package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/linspacestrom/go-project/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CheckUserExistsByEmail(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, email, passwordHash, role string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.UserCredentials, error)
}

type TokenGenerator interface {
	Generate(userID uuid.UUID, role string) (string, error)
}

type Service struct {
	tokenGen TokenGenerator
	userRepo UserRepository
}

func NewService(tokenGen TokenGenerator, userRepo UserRepository) *Service {
	return &Service{tokenGen: tokenGen, userRepo: userRepo}
}

func (s *Service) Register(ctx context.Context, email, password, role string) (*domain.User, error) {
	exists, err := s.userRepo.CheckUserExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrEmailExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	return s.userRepo.CreateUser(ctx, email, string(hash), role)
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	creds, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(creds.PasswordHash), []byte(password)); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	return s.tokenGen.Generate(creds.ID, creds.Role)
}
