package v1

import (
	"autopilot/backends/api/internal/handler"
	"autopilot/backends/api/internal/validator"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

// SessionUser is the user object for the session.
type SessionUser struct {
	ID                 string     `json:"id" doc:"The user's ID"`
	Email              string     `json:"email" doc:"The user's email address"`
	Name               string     `json:"name" doc:"The user's name"`
	Image              *string    `json:"image,omitempty" doc:"The user's profile image URL"`
	IsVerified         bool       `json:"isVerified" doc:"Whether the user's email is verified"`
	IsTwoFactorEnabled bool       `json:"isTwoFactorEnabled" doc:"Whether two-factor authentication is enabled"`
	LastActiveAt       *time.Time `json:"lastActiveAt,omitempty" doc:"The user's last activity time"`
	LastLoggedInAt     *time.Time `json:"lastLoggedInAt,omitempty" doc:"The user's last login time"`
	SessionExpiresAt   time.Time  `json:"sessionExpiresAt" doc:"When the current session will expire"`
}

// MeRequest is the request body for the get session endpoint.
type MeRequest struct {
	Session http.Cookie `cookie:"session" doc:"The session cookie"`
}

// MeResponse is the response body for the get session endpoint.
type MeResponse struct {
	Body struct {
		IsTwoFactorPending bool        `json:"isTwoFactorPending" doc:"Whether two-factor authentication is pending"`
		User               SessionUser `json:"user" doc:"The user object"`
	}
}

// Me is the handler for the get session endpoint.
func (v *V1) Me(ctx context.Context, input *MeRequest) (*MeResponse, error) {
	session, err := v.service.Session.GetByToken(ctx, input.Session.Value)
	if err != nil {
		v.Logger.Error("Failed to get session", "error", err.Error())
		return nil, huma.Error401Unauthorized("Invalid session")
	}

	user, err := v.service.User.GetByID(ctx, session.UserID)
	if err != nil {
		v.Logger.Error("Failed to get user", "error", err.Error())
		return nil, huma.Error500InternalServerError("Failed to get user")
	}

	response := &MeResponse{}
	response.Body.User.ID = session.UserID
	response.Body.User.Email = user.Email
	response.Body.User.Name = user.Name
	response.Body.User.Image = user.Image
	response.Body.User.IsVerified = user.IsEmailVerified()
	response.Body.User.IsTwoFactorEnabled = false
	response.Body.User.LastActiveAt = user.LastActiveAt
	response.Body.User.LastLoggedInAt = user.LastLoggedInAt
	response.Body.User.SessionExpiresAt = session.ExpiresAt

	return response, nil
}

// SignInRequest is the request body for the sign in endpoint.
type SignInRequest struct {
	Body struct {
		CfTurnstileToken validator.TurnstileToken `json:"cfTurnstileToken" required:"true" doc:"The Cloudflare Turnstile token" example:"XXX.DUMMY.TOKEN"`
		Email            string                   `json:"email" required:"true" doc:"The user's email address" format:"email" example:"john_doe@example.com"`
		Password         string                   `json:"password" required:"true" doc:"The user's password" example:"password123"`
	}
}

// SignInResponse is the response body for the sign in endpoint.
type SignInResponse struct {
	Body struct {
		IsTwoFactorPending bool        `json:"isTwoFactorPending" doc:"Whether two-factor authentication is pending"`
		User               SessionUser `json:"user" doc:"The user object"`
	}

	SetCookies []http.Cookie `header:"Set-Cookie"`
}

// SignIn is the handler for the sign in endpoint.
func (v *V1) SignIn(ctx context.Context, input *SignInRequest) (*SignInResponse, error) {
	session, err := v.service.Session.Create(ctx, input.Body.Email, input.Body.Password)
	if err != nil {
		v.Logger.Error("Failed to sign in", "error", err.Error())
		return nil, handler.ConvertServiceError(401, err)
	}

	user, err := v.service.User.GetByID(ctx, session.UserID)
	if err != nil {
		v.Logger.Error("Failed to get user", "error", err.Error())
		return nil, handler.ConvertServiceError(500, err)
	}

	response := &SignInResponse{
		SetCookies: []http.Cookie{
			v.newSessionCookie(
				session.Token,
				int(time.Until(session.ExpiresAt).Seconds()),
				session.ExpiresAt,
			),
			v.newRefreshCookie(
				session.RefreshToken,
				int(time.Until(session.RefreshExpiresAt).Seconds()),
				session.RefreshExpiresAt,
			),
		},
	}

	response.Body.User.ID = session.UserID
	response.Body.User.Email = user.Email
	response.Body.User.Name = user.Name
	response.Body.User.Image = user.Image
	response.Body.User.IsVerified = user.IsEmailVerified()
	response.Body.User.IsTwoFactorEnabled = false
	response.Body.User.LastActiveAt = user.LastActiveAt
	response.Body.User.LastLoggedInAt = user.LastLoggedInAt
	response.Body.User.SessionExpiresAt = session.ExpiresAt

	return response, nil
}

// SignOutRequest is the request body for the sign out endpoint.
type SignOutRequest struct {
	Session http.Cookie `cookie:"session" doc:"The session cookie"`
}

// SignOutResponse is the response body for the sign out endpoint.
type SignOutResponse struct {
	SetCookies []http.Cookie `header:"Set-Cookie"`
}

// SignOut is the handler for the sign out endpoint.
func (v *V1) SignOut(ctx context.Context, input *SignOutRequest) (*SignOutResponse, error) {
	if err := v.service.Session.Invalidate(ctx, input.Session.Value); err != nil {
		v.Logger.Error("Failed to sign out", "error", err.Error())
		return nil, handler.ConvertServiceError(401, err)
	}

	response := &SignOutResponse{
		SetCookies: []http.Cookie{
			v.newSessionCookie("", -1, time.Time{}),
			v.newRefreshCookie("", -1, time.Time{}),
		},
	}

	return response, nil
}

// newSessionCookie creates a new session cookie with standard configuration
func (v *V1) newSessionCookie(value string, maxAge int, expiresAt time.Time) http.Cookie {
	return http.Cookie{
		Name:     "session",
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   strings.HasPrefix(v.Config.App.BaseURL, "https://"),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
		Expires:  expiresAt,
	}
}

// newRefreshCookie creates a new refresh token cookie with standard configuration
func (v *V1) newRefreshCookie(value string, maxAge int, expiresAt time.Time) http.Cookie {
	return http.Cookie{
		Name:     "refresh_token",
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   strings.HasPrefix(v.Config.App.BaseURL, "https://"),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
		Expires:  expiresAt,
	}
}
