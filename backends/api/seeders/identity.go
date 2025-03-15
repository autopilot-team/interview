package seeders

import (
	"autopilot/backends/api/internal/identity/model"
	"autopilot/backends/api/internal/identity/store"
	"autopilot/backends/internal/core"
	"autopilot/backends/internal/types"
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Identity seeds the identity database
func Identity(ctx context.Context, db core.DBer) error {
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

	storeManager := store.NewManager(tx)

	type seedUser struct {
		id    string
		name  string
		email string
		role  types.Role
	}

	type seedData struct {
		entity   model.Entity
		users    []seedUser
		children []seedData
	}

	seeds := []seedData{
		{
			entity: model.Entity{
				Name:   "Acme",
				Slug:   "acme",
				Status: model.EntityStatusActive,
				Type:   model.EntityTypeAccount,
			},
			users: []seedUser{
				{
					name:  "Administrator",
					email: "admin@acme.com",
					role:  types.RoleAdmin,
				},
			},
		},
	}

	// Create entities and their associated users recursively
	var createEntity func(seed seedData, parentID *string) error
	password := "Strongpa$$w0rd!"

	createEntity = func(seed seedData, parentID *string) error {
		// Set parent ID if provided
		seed.entity.ParentID = parentID

		// Create entity
		newEntity, err := storeManager.Entity.Create(ctx, &seed.entity)
		if err != nil {
			return err
		}
		seed.entity = *newEntity

		// Create associated users
		for i, u := range seed.users {
			passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), model.PasswordHashBcryptCost)
			if err != nil {
				return err
			}

			now := time.Now()
			passwordHashStr := string(passwordHash)
			user, err := storeManager.User.Create(ctx, &model.User{
				Name:            u.name,
				Email:           u.email,
				EmailVerifiedAt: &now,
				PasswordHash:    &passwordHashStr,
			})
			if err != nil {
				return err
			}
			seed.users[i].id = user.ID

			// Create member record to associate user with entity
			_, err = storeManager.Membership.Create(ctx, &model.Membership{
				EntityID: &seed.entity.ID,
				UserID:   user.ID,
				Role:     u.role,
			})
			if err != nil {
				return err
			}
		}
		return nil
	}

	// Create all entities with their hierarchies
	for _, seed := range seeds {
		if err := createEntity(seed, nil); err != nil {
			return err
		}
	}

	return nil
}
