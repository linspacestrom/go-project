package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/linspacestrom/go-project/internal/domain"
	"github.com/linspacestrom/go-project/internal/repository/postgres"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) GetMe(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *Service) UpdateMe(ctx context.Context, userID uuid.UUID, fullName *string, birthDate *time.Time) (*domain.User, error) {
	return s.repo.UpdateUserProfile(ctx, userID, fullName, birthDate)
}

func (s *Service) RegisterMentor(ctx context.Context, email, password, fullName string, cityID uuid.UUID, description, title *string) (*domain.User, error) {
	exists, err := s.repo.CheckUserExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrEmailExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash mentor password: %w", err)
	}

	user, err := s.repo.CreateUser(ctx, postgres.CreateUserParams{
		Email:        email,
		PasswordHash: string(hash),
		Role:         domain.RoleMentor,
		FullName:     fullName,
		CityID:       &cityID,
	})
	if err != nil {
		return nil, err
	}
	if err = s.repo.CreateMentorProfile(ctx, user.ID, description, title); err != nil {
		return nil, err
	}
	if err = s.publishEvent(ctx, "user", &user.ID, domain.EventMentorRegistered, map[string]any{
		"user_id": user.ID,
		"email":   user.Email,
		"city_id": cityID,
	}); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) UpdateUserCity(ctx context.Context, userID, cityID uuid.UUID) error {
	if err := s.repo.UpdateUserCity(ctx, userID, cityID); err != nil {
		return err
	}

	payload := map[string]any{"user_id": userID, "city_id": cityID}
	if err := s.publishEvent(ctx, "user", &userID, domain.EventAdminCityChanged, payload); err != nil {
		return err
	}

	return nil
}
