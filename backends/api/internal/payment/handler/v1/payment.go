package v1

import (
	"autopilot/backends/api/internal/payment/model"
	"autopilot/backends/api/pkg/httpx"
	"context"
	"time"
)

type Payment struct {
	ID          string                  `json:"id" db:"id"`
	MerchantID  string                  `json:"merchant_id" db:"merchant_id"`
	Amount      int64                   `json:"amount" db:"amount"`     // Amount in cents
	Currency    string                  `json:"currency" db:"currency"` // ISO 4217
	Status      model.PaymentStatus     `json:"status" db:"status"`
	Provider    string                  `json:"provider" db:"provider"`
	Method      model.PaymentMethodType `json:"method" db:"method"`
	Description string                  `json:"description" db:"description"`
	Metadata    map[string]any          `json:"metadata" db:"metadata"`
	CreatedAt   time.Time               `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at" db:"updated_at"`
	CompletedAt time.Time               `json:"completed_at,omitempty" db:"completed_at"`
}

// CreatePaymentRequest
type CreatePaymentRequest struct {
	Body struct {
		MerchantID  string
		Amount      httpx.Money
		Currency    httpx.Currency
		Provider    string
		Method      string
		Description string
		Metadata    map[string]any
	}
}

// CreatePaymentResponse is the response body for the update user endpoint.
type CreatePaymentResponse struct {
	Body Payment
}

// UpdateUser is the handler for the update user endpoint.
func (v *V1) CreatePayment(ctx context.Context, input *CreatePaymentRequest) (*CreatePaymentResponse, error) {
	return nil, nil
}
