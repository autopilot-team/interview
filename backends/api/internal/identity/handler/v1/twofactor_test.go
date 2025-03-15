package v1

import (
	"autopilot/backends/api/internal/identity"
	"autopilot/backends/api/internal/identity/model"
	"autopilot/backends/api/internal/identity/service"
	"autopilot/backends/api/pkg/httpx"
	"autopilot/backends/api/pkg/testutil"
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifyTwoFactor(t *testing.T) {
	t.Parallel()
	verifyPath := BasePath("/identity/verify-two-factor")
	api, container, mods := testutil.Container(t)

	ctx := context.Background()
	auth := identity.NewAuthentication(container, api, mods.Identity.Service)
	err := AddRoutes(container, api, mods.Identity.Service, auth)
	require.NoError(t, err)

	tests := []struct {
		name           string
		setup          func(t *testing.T) (string, string) // Returns session token and valid code
		payload        map[string]any
		expectedStatus int
		err            error
		checkCookies   bool
	}{
		{
			name: "should verify two-factor authentication successfully",
			setup: func(t *testing.T) (string, string) {
				_, session, twoFactorSetup := setupTestUserWithTwoFactor(t, ctx, mods.Identity.Service)
				return session.Token, twoFactorSetup.BackupCodes[0]
			},
			payload: map[string]any{
				"code": "", // Will be set from setup
			},
			expectedStatus: http.StatusNoContent,
			checkCookies:   true,
		},
		{
			name: "should reject missing session cookie",
			setup: func(t *testing.T) (string, string) {
				_, _, twoFactorSetup := setupTestUserWithTwoFactor(t, ctx, mods.Identity.Service)
				return "", twoFactorSetup.BackupCodes[0]
			},
			payload: map[string]any{
				"code": "", // Will be set from setup
			},
			expectedStatus: http.StatusUnauthorized,
			err:            httpx.ErrUnauthenticated,
		},
		{
			name: "should reject invalid session",
			setup: func(t *testing.T) (string, string) {
				_, _, twoFactorSetup := setupTestUserWithTwoFactor(t, ctx, mods.Identity.Service)
				return "invalid-session", twoFactorSetup.BackupCodes[0]
			},
			payload: map[string]any{
				"code": "", // Will be set from setup
			},
			expectedStatus: http.StatusUnauthorized,
			err:            httpx.ErrUnauthenticated,
		},
		{
			name: "should reject invalid 2FA code",
			setup: func(t *testing.T) (string, string) {
				_, session, _ := setupTestUserWithTwoFactor(t, ctx, mods.Identity.Service)
				return session.Token, "invalid"
			},
			payload: map[string]any{
				"code": "", // Will be set from setup
			},
			expectedStatus: http.StatusUnauthorized,
			err:            httpx.ErrInvalidTwoFactorCode,
		},
		{
			name: "should reject missing code",
			setup: func(t *testing.T) (string, string) {
				_, session, _ := setupTestUserWithTwoFactor(t, ctx, mods.Identity.Service)
				return session.Token, ""
			},
			payload:        map[string]any{},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject used backup code",
			setup: func(t *testing.T) (string, string) {
				user, session, twoFactorSetup := setupTestUserWithTwoFactor(t, ctx, mods.Identity.Service)

				// Use the backup code once
				err := mods.Identity.Service.TwoFactor.Verify(ctx, user.ID, twoFactorSetup.BackupCodes[0])
				require.NoError(t, err)

				// Try to use the same code again
				return session.Token, twoFactorSetup.BackupCodes[0]
			},
			payload: map[string]any{
				"code": "", // Will be set from setup
			},
			expectedStatus: http.StatusUnauthorized,
			err:            httpx.ErrInvalidTwoFactorCode,
		},
		{
			name: "should reject expired session",
			setup: func(t *testing.T) (string, string) {
				_, session, twoFactorSetup := setupTestUserWithTwoFactor(t, ctx, mods.Identity.Service)

				// Force expire the session by invalidating it
				err := mods.Identity.Service.Session.Invalidate(ctx, session.Token)
				require.NoError(t, err)

				return session.Token, twoFactorSetup.BackupCodes[0]
			},
			payload: map[string]any{
				"code": "", // Will be set from setup
			},
			expectedStatus: http.StatusUnauthorized,
			err:            httpx.ErrUnauthenticated,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Run setup to get fresh session and code
			sessionToken, code := tc.setup(t)
			if code != "" {
				tc.payload["code"] = code
			}

			// Prepare request with payload and cookie header
			args := []any{tc.payload}
			if sessionToken != "" {
				args = append(args, fmt.Sprintf("Cookie: session=%s", sessionToken))
			}

			resp := api.Post(verifyPath, args...)
			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.err != nil {
				httpx.AssertErr(t, tc.err, resp.Body)
				return
			}

			if tc.checkCookies {
				cookies := resp.Header().Values("Set-Cookie")
				assert.Len(t, cookies, 2) // Session and refresh cookies

				var foundSession, foundRefresh bool
				for _, cookie := range cookies {
					if strings.Contains(cookie, "session=") {
						assert.Contains(t, cookie, "HttpOnly")
						assert.Contains(t, cookie, "SameSite=Lax")
						foundSession = true
					} else if strings.Contains(cookie, "refresh_token=") {
						assert.Contains(t, cookie, "HttpOnly")
						assert.Contains(t, cookie, "SameSite=Lax")
						foundRefresh = true
					}
				}

				assert.True(t, foundSession, "Session cookie not found")
				assert.True(t, foundRefresh, "Refresh cookie not found")
			}
		})
	}
}

func setupTestUserWithTwoFactor(t *testing.T, ctx context.Context, service *service.Manager) (*model.User, *model.Session, *service.TwoFactorSetupData) {
	// Create a test user
	now := time.Now()
	user, err := service.User.Create(ctx, &model.User{
		Email:           fmt.Sprintf("test-%d@test.com", time.Now().UnixNano()),
		Name:            "Test User",
		EmailVerifiedAt: &now,
	}, "StrongPass123!")
	require.NoError(t, err)

	// Create a test session with 2FA pending
	session, err := service.Session.Create(ctx, user.Email, "StrongPass123!")
	require.NoError(t, err)

	// Enable 2FA for the user and get the backup codes
	twoFactorSetup, err := service.TwoFactor.Setup(ctx, user.ID)
	require.NoError(t, err)

	// Mark session as 2FA pending
	err = service.Session.UpdateTwoFactorStatus(ctx, session.Token, true)
	require.NoError(t, err)

	return user, session, twoFactorSetup
}
