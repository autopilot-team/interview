package v1

import (
	"autopilot/backends/api/internal/identity"
	"autopilot/backends/api/internal/identity/model"
	"autopilot/backends/api/pkg/app/mocks"
	"autopilot/backends/api/pkg/httpx"
	"autopilot/backends/api/pkg/testutil"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSignIn(t *testing.T) {
	t.Parallel()
	signInPath := BasePath("/identity/sign-in")
	api, container, mods := testutil.Container(t)

	ctx := context.Background()
	auth := identity.NewAuthentication(container, api, mods.Identity.Service)
	err := AddRoutes(container, api, mods.Identity.Service, auth)
	require.NoError(t, err)

	// Create a test user with unverified email
	unverifiedUser, err := mods.Identity.Service.User.Create(ctx, &model.User{
		Email: "unverified@test.com",
		Name:  "Unverified User",
	}, "StrongPass123!")
	require.NoError(t, err)

	// Get the initial verification
	initialVerification, err := mods.Identity.Store.User.GetVerificationByValue(ctx, model.VerificationContextEmailVerification, unverifiedUser.Email)
	require.NoError(t, err)

	// Create a test user with verified email
	now := time.Now()
	_, err = mods.Identity.Service.User.Create(ctx, &model.User{
		Email:           "test@test.com",
		Name:            "Test User",
		EmailVerifiedAt: &now,
	}, "StrongPass123!")
	require.NoError(t, err)

	// Get the mock Turnstile instance
	mockTurnstile := container.Turnstile.(*mocks.MockTurnstiler)
	// Set up default success behavior for Turnstile
	mockTurnstile.On("Verify", mock.Anything, "XXX.DUMMY.TOKEN", "").Return(true, nil)
	// Set up failure behavior for invalid token
	mockTurnstile.On("Verify", mock.Anything, "INVALID.TOKEN", "").Return(false, nil)
	// Set up error behavior for error token
	mockTurnstile.On("Verify", mock.Anything, "ERROR.TOKEN", "").Return(false, assert.AnError)

	// Test expired verification
	t.Run("unverified email with expired verification", func(t *testing.T) {
		// Delete the initial verification to simulate expiration
		err := mods.Identity.Store.User.DeleteVerification(ctx, initialVerification.ID)
		require.NoError(t, err)

		resp := api.Post(signInPath, map[string]any{
			"cfTurnstileToken": "XXX.DUMMY.TOKEN",
			"email":            "unverified@test.com",
			"password":         "StrongPass123!",
		})
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		httpx.AssertErr(t, httpx.ErrEmailNotVerified, resp.Body)

		// Verify a new verification was created
		newVerification, err := mods.Identity.Store.User.GetVerificationByValue(ctx, model.VerificationContextEmailVerification, unverifiedUser.Email)
		assert.NoError(t, err)
		assert.NotNil(t, newVerification)
		assert.NotEqual(t, initialVerification.ID, newVerification.ID)
		assert.True(t, newVerification.ExpiresAt.After(time.Now()))
	})

	tests := []struct {
		name           string
		payload        map[string]any
		expectedStatus int
		err            error
		checkCookie    bool
	}{
		{
			name: "should sign in successfully",
			payload: map[string]any{
				"cfTurnstileToken": "XXX.DUMMY.TOKEN",
				"email":            "test@test.com",
				"password":         "StrongPass123!",
			},
			expectedStatus: http.StatusOK,
			checkCookie:    true,
		},
		{
			name: "should reject invalid email format",
			payload: map[string]any{
				"cfTurnstileToken": "XXX.DUMMY.TOKEN",
				"email":            "invalid-email",
				"password":         "StrongPass123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject missing required fields",
			payload: map[string]any{
				"cfTurnstileToken": "XXX.DUMMY.TOKEN",
				"email":            "test@test.com",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject invalid credentials",
			payload: map[string]any{
				"cfTurnstileToken": "XXX.DUMMY.TOKEN",
				"email":            "test@test.com",
				"password":         "WrongPass123!",
			},
			expectedStatus: http.StatusUnauthorized,
			err:            httpx.ErrInvalidCredentials,
		},
		{
			name: "should reject invalid turnstile token",
			payload: map[string]any{
				"cfTurnstileToken": "INVALID.TOKEN",
				"email":            "test@test.com",
				"password":         "StrongPass123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject turnstile verification error",
			payload: map[string]any{
				"cfTurnstileToken": "ERROR.TOKEN",
				"email":            "test@test.com",
				"password":         "StrongPass123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := api.Post(signInPath, tc.payload)
			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.err != nil {
				httpx.AssertErr(t, tc.err, resp.Body)
				return
			}

			if tc.checkCookie {
				cookie := resp.Header().Get("Set-Cookie")
				assert.Contains(t, cookie, "session=")
				assert.Contains(t, cookie, "HttpOnly")
				assert.Contains(t, cookie, "SameSite=Lax")
			}
		})
	}
}

func TestSignOut(t *testing.T) {
	t.Parallel()
	signOutPath := BasePath("/identity/sign-out")
	api, container, mods := testutil.Container(t)

	ctx := context.Background()
	auth := identity.NewAuthentication(container, api, mods.Identity.Service)
	err := AddRoutes(container, api, mods.Identity.Service, auth)
	assert.NoError(t, err)

	// Create a test user and session
	now := time.Now()
	_, err = mods.Identity.Service.User.Create(ctx, &model.User{
		Email:           "test@test.com",
		Name:            "Test User",
		EmailVerifiedAt: &now,
	}, "StrongPass123!")
	require.NoError(t, err)

	session, err := mods.Identity.Service.Session.Create(ctx, "test@test.com", "StrongPass123!")
	require.NoError(t, err)

	// Verify session was created successfully
	_, err = mods.Identity.Service.Session.GetByToken(ctx, session.Token)
	require.NoError(t, err)

	tests := []struct {
		name           string
		cookie         string
		expectedStatus int
		err            error
		checkCookie    bool
	}{
		{
			name:           "should sign out successfully",
			cookie:         session.Token,
			expectedStatus: http.StatusNoContent,
			checkCookie:    true,
		},
		{
			name:           "should reject invalid session token",
			cookie:         "invalid-token",
			expectedStatus: http.StatusUnauthorized,
			err:            httpx.ErrUnauthenticated,
		},
		{
			name:           "should reject missing session token",
			cookie:         "",
			expectedStatus: http.StatusUnauthorized,
			err:            httpx.ErrUnauthenticated,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cookie := ""
			if tc.cookie != "" {
				cookie = fmt.Sprintf("Cookie: session=%s", tc.cookie)
			}

			resp := api.Delete(signOutPath, cookie)
			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.err != nil {
				httpx.AssertErr(t, tc.err, resp.Body)
				return
			}

			if tc.checkCookie {
				cookie := resp.Header().Get("Set-Cookie")
				assert.Contains(t, cookie, "session=;")
				assert.Contains(t, cookie, "HttpOnly;")
				assert.Contains(t, cookie, "SameSite=Lax")
				assert.Contains(t, cookie, "Max-Age=0")
			}
		})
	}
}

func TestDeleteAllSessions(t *testing.T) {
	t.Parallel()
	deleteAllSessionsPath := BasePath("/identity/sessions")
	api, container, mods := testutil.Container(t)

	ctx := context.Background()
	auth := identity.NewAuthentication(container, api, mods.Identity.Service)
	err := AddRoutes(container, api, mods.Identity.Service, auth)
	assert.NoError(t, err)

	// Create a test user and multiple sessions
	now := time.Now()
	_, err = mods.Identity.Service.User.Create(ctx, &model.User{
		Email:           "test@test.com",
		Name:            "Test User",
		EmailVerifiedAt: &now,
	}, "StrongPass123!")
	require.NoError(t, err)

	// Create three sessions for the same user
	session1, err := mods.Identity.Service.Session.Create(ctx, "test@test.com", "StrongPass123!")
	require.NoError(t, err)

	session2, err := mods.Identity.Service.Session.Create(ctx, "test@test.com", "StrongPass123!")
	require.NoError(t, err)

	session3, err := mods.Identity.Service.Session.Create(ctx, "test@test.com", "StrongPass123!")
	require.NoError(t, err)

	tests := []struct {
		name           string
		cookie         string
		expectedStatus int
		err            error
		checkSessions  bool
	}{
		{
			name:           "should delete all other sessions successfully",
			cookie:         session1.Token,
			expectedStatus: http.StatusNoContent,
			checkSessions:  true,
		},
		{
			name:           "should reject invalid session token",
			cookie:         session2.Token,
			expectedStatus: http.StatusUnauthorized,
			err:            httpx.ErrUnauthenticated,
		},
		{
			name:           "should reject missing session token",
			cookie:         "",
			expectedStatus: http.StatusUnauthorized,
			err:            httpx.ErrUnauthenticated,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cookie := ""
			if tc.cookie != "" {
				cookie = fmt.Sprintf("Cookie: session=%s", tc.cookie)
			}

			resp := api.Delete(deleteAllSessionsPath, cookie)
			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.err != nil {
				httpx.AssertErr(t, tc.err, resp.Body)
				return
			}

			if tc.checkSessions {
				// Verify session2 is invalidated
				session, err := mods.Identity.Service.Session.GetByToken(ctx, session1.Token)
				assert.Equal(t, session, session1)
				require.NoError(t, err)

				_, err = mods.Identity.Service.Session.GetByToken(ctx, session2.Token)
				assert.Error(t, err)

				_, err = mods.Identity.Service.Session.GetByToken(ctx, session3.Token)
				assert.Error(t, err)
			}
		})
	}
}

func TestDeleteSession(t *testing.T) {
	t.Parallel()
	deleteSessionPath := func(id string) string { return BasePath("/identity/sessions/" + id) }
	api, container, mods := testutil.Container(t)

	ctx := context.Background()
	auth := identity.NewAuthentication(container, api, mods.Identity.Service)
	err := AddRoutes(container, api, mods.Identity.Service, auth)
	assert.NoError(t, err)

	// Create test users and multiple sessions
	now := time.Now()
	_, err = mods.Identity.Service.User.Create(ctx, &model.User{
		Email:           "test@test.com",
		Name:            "Test User",
		EmailVerifiedAt: &now,
	}, "StrongPass123!")
	require.NoError(t, err)

	_, err = mods.Identity.Service.User.Create(ctx, &model.User{
		Email:           "test2@test.com",
		Name:            "Test User 2",
		EmailVerifiedAt: &now,
	}, "StrongPass123!")
	require.NoError(t, err)

	// Create two sessions for the same user
	session1, err := mods.Identity.Service.Session.Create(ctx, "test@test.com", "StrongPass123!")
	require.NoError(t, err)

	session2, err := mods.Identity.Service.Session.Create(ctx, "test@test.com", "StrongPass123!")
	require.NoError(t, err)

	session3, err := mods.Identity.Service.Session.Create(ctx, "test@test.com", "StrongPass123!")
	require.NoError(t, err)

	session4, err := mods.Identity.Service.Session.Create(ctx, "test2@test.com", "StrongPass123!")
	require.NoError(t, err)

	tests := []struct {
		name           string
		cookie         string
		session        string
		expectedStatus int
		err            error
		checkRemoved   string
		checkPresent   string
	}{
		{
			name:           "should delete other session successfully",
			cookie:         session1.Token,
			session:        session2.ID,
			expectedStatus: http.StatusNoContent,
			checkRemoved:   session2.Token,
			checkPresent:   session3.Token,
		},
		{
			name:           "should no-op on deleting invalid session",
			cookie:         session1.Token,
			session:        session2.ID,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "should reject missing session token",
			cookie:         "",
			session:        session3.ID,
			expectedStatus: http.StatusUnauthorized,
			err:            httpx.ErrUnauthenticated,
		},
		{
			name:           "should no-op on deleting session of another user",
			cookie:         session1.Token,
			session:        session4.ID,
			expectedStatus: http.StatusNoContent,
			checkPresent:   session4.Token,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cookie := ""
			if tc.cookie != "" {
				cookie = fmt.Sprintf("Cookie: session=%s", tc.cookie)
			}

			resp := api.Delete(deleteSessionPath(tc.session), cookie)
			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.err != nil {
				httpx.AssertErr(t, tc.err, resp.Body)
				return
			}

			if tc.checkRemoved != "" {
				// Verify session is invalidated
				_, err := mods.Identity.Service.Session.GetByToken(ctx, tc.checkRemoved)
				assert.Error(t, err)
			}

			if tc.checkPresent != "" {
				// Verify session is valid
				session, err := mods.Identity.Service.Session.GetByToken(ctx, tc.checkPresent)
				assert.Equal(t, tc.checkPresent, session.Token)
				require.NoError(t, err)

			}
		})
	}
}

func TestRefreshSession(t *testing.T) {
	t.Parallel()
	refreshSessionPath := BasePath("/identity/refresh-session")
	api, container, mods := testutil.Container(t)

	ctx := context.Background()
	auth := identity.NewAuthentication(container, api, mods.Identity.Service)
	err := AddRoutes(container, api, mods.Identity.Service, auth)
	require.NoError(t, err)

	// Create a test user and session
	now := time.Now()
	_, err = mods.Identity.Service.User.Create(ctx, &model.User{
		Email:           "test@test.com",
		Name:            "Test User",
		EmailVerifiedAt: &now,
	}, "StrongPass123!")
	require.NoError(t, err)

	session, err := mods.Identity.Service.Session.Create(ctx, "test@test.com", "StrongPass123!")
	require.NoError(t, err)

	tests := []struct {
		name           string
		cookie         string
		expectedStatus int
		err            error
		checkCookie    bool
	}{
		{
			name:           "should refresh session successfully",
			cookie:         session.RefreshToken,
			expectedStatus: http.StatusNoContent,
			checkCookie:    true,
		},
		{
			name:           "should reject invalid refresh token",
			cookie:         "invalid-token",
			expectedStatus: http.StatusUnauthorized,
			err:            httpx.ErrInvalidRefreshToken,
		},
		{
			name:           "should reject missing refresh token",
			cookie:         "",
			expectedStatus: http.StatusUnauthorized,
			err:            httpx.ErrInvalidRefreshToken,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cookie := ""
			if tc.cookie != "" {
				cookie = fmt.Sprintf("Cookie: refresh_token=%s", tc.cookie)
			}

			resp := api.Post(refreshSessionPath, cookie)
			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.err != nil {
				httpx.AssertErr(t, tc.err, resp.Body)
				return
			}

			if tc.checkCookie {
				cookie := resp.Header().Get("Set-Cookie")
				assert.Contains(t, cookie, "session=")
				assert.Contains(t, cookie, "HttpOnly")
				assert.Contains(t, cookie, "SameSite=Lax")

				// Verify old session is invalidated
				_, err := mods.Identity.Service.Session.GetByToken(ctx, session.Token)
				assert.Error(t, err)
			}
		})
	}
}

