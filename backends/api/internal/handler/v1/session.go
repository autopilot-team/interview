package v1

import "context"

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignInResponse struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *V1) SignIn(ctx context.Context, input *SignInRequest) (*SignInResponse, error) {
	return nil, nil
}
