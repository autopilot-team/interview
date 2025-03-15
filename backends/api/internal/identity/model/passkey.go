package model

import "time"

// Passkey represents a user's passkey for authentication
type Passkey struct {
	ID           string    `db:"id"`
	BackedUpAt   time.Time `db:"backed_up_at"`
	Counter      int       `db:"counter"`
	CredentialID string    `db:"credential_id"`
	DeviceType   string    `db:"device_type"`
	Name         *string   `db:"name"`
	PublicKey    string    `db:"public_key"`
	Transports   *string   `db:"transports"`
	UserID       string    `db:"user_id"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
