package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/linspacestrom/go-project/internal/domain"
)

type CreateUserParams struct {
	Email        string
	PasswordHash string
	Role         string
	FullName     string
	BirthDate    *time.Time
	CityID       *uuid.UUID
}

func (r *Repository) CheckUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	query, args, err := psql.Select("1").
		Prefix("SELECT EXISTS (").
		From("users").
		Where(squirrel.Eq{"email": email}).
		Suffix(")").
		ToSql()

	if err != nil {
		return false, fmt.Errorf("build check user query: %w", err)
	}

	var exists bool
	err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

func (r *Repository) CreateUser(ctx context.Context, params CreateUserParams) (*domain.User, error) {
	roleExpr := squirrel.Expr(
		"COALESCE((SELECT id FROM user_roles WHERE name = ?), (SELECT id FROM user_roles WHERE name = 'user' AND ? = 'student'))",
		params.Role,
		params.Role,
	)

	query, args, err := psql.Insert("users").
		Columns("email", "password_hash", "role_id", "full_name", "birth_date", "city_id").
		Values(
			params.Email,
			params.PasswordHash,
			roleExpr,
			params.FullName,
			params.BirthDate,
			params.CityID,
		).
		Suffix("RETURNING id, email, full_name, birth_date, city_id, created_at, updated_at, is_active").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build insert query: %w", err)
	}

	user := &domain.User{Role: params.Role}
	if err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.FullName,
		&user.BirthDate,
		&user.CityID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
	); err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}

	return user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.UserCredentials, error) {
	query, args, err := psql.Select(
		"u.id", "u.email", "u.password_hash", "u.full_name", "u.birth_date", "u.city_id", "u.is_active", "u.created_at", "u.updated_at",
		"CASE WHEN r.name = 'user' THEN 'student' ELSE r.name END",
	).From("users u").
		Join("user_roles r ON u.role_id = r.id").
		Where(squirrel.Eq{"u.email": email}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build login query: %w", err)
	}

	var u domain.UserCredentials
	err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.FullName,
		&u.BirthDate,
		&u.CityID,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.Role,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrInvalidCredentials
		}

		return nil, fmt.Errorf("scan user credentials: %w", err)
	}

	return &u, nil
}

func (r *Repository) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	query, args, err := psql.Select(
		"u.id", "u.email", "u.full_name", "u.birth_date", "u.city_id", "u.is_active", "u.created_at", "u.updated_at",
		"CASE WHEN r.name = 'user' THEN 'student' ELSE r.name END",
	).From("users u").
		Join("user_roles r ON u.role_id = r.id").
		Where(squirrel.Eq{"u.id": userID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get user by id query: %w", err)
	}

	var u domain.User
	err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(
		&u.ID,
		&u.Email,
		&u.FullName,
		&u.BirthDate,
		&u.CityID,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.Role,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUnauthorized
		}

		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &u, nil
}

func (r *Repository) UpdateUserProfile(ctx context.Context, userID uuid.UUID, fullName *string, birthDate *time.Time) (*domain.User, error) {
	builder := psql.Update("users").
		Set("updated_at", squirrel.Expr("now()"))
	if fullName != nil {
		builder = builder.Set("full_name", *fullName)
	}
	if birthDate != nil {
		builder = builder.Set("birth_date", *birthDate)
	}

	query, args, err := builder.
		Where(squirrel.Eq{"id": userID}).
		Suffix("RETURNING id, email, full_name, birth_date, city_id, is_active, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build update profile query: %w", err)
	}

	var u domain.User
	err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(
		&u.ID,
		&u.Email,
		&u.FullName,
		&u.BirthDate,
		&u.CityID,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUnauthorized
		}

		return nil, fmt.Errorf("update profile: %w", err)
	}

	roleQ, roleArgs, err := psql.Select("CASE WHEN r.name = 'user' THEN 'student' ELSE r.name END").
		From("users u").
		Join("user_roles r ON r.id=u.role_id").
		Where(squirrel.Eq{"u.id": userID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build role query: %w", err)
	}
	if err = r.GetConn(ctx).QueryRow(ctx, roleQ, roleArgs...).Scan(&u.Role); err != nil {
		return nil, fmt.Errorf("scan role: %w", err)
	}

	return &u, nil
}

func (r *Repository) UpdateUserCity(ctx context.Context, userID, cityID uuid.UUID) error {
	query, args, err := psql.Update("users").
		Set("city_id", cityID).
		Set("updated_at", squirrel.Expr("now()")).
		Where(squirrel.Eq{"id": userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build update city query: %w", err)
	}

	tag, err := r.GetConn(ctx).Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update user city: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrUnauthorized
	}

	return nil
}

func (r *Repository) CreateStudentProfile(ctx context.Context, userID uuid.UUID, university string, course int, degreeType string) error {
	query, args, err := psql.Insert("student_profiles").
		Columns("user_id", "university", "course", "degree_type").
		Values(userID, university, course, degreeType).
		ToSql()
	if err != nil {
		return fmt.Errorf("build create student profile query: %w", err)
	}

	if _, err = r.GetConn(ctx).Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("create student profile: %w", err)
	}

	return nil
}

func (r *Repository) CreateMentorProfile(ctx context.Context, userID uuid.UUID, description, title *string) error {
	query, args, err := psql.Insert("mentor_profiles").
		Columns("user_id", "description", "title", "is_verified").
		Values(userID, description, title, true).
		ToSql()
	if err != nil {
		return fmt.Errorf("build create mentor profile query: %w", err)
	}

	if _, err = r.GetConn(ctx).Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("create mentor profile: %w", err)
	}

	return nil
}

type RefreshToken struct {
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	RevokedAt *time.Time
}

func (r *Repository) SaveRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	query, args, err := psql.Insert("refresh_tokens").
		Columns("user_id", "token", "expires_at").
		Values(userID, token, expiresAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("build save refresh token query: %w", err)
	}
	if _, err = r.GetConn(ctx).Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("save refresh token: %w", err)
	}

	return nil
}

func (r *Repository) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	query, args, err := psql.Select("user_id", "token", "expires_at", "revoked_at").
		From("refresh_tokens").
		Where(squirrel.Eq{"token": token}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get refresh token query: %w", err)
	}

	var rt RefreshToken
	err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(&rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.RevokedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrInvalidCredentials
		}

		return nil, fmt.Errorf("get refresh token: %w", err)
	}

	return &rt, nil
}

func (r *Repository) RevokeRefreshToken(ctx context.Context, token string) error {
	query, args, err := psql.Update("refresh_tokens").
		Set("revoked_at", squirrel.Expr("now()")).
		Where(squirrel.Eq{"token": token}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build revoke refresh token query: %w", err)
	}
	if _, err = r.GetConn(ctx).Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}

	return nil
}
