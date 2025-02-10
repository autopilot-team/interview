package seeders

import (
	"autopilot/backends/api/internal/model"
	"autopilot/backends/api/internal/store"
	"autopilot/backends/internal/core"
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
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

	// Define users to be created with roles
	password := "Strongpa$$w0rd!"
	users := []struct {
		name     string
		email    string
		password string
	}{
		{
			name:     "Financial Administrator",
			email:    "admin@example.com",
			password: password,
		},
	}

	storeManager := store.NewManager(db)

	// Create users
	for _, u := range users {
		// Hash password
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.password), model.PasswordHashBcryptCostBcryptCost)
		if err != nil {
			return err
		}

		now := time.Now()
		passwordHashStr := string(passwordHash)
		user := &model.User{
			Name:            u.name,
			Email:           u.email,
			EmailVerifiedAt: &now,
			PasswordHash:    &passwordHashStr,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		if err = storeManager.User.Create(ctx, user, tx); err != nil {
			return err
		}
	}

	return nil
}
