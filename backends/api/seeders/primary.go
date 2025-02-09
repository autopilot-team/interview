package seeders

import (
	"autopilot/backends/internal/core"
	"context"
)

// Primary seeds the primary database
func Primary(ctx context.Context, db core.DBer) error {
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
