package v1

import (
	"autopilot/backends/api/internal/identity"
	"autopilot/backends/api/internal/identity/model"
	"autopilot/backends/api/pkg/app/mocks"
	"autopilot/backends/api/pkg/httpx"
	"autopilot/backends/api/pkg/testutil"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	_ "embed"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/sample.svg
var badProfile []byte

//go:embed testdata/profile.png
var goodProfile []byte

func TestGetUser(t *testing.T) {
	t.Parallel()
	getUserPath := BasePath("/users/")
	api, container, mods := testutil.Container(t)

	ctx := context.Background()
	auth := identity.NewAuthentication(container, api, mods.Identity.Service)
	err := AddRoutes(container, api, mods.Identity.Service, auth)
	require.NoError(t, err)

	now := time.Now()
	user, err := mods.Identity.Service.User.Create(ctx, &model.User{
		Email:           "test@test.com",
		Name:            "Test User",
		EmailVerifiedAt: &now,
	}, "StrongPass123!")
	require.NoError(t, err)

	session, err := mods.Identity.Service.Session.Create(ctx, "test@test.com", "StrongPass123!")
	require.NoError(t, err)

	tests := []struct {
		name           string
		noCookie       bool
		userID         string
		expectedStatus int
		err            error
	}{
		{
			name:           "should return user by id successfully",
			userID:         user.ID,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "should return user by @me successfully",
			userID:         "@me",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "should fail for invalid user ID",
			userID:         "bad",
			expectedStatus: http.StatusNotFound,
			err:            httpx.ErrUserNotFound,
		},
		{
			name:           "should fail for invalid token",
			noCookie:       true,
			userID:         "@me",
			expectedStatus: http.StatusUnauthorized,
			err:            httpx.ErrUnauthenticated,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var args []any
			if !tc.noCookie {
				args = append(args, "Cookie: session="+session.Token)
			}
			resp := api.Get(getUserPath+tc.userID, args...)
			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.err != nil {
				httpx.AssertErr(t, tc.err, resp.Body)
				return
			}

			var response User
			err := json.NewDecoder(resp.Body).Decode(&response)
			assert.NoError(t, err)

			assert.Equal(t, user.ID, response.ID)
			assert.Equal(t, user.Email, response.Email)
			assert.Equal(t, user.Name, response.Name)
			assert.Equal(t, user.Image, response.Image)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	getUserPath := BasePath("/users/")
	api, container, mods := testutil.Container(t)

	ctx := context.Background()
	auth := identity.NewAuthentication(container, api, mods.Identity.Service)
	err := AddRoutes(container, api, mods.Identity.Service, auth)
	require.NoError(t, err)

	now := time.Now()
	user, err := mods.Identity.Service.User.Create(ctx, &model.User{
		Email:           "test@test.com",
		Name:            "Test User",
		EmailVerifiedAt: &now,
	}, "StrongPass123!")
	require.NoError(t, err)

	session, err := mods.Identity.Service.Session.Create(ctx, "test@test.com", "StrongPass123!")
	require.NoError(t, err)

	user.Name = "New Name"
	tests := []struct {
		name           string
		noCookie       bool
		payload        map[string]any
		userID         string
		expectedStatus int
		err            error
	}{
		{
			name:   "should succeed for valid update",
			userID: user.ID,
			payload: map[string]any{
				"name": "New Name",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "should succeed for valid update",
			userID: "@me",
			payload: map[string]any{
				"name": "New Name",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "should fail for valid name",
			userID: "@me",
			payload: map[string]any{
				"name": "T",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name:           "should fail for invalid user ID",
			userID:         "bad",
			payload:        map[string]any{"name": "fail"},
			expectedStatus: http.StatusNotFound,
			err:            httpx.ErrUserNotFound,
		},
		{
			name:           "should fail for invalid token",
			noCookie:       true,
			userID:         "@me",
			payload:        map[string]any{"name": "fail"},
			expectedStatus: http.StatusUnauthorized,
			err:            httpx.ErrUnauthenticated,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			args := []any{tc.payload}
			if !tc.noCookie {
				args = append(args, "Cookie: session="+session.Token)
			}
			resp := api.Put(getUserPath+tc.userID, args...)
			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.err != nil {
				httpx.AssertErr(t, tc.err, resp.Body)
				return
			}

			var response User
			err := json.NewDecoder(resp.Body).Decode(&response)
			assert.NoError(t, err)

			assert.Equal(t, user.ID, response.ID)
			assert.Equal(t, user.Email, response.Email)
			assert.Equal(t, user.Name, response.Name)
			assert.Equal(t, user.Image, response.Image)
		})
	}
}

func TestUpdateUserImage(t *testing.T) {
	t.Parallel()
	getUserPath := BasePath("/users/")
	api, container, mods := testutil.Container(t)

	ctx := context.Background()
	auth := identity.NewAuthentication(container, api, mods.Identity.Service)
	err := AddRoutes(container, api, mods.Identity.Service, auth)
	require.NoError(t, err)

	now := time.Now()
	user, err := mods.Identity.Service.User.Create(ctx, &model.User{
		Email:           "test@test.com",
		Name:            "Test User",
		EmailVerifiedAt: &now,
	}, "StrongPass123!")
	require.NoError(t, err)

	session, err := mods.Identity.Service.Session.Create(ctx, "test@test.com", "StrongPass123!")
	require.NoError(t, err)

	user.Name = "New Name"
	tests := []struct {
		name           string
		noCookie       bool
		fileName       string
		payload        []byte
		userID         string
		expectedStatus int
		err            error
	}{
		{
			name:           "should succeed for valid update",
			userID:         user.ID,
			payload:        goodProfile,
			fileName:       "profile.png",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "should fail for no file name",
			userID:         user.ID,
			payload:        goodProfile,
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name:           "should fail for invalid file input",
			userID:         user.ID,
			payload:        badProfile,
			fileName:       "profile.png",
			expectedStatus: http.StatusBadRequest,
			err:            httpx.ErrInvalidImageFormat,
		},
		{
			name:           "should fail for invalid user ID",
			userID:         "bad",
			payload:        goodProfile,
			fileName:       "profile.png",
			expectedStatus: http.StatusNotFound,
			err:            httpx.ErrUserNotFound,
		},
		{
			name:           "should fail for invalid token",
			noCookie:       true,
			userID:         "@me",
			payload:        goodProfile,
			fileName:       "profile.png",
			expectedStatus: http.StatusUnauthorized,
			err:            httpx.ErrUnauthenticated,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			args := []any{bytes.NewReader(tc.payload)}
			if !tc.noCookie {
				args = append(args, "Cookie: session="+session.Token)
			}
			if tc.fileName != "" {
				args = append(args, "X-File-Name: "+tc.fileName)
			}
			resp := api.Post(getUserPath+tc.userID+"/image", args...)
			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.err != nil {
				httpx.AssertErr(t, tc.err, resp.Body)
				return
			}

			assert.Empty(t, resp.Body.String())
		})
	}
}

func TestSignUp(t *testing.T) {
	t.Parallel()
	signUpPath := BasePath("/identity/sign-up")
	api, container, mods := testutil.Container(t)

	auth := identity.NewAuthentication(container, api, mods.Identity.Service)
	err := AddRoutes(container, api, mods.Identity.Service, auth)
	assert.NoError(t, err)

	mockWorker := container.Worker.(*testutil.MockWorker)

	mockTurnstile := container.Turnstile.(*mocks.MockTurnstiler)
	mockTurnstile.On("Verify", mock.Anything, "XXX.DUMMY.TOKEN", "").Return(true, nil)
	mockTurnstile.On("Verify", mock.Anything, "INVALID.TOKEN", "").Return(false, nil)
	mockTurnstile.On("Verify", mock.Anything, "ERROR.TOKEN", "").Return(false, assert.AnError)

	tests := []struct {
		name           string
		payload        map[string]any
		expectedStatus int
		err            error
	}{
		{
			name: "should sign up successfully",
			payload: map[string]any{
				"cfTurnstileToken": "XXX.DUMMY.TOKEN",
				"email":            "test@test.com",
				"name":             "Test User",
				"password":         "StrongPass123!",
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "should reject invalid email format",
			payload: map[string]any{
				"cfTurnstileToken": "XXX.DUMMY.TOKEN",
				"email":            "invalid-email",
				"name":             "Test User",
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
			name: "should reject password too short",
			payload: map[string]any{
				"cfTurnstileToken": "XXX.DUMMY.TOKEN",
				"email":            "test@test.com",
				"name":             "Test User",
				"password":         "123",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject name too short",
			payload: map[string]any{
				"cfTurnstileToken": "XXX.DUMMY.TOKEN",
				"email":            "test@test.com",
				"name":             "T",
				"password":         "StrongPass123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject name too long",
			payload: map[string]any{
				"cfTurnstileToken": "XXX.DUMMY.TOKEN",
				"email":            "test@test.com",
				"name":             strings.Repeat("a", 101),
				"password":         "StrongPass123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject password missing number",
			payload: map[string]any{
				"cfTurnstileToken": "XXX.DUMMY.TOKEN",
				"email":            "test@test.com",
				"name":             "Test User",
				"password":         "StrongPass!!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject password missing special char",
			payload: map[string]any{
				"cfTurnstileToken": "XXX.DUMMY.TOKEN",
				"email":            "test@test.com",
				"name":             "Test User",
				"password":         "StrongPass123",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject password missing uppercase",
			payload: map[string]any{
				"cfTurnstileToken": "XXX.DUMMY.TOKEN",
				"email":            "test@test.com",
				"name":             "Test User",
				"password":         "strongpass123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject password missing lowercase",
			payload: map[string]any{
				"cfTurnstileToken": "XXX.DUMMY.TOKEN",
				"email":            "test@test.com",
				"name":             "Test User",
				"password":         "STRONGPASS123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject invalid turnstile token",
			payload: map[string]any{
				"cfTurnstileToken": "INVALID.TOKEN",
				"email":            "test@test.com",
				"name":             "Test User",
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
				"name":             "Test User",
				"password":         "StrongPass123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockWorker.Reset()
			resp := api.Post(signUpPath, tc.payload)
			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.err != nil {
				httpx.AssertErr(t, tc.err, resp.Body)
				return
			}
			assert.Empty(t, resp.Body.String())

			// verify successful email
			mailerRequests := mockWorker.MailRequests
			assert.Len(t, mailerRequests, 1)
			req := mailerRequests[0]
			assert.Equal(t, req.Email, "test@test.com")

			newVerification, err := mods.Identity.Store.User.GetVerificationByValue(context.Background(), model.VerificationContextEmailVerification, req.Email)
			require.NoError(t, err)

			assert.Contains(t, req.Data["VerificationURL"], newVerification.ID)
		})
	}

	// Test duplicate email separately as it requires two API calls
	t.Run("duplicate email", func(t *testing.T) {
		payload := map[string]any{
			"cfTurnstileToken": "XXX.DUMMY.TOKEN",
			"email":            "duplicate@test.com",
			"name":             "Test User",
			"password":         "StrongPass123!",
		}

		// First signup should succeed
		resp := api.Post(signUpPath, payload)
		assert.Equal(t, http.StatusNoContent, resp.Code)
		assert.Empty(t, resp.Body.String())

		// Second signup with same email should fail
		resp = api.Post(signUpPath, payload)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
		httpx.AssertErr(t, httpx.ErrEmailExists, resp.Body)
	})

	// Verify all Turnstile mock expectations were met
	mockTurnstile.AssertExpectations(t)
}

func TestVerifyEmail(t *testing.T) {
	t.Parallel()
	verifyEmailPath := BasePath("/identity/verify-email")
	api, container, mods := testutil.Container(t)

	auth := identity.NewAuthentication(container, api, mods.Identity.Service)
	err := AddRoutes(container, api, mods.Identity.Service, auth)
	require.NoError(t, err)

	// Create a test user
	ctx := context.Background()
	user, err := mods.Identity.Service.User.Create(ctx, &model.User{
		Email: "test@test.com",
		Name:  "Test User",
	}, "StrongPass123!")
	require.NoError(t, err)

	// Create a verification for the test user
	verification, err := mods.Identity.Store.Verification.GetByValue(ctx, model.VerificationContextEmailVerification, user.Email)
	require.NoError(t, err)

	tests := []struct {
		name           string
		payload        map[string]any
		expectedStatus int
		err            error
	}{
		{
			name: "should verify email successfully",
			payload: map[string]any{
				"token": verification.ID,
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "should reject invalid token",
			payload: map[string]any{
				"token": "invalid-token",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject wrong token",
			payload: map[string]any{
				"token": uuid.NewString(),
			},
			expectedStatus: http.StatusUnauthorized,
			err:            httpx.ErrInvalidOrExpiredToken,
		},
		{
			name: "should reject already verified token",
			payload: map[string]any{
				"token": verification.ID,
			},
			expectedStatus: http.StatusUnauthorized,
			err:            httpx.ErrInvalidOrExpiredToken,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := api.Post(verifyEmailPath, tc.payload)
			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.err != nil {
				httpx.AssertErr(t, tc.err, resp.Body)
			}
		})
	}
}

func TestForgotPassword(t *testing.T) {
	t.Parallel()
	forgotPasswordPath := BasePath("/identity/forgot-password")
	api, container, mods := testutil.Container(t)

	auth := identity.NewAuthentication(container, api, mods.Identity.Service)
	err := AddRoutes(container, api, mods.Identity.Service, auth)
	assert.NoError(t, err)

	// Get the mock Turnstile instance
	mockTurnstile := container.Turnstile.(*mocks.MockTurnstiler)
	// Set up default success behavior for Turnstile
	mockTurnstile.On("Verify", mock.Anything, "XXX.DUMMY.TOKEN", "").Return(true, nil)
	// Set up failure behavior for invalid token
	mockTurnstile.On("Verify", mock.Anything, "INVALID.TOKEN", "").Return(false, nil)
	// Set up error behavior for error token
	mockTurnstile.On("Verify", mock.Anything, "ERROR.TOKEN", "").Return(false, assert.AnError)

	tests := []struct {
		name           string
		payload        map[string]any
		expectedStatus int
		err            error
	}{
		{
			name: "should request password reset successfully",
			payload: map[string]any{
				"cfTurnstileToken": "XXX.DUMMY.TOKEN",
				"email":            "test@test.com",
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "should reject invalid email format",
			payload: map[string]any{
				"cfTurnstileToken": "XXX.DUMMY.TOKEN",
				"email":            "invalid-email",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject missing email",
			payload: map[string]any{
				"cfTurnstileToken": "XXX.DUMMY.TOKEN",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject invalid turnstile token",
			payload: map[string]any{
				"cfTurnstileToken": "INVALID.TOKEN",
				"email":            "test@test.com",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject turnstile verification error",
			payload: map[string]any{
				"cfTurnstileToken": "ERROR.TOKEN",
				"email":            "test@test.com",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should handle non-existent email gracefully",
			payload: map[string]any{
				"cfTurnstileToken": "XXX.DUMMY.TOKEN",
				"email":            "nonexistent@test.com",
			},
			expectedStatus: http.StatusNoContent, // Should still return success to prevent email enumeration
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := api.Post(forgotPasswordPath, tc.payload)
			assert.Equal(t, tc.expectedStatus, resp.Code)
			if tc.err != nil {
				httpx.AssertErr(t, tc.err, resp.Body)
				return
			}
		})
	}

	// Verify all Turnstile mock expectations were met
	mockTurnstile.AssertExpectations(t)
}

func TestResetPassword(t *testing.T) {
	t.Parallel()
	resetPasswordPath := BasePath("/identity/reset-password")
	api, container, mods := testutil.Container(t)

	auth := identity.NewAuthentication(container, api, mods.Identity.Service)
	err := AddRoutes(container, api, mods.Identity.Service, auth)
	require.NoError(t, err)

	// Create a test user
	ctx := context.Background()
	user, err := mods.Identity.Service.User.Create(ctx, &model.User{
		Email: "test@test.com",
		Name:  "Test User",
	}, "OldPass123!")
	require.NoError(t, err)

	// Get the email verification record
	verification, err := mods.Identity.Store.Verification.GetByValue(ctx, model.VerificationContextEmailVerification, user.Email)
	require.NoError(t, err)

	// Verify the user's email
	err = mods.Identity.Service.User.VerifyEmail(ctx, verification.ID)
	require.NoError(t, err)

	// Create a password reset verification
	now := time.Now()
	passwordResetVerification := &model.Verification{
		Context:   model.VerificationContextPasswordReset,
		Value:     user.Email,
		ExpiresAt: now.Add(model.PasswordResetDuration),
		CreatedAt: now,
		UpdatedAt: now,
	}
	newPasswordResetVerification, err := mods.Identity.Store.User.CreateVerification(ctx, passwordResetVerification)
	require.NoError(t, err)

	// Create an expired verification for testing
	expiredVerification := &model.Verification{
		Context:   model.VerificationContextPasswordReset,
		Value:     user.Email,
		ExpiresAt: now.Add(-time.Hour), // Expired 1 hour ago
		CreatedAt: now.Add(-time.Hour),
		UpdatedAt: now.Add(-time.Hour),
	}
	newExpiredVerification, err := mods.Identity.Store.User.CreateVerification(ctx, expiredVerification)
	require.NoError(t, err)

	tests := []struct {
		name           string
		payload        map[string]any
		expectedStatus int
		err            error
	}{
		{
			name: "should reset password successfully",
			payload: map[string]any{
				"token":       newPasswordResetVerification.ID,
				"newPassword": "NewStrongPass123!",
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "should reject invalid token",
			payload: map[string]any{
				"token":       "invalid-token",
				"newPassword": "NewStrongPass123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject expired token",
			payload: map[string]any{
				"token":       newExpiredVerification.ID,
				"newPassword": "NewStrongPass123!",
			},
			expectedStatus: http.StatusUnauthorized,
			err:            httpx.ErrInvalidOrExpiredToken,
		},
		{
			name: "should reject missing token",
			payload: map[string]any{
				"newPassword": "NewStrongPass123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject missing new password",
			payload: map[string]any{
				"token": newPasswordResetVerification.ID,
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject weak password - too short",
			payload: map[string]any{
				"token":       newPasswordResetVerification.ID,
				"newPassword": "Weak1!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject weak password - missing uppercase",
			payload: map[string]any{
				"token":       newPasswordResetVerification.ID,
				"newPassword": "weakpass123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject weak password - missing lowercase",
			payload: map[string]any{
				"token":       newPasswordResetVerification.ID,
				"newPassword": "STRONGPASS123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject weak password - missing number",
			payload: map[string]any{
				"token":       newPasswordResetVerification.ID,
				"newPassword": "StrongPass!!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject weak password - missing special character",
			payload: map[string]any{
				"token":       newPasswordResetVerification.ID,
				"newPassword": "StrongPass123",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := api.Post(resetPasswordPath, tc.payload)
			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.err != nil {
				httpx.AssertErr(t, tc.err, resp.Body)
				return
			}
		})
	}

	// Get the mock Turnstile instance
	mockTurnstile := container.Turnstile.(*mocks.MockTurnstiler)
	// Set up default success behavior for Turnstile
	mockTurnstile.On("Verify", mock.Anything, "XXX.DUMMY.TOKEN", "").Return(true, nil)

	// Test that the password was actually changed
	t.Run("verify password was changed", func(t *testing.T) {
		// Try to sign in with old password
		resp := api.Post(BasePath("/identity/sign-in"), map[string]any{
			"cfTurnstileToken": "XXX.DUMMY.TOKEN",
			"email":            user.Email,
			"password":         "OldPass123!",
		})
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		httpx.AssertErr(t, httpx.ErrInvalidCredentials, resp.Body)

		// Try to sign in with new password
		resp = api.Post(BasePath("/identity/sign-in"), map[string]any{
			"cfTurnstileToken": "XXX.DUMMY.TOKEN",
			"email":            user.Email,
			"password":         "NewStrongPass123!",
		})
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	// Verify all Turnstile mock expectations were met
	mockTurnstile.AssertExpectations(t)
}

func TestUpdatePassword(t *testing.T) {
	t.Parallel()
	updatePasswordPath := BasePath("/identity/update-password")
	api, container, mods := testutil.Container(t)

	ctx := context.Background()
	auth := identity.NewAuthentication(container, api, mods.Identity.Service)
	err := AddRoutes(container, api, mods.Identity.Service, auth)
	require.NoError(t, err)

	// Create a test user with verified email
	now := time.Now()
	newUser, err := mods.Identity.Service.User.Create(ctx, &model.User{
		Email:           "test@test.com",
		Name:            "Test User",
		EmailVerifiedAt: &now,
	}, "CurrentPass123!")
	require.NoError(t, err)

	// Create a session for the test user
	newSession, err := mods.Identity.Service.Session.Create(ctx, "test@test.com", "CurrentPass123!")
	require.NoError(t, err)

	tests := []struct {
		name           string
		cookie         string
		payload        map[string]any
		expectedStatus int
		err            error
	}{
		{
			name:   "should update password successfully",
			cookie: newSession.Token,
			payload: map[string]any{
				"currentPassword": "CurrentPass123!",
				"newPassword":     "NewStrongPass123!",
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:   "should reject invalid current password",
			cookie: newSession.Token,
			payload: map[string]any{
				"currentPassword": "WrongPass123!",
				"newPassword":     "NewStrongPass123!",
			},
			expectedStatus: http.StatusUnauthorized,
			err:            httpx.ErrInvalidCredentials,
		},
		{
			name:   "should reject weak new password - too short",
			cookie: newSession.Token,
			payload: map[string]any{
				"currentPassword": "CurrentPass123!",
				"newPassword":     "Weak1!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name:   "should reject weak new password - missing uppercase",
			cookie: newSession.Token,
			payload: map[string]any{
				"currentPassword": "CurrentPass123!",
				"newPassword":     "weakpass123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name:   "should reject weak new password - missing lowercase",
			cookie: newSession.Token,
			payload: map[string]any{
				"currentPassword": "CurrentPass123!",
				"newPassword":     "STRONGPASS123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name:   "should reject weak new password - missing number",
			cookie: newSession.Token,
			payload: map[string]any{
				"currentPassword": "CurrentPass123!",
				"newPassword":     "StrongPass!!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name:   "should reject weak new password - missing special character",
			cookie: newSession.Token,
			payload: map[string]any{
				"currentPassword": "CurrentPass123!",
				"newPassword":     "StrongPass123",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name:   "should reject missing current password",
			cookie: newSession.Token,
			payload: map[string]any{
				"newPassword": "NewStrongPass123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name:   "should reject missing new password",
			cookie: newSession.Token,
			payload: map[string]any{
				"currentPassword": "CurrentPass123!",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name:   "should reject invalid session",
			cookie: "invalid-token",
			payload: map[string]any{
				"currentPassword": "CurrentPass123!",
				"newPassword":     "NewStrongPass123!",
			},
			expectedStatus: http.StatusUnauthorized,
			err:            httpx.ErrUnauthenticated,
		},
		{
			name: "should reject missing session",
			payload: map[string]any{
				"currentPassword": "CurrentPass123!",
				"newPassword":     "NewStrongPass123!",
			},
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

			args := []any{tc.payload}
			if cookie != "" {
				args = append(args, cookie)
			}

			resp := api.Post(updatePasswordPath, args...)
			assert.Equal(t, tc.expectedStatus, resp.Code)
			if tc.err != nil {
				httpx.AssertErr(t, tc.err, resp.Body)
				return
			}
		})
	}

	// Get the mock Turnstile instance
	mockTurnstile := container.Turnstile.(*mocks.MockTurnstiler)
	// Set up default success behavior for Turnstile
	mockTurnstile.On("Verify", mock.Anything, "XXX.DUMMY.TOKEN", "").Return(true, nil)

	// Test that the password was actually changed
	t.Run("verify password was changed", func(t *testing.T) {
		// Try to sign in with old password
		resp := api.Post(BasePath("/identity/sign-in"), map[string]any{
			"cfTurnstileToken": "XXX.DUMMY.TOKEN",
			"email":            newUser.Email,
			"password":         "CurrentPass123!",
		})
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		httpx.AssertErr(t, httpx.ErrInvalidCredentials, resp.Body)

		// Try to sign in with new password
		resp = api.Post(BasePath("/identity/sign-in"), map[string]any{
			"cfTurnstileToken": "XXX.DUMMY.TOKEN",
			"email":            newUser.Email,
			"password":         "NewStrongPass123!",
		})
		assert.Equal(t, http.StatusOK, resp.Code)
	})
}
