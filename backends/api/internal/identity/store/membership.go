package store

import (
	"autopilot/backends/api/internal/identity/model"
	"autopilot/backends/internal/core"
	"context"
)

// Membershiper defines the interface for membership store operations
type Membershiper interface {
	Create(ctx context.Context, membership *model.Membership) (*model.Membership, error)
	GetByUserID(ctx context.Context, userID string) ([]*model.Membership, error)
	GetByEntityID(ctx context.Context, entityID string) ([]*model.Membership, error)
	GetByEntityIDWithInheritance(ctx context.Context, userID, entityID string) ([]*model.Membership, error)
	GetByUserIDWithInheritance(ctx context.Context, userID string) ([]*model.Membership, error)
	WithQuerier(q core.Querier) Membershiper
}

func (s *Membership) WithQuerier(q core.Querier) Membershiper {
	return &Membership{q}
}

// Membership implements Memberer interface
type Membership struct {
	core.Querier
}

// NewMembership creates a new membership store
func NewMembership(db core.Querier) Membershiper {
	return &Membership{
		db,
	}
}

// Create creates a new membership
func (s *Membership) Create(ctx context.Context, membership *model.Membership) (*model.Membership, error) {
	query := `
		INSERT INTO memberships (
			entity_id,
			role,
			user_id
		) VALUES (
			$1, $2, $3
		) RETURNING id
	`

	err := s.QueryRowContext(
		ctx,
		query,
		membership.EntityID,
		membership.Role,
		membership.UserID,
	).Scan(&membership.ID)
	if err != nil {
		return nil, err
	}

	return membership, nil
}

// GetByUserID retrieves all memberships for a user
func (s *Membership) GetByUserID(ctx context.Context, userID string) ([]*model.Membership, error) {
	query := `
		SELECT
			id,
			entity_id,
			role,
			user_id,
			created_at,
			updated_at
		FROM
			memberships
		WHERE
			user_id = $1
	`

	rows, err := s.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memberships []*model.Membership
	for rows.Next() {
		membership := &model.Membership{}
		if err := rows.Scan(
			&membership.ID,
			&membership.EntityID,
			&membership.Role,
			&membership.UserID,
			&membership.CreatedAt,
			&membership.UpdatedAt,
		); err != nil {
			return nil, err
		}

		memberships = append(memberships, membership)
	}

	return memberships, rows.Err()
}

// GetByEntityID retrieves all memberships for an entity
func (s *Membership) GetByEntityID(ctx context.Context, entityID string) ([]*model.Membership, error) {
	query := `
		SELECT
			id,
			entity_id,
			role,
			user_id,
			created_at,
			updated_at
		FROM memberships
		WHERE entity_id = $1
	`

	rows, err := s.QueryContext(ctx, query, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memberships []*model.Membership
	for rows.Next() {
		membership := &model.Membership{}
		if err := rows.Scan(
			&membership.ID,
			&membership.EntityID,
			&membership.Role,
			&membership.UserID,
			&membership.CreatedAt,
			&membership.UpdatedAt,
		); err != nil {
			return nil, err
		}

		memberships = append(memberships, membership)
	}

	return memberships, rows.Err()
}

// GetByEntityIDWithInheritance retrieves a membership with the specified
// entityID. The function also ensures that a valid membership, via inheritence
// is indeed present.
func (s *Membership) GetByEntityIDWithInheritance(ctx context.Context, userID, entityID string) ([]*model.Membership, error) {
	query := `
		WITH RECURSIVE entity_hierarchy AS (
			SELECT
				e.id,
				e.parent_id,
				m.id as member_id,
				m.role,
				m.user_id
			FROM entities e
			INNER JOIN memberships m ON m.entity_id = e.id
			WHERE m.user_id = $1

			UNION ALL

			SELECT
				e.id,
				e.parent_id,
				eh.member_id,
				eh.role,
				eh.user_id
			FROM entities e
			INNER JOIN entity_hierarchy eh ON e.parent_id = eh.id
		)
		SELECT DISTINCT ON (eh.id)
			m.id,
			eh.id as entity_id,
			eh.role,
			m.user_id,
			m.created_at,
			m.updated_at
		FROM entity_hierarchy eh
		INNER JOIN memberships m ON m.id = eh.member_id
		WHERE eh.id = $2`

	rows, err := s.QueryContext(ctx, query, userID, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memberships []*model.Membership
	for rows.Next() {
		membership := &model.Membership{}
		if err := rows.Scan(
			&membership.ID,
			&membership.EntityID,
			&membership.Role,
			&membership.UserID,
			&membership.CreatedAt,
			&membership.UpdatedAt,
		); err != nil {
			return nil, err
		}

		memberships = append(memberships, membership)
	}

	return memberships, rows.Err()
}

// GetByUserIDWithInheritance retrieves all memberships for a user including inherited ones through entity hierarchy
func (s *Membership) GetByUserIDWithInheritance(ctx context.Context, userID string) ([]*model.Membership, error) {
	query := `
		WITH RECURSIVE entity_hierarchy AS (
			-- Get all entities the user has direct membership to
			SELECT
				e.id,
				e.parent_id,
				m.id as member_id,
				m.role,
				m.user_id,
				1 as level
			FROM entities e
			INNER JOIN memberships m ON m.entity_id = e.id
			WHERE m.user_id = $1

			UNION ALL

			-- Get all child entities
			SELECT
				e.id,
				e.parent_id,
				eh.member_id,
				eh.role,
				eh.user_id,
				eh.level + 1
			FROM entities e
			INNER JOIN entity_hierarchy eh ON e.parent_id = eh.id
		)
		SELECT DISTINCT ON (eh.id)
			m.id,
			eh.id as entity_id,
			eh.role,
			m.user_id,
			m.created_at,
			m.updated_at
		FROM entity_hierarchy eh
		INNER JOIN memberships m ON m.id = eh.member_id
		ORDER BY eh.id, eh.level ASC
	`

	rows, err := s.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memberships []*model.Membership
	for rows.Next() {
		membership := &model.Membership{}
		entityID := ""
		if err := rows.Scan(
			&membership.ID,
			&entityID,
			&membership.Role,
			&membership.UserID,
			&membership.CreatedAt,
			&membership.UpdatedAt,
		); err != nil {
			return nil, err
		}
		membership.EntityID = &entityID
		memberships = append(memberships, membership)
	}

	return memberships, rows.Err()
}
