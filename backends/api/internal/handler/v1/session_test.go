package v1

import (
	"autopilot/backends/api/internal/model"
	"autopilot/backends/api/internal/testutil"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignIn(t *testing.T) {
	var signInPath = BasePath("/identity/sign-in")
	api, container, service, storeManager := testutil.NewMocks(t)
	defer container.Close()

	ctx := context.Background()
	err := AddRoutes(container, api, service)
	assert.Nil(t, err)

	// Create a test user with unverified email
	unverifiedUser, err := service.User.Create(ctx, &model.User{
		Email: "unverified@test.com",
		Name:  "Unverified User",
	}, "StrongPass123!")
	assert.NotNil(t, unverifiedUser)
	assert.Nil(t, err)

	// Get the initial verification
	initialVerification, err := storeManager.User.GetVerificationByValue(ctx, model.VerificationContextEmailVerification, unverifiedUser.Email)
	assert.NotNil(t, initialVerification)
	assert.Nil(t, err)

	// Create a test user with verified email
	now := time.Now()
	user, err := service.User.Create(ctx, &model.User{
		Email:           "test@test.com",
		Name:            "Test User",
		EmailVerifiedAt: &now,
	}, "StrongPass123!")
	assert.NotNil(t, user)
	assert.Nil(t, err)

	// Get the mock Turnstile instance
	mockTurnstile := container.Turnstile.(*testutil.MockTurnstile)
	// Set up default success behavior for Turnstile
	mockTurnstile.On("Verify", mock.Anything, "XXX.DUMMY.TOKEN", "").Return(true, nil)
	// Set up failure behavior for invalid token
	mockTurnstile.On("Verify", mock.Anything, "INVALID.TOKEN", "").Return(false, nil)
	// Set up error behavior for error token
	mockTurnstile.On("Verify", mock.Anything, "ERROR.TOKEN", "").Return(false, assert.AnError)

	// Test expired verification
	t.Run("unverified email with expired verification", func(t *testing.T) {
		// Delete the initial verification to simulate expiration
		err := storeManager.User.DeleteVerification(ctx, initialVerification.ID)
		assert.Nil(t, err)

		resp := api.Post(signInPath, map[string]any{
			"cfTurnstileToken": "XXX.DUMMY.TOKEN",
			"email":            "unverified@test.com",
			"password":         "StrongPass123!",
		})
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Equal(t, `{"code":"identity.email_not_verified","errors":[],"message":"Email verification required"}`, strings.Trim(resp.Body.String(), "\n"))

		// Verify a new verification was created
		newVerification, err := storeManager.User.GetVerificationByValue(ctx, model.VerificationContextEmailVerification, unverifiedUser.Email)
		assert.Nil(t, err)
		assert.NotNil(t, newVerification)
		assert.NotEqual(t, initialVerification.ID, newVerification.ID)
		assert.True(t, newVerification.ExpiresAt.After(time.Now()))
	})

	tests := []struct {
		name           string
		payload        map[string]any
		expectedStatus int
		expectedBody   string
		checkCookie    bool
	}{
		{
			name: "successful sign in",
			payload: map[string]any{
				"cfTurnstileToken": "XXX.DUMMY.TOKEN",
				"email":            "test@test.com",
				"password":         "StrongPass123!",
			},
			expectedStatus: http.StatusOK,
			checkCookie:    true,
		},
		{
			name: "invalid email format",
			payload: map[string]any{
				"cfTurnstileToken": "XXX.DUMMY.TOKEN",
				"email":            "invalid-email",
				"password":         "StrongPass123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"errors":[{"code":"INVALID_EMAIL","location":"body.email","message":"expected string to be RFC 5322 email: mail: missing '@' or angle-addr"}],"message":"validation failed"}`,
		},
		{
			name: "missing required fields",
			payload: map[string]any{
				"cfTurnstileToken": "XXX.DUMMY.TOKEN",
				"email":            "test@test.com",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"errors":[{"code":"REQUIRED","location":"body.password","message":"expected required property password to be present"}],"message":"validation failed"}`,
		},
		{
			name: "invalid credentials",
			payload: map[string]any{
				"cfTurnstileToken": "XXX.DUMMY.TOKEN",
				"email":            "test@test.com",
				"password":         "WrongPass123!",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"code":"identity.invalid_credentials","errors":[],"message":"Invalid credentials"}`,
		},
		{
			name: "invalid turnstile token",
			payload: map[string]any{
				"cfTurnstileToken": "INVALID.TOKEN",
				"email":            "test@test.com",
				"password":         "StrongPass123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"errors":[{"code":"FAILED_TO_VERIFY_TURNSTILE_TOKEN","location":"body.cfTurnstileToken","message":"failed to verify Turnstile token"}],"message":"validation failed"}`,
		},
		{
			name: "turnstile verification error",
			payload: map[string]any{
				"cfTurnstileToken": "ERROR.TOKEN",
				"email":            "test@test.com",
				"password":         "StrongPass123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"errors":[{"code":"UNABLE_TO_VERIFY_TURNSTILE_TOKEN","location":"body.cfTurnstileToken","message":"unable to verify Turnstile token"},{"code":"FAILED_TO_VERIFY_TURNSTILE_TOKEN","location":"body.cfTurnstileToken","message":"failed to verify Turnstile token"}],"message":"validation failed"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := api.Post(signInPath, tc.payload)
			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.expectedBody != "" {
				assert.Equal(t, tc.expectedBody, strings.Trim(resp.Body.String(), "\n"))
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
	var signOutPath = BasePath("/identity/sign-out")
	api, container, service, _ := testutil.NewMocks(t)
	defer container.Close()

	ctx := context.Background()
	err := AddRoutes(container, api, service)
	assert.NoError(t, err)

	// Create a test user and session
	now := time.Now()
	user, err := service.User.Create(ctx, &model.User{
		Email:           "test@test.com",
		Name:            "Test User",
		EmailVerifiedAt: &now,
	}, "StrongPass123!")
	assert.NotNil(t, user)
	assert.Nil(t, err)

	session, err := service.Session.Create(ctx, "test@test.com", "StrongPass123!")
	assert.NotNil(t, session)
	assert.Nil(t, err)

	// Verify session was created successfully
	verifySession, err := service.Session.GetByToken(ctx, session.Token)
	assert.NotNil(t, verifySession)
	assert.Nil(t, err)

	tests := []struct {
		name           string
		cookie         string
		expectedStatus int
		expectedBody   string
		checkCookie    bool
	}{
		{
			name:           "successful sign out",
			cookie:         session.Token,
			expectedStatus: http.StatusNoContent,
			checkCookie:    true,
		},
		{
			name:           "invalid session token",
			cookie:         "invalid-token",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"code":"identity.invalid_session","errors":[],"message":"Invalid or expired session"}`,
		},
		{
			name:           "missing session token",
			cookie:         "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"code":"identity.invalid_session","errors":[],"message":"Invalid or expired session"}`,
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

			if tc.expectedBody != "" {
				assert.Equal(t, tc.expectedBody, strings.Trim(resp.Body.String(), "\n"))
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

func TestMe(t *testing.T) {
	var mePath = BasePath("/identity/me")
	api, container, service, _ := testutil.NewMocks(t)
	defer container.Close()

	ctx := context.Background()
	err := AddRoutes(container, api, service)
	assert.NoError(t, err)

	// Create a test user and session
	now := time.Now()
	user, err := service.User.Create(ctx, &model.User{
		Email:           "test@test.com",
		Name:            "Test User",
		EmailVerifiedAt: &now,
		LastActiveAt:    &now,
		LastLoggedInAt:  &now,
	}, "StrongPass123!")
	assert.NotNil(t, user)
	assert.Nil(t, err)

	session, err := service.Session.Create(ctx, "test@test.com", "StrongPass123!")
	assert.NotNil(t, session)
	assert.Nil(t, err)

	tests := []struct {
		name           string
		cookie         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful me",
			cookie:         session.Token,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid session token",
			cookie:         "invalid-token",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"errors":[],"message":"Invalid session"}`,
		},
		{
			name:           "missing session token",
			cookie:         "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"errors":[],"message":"Invalid session"}`,
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

			if tc.expectedBody != "" {
				assert.Equal(t, tc.expectedBody, strings.Trim(resp.Body.String(), "\n"))
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
