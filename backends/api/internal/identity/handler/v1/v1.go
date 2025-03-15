package v1

import (
	"autopilot/backends/api/internal/identity/service"
	"autopilot/backends/api/pkg/app"
	"autopilot/backends/api/pkg/httpx"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// V1 is the v1 API handler
type V1 struct {
	*app.Container
	identity *service.Manager
}

var TagIdentity = huma.Tag{
	Name:        "Identity",
	Description: `Identity & Access Management`,
}

func BasePath(path string) string {
	return fmt.Sprintf("/v1%s", path)
}

// AddRoutes adds the v1 API docs/routes to the http server
func AddRoutes(container *app.Container, humaAPI huma.API, identity *service.Manager, auth httpx.Authenticator) error {
	api := httpx.InitHandler(humaAPI, container.Mode, auth)
	api.AddTags(&TagIdentity)

	v1 := &V1{
		Container: container,
		identity:  identity,
	}

	// API Key access routes
	// All API key info endpoints are scoped to user session only.

	// User Routes
	httpx.Register(api, huma.Operation{
		Method:      http.MethodGet,
		OperationID: "get-user",
		Path:        BasePath("/users/{id}"),
		Summary:     "Get user",
		Tags:        []string{TagIdentity.Name},
	}, v1.GetUser, api.WithUserSession())

	httpx.Register(api, huma.Operation{
		Method:      http.MethodPut,
		OperationID: "update-user",
		Path:        BasePath("/users/{id}"),
		Summary:     "Update user",
		Tags:        []string{TagIdentity.Name},
	}, v1.UpdateUser, api.WithUserSession())

	httpx.Register(api, huma.Operation{
		Method:       http.MethodPost,
		OperationID:  "update-user-image",
		Path:         BasePath("/users/{id}/image"),
		Summary:      "Update user profile image",
		Tags:         []string{TagIdentity.Name},
		MaxBodyBytes: 5 * 1024 * 1024, // 5 MiB max

	}, v1.UpdateUserImage, api.WithUserSession())

	// Identity Routes

	// Allow me endpoint with API key
	httpx.Register(api, huma.Operation{
		Method:      http.MethodGet,
		OperationID: "get-session",
		Path:        BasePath("/identity/me"),
		Summary:     "Get current user session",
		Tags:        []string{TagIdentity.Name},
	}, v1.Me, api.WithUserSession())

	// Identity endpoints which are unauthenticated.

	httpx.Register(api, huma.Operation{
		Method:      http.MethodPost,
		OperationID: "sign-in",
		Path:        BasePath("/identity/sign-in"),
		Summary:     "Authenticate and create a new session",
		Tags:        []string{TagIdentity.Name},
	}, v1.SignIn, api.WithUnauthenticated())

	httpx.Register(api, huma.Operation{
		Method:      http.MethodPost,
		OperationID: "sign-up",
		Path:        BasePath("/identity/sign-up"),
		Summary:     "Create a new user account",
		Tags:        []string{TagIdentity.Name},
	}, v1.SignUp, api.WithUnauthenticated())

	httpx.Register(api, huma.Operation{
		Method:      http.MethodPost,
		OperationID: "forgot-password",
		Path:        BasePath("/identity/forgot-password"),
		Summary:     "Initiate password reset process",
		Tags:        []string{TagIdentity.Name},
	}, v1.ForgotPassword, api.WithUnauthenticated())

	httpx.Register(api, huma.Operation{
		Method:      http.MethodPost,
		OperationID: "reset-password",
		Path:        BasePath("/identity/reset-password"),
		Summary:     "Complete password reset process",
		Tags:        []string{TagIdentity.Name},
	}, v1.ResetPassword, api.WithUnauthenticated())

	httpx.Register(api, huma.Operation{
		Method:      http.MethodPost,
		OperationID: "verify-email",
		Path:        BasePath("/identity/verify-email"),
		Summary:     "Confirm user email address",
		Tags:        []string{TagIdentity.Name},
	}, v1.VerifyEmail, api.WithUnauthenticated())

	// Private identity endpoints with rate limits

	httpx.Register(api, huma.Operation{
		Method:      http.MethodPost,
		OperationID: "update-password",
		Path:        BasePath("/identity/update-password"),
		Summary:     "Update user password",
		Tags:        []string{TagIdentity.Name},
	}, v1.UpdatePassword, api.WithUserSession())

	httpx.Register(api, huma.Operation{
		Method:      http.MethodPost,
		OperationID: "setup-two-factor",
		Path:        BasePath("/identity/setup-two-factor"),
		Summary:     "Setup two-factor authentication",
		Tags:        []string{TagIdentity.Name},
	}, v1.SetupTwoFactor, api.WithUserSession())

	httpx.Register(api, huma.Operation{
		Method:      http.MethodPost,
		OperationID: "enable-two-factor",
		Path:        BasePath("/identity/enable-two-factor"),
		Summary:     "Enable two-factor authentication after setup",
		Tags:        []string{TagIdentity.Name},
	}, v1.EnableTwoFactor, api.WithUserSession())

	httpx.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		OperationID: "disable-two-factor",
		Path:        BasePath("/identity/disable-two-factor"),
		Summary:     "Disable two-factor authentication",
		Tags:        []string{TagIdentity.Name},
	}, v1.DisableTwoFactor, api.WithUserSession())

	httpx.Register(api, huma.Operation{
		Method:      http.MethodPost,
		OperationID: "verify-two-factor",
		Path:        BasePath("/identity/verify-two-factor"),
		Summary:     "Verify two-factor authentication code during sign-in",
		Tags:        []string{TagIdentity.Name},
	}, v1.VerifyTwoFactor, api.WithUnauthenticated())

	httpx.Register(api, huma.Operation{
		Method:      http.MethodPost,
		OperationID: "verify-password",
		Path:        BasePath("/identity/verify-password"),
		Summary:     "Verify current password for sensitive operations",
		Tags:        []string{TagIdentity.Name},
	}, v1.VerifyPassword, api.WithUserSession())

	// Non-rate-limited endpoints
	httpx.Register(api, huma.Operation{
		Method:      http.MethodPost,
		OperationID: "refresh-session",
		Path:        BasePath("/identity/refresh-session"),
		Summary:     "Extend current session validity",
		Tags:        []string{TagIdentity.Name},
	}, v1.RefreshSession, api.WithUnauthenticated())

	httpx.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		OperationID: "sign-out",
		Path:        BasePath("/identity/sign-out"),
		Summary:     "Terminate current session",
		Tags:        []string{TagIdentity.Name},
	}, v1.SignOut, api.WithUserSession())

	httpx.Register(api, huma.Operation{
		Method:      http.MethodGet,
		OperationID: "get-all-sessions",
		Path:        BasePath("/identity/sessions"),
		Summary:     "Fetch all active sessions",
		Tags:        []string{TagIdentity.Name},
	}, v1.GetAllSessions, api.WithUserSession())

	httpx.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		OperationID: "delete-all-sessions",
		Path:        BasePath("/identity/sessions"),
		Summary:     "Terminate all active sessions",
		Tags:        []string{TagIdentity.Name},
	}, v1.DeleteAllSessions, api.WithUserSession())

	httpx.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		OperationID: "delete-session",
		Path:        BasePath("/identity/sessions/{id}"),
		Summary:     "Terminate session by id",
		Tags:        []string{TagIdentity.Name},
	}, v1.DeleteSession, api.WithUserSession())

	httpx.Register(api, huma.Operation{
		Method:      http.MethodPost,
		OperationID: "regenerate-qr-code",
		Path:        BasePath("/identity/regenerate-qr-code"),
		Summary:     "Regenerate QR code for existing two-factor authentication setup",
		Tags:        []string{TagIdentity.Name},
	}, v1.RegenerateQRCode, api.WithUserSession())

	return nil
}
