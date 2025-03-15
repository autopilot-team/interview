package store

import (
	"autopilot/backends/api/internal/identity/model"
	"autopilot/backends/internal/core"
	"context"
	"database/sql"
)

// Verificationer is the store for verification operations.
type Verificationer interface {
	Delete(ctx context.Context, id string) error
	GetByValue(ctx context.Context, context string, value string) (*model.Verification, error)
	WithQuerier(q core.Querier) Verificationer
}

// Verification is the store for verification operations.
type Verification struct {
	core.Querier
}

func (v Verification) WithQuerier(q core.Querier) Verificationer {
	return &Verification{q}
}

// NewVerification creates a new Verification.
func NewVerification(db core.Querier) *Verification {
	return &Verification{db}
}

// GetByValue gets a verification by value.
func (s *Verification) GetByValue(ctx context.Context, context string, value string) (*model.Verification, error) {
	var verification model.Verification
	query := `SELECT * FROM verifications WHERE context = $1 AND value = $2`

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

// Delete deletes a verification by ID
func (s *Verification) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM verifications WHERE id = $1`

	_, err := s.ExecContext(ctx, query, id)
	return err
}
