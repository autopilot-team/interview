package store

import (
	"autopilot/backends/api/internal/identity/model"
	"autopilot/backends/internal/core"
	"context"
	"database/sql"
)

// Userer is the store for user operations.
type Userer interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	CreateVerification(ctx context.Context, verification *model.Verification) (*model.Verification, error)
	DeleteVerification(ctx context.Context, id string) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetVerification(ctx context.Context, context string, id string) (*model.Verification, error)
	GetVerificationByValue(ctx context.Context, context string, value string) (*model.Verification, error)
	Update(ctx context.Context, user *model.User) error
	WithQuerier(q core.Querier) Userer
}

// User is the store for user operations.
type User struct {
	core.Querier
}

func (s *User) WithQuerier(q core.Querier) Userer {
	return &User{q}
}

// NewUser creates a new User.
func NewUser(db core.Querier) Userer {
	return &User{db}
}

// Create creates a new user.
func (s *User) Create(ctx context.Context, user *model.User) (*model.User, error) {
	query := `
		INSERT INTO users (
			name, email, email_verified_at, failed_login_attempts,
			image, last_active_at, last_logged_in_at, locked_at,
			password_changed_at, password_hash
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8,
			$9, $10
		) RETURNING
			id, name, email, email_verified_at, failed_login_attempts,
			image, last_active_at, last_logged_in_at, locked_at,
			password_changed_at, password_hash, created_at, updated_at
	`

	var created model.User
	err := s.QueryRowContext(
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
	).Scan(
		&created.ID,
		&created.Name,
		&created.Email,
		&created.EmailVerifiedAt,
		&created.FailedLoginAttempts,
		&created.Image,
		&created.LastActiveAt,
		&created.LastLoggedInAt,
		&created.LockedAt,
		&created.PasswordChangedAt,
		&created.PasswordHash,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

// CreateVerification creates a new verification.
func (s *User) CreateVerification(ctx context.Context, verification *model.Verification) (*model.Verification, error) {
	query := `
		INSERT INTO verifications (
			context, value, expires_at
		) VALUES (
			$1, $2, $3
		) RETURNING
			id, context, value, expires_at, created_at, updated_at
	`

	var created model.Verification
	err := s.QueryRowContext(
		ctx,
		query,
		verification.Context,
		verification.Value,
		verification.ExpiresAt,
	).Scan(
		&created.ID,
		&created.Context,
		&created.Value,
		&created.ExpiresAt,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

// DeleteVerification deletes a verification.
func (s *User) DeleteVerification(ctx context.Context, id string) error {
	query := `DELETE FROM verifications WHERE id = $1`

	result, err := s.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return err
	}

	return nil
}

// ExistsByEmail checks if a user exists by email.
func (s *User) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	err := s.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetByID gets a user by ID.
func (s *User) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE id = $1`

	err := s.QueryRowContext(ctx, query, id).Scan(
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
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByEmail gets a user by email.
func (s *User) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	query := `SELECT
			id, name, email, email_verified_at, failed_login_attempts,
			image, last_active_at, last_logged_in_at, locked_at,
			password_changed_at, password_hash, created_at, updated_at
		FROM
			users
		WHERE
			email = $1`

	err := s.QueryRowContext(ctx, query, email).Scan(
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
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates a user.
func (s *User) Update(ctx context.Context, user *model.User) error {
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
			updated_at = $11
		WHERE id = $12
	`

	result, err := s.ExecContext(
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
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return err
	}

	return nil
}

// GetVerification gets a verification by ID and context.
func (s *User) GetVerification(ctx context.Context, context string, id string) (*model.Verification, error) {
	var verification model.Verification
	query := `SELECT
			*
		FROM
			verifications
		WHERE
			context = $1 AND id = $2`

	err := s.QueryRowContext(ctx, query, context, id).Scan(
		&verification.ID,
		&verification.Context,
		&verification.Value,
		&verification.ExpiresAt,
		&verification.CreatedAt,
		&verification.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &verification, nil
}

// GetVerificationByValue gets a verification by value and context.
func (s *User) GetVerificationByValue(ctx context.Context, context string, value string) (*model.Verification, error) {
	query := `SELECT
			*
		FROM
			verifications
		WHERE
			context = $1 AND value = $2`

	var verification model.Verification
	err := s.QueryRowContext(ctx, query, context, value).Scan(
		&verification.ID,
		&verification.Context,
		&verification.Value,
		&verification.ExpiresAt,
		&verification.CreatedAt,
		&verification.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &verification, nil
}
