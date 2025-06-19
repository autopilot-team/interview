package store

import (
	"autopilot/backends/api/internal/payment/model"
	"autopilot/backends/internal/core"
	"context"
)

type Refunder interface {
	Create(ctx context.Context, refund *model.Refund) (*model.Refund, error)
}

type Refund struct {
	core.Querier
}

func NewRefund(q core.Querier) Refunder {
	return &Refund{q}
}

func (s *Refund) WithQuerier(q core.Querier) Refunder {
	return &Refund{q}
}

func (s *Refund) Create(ctx context.Context, refund *model.Refund) (*model.Refund, error) {
	query := `
		INSERT INTO refunds (
		) VALUES (
		 )
		RETURNING
			id, created_at, updated_at
	`
	created := *refund
	err := s.QueryRowContext(
		ctx,
		query,
	).Scan(
		&created.ID,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &created, nil
}
