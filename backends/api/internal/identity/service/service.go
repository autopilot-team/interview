package service

import (
	"autopilot/backends/api/internal/identity/model"
	"autopilot/backends/api/internal/identity/store"
	"autopilot/backends/api/pkg/app"
	"autopilot/backends/api/pkg/httpx"
	"autopilot/backends/api/pkg/middleware"
	"autopilot/backends/internal/types"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// Manager is a collection of services used by the handlers/workers.
type Manager struct {
	Entity     Entityer
	Membership Membershiper
	Session    Sessioner
	TwoFactor  TwoFactorer
	User       Userer
}

// NewManager creates a new service manager
func NewManager(container *app.Container, store *store.Manager) *Manager {
	twoFactorService := NewTwoFactor(container, store)
	sessionService := NewSession(container, store, twoFactorService)
	membershipService := NewMembership(container, store)
	entityService := NewEntity(container, store)

	return &Manager{
		Entity:     entityService,
		Membership: membershipService,
		Session:    sessionService,
		TwoFactor:  twoFactorService,
		User:       NewUser(container, store),
	}
}

// auditLog is a helper function to create audit logs consistently across services
func auditLog(ctx context.Context, store *store.Manager, resourceType types.Resource, action types.Action, resourceID, userID string, metadata map[string]any) error {
	auditLog := &model.AuditLog{
		Action:       action,
		ResourceID:   resourceID,
		ResourceType: resourceType,
		UserID:       userID,
	}

	// Get request metadata for IP and user agent
	reqMetadata := middleware.GetRequestMetadata(ctx)
	if reqMetadata != nil {
		auditLog.IPAddress = &reqMetadata.IPAddress
		auditLog.UserAgent = &reqMetadata.UserAgent
	}

	// Convert metadata map to JSON if provided
	if metadata != nil {
		metadataJSON, err := json.Marshal(metadata)
		if err == nil {
			auditLog.Metadata = metadataJSON
		}
	}

	if _, err := store.AuditLog.Create(ctx, auditLog); err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}

	return nil
}

// createEmailVerification creates a new email verification record and sends the verification email
func createEmailVerification(ctx context.Context, store *store.Manager, container *app.Container, email string, user *model.User) error {
	now := time.Now()
	verification := &model.Verification{
		Context:   model.VerificationContextEmailVerification,
		Value:     email,
		ExpiresAt: now.Add(model.EmailVerificationDuration),
		CreatedAt: now,
		UpdatedAt: now,
	}

	verification, err := store.User.CreateVerification(ctx, verification)
	if err != nil {
		return httpx.ErrUnknown.WithInternal(err)
	}

	// Queue verification email
	locale := middleware.GetLocale(ctx)
	t := middleware.GetT(ctx)
	if t == nil {
		t = i18n.NewLocalizer(container.I18nBundle.Bundle, locale)
	}

	subject, err := t.Localize(&i18n.LocalizeConfig{
		MessageID: "welcome.title",
		TemplateData: map[string]any{
			"AppName": container.Config.App.Name,
		},
	})
	if err != nil {
		container.Logger.Error("Failed to localize email subject", "error", err)
		subject = fmt.Sprintf("Welcome to %s", container.Config.App.Name)
	}

	if _, err := container.Worker.Insert(ctx, MailerArgs{
		Data: map[string]any{
			"AssetsURL":       container.Config.App.AssetsURL,
			"AppName":         container.Config.App.Name,
			"Duration":        model.EmailVerificationDuration.Hours(),
			"Email":           user.Email,
			"Name":            user.Name,
			"VerificationURL": fmt.Sprintf("%s/verify-email?token=%s", container.Config.App.DashboardURL, verification.ID),
		},
		Email:    user.Email,
		Locale:   locale,
		Subject:  subject,
		Template: "welcome",
	}, nil); err != nil {
		container.Logger.Error("Failed to queue verification email", "error", err)
	}

	return nil
}

// generateSecureToken generates a secure random token of the specified length
func generateSecureToken(length int) (string, error) {
	token := make([]byte, length)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(token), nil
}