func TestMe(t *testing.T) {
	t.Parallel()
	mePath := BasePath("/identity/me")
	api, container, mods := testutil.Container(t)

	ctx := context.Background()
	auth := identity.NewAuthentication(container, api, mods.Identity.Service)
	err := AddRoutes(container, api, mods.Identity.Service, auth)
	assert.NoError(t, err)

	// Create a test user and session
	now := time.Now()
	user, err := mods.Identity.Service.User.Create(ctx, &model.User{
		Email:           "test@test.com",
		Name:            "Test User",
		EmailVerifiedAt: &now,
		LastActiveAt:    &now,
		LastLoggedInAt:  &now,
	}, "StrongPass123!")
	require.NoError(t, err)

	session, err := mods.Identity.Service.Session.Create(ctx, "test@test.com", "StrongPass123!")
	require.NoError(t, err)

	tests := []struct {
		name           string
		cookie         string
		expectedStatus int
		err            error
	}{
		{
			name:           "should get current user successfully",
			cookie:         session.Token,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "should reject invalid session token",
			cookie:         "invalid-token",
			expectedStatus: http.StatusUnauthorized,
			err:            httpx.ErrUnauthenticated,
		},
		{
			name:           "should reject missing session token",
			cookie:         "",
			expectedStatus: http.StatusUnauthorized,
			err:            httpx.ErrUnauthenticated,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cookie := ""
			if tc.cookie != "" {
				cookie = fmt.Sprintf("Cookie: session=%s", tc.cookie)
			}

			resp := api.Get(mePath, cookie)
			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.err != nil {
				httpx.AssertErr(t, tc.err, resp.Body)
				return
			}

			if tc.expectedStatus == http.StatusOK {
				var response struct {
					User struct {
						ID                 string     `json:"id"`
						Email              string     `json:"email"`
						Name               string     `json:"name"`
						Image              *string    `json:"image,omitempty"`
						IsVerified         bool       `json:"isVerified"`
						IsTwoFactorEnabled bool       `json:"isTwoFactorEnabled"`
						LastActiveAt       *time.Time `json:"lastActiveAt,omitempty"`
						LastLoggedInAt     *time.Time `json:"lastLoggedInAt,omitempty"`
						SessionExpiresAt   time.Time  `json:"sessionExpiresAt"`
					} `json:"user"`
				}
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)

				assert.Equal(t, user.ID, response.User.ID)
				assert.Equal(t, user.Email, response.User.Email)
				assert.Equal(t, user.Name, response.User.Name)
				assert.Equal(t, user.Image, response.User.Image)
				assert.True(t, response.User.IsVerified)
				assert.False(t, response.User.IsTwoFactorEnabled) // No 2FA in test setup
			}
		})
	}
}
