package seeders

import (
	"autopilot/backends/internal/core"
	"context"
)

// Payment seeds the payment database
func Payment(ctx context.Context, db core.DBer) error {
	tx, err := db.Writer().Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		_ = tx.Commit()
	}()

	return nil
}
