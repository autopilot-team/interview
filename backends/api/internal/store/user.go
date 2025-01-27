package store

import (
	"autopilot/backends/api/internal/app"
	"autopilot/backends/api/internal/model"
	"autopilot/backends/internal/core"
	"context"
	"database/sql"
	"fmt"
)

// User is the store for user operations.
type User struct {
	core.DBer
}

// NewUser creates a new User.
func NewUser(container *app.Container) *User {
	return &User{
		container.DB.Primary,
	}
}

// Create creates a new user.
func (s *User) Create(ctx context.Context, tx *sql.Tx, user *model.User) error {
	query := `
		INSERT INTO users (
			id, name, email, email_verified_at, failed_login_attempts,
			image, last_active_at, last_logged_in_at, locked_at,
			password_changed_at, password_hash, two_factor_enabled,
			created_at, updated_at
		) VALUES (
			uuid7(), $1, $2, $3, $4,
			$5, $6, $7, $8,
			$9, $10, $11,
			$12, $13
		) RETURNING id
	`

	err := tx.QueryRowContext(
		ctx,
		query,
		user.Name,
		user.Email,
		user.EmailVerifiedAt,
		user.FailedLoginAttempts,
		user.Image,
		user.LastActiveAt,
		user.LastLoggedInAt,
		user.LockedAt,
		user.PasswordChangedAt,
		user.PasswordHash,
		user.TwoFactorEnabled,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("inserting user: %w", err)
	}

	return nil
}

// Update updates a user.
func (s *User) Update(ctx context.Context, tx *sql.Tx, user *model.User) error {
	query := `
		UPDATE users SET
			name = $1,
			email = $2,
			email_verified_at = $3,
			failed_login_attempts = $4,
			image = $5,
			last_active_at = $6,
			last_logged_in_at = $7,
			locked_at = $8,
			password_changed_at = $9,
			password_hash = $10,
			two_factor_enabled = $11,
			updated_at = $12
		WHERE id = $13
	`

	result, err := tx.ExecContext(
		ctx,
		query,
		user.Name,
		user.Email,
		user.EmailVerifiedAt,
		user.FailedLoginAttempts,
		user.Image,
		user.LastActiveAt,
		user.LastLoggedInAt,
		user.LockedAt,
		user.PasswordChangedAt,
		user.PasswordHash,
		user.TwoFactorEnabled,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("updating user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// GetByEmail gets a user by email.
func (s *User) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	query := `
		SELECT
			id, name, email, email_verified_at, failed_login_attempts,
			image, last_active_at, last_logged_in_at, locked_at,
			password_changed_at, password_hash, two_factor_enabled,
			created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := s.Writer().QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.EmailVerifiedAt,
		&user.FailedLoginAttempts,
		&user.Image,
		&user.LastActiveAt,
		&user.LastLoggedInAt,
		&user.LockedAt,
		&user.PasswordChangedAt,
		&user.PasswordHash,
		&user.TwoFactorEnabled,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getting user by email: %w", err)
	}

	return &user, nil
}

// CreateVerification creates a new verification.
func (s *User) CreateVerification(ctx context.Context, tx *sql.Tx, verification *model.Verification) error {
	query := `
		INSERT INTO verifications (
			id, value, expires_at, created_at, updated_at
		) VALUES (
			uuid7(), $1, $2, $3, $4
		) RETURNING id
	`

	err := tx.QueryRowContext(
		ctx,
		query,
		verification.Value,
		verification.ExpiresAt,
		verification.CreatedAt,
		verification.UpdatedAt,
	).Scan(&verification.ID)

	if err != nil {
		return fmt.Errorf("inserting verification: %w", err)
	}

	return nil
}

// GetVerification gets a verification by ID.
func (s *User) GetVerification(ctx context.Context, id string) (*model.Verification, error) {
	var verification model.Verification
	query := `
		SELECT
			id, value, expires_at, created_at, updated_at
		FROM verifications
		WHERE id = $1
	`

	err := s.Writer().QueryRowContext(ctx, query, id).Scan(
		&verification.ID,
		&verification.Value,
		&verification.ExpiresAt,
		&verification.CreatedAt,
		&verification.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("getting verification: %w", err)
	}

	return &verification, nil
}

// DeleteVerification deletes a verification.
func (s *User) DeleteVerification(ctx context.Context, tx *sql.Tx, id string) error {
	query := `DELETE FROM verifications WHERE id = $1`

	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting verification: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("verification not found")
	}

	return nil
}

// ExistsByEmail checks if a user exists by email.
func (s *User) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	err := s.Writer().QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking email existence: %w", err)
	}

	return exists, nil
}
