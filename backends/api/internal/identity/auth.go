package identity

import (
	"autopilot/backends/api/internal/identity/service"
	"autopilot/backends/api/pkg/app"
	"autopilot/backends/api/pkg/httpx"
	"autopilot/backends/api/pkg/middleware"
	"autopilot/backends/internal/types"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type Authentication struct {
	*app.Container
	API     huma.API
	Session service.Sessioner
}

func NewAuthentication(container *app.Container, api huma.API, manager *service.Manager) httpx.Authenticator {
	return &Authentication{
		Container: container,
		API:       api,
		Session:   manager.Session,
	}
}

func (s *Authentication) RequireUserSession(ctx huma.Context, next func(huma.Context)) {
	cookie, err := huma.ReadCookie(ctx, "session")
	if err != nil || cookie.Value == "" {
		_ = huma.WriteErr(s.API, ctx, http.StatusUnauthorized, "Unauthenticated", httpx.ErrUnauthenticated)
		return
	}

	session, err := s.Session.GetByToken(ctx.Context(), cookie.Value)
	if err != nil {
		s.Logger.Error("Failed to get session", "error", err)
		_ = huma.WriteErr(s.API, ctx, http.StatusUnauthorized, "Unauthenticated", httpx.ErrUnauthenticated)
		return
	}

	entityID := middleware.GetActiveEntity(ctx.Context())
	mode := types.GetOperationMode(ctx.Context())
	next(httpx.WithAuthInfo(ctx, httpx.AuthInfo{
		Authenticated: true,
		EntityID:      entityID,
		UserID:        session.UserID,
		Mode:          mode,
	}))
}

func (s *Authentication) RequireSecretKey(ctx huma.Context, next func(huma.Context)) {
	token := ctx.Header("X-Api-Key")
	if token == "" {
		_ = huma.WriteErr(s.API, ctx, http.StatusUnauthorized, "Unauthenticated", httpx.ErrUnauthenticated)
		return
	}
}

func (s *Authentication) RequirePublishableKey(ctx huma.Context, next func(huma.Context)) {
	_ = huma.WriteErr(s.API, ctx, http.StatusInternalServerError, "Unimplemented")
}

func (s *Authentication) RequireAuthenticated(ctx huma.Context, next func(huma.Context)) {
	cookie, err := huma.ReadCookie(ctx, "session")
	if err == nil && cookie.Value != "" {
		s.cookie(ctx, cookie.Value, next)
		return
	}

	token := ctx.Header("X-Api-Key")
	if token == "" {
		_ = huma.WriteErr(s.API, ctx, http.StatusUnauthorized, "Unauthenticated", httpx.ErrUnauthenticated)
		return
	}

	s.secretKey(ctx, token, next)
}

func (s *Authentication) cookie(ctx huma.Context, cookie string, next func(huma.Context)) {
	session, err := s.Session.GetByToken(ctx.Context(), cookie)
	if err != nil {
		s.Logger.Error("Failed to get session", "error", err)
		_ = huma.WriteErr(s.API, ctx, http.StatusUnauthorized, "Unauthenticated", httpx.ErrUnauthenticated)
		return
	}

	entityID := middleware.GetActiveEntity(ctx.Context())
	mode := types.GetOperationMode(ctx.Context())
	next(httpx.WithAuthInfo(ctx, httpx.AuthInfo{
		Authenticated: true,
		EntityID:      entityID,
		UserID:        session.UserID,
		Mode:          mode,
		EntityRole:    session.Role(entityID),
	}))
}

func (s *Authentication) secretKey(ctx huma.Context, _ string, _ func(huma.Context)) {
	_ = huma.WriteErr(s.API, ctx, http.StatusUnauthorized, "Unauthenticated")
}
