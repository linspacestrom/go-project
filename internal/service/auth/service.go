package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/linspacestrom/go-project/internal/domain"
	"github.com/linspacestrom/go-project/internal/repository/postgres"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CheckUserExistsByEmail(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, params postgres.CreateUserParams) (*domain.User, error)
	CreateStudentProfile(ctx context.Context, userID uuid.UUID, university string, course int, degreeType string) error
	CreateMentorProfile(ctx context.Context, userID uuid.UUID, description, title *string) error
	GetUserByEmail(ctx context.Context, email string) (*domain.UserCredentials, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	SaveRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, token string) (*postgres.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	CreateOutboxEvent(ctx context.Context, aggregateType string, aggregateID *uuid.UUID, eventType string, payload []byte) error
}

type TokenGenerator interface {
	Generate(userID uuid.UUID, role string) (string, error)
}

type Service struct {
	tokenGen        TokenGenerator
	userRepo        UserRepository
	refreshTokenTTL time.Duration
}

type RegisterParams struct {
	Email       string
	Password    string
	Role        string
	FullName    string
	BirthDate   *time.Time
	University  string
	Course      int
	DegreeType  string
	CityID      *uuid.UUID
	Description *string
	Title       *string
}

func NewService(tokenGen TokenGenerator, userRepo UserRepository, refreshTokenTTL time.Duration) *Service {
	return &Service{tokenGen: tokenGen, userRepo: userRepo, refreshTokenTTL: refreshTokenTTL}
}

func (s *Service) Register(ctx context.Context, params RegisterParams) (*domain.User, string, string, error) {
	exists, err := s.userRepo.CheckUserExistsByEmail(ctx, params.Email)
	if err != nil {
		return nil, "", "", err
	}
	if exists {
		return nil, "", "", domain.ErrEmailExists
	}
	if params.Role == "" {
		params.Role = domain.RoleStudent
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", "", fmt.Errorf("hashing password: %w", err)
	}

	created, err := s.userRepo.CreateUser(ctx, postgres.CreateUserParams{
		Email:        params.Email,
		PasswordHash: string(hash),
		Role:         params.Role,
		FullName:     params.FullName,
		BirthDate:    params.BirthDate,
		CityID:       params.CityID,
	})
	if err != nil {
		return nil, "", "", err
	}

	switch params.Role {
	case domain.RoleStudent:
		if params.University != "" && params.Course > 0 && params.DegreeType != "" {
			if err = s.userRepo.CreateStudentProfile(ctx, created.ID, params.University, params.Course, params.DegreeType); err != nil {
				return nil, "", "", err
			}
		}
	case domain.RoleMentor:
		if err = s.userRepo.CreateMentorProfile(ctx, created.ID, params.Description, params.Title); err != nil {
			return nil, "", "", err
		}
	}

	accessToken, err := s.tokenGen.Generate(created.ID, created.Role)
	if err != nil {
		return nil, "", "", err
	}
	refreshToken := xid.New().String()
	expiresAt := time.Now().UTC().Add(s.refreshTokenTTL)
	if err = s.userRepo.SaveRefreshToken(ctx, created.ID, refreshToken, expiresAt); err != nil {
		return nil, "", "", err
	}

	eventType := domain.EventUserRegistered
	if params.Role == domain.RoleMentor {
		eventType = domain.EventMentorRegistered
	}
	_ = s.userRepo.CreateOutboxEvent(ctx, "user", &created.ID, eventType, []byte(fmt.Sprintf(`{"user_id":"%s","email":"%s"}`, created.ID, created.Email)))

	return created, accessToken, refreshToken, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, string, error) {
	creds, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(creds.PasswordHash), []byte(password)); err != nil {
		return "", "", domain.ErrInvalidCredentials
	}

	accessToken, err := s.tokenGen.Generate(creds.ID, creds.Role)
	if err != nil {
		return "", "", err
	}

	refreshToken := xid.New().String()
	expiresAt := time.Now().UTC().Add(s.refreshTokenTTL)
	if err = s.userRepo.SaveRefreshToken(ctx, creds.ID, refreshToken, expiresAt); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *Service) Refresh(ctx context.Context, token string) (string, string, error) {
	rt, err := s.userRepo.GetRefreshToken(ctx, token)
	if err != nil {
		return "", "", err
	}
	if rt.RevokedAt != nil || time.Now().UTC().After(rt.ExpiresAt) {
		return "", "", domain.ErrUnauthorized
	}
	user, err := s.userRepo.GetUserByID(ctx, rt.UserID)
	if err != nil {
		return "", "", err
	}
	accessToken, err := s.tokenGen.Generate(user.ID, user.Role)
	if err != nil {
		return "", "", err
	}
	newRefresh := xid.New().String()
	if err = s.userRepo.RevokeRefreshToken(ctx, token); err != nil {
		return "", "", err
	}
	if err = s.userRepo.SaveRefreshToken(ctx, user.ID, newRefresh, time.Now().UTC().Add(s.refreshTokenTTL)); err != nil {
		return "", "", err
	}

	return accessToken, newRefresh, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	if token == "" {
		return errors.New("empty refresh token")
	}

	return s.userRepo.RevokeRefreshToken(ctx, token)
}
