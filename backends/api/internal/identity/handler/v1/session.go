package v1

import (
	"autopilot/backends/api/internal/identity/model"
	"autopilot/backends/api/pkg/httpx"
	"autopilot/backends/api/pkg/middleware"
	"autopilot/backends/internal/types"
	"context"
	"errors"
	"net/http"
	"strings"
	"time"
)

// SessionUser is the user object for the session.
type SessionUser struct {
	ID                 string        `json:"id" doc:"The user's ID"`
	Email              string        `json:"email" doc:"The user's email address"`
	Name               string        `json:"name" doc:"The user's name"`
	Image              *string       `json:"image,omitempty" doc:"The user's profile image URL"`
	IsVerified         bool          `json:"isVerified" doc:"Whether the user's email is verified"`
	IsTwoFactorEnabled bool          `json:"isTwoFactorEnabled" doc:"Whether two-factor authentication is enabled"`
	LastActiveAt       *time.Time    `json:"lastActiveAt,omitempty" doc:"The user's last activity time"`
	LastLoggedInAt     *time.Time    `json:"lastLoggedInAt,omitempty" doc:"The user's last login time"`
	SessionExpiresAt   time.Time     `json:"sessionExpiresAt" doc:"When the current session will expire"`
	Memberships        []*Membership `json:"memberships,omitempty" doc:"The user's entity memberships"`
}

// Membership represents a user's membership in an entity
type Membership struct {
	ID       string  `json:"id" doc:"The membership ID"`
	EntityID *string `json:"entityId" doc:"The entity ID"`
	Role     string  `json:"role" doc:"The user's role in the entity"`
	Entity   *Entity `json:"entity" doc:"The entity details"`
}

// Entity represents an entity in the session
type Entity struct {
	ID       string  `json:"id" doc:"The entity's ID"`
	Name     string  `json:"name" doc:"The entity's name"`
	Slug     string  `json:"slug" doc:"The entity's slug"`
	Type     string  `json:"type" doc:"The entity's type"`
	Status   string  `json:"status" doc:"The entity's status"`
	ParentID *string `json:"parentId,omitempty" doc:"The parent entity's ID"`
	Logo     *string `json:"logo,omitempty" doc:"The entity's logo URL"`
	Domain   *string `json:"domain,omitempty" doc:"The entity's domain"`
}

type (
	EntityResource = types.Resource
	EntityAction   = types.Action
)

type EntityRole struct {
	Name   string                            `json:"name" doc:"The name of the user's role in the entity"`
	Access map[EntityResource][]EntityAction `json:"access" doc:"The role's access to resources in the entity"`
}

// Session is the object describing an active user session.
type Session struct {
	ID           string    `json:"id" doc:"The session ID"`
	UserID       string    `json:"userId" doc:"The user's ID"`
	Current      bool      `json:"current" doc:"Whether this is the user's current session"`
	IPAddress    *string   `json:"ipAddress" doc:"The last seen IP address"`
	Country      *string   `json:"country" doc:"The last seen IP country"`
	UserAgent    *string   `json:"userAgent" doc:"The last seen user agent"`
	CreatedAt    time.Time `json:"createdAt" doc:"The session creation time"`
	LastActiveAt time.Time `json:"updatedAt" doc:"The session's last activity time"`
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
		ActiveEntity       *Entity     `json:"activeEntity,omitempty" doc:"The currently active entity"`
		EntityRole         *EntityRole `json:"entityRole,omitempty" doc:"The permissions for the currently active entity"`
	}
}

