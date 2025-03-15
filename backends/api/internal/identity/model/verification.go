package model

import "time"

const (
	// VerificationContextEmailVerification represents email verification context
	VerificationContextEmailVerification = "email_verification"

	// VerificationContextPasswordReset represents password reset context
	VerificationContextPasswordReset = "password_reset"

	// EmailVerificationDuration is the duration for which email verification links are valid
	EmailVerificationDuration = 24 * time.Hour

	// PasswordResetDuration is the duration for which password reset links are valid
	PasswordResetDuration = 1 * time.Hour
)

// Verification represents an email or other verification process
type Verification struct {
	ID        string    `db:"id"`
	Context   string    `db:"context"`
	Value     string    `db:"value"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// IsExpired checks if the verification has expired
func (v *Verification) IsExpired() bool {
	return time.Now().After(v.ExpiresAt)
}
