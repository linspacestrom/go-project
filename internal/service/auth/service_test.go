package auth

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/linspacestrom/go-project/internal/domain"
	"github.com/linspacestrom/go-project/internal/repository/postgres"
)

type fakeTokenGen struct{}

func (fakeTokenGen) Generate(userID uuid.UUID, role string) (string, error) {
	return userID.String() + ":" + role, nil
}

type fakeRepo struct {
	users        map[string]*domain.UserCredentials
	refresh      map[string]*postgres.RefreshToken
	createUserFn func(params postgres.CreateUserParams) *domain.User
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{users: map[string]*domain.UserCredentials{}, refresh: map[string]*postgres.RefreshToken{}}
}

func (f *fakeRepo) CheckUserExistsByEmail(_ context.Context, email string) (bool, error) {
	_, ok := f.users[email]
	return ok, nil
}

func (f *fakeRepo) CreateUser(_ context.Context, params postgres.CreateUserParams) (*domain.User, error) {
	user := &domain.User{ID: uuid.New(), Email: params.Email, Role: params.Role, FullName: params.FullName, IsActive: true, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
	f.users[params.Email] = &domain.UserCredentials{User: *user, PasswordHash: params.PasswordHash}
	return user, nil
}

func (f *fakeRepo) CreateStudentProfile(_ context.Context, _ uuid.UUID, _ string, _ int, _ string) error {
	return nil
}

func (f *fakeRepo) CreateMentorProfile(_ context.Context, _ uuid.UUID, _, _ *string) error {
	return nil
}

func (f *fakeRepo) GetUserByEmail(_ context.Context, email string) (*domain.UserCredentials, error) {
	if u, ok := f.users[email]; ok {
		return u, nil
	}
	return nil, domain.ErrInvalidCredentials
}

func (f *fakeRepo) GetUserByID(_ context.Context, userID uuid.UUID) (*domain.User, error) {
	for _, u := range f.users {
		if u.ID == userID {
			copy := u.User
			return &copy, nil
		}
	}
	return nil, domain.ErrUnauthorized
}

func (f *fakeRepo) SaveRefreshToken(_ context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	f.refresh[token] = &postgres.RefreshToken{UserID: userID, Token: token, ExpiresAt: expiresAt}
	return nil
}

func (f *fakeRepo) GetRefreshToken(_ context.Context, token string) (*postgres.RefreshToken, error) {
	if rt, ok := f.refresh[token]; ok {
		return rt, nil
	}
	return nil, domain.ErrInvalidCredentials
}

func (f *fakeRepo) RevokeRefreshToken(_ context.Context, token string) error {
	if rt, ok := f.refresh[token]; ok {
		now := time.Now().UTC()
		rt.RevokedAt = &now
	}
	return nil
}

func (f *fakeRepo) CreateOutboxEvent(_ context.Context, _ string, _ *uuid.UUID, _ string, _ []byte) error {
	return nil
}

func TestRegisterAndLogin(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(fakeTokenGen{}, repo, 24*time.Hour)

	user, access, refresh, err := svc.Register(context.Background(), RegisterParams{
		Email:    "student@example.com",
		Password: "password123",
		FullName: "Student",
	})
	if err != nil {
		t.Fatalf("register returned error: %v", err)
	}
	if user.Role != domain.RoleStudent {
		t.Fatalf("expected default role student, got %s", user.Role)
	}
	if access == "" || refresh == "" {
		t.Fatalf("expected non-empty tokens")
	}

	loginAccess, loginRefresh, err := svc.Login(context.Background(), "student@example.com", "password123")
	if err != nil {
		t.Fatalf("login returned error: %v", err)
	}
	if loginAccess == "" || loginRefresh == "" {
		t.Fatalf("expected non-empty login tokens")
	}
}

func TestRefreshRotatesToken(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(fakeTokenGen{}, repo, 24*time.Hour)
	_, _, refresh, err := svc.Register(context.Background(), RegisterParams{
		Email:    "mentor@example.com",
		Password: "password123",
		FullName: "Mentor",
		Role:     domain.RoleMentor,
	})
	if err != nil {
		t.Fatalf("register returned error: %v", err)
	}

	access2, refresh2, err := svc.Refresh(context.Background(), refresh)
	if err != nil {
		t.Fatalf("refresh returned error: %v", err)
	}
	if access2 == "" || refresh2 == "" || refresh2 == refresh {
		t.Fatalf("expected rotated refresh token")
	}
}
