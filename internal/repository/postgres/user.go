package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/linspacestrom/go-project/internal/domain"
)

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

func (r *Repository) CreateUser(ctx context.Context, email, passwordHash, role string) (*domain.User, error) {
	query, args, err := psql.Insert("users").
		Columns("email", "password_hash", "role_id").
		Values(
			email,
			passwordHash,
			squirrel.Expr("(SELECT id FROM user_roles WHERE name = ?)", role),
		).
		Suffix("RETURNING id, created_at").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build insert query: %w", err)
	}

	user := &domain.User{Email: email, Role: role}

	err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)
	if err == nil {
		user.CreatedAt = user.CreatedAt.UTC()
	}
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}

	return user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.UserCredentials, error) {
	query, args, err := psql.Select("u.id", "u.email", "u.password_hash", "r.name", "u.created_at").
		From("users u").
		Join("user_roles r ON u.role_id = r.id").
		Where(squirrel.Eq{"u.email": email}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build login query: %w", err)
	}

	var u domain.UserCredentials
	err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err == nil {
		u.CreatedAt = u.CreatedAt.UTC()
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrInvalidCredentials
		}

		return nil, fmt.Errorf("scan user credentials: %w", err)
	}

	return &u, nil
}
