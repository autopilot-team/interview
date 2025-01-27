package model

import "time"

// Verification represents an email or other verification process
type Verification struct {
	ID        string    `db:"id"`
	Value     string    `db:"value"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// IsExpired checks if the verification has expired
func (v *Verification) IsExpired() bool {
	return time.Now().After(v.ExpiresAt)
}
