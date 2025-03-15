package store

import (
	"autopilot/backends/api/internal/identity/model"
	"autopilot/backends/internal/core"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// TwoFactorer is the store for two-factor authentication operations.
type TwoFactorer interface {
	Create(ctx context.Context, twoFactor *model.TwoFactor) (*model.TwoFactor, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*model.TwoFactor, error)
	GetByUserID(ctx context.Context, userID string) (*model.TwoFactor, error)
	Update(ctx context.Context, twoFactor *model.TwoFactor) error
	WithQuerier(q core.Querier) TwoFactorer
}

// TwoFactor implements TwoFactorer interface
type TwoFactor struct {
	core.Querier
}

func (s *TwoFactor) WithQuerier(q core.Querier) TwoFactorer {
	return &TwoFactor{q}
}

// NewTwoFactor creates a new TwoFactor store
func NewTwoFactor(db core.Querier) *TwoFactor {
	return &TwoFactor{db}
}

// Create creates a new two-factor authentication record
func (s *TwoFactor) Create(ctx context.Context, twoFactor *model.TwoFactor) (*model.TwoFactor, error) {
	query := `
		INSERT INTO two_factors (
			backup_codes, secret, user_id, enabled_at
		) VALUES (
			$1, $2, $3, $4
		) RETURNING
			id, backup_codes, secret, user_id, enabled_at, created_at, updated_at
	`

	var (
		backupCodesJSON []byte // temporary holder for JSONB data
		created         model.TwoFactor
	)
	err := s.QueryRowContext(
		ctx,
		query,
		twoFactor.BackupCodes,
		twoFactor.Secret,
		twoFactor.UserID,
		twoFactor.EnabledAt,
	).Scan(
		&created.ID,
		&backupCodesJSON,
		&created.Secret,
		&created.UserID,
		&created.EnabledAt,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSONB data into the BackupCodes slice
	if err := json.Unmarshal(backupCodesJSON, &created.BackupCodes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal backup codes: %w", err)
	}

	return &created, nil
}

// Delete deletes a two-factor authentication record
func (s *TwoFactor) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM two_factors WHERE id = $1"

	_, err := s.ExecContext(ctx, query, id)
	return err
}

// GetByID retrieves a two-factor authentication record by ID
func (s *TwoFactor) GetByID(ctx context.Context, id string) (*model.TwoFactor, error) {
	var twoFactor model.TwoFactor
	query := "SELECT * FROM two_factors WHERE id = $1"

	err := s.QueryRowContext(ctx, query, id).Scan(
		&twoFactor.ID,
		&twoFactor.BackupCodes,
		&twoFactor.Secret,
		&twoFactor.UserID,
		&twoFactor.EnabledAt,
		&twoFactor.CreatedAt,
		&twoFactor.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &twoFactor, nil
}

// GetByUserID retrieves a two-factor authentication record by user ID
func (s *TwoFactor) GetByUserID(ctx context.Context, userID string) (*model.TwoFactor, error) {
	var twoFactor model.TwoFactor
	var backupCodesJSON []byte // temporary holder for JSONB data

	err := s.QueryRowContext(
		ctx,
		"SELECT * FROM two_factors WHERE user_id = $1",
		userID,
	).Scan(
		&twoFactor.ID,
		&backupCodesJSON, // scan JSONB into []byte first
		&twoFactor.Secret,
		&twoFactor.UserID,
		&twoFactor.FailedAttempts,
		&twoFactor.LastFailedAttemptAt,
		&twoFactor.LockedUntil,
		&twoFactor.EnabledAt,
		&twoFactor.CreatedAt,
		&twoFactor.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	// Unmarshal the JSONB data into the BackupCodes slice
	if err := json.Unmarshal(backupCodesJSON, &twoFactor.BackupCodes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal backup codes: %w", err)
	}

	return &twoFactor, nil
}

// Update updates a two-factor authentication record
func (s *TwoFactor) Update(ctx context.Context, twoFactor *model.TwoFactor) error {
	twoFactor.UpdatedAt = time.Now()
	query := `
		UPDATE two_factors
		SET backup_codes = $1,
			secret = $2,
			enabled_at = $3,
			updated_at = $4
		WHERE id = $5
	`

	_, err := s.ExecContext(
		ctx,
		query,
		twoFactor.BackupCodes,
		twoFactor.Secret,
		twoFactor.EnabledAt,
		twoFactor.UpdatedAt,
		twoFactor.ID,
	)

	return err
}
