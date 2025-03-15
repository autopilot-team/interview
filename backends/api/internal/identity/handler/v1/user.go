package v1

import (
	"autopilot/backends/api/internal/identity/model"
	"autopilot/backends/api/pkg/httpx"
	"autopilot/backends/internal/core"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"image"
	"time"

	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"
)

// User is object representing a user.
type User struct {
	ID              string     `json:"id,omitempty"`
	Name            string     `json:"name,omitempty"`
	Email           string     `json:"email,omitempty"`
	EmailVerifiedAt *time.Time `json:"emailVerifiedAt,omitempty"`
	// FailedLoginAttempts int       `json:"failedLoginAttempts,omitempty"`
	Image          *string    `json:"image,omitempty"`
	LastActiveAt   *time.Time `json:"lastActiveAt,omitempty"`
	LastLoggedInAt *time.Time `json:"lastLoggedInAt,omitempty"`
	// LockedAt            time.Time `json:"lockedAt,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// GetUserRequest is the request body for the get user endpoint.
type GetUserRequest struct {
	ID string `path:"id" required:"true" doc:"The ID of the user. Use @me to refer to the current user." example:"@me"`
}

// GetUserResponse is the response body for the get user endpoint.
type GetUserResponse struct {
	Body User
}

// GetUser is the handler for the get user endpoint.
func (v *V1) GetUser(ctx context.Context, input *GetUserRequest) (*GetUserResponse, error) {
	auth := httpx.GetAuthInfo(ctx)
	if input.ID == "@me" || input.ID == "%40me" {
		input.ID = auth.UserID
	}
	if input.ID != auth.UserID {
		return nil, httpx.ErrUserNotFound
	}

	user, err := v.identity.User.GetByID(ctx, auth.UserID)
	if err != nil {
		v.Logger.Error("Failed to get user", "error", err)
		return nil, err
	}
	return &GetUserResponse{
		Body: User{
			ID:              user.ID,
			Name:            user.Name,
			Email:           user.Email,
			EmailVerifiedAt: user.EmailVerifiedAt,
			Image:           user.Image,
			LastActiveAt:    user.LastActiveAt,
			LastLoggedInAt:  user.LastLoggedInAt,
			CreatedAt:       user.CreatedAt,
			UpdatedAt:       user.UpdatedAt,
		},
	}, nil
}

// UpdateUserRequest is the request body for the update user endpoint.
type UpdateUserRequest struct {
	ID   string `path:"id" required:"true" doc:"The ID of the user. Use @me to refer to the current user." example:"@me"`
	Body struct {
		Name string `json:"name" required:"true" minLength:"2" maxLength:"100" doc:"The user's full name" example:"John Doe"`
	}
}

// UpdateUserResponse is the response body for the update user endpoint.
type UpdateUserResponse struct {
	Body User
}

// UpdateUser is the handler for the update user endpoint.
func (v *V1) UpdateUser(ctx context.Context, input *UpdateUserRequest) (*UpdateUserResponse, error) {
	auth := httpx.GetAuthInfo(ctx)
	if input.ID == "@me" || input.ID == "%40me" {
		input.ID = auth.UserID
	}
	if input.ID != auth.UserID {
		return nil, httpx.ErrUserNotFound
	}

	user, err := v.identity.User.Update(ctx, &model.User{
		ID:   input.ID,
		Name: input.Body.Name,
	})
	if err != nil {
		v.Logger.Error("Failed to update user", "error", err)
		return nil, err
	}
	return &UpdateUserResponse{
		Body: User{
			ID:              user.ID,
			Name:            user.Name,
			Email:           user.Email,
			EmailVerifiedAt: user.EmailVerifiedAt,
			Image:           user.Image,
			LastActiveAt:    user.LastActiveAt,
			LastLoggedInAt:  user.LastLoggedInAt,
			CreatedAt:       user.CreatedAt,
			UpdatedAt:       user.UpdatedAt,
		},
	}, nil
}

// UpdateUserImageRequest is the request body for the update user image endpoint.
type UpdateUserImageRequest struct {
	ID       string `path:"id" required:"true" doc:"The ID of the user. Use @me to refer to the current user." example:"@me"`
	FileName string `header:"X-File-Name" required:"true" doc:"The original file name."`
	RawBody  []byte `contentType:"image/*" doc:"Supports png, jpeg and webp file types."`
}

// UpdateUserImageResponse is the response body for the update user image endpoint.
type UpdateUserImageResponse struct{}

// UpdateUserImage is the handler for the update user image endpoint.
func (v *V1) UpdateUserImage(ctx context.Context, input *UpdateUserImageRequest) (*UpdateUserImageResponse, error) {
	auth := httpx.GetAuthInfo(ctx)
	if input.ID == "@me" || input.ID == "%40me" {
		input.ID = auth.UserID
	}
	if input.ID != auth.UserID {
		return nil, httpx.ErrUserNotFound
	}

	key, format, err := validateImage(input.RawBody)
	if err != nil {
		return nil, err
	}

	_, err = v.Storage.Identity.Upload(ctx, key, bytes.NewReader(input.RawBody), &core.ObjectMetadata{
		Key:            key,
		Size:           int64(len(input.RawBody)),
		ContentType:    "image/" + format,
		LastModified:   time.Now(),
		CustomMetadata: map[string]string{"name": input.FileName},
	})
	if err != nil {
		v.Logger.Error("error uploading image to storage", "error", err)
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	info, err := v.Storage.Identity.GenerateDownloadURL(ctx, key, 0)
	if err != nil {
		v.Logger.Error("error generating link to storage object", "error", err)
		return nil, httpx.ErrUnknown.WithInternal(err)
	}

	if _, err := v.identity.User.Update(ctx, &model.User{
		ID:    input.ID,
		Image: &info.URL,
	}); err != nil {
		return nil, err
	}

	return &UpdateUserImageResponse{}, nil
}

// ForgotPasswordRequest is the request body for the forgot password endpoint.
type ForgotPasswordRequest struct {
	Body struct {
		CfTurnstileToken httpx.TurnstileToken `json:"cfTurnstileToken" required:"true" doc:"The Cloudflare Turnstile token" example:"XXX.DUMMY.TOKEN"`
		Email            string               `json:"email" required:"true" format:"email" doc:"The user's email address" example:"user@example.com"`
	}
}

// ForgotPasswordResponse is the response body for the forgot password endpoint.
type ForgotPasswordResponse struct{}

// ForgotPassword initiates the password reset process.
func (v *V1) ForgotPassword(ctx context.Context, input *ForgotPasswordRequest) (*ForgotPasswordResponse, error) {
	err := v.identity.User.InitiatePasswordReset(ctx, input.Body.Email)
	if err != nil {
		v.Logger.Error("Failed to initiate password reset", "error", err)
		return nil, err
	}

	return &ForgotPasswordResponse{}, nil
}

// ResetPasswordRequest is the request body for the reset password endpoint.
type ResetPasswordRequest struct {
	Body struct {
		Token       string         `json:"token" format:"uuid" required:"true" doc:"The password reset token" example:"abc123"`
		NewPassword httpx.Password `json:"newPassword" doc:"The new password" example:"NewStrongPass123!"`
	}
}

// ResetPasswordResponse is the response body for the reset password endpoint.
type ResetPasswordResponse struct{}

// ResetPassword completes the password reset process.
func (v *V1) ResetPassword(ctx context.Context, input *ResetPasswordRequest) (*ResetPasswordResponse, error) {
	err := v.identity.User.ResetPassword(ctx, input.Body.Token, string(input.Body.NewPassword))
	if err != nil {
		v.Logger.Error("Failed to reset password", "error", err)
		return nil, err
	}

	return &ResetPasswordResponse{}, nil
}

// SignUpRequest is the request body for the sign up endpoint.
type SignUpRequest struct {
	Body struct {
		CfTurnstileToken httpx.TurnstileToken `json:"cfTurnstileToken" required:"true" doc:"The Cloudflare Turnstile token" example:"XXX.DUMMY.TOKEN"`
		Email            string               `json:"email" required:"true" format:"email" doc:"The user's email address" example:"user@example.com"`
		Name             string               `json:"name" required:"true" minLength:"2" maxLength:"100" doc:"The user's full name" example:"John Doe"`
		Password         httpx.Password       `json:"password" doc:"The user's password" example:"StrongPass123!"`
	}
}

// SignUpResponse is the response body for the sign up endpoint.
type SignUpResponse struct{}

// SignUp is the handler for the sign up endpoint.
func (v *V1) SignUp(ctx context.Context, input *SignUpRequest) (*SignUpResponse, error) {
	_, err := v.identity.User.Create(ctx, &model.User{
		Email: input.Body.Email,
		Name:  input.Body.Name,
	}, string(input.Body.Password))
	if err != nil {
		v.Logger.Error("Failed to sign up user", "error", err)
		return nil, err
	}

	return &SignUpResponse{}, nil
}

// VerifyEmailRequest is the request body for the verify email endpoint.
type VerifyEmailRequest struct {
	Body struct {
		Token string `json:"token" format:"uuid" required:"true" doc:"The verification token"`
	}
}

// VerifyEmailResponse is the response body for the verify email endpoint.
type VerifyEmailResponse struct{}

// VerifyEmail is the handler for the verify email endpoint.
func (v *V1) VerifyEmail(ctx context.Context, input *VerifyEmailRequest) (*VerifyEmailResponse, error) {
	err := v.identity.User.VerifyEmail(ctx, input.Body.Token)
	if err != nil {
		v.Logger.Error("Failed to verify email", "error", err)
		return nil, err
	}

	return &VerifyEmailResponse{}, nil
}

// VerifyPasswordRequest is the request body for the verify password endpoint.
type VerifyPasswordRequest struct {
	Body struct {
		Password string `json:"password" required:"true" doc:"The current password to verify" example:"current-password"`
	}
}

// VerifyPasswordResponse is the response body for the verify password endpoint.
type VerifyPasswordResponse struct {
	Body struct {
		Verified bool `json:"verified" doc:"Whether the password was verified successfully"`
	}
}

// VerifyPassword verifies the current password for sensitive operations
func (v *V1) VerifyPassword(ctx context.Context, input *VerifyPasswordRequest) (*VerifyPasswordResponse, error) {
	auth := httpx.GetAuthInfo(ctx)
	user, err := v.identity.User.GetByID(ctx, auth.UserID)
	if err != nil {
		v.Logger.Error("Failed to get user", "error", err)
		return nil, err
	}

	if !user.VerifyPassword(input.Body.Password) {
		return nil, httpx.ErrInvalidCredentials
	}

	response := &VerifyPasswordResponse{}
	response.Body.Verified = true

	return response, nil
}

// UpdatePasswordRequest is the request body for the update password endpoint.
type UpdatePasswordRequest struct {
	Body struct {
		CurrentPassword httpx.Password `json:"currentPassword" doc:"The current password" example:"CurrentPass123!"`
		NewPassword     httpx.Password `json:"newPassword" doc:"The new password" example:"NewStrongPass123!"`
	}
}

// UpdatePasswordResponse is the response body for the update password endpoint.
type UpdatePasswordResponse struct{}

// UpdatePassword handles password update requests
func (v *V1) UpdatePassword(ctx context.Context, input *UpdatePasswordRequest) (*UpdatePasswordResponse, error) {
	auth := httpx.GetAuthInfo(ctx)
	err := v.identity.User.UpdatePassword(ctx, auth.UserID, string(input.Body.CurrentPassword), string(input.Body.NewPassword))
	if err != nil {
		v.Logger.Error("Failed to update password", "error", err)
		return nil, err
	}

	return &UpdatePasswordResponse{}, nil
}

func validateImage(img []byte) (key string, format string, err error) {
	_, format, err = image.Decode(bytes.NewBuffer(img))
	if err != nil {
		return "", "", httpx.ErrInvalidImageFormat
	}
	switch format {
	case "png", "jpeg", "webp":
	default:
		return "", "", httpx.ErrInvalidImageFormat
	}
	hash := sha256.Sum256(img)
	key = base64.RawURLEncoding.EncodeToString(hash[:])
	return key, format, nil
}
