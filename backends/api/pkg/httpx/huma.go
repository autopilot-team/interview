package httpx

import (
	"autopilot/backends/internal/types"
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type HandlerOption func(*huma.Operation)

var tooManyRequestsRef = &huma.Response{
	Description: "Too many requests - rate limit exceeded",
	Ref:         "#/components/responses/TooManyRequests",
}

var securityApiKey = map[string][]string{"API Key Authentication": {}}

type API struct {
	huma.API

	mode          types.Mode
	authenticator Authenticator
}

func InitHandler(api huma.API, mode types.Mode, authenticator Authenticator) API {
	if api.OpenAPI().Components == nil {
		api.OpenAPI().Components = &huma.Components{}
	}

	if api.OpenAPI().Components.Extensions == nil {
		api.OpenAPI().Components.Extensions = make(map[string]any)
	}

	// Add rate limit documentation to OpenAPI spec
	api.OpenAPI().Components.Extensions["x-rate-limiting"] = map[string]any{
		"public": map[string]any{
			"description": "Rate limits for public endpoints that don't require authentication",
			"rate": map[string]any{
				"requests": 30,
				"window":   "5 minutes",
			},
			"scope": "by IP address",
		},
		"private": map[string]any{
			"description": "Rate limits for authenticated endpoints",
			"rate": map[string]any{
				"requests": 300,
				"window":   "1 minute",
			},
			"scope": "by IP address and endpoint",
		},
	}

	// Initialize responses map if nil
	if api.OpenAPI().Components.Responses == nil {
		api.OpenAPI().Components.Responses = make(map[string]*huma.Response)
	}

	// Add rate limit response to components
	api.OpenAPI().Components.Responses["TooManyRequests"] = &huma.Response{
		Description: "Too many requests - rate limit exceeded",
		Content: map[string]*huma.MediaType{
			"application/json": {
				Schema: &huma.Schema{
					Type: "object",
					Properties: map[string]*huma.Schema{
						"error": {
							Type:        "string",
							Description: "Error message",
							Examples:    []any{"Rate limit exceeded for API operations"},
						},
					},
				},
			},
		},
		Headers: map[string]*huma.Param{
			"X-RateLimit-Limit": {
				Description: "The number of allowed requests in the current period",
				Schema: &huma.Schema{
					Type: "integer",
				},
			},
			"X-RateLimit-Remaining": {
				Description: "The number of remaining requests in the current period",
				Schema: &huma.Schema{
					Type: "integer",
				},
			},
			"X-RateLimit-Reset": {
				Description: "The remaining window before the rate limit resets in UTC epoch seconds",
				Schema: &huma.Schema{
					Type: "integer",
				},
			},
		},
	}

	return API{
		API:           api,
		mode:          mode,
		authenticator: authenticator,
	}
}

func (a API) AddTags(tags ...*huma.Tag) {
	oapi := a.OpenAPI()
	oapi.Tags = append(oapi.Tags, tags...)
}

func RegisterPublish[I, O any](api API, op huma.Operation, handler func(context.Context, *I) (*O, error), opts ...HandlerOption) {
	opts = append(opts, WithPublish())
	Register(api, op, handler, opts...)
}

func Register[I, O any](api API, op huma.Operation, handler func(context.Context, *I) (*O, error), opts ...HandlerOption) {
	// Hide endpoints by default
	op.Hidden = api.mode == types.ReleaseMode

	if op.Responses == nil {
		op.Responses = make(map[string]*huma.Response, 1)
	}
	op.Responses["429"] = tooManyRequestsRef

	for _, opt := range opts {
		opt(&op)
	}

	// Default authenticated endpoint
	if len(op.Middlewares) == 0 {
		op.Middlewares = append(op.Middlewares, api.authenticator.RequireAuthenticated)
		op.Security = append(op.Security, securityApiKey)
	}

	huma.Register(api, op, handler)
}

func WithoutRateLimit() HandlerOption {
	return func(op *huma.Operation) {
		delete(op.Responses, "429")
	}
}

func WithPublish() HandlerOption {
	return func(op *huma.Operation) {
		op.Hidden = false
	}
}

func (a API) WithUnauthenticated() HandlerOption {
	return func(op *huma.Operation) {
		// insert empty middleware
		op.Middlewares = append(op.Middlewares, func(ctx huma.Context, next func(huma.Context)) {
			next(ctx)
		})
	}
}

func (a API) WithUserSession() HandlerOption {
	return func(op *huma.Operation) {
		op.Middlewares = append(op.Middlewares, a.authenticator.RequireUserSession)
	}
}

func (a API) WithSecretKey() HandlerOption {
	return func(op *huma.Operation) {
		op.Middlewares = append(op.Middlewares, a.authenticator.RequireSecretKey)
		op.Security = append(op.Security, securityApiKey)
	}
}

func (a API) WithPermission(resource types.Resource, action types.Action) HandlerOption {
	return func(op *huma.Operation) {
		if len(op.Middlewares) == 0 {
			op.Middlewares = append(op.Middlewares, a.authenticator.RequireAuthenticated)
			op.Security = append(op.Security, securityApiKey)
		}

		op.Middlewares = append(op.Middlewares, func(ctx huma.Context, next func(huma.Context)) {
			auth := GetAuthInfo(ctx.Context())
			if !auth.Authenticated {
				_ = huma.WriteErr(a.API, ctx, http.StatusUnauthorized, "Unauthenticated")
				return
			}

			if auth.EntityID == "" {
				_ = huma.WriteErr(a.API, ctx, http.StatusBadRequest, "No active entity selected", ErrEntityNotFound)
				return
			}

			if !auth.EntityRole.HasPermission(resource, action) {
				_ = huma.WriteErr(a.API, ctx, http.StatusForbidden, "Insufficient permissions", ErrInsufficientPermissions)
				return
			}

			next(ctx)
		})
	}
}
