package auth

import (
	"context"
	"database/sql"
	"errors"
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type authRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &authRepository{db}
}

func (r *authRepository) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, role_id, branch_id, name, email, password_hash, phone)
		VALUES (?, (SELECT id FROM roles WHERE name = 'Customer' LIMIT 1), ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.BranchID, user.Name, user.Email, user.PasswordHash, user.Phone)
	return err
}

func (r *authRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, role_id, branch_id, name, email, password_hash, phone, profile_picture, is_active, created_at, updated_at
		FROM users WHERE email = ?
	`
	row := r.db.QueryRowContext(ctx, query, email)

	var u User
	err := row.Scan(
		&u.ID, &u.RoleID, &u.BranchID, &u.Name, &u.Email, &u.PasswordHash,
		&u.Phone, &u.ProfilePicture, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &u, nil
}