// Me is the handler for the get session endpoint.
func (v *V1) Me(ctx context.Context, input *MeRequest) (*MeResponse, error) {
	session, err := v.identity.Session.GetByTokenFull(ctx, input.Session.Value)
	if err != nil {
		v.Logger.Error("Failed to get session", "error", err)
		return nil, httpx.ErrUnauthenticated
	}

	user, err := v.identity.User.GetByID(ctx, session.UserID)
	if err != nil {
		v.Logger.Error("Failed to get user", "error", err)
		return nil, err
	}

	// Check if 2FA is enabled
	twoFactor, err := v.identity.TwoFactor.GetByUserID(ctx, user.ID)
	if err != nil && !errors.Is(err, httpx.ErrTwoFactorNotEnabled) {
		v.Logger.Error("Failed to get two-factor status", "error", err)
		return nil, err
	}
	isTwoFactorEnabled := twoFactor != nil && twoFactor.EnabledAt != nil

	// Convert model memberships to response memberships
	responseMemberships := make([]*Membership, len(session.Memberships))
	for i, m := range session.Memberships {
		// Get full entity details for each membership
		entity, err := v.identity.Entity.Get(ctx, *m.EntityID)
		if err != nil {
			v.Logger.Error("Failed to get entity for membership", "error", err)
			return nil, err
		}

		responseMemberships[i] = &Membership{
			ID:       m.ID,
			EntityID: m.EntityID,
			Role:     string(m.Role),
			Entity: &Entity{
				ID:       entity.ID,
				Name:     entity.Name,
				Slug:     entity.Slug,
				Type:     string(entity.Type),
				Status:   string(entity.Status),
				ParentID: entity.ParentID,
				Logo:     entity.Logo,
				Domain:   entity.Domain,
			},
		}
	}

	response := &MeResponse{}
	response.Body.User.ID = session.UserID
	response.Body.User.Email = user.Email
	response.Body.User.Name = user.Name
	response.Body.User.Image = user.Image
	response.Body.User.IsVerified = user.IsEmailVerified()
	response.Body.User.IsTwoFactorEnabled = isTwoFactorEnabled
	response.Body.User.LastActiveAt = user.LastActiveAt
	response.Body.User.LastLoggedInAt = user.LastLoggedInAt
	response.Body.User.SessionExpiresAt = session.ExpiresAt
	response.Body.User.Memberships = responseMemberships

	// If there's an active entity, get its details
	entityID := middleware.GetActiveEntity(ctx)
	var entity *model.Entity
	if entityID != "" {
		e, err := v.identity.Entity.Get(ctx, entityID)
		if err == nil && session.Role(e.ID).HasPermission(types.ResourceEntity, types.ActionRead) {
			entity = e
		}
	}
	if entity == nil && len(session.Memberships) != 0 {
		entity, err = v.identity.Entity.Get(ctx, *session.Memberships[0].EntityID)
		if err != nil {
			v.Logger.Error("Failed to get active entity", "error", err)
			return nil, err
		}
	}

	if entity != nil {
		response.Body.ActiveEntity = &Entity{
			ID:       entity.ID,
			Name:     entity.Name,
			Slug:     entity.Slug,
			Type:     string(entity.Type),
			Status:   string(entity.Status),
			ParentID: entity.ParentID,
			Logo:     entity.Logo,
			Domain:   entity.Domain,
		}
		role := session.Role(entity.ID)
		if perms, ok := types.RolePermissions[role]; ok {
			response.Body.EntityRole = &EntityRole{
				Name:   role.String(),
				Access: perms,
			}
		}
	}

	return response, nil
}

// SignInRequest is the request body for the sign in endpoint.
type SignInRequest struct {
	Body struct {
		CfTurnstileToken httpx.TurnstileToken `json:"cfTurnstileToken" required:"true" doc:"The Cloudflare Turnstile token" example:"XXX.DUMMY.TOKEN"`
		Email            string               `json:"email" required:"true" doc:"The user's email address" format:"email" example:"john_doe@example.com"`
		Password         string               `json:"password" required:"true" doc:"The user's password" example:"password123"`
	}
}

// SignInResponse is the response body for the sign in endpoint.
type SignInResponse struct {
	Body struct {
		IsTwoFactorPending bool `json:"isTwoFactorPending" doc:"Whether two-factor authentication is pending"`
	}

	SetCookies []http.Cookie `header:"Set-Cookie"`
}

