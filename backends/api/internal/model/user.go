package model

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID                  string     `db:"id"`
	Name                string     `db:"name"`
	Email               string     `db:"email"`
	EmailVerifiedAt     *time.Time `db:"email_verified_at"`
	FailedLoginAttempts int        `db:"failed_login_attempts"`
	Image               *string    `db:"image"`
	LastActiveAt        *time.Time `db:"last_active_at"`
	LastLoggedInAt      *time.Time `db:"last_logged_in_at"`
	LockedAt            *time.Time `db:"locked_at"`
	PasswordChangedAt   *time.Time `db:"password_changed_at"`
	PasswordHash        *string    `db:"password_hash"`
	TwoFactorEnabled    bool       `db:"two_factor_enabled"`
	CreatedAt           time.Time  `db:"created_at"`
	UpdatedAt           time.Time  `db:"updated_at"`
}

// IsLocked checks if the user account is locked at the given time
func (u *User) IsLocked(at time.Time) bool {
	return u.LockedAt != nil && u.LockedAt.After(at)
}

// IsEmailVerified checks if the user's email is verified
func (u *User) IsEmailVerified() bool {
	return u.EmailVerifiedAt != nil
}

// HasPassword checks if the user has a password set
func (u *User) HasPassword() bool {
	return u.PasswordHash != nil
}

// HasLoggedIn checks if the user has ever logged in
func (u *User) HasLoggedIn() bool {
	return u.LastLoggedInAt != nil
}