// SignIn is the handler for the sign in endpoint.
func (v *V1) SignIn(ctx context.Context, input *SignInRequest) (*SignInResponse, error) {
	session, err := v.identity.Session.Create(ctx, input.Body.Email, input.Body.Password)
	if err != nil {
		// Special handling for 2FA pending case
		if errors.Is(err, httpx.ErrTwoFactorPending) && session != nil {
			// Return success with temporary session and 2FA pending flag
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

			response.Body.IsTwoFactorPending = true

			return response, nil
		}

		v.Logger.Error("Failed to sign in", "error", err)
		return nil, err
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
	if err := v.identity.Session.Invalidate(ctx, input.Session.Value); err != nil {
		v.Logger.Error("Failed to sign out", "error", err)
		return nil, err
	}

	response := &SignOutResponse{
		SetCookies: []http.Cookie{
			v.newSessionCookie("", -1, time.Time{}),
			v.newRefreshCookie("", -1, time.Time{}),
		},
	}

	return response, nil
}

// GetAllSessionsRequest is the request body for listing all sessions endpoint.
type GetAllSessionsRequest struct {
	Session http.Cookie `cookie:"session" doc:"The session cookie"`
}

// GetAllSessionsResponse is the response body for the list all sessions endpoint.
type GetAllSessionsResponse struct {
	Body struct {
		Sessions []Session `json:"sessions" doc:"The current active user sessions"`
	}
}

func (v *V1) GetAllSessions(ctx context.Context, input *GetAllSessionsRequest) (*GetAllSessionsResponse, error) {
	// Get active sessions for current user ID
	sessions, err := v.identity.Session.ListByToken(ctx, input.Session.Value)
	if err != nil {
		v.Logger.Error("Failed to list user sessions", "error", err)
		return nil, err
	}
	response := &GetAllSessionsResponse{}
	for _, s := range sessions {
		response.Body.Sessions = append(response.Body.Sessions, Session{
			ID:           s.ID,
			UserID:       s.UserID,
			Current:      s.Token == input.Session.Value,
			IPAddress:    s.IPAddress,
			Country:      s.Country,
			UserAgent:    s.UserAgent,
			CreatedAt:    s.CreatedAt,
			LastActiveAt: s.UpdatedAt,
		})
	}
	return response, nil
}

// DeleteAllSessionsRequest is the request body for the delete all sessions endpoint.
type DeleteAllSessionsRequest struct {
	Session http.Cookie `cookie:"session" doc:"The session cookie"`
}

// DeleteAllSessionsResponse is the response body for the delete all sessions endpoint.
type DeleteAllSessionsResponse struct{}

// DeleteAllSessions is the handler for the delete all sessions endpoint.
func (v *V1) DeleteAllSessions(ctx context.Context, input *DeleteAllSessionsRequest) (*DeleteAllSessionsResponse, error) {
	// Get current session to get user ID
	session, err := v.identity.Session.GetByToken(ctx, input.Session.Value)
	if err != nil {
		v.Logger.Error("Failed to get session", "error", err)
		return nil, err
	}

	if err := v.identity.Session.InvalidateAllSessions(ctx, session.UserID, input.Session.Value); err != nil {
		v.Logger.Error("Failed to delete all sessions", "error", err)
		return nil, err
	}

	response := &DeleteAllSessionsResponse{}

	return response, nil
}

// DeleteSessionRequest is the request body for the delete session endpoint.
type DeleteSessionRequest struct {
	Session   http.Cookie `cookie:"session" doc:"The session cookie"`
	SessionID string      `path:"id" doc:"The session id"`
}

// DeleteSessionResponse is the response data for the delete session endpoint.
type DeleteSessionResponse struct{}

// DeleteSession is the handler for the delete session request.
func (v *V1) DeleteSession(ctx context.Context, input *DeleteSessionRequest) (*DeleteSessionResponse, error) {
	if err := v.identity.Session.InvalidateByID(ctx, input.Session.Value, input.SessionID); err != nil {
		v.Logger.Error("Failed to delete session", "session_id", input.SessionID, "error", err)
		return nil, err
	}

	response := &DeleteSessionResponse{}
	return response, nil
}

// RefreshSessionRequest is the request body for the refresh session endpoint.
type RefreshSessionRequest struct {
	RefreshToken http.Cookie `cookie:"refresh_token" doc:"The refresh token cookie"`
}

// RefreshSessionResponse is the response body for the refresh session endpoint.
type RefreshSessionResponse struct {
	SetCookies []http.Cookie `header:"Set-Cookie"`
}

// RefreshSession is the handler for the refresh session endpoint.
func (v *V1) RefreshSession(ctx context.Context, input *RefreshSessionRequest) (*RefreshSessionResponse, error) {
	session, err := v.identity.Session.Refresh(ctx, input.RefreshToken.Value)
	if err != nil {
		v.Logger.Error("Failed to refresh session", "error", err)
		return nil, err
	}

	response := &RefreshSessionResponse{
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

	return response, nil
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
