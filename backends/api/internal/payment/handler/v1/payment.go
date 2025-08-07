package v1

import (
	"autopilot/backends/api/internal/payment/model"
	"autopilot/backends/api/pkg/httpx"
	"context"
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID          string                  `json:"id" db:"id"`
	MerchantID  string                  `json:"merchantId" db:"merchant_id"`
	Amount      int64                   `json:"amount" db:"amount"`     // Amount in cents
	Currency    string                  `json:"currency" db:"currency"` // ISO 4217
	Status      model.PaymentStatus     `json:"status" db:"status"`
	Provider    string                  `json:"provider" db:"provider"`
	Method      model.PaymentMethodType `json:"method" db:"method"`
	Description string                  `json:"description" db:"description"`
	Metadata    map[string]any          `json:"metadata" db:"metadata"`
	CreatedAt   time.Time               `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time               `json:"updatedAt" db:"updated_at"`
	CompletedAt time.Time               `json:"completedAt,omitempty" db:"completed_at"`
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
	inputPayment := &model.Payment{
		MerchantID:   uuid.MustParse(input.Body.MerchantID),
		Amount:       int64(input.Body.Amount),
		Currency:     input.Body.Currency.Code,
		Provider:     input.Body.Provider,
		Method:       model.PaymentMethodType(input.Body.Method),
		Description:  input.Body.Description,
		Metadata:     input.Body.Metadata,
		Status:       model.PaymentStatusPending,
		ErrorMessage: nil,
		CompletedAt:  nil,
	}
	payment, err := v.paymentService.Payment.Create(ctx, inputPayment)
	if err != nil {
		return nil, err
	}

	response := &CreatePaymentResponse{
		Body: Payment{
			ID:          payment.ID.String(),
			MerchantID:  payment.MerchantID.String(),
			Amount:      payment.Amount,
			Currency:    payment.Currency,
			Status:      payment.Status,
			Provider:    payment.Provider,
			Method:      payment.Method,
			Description: payment.Description,
			Metadata:    payment.Metadata,
			CreatedAt:   payment.CreatedAt,
			UpdatedAt:   payment.UpdatedAt,
		},
	}
	return response, nil
}

	return nil, nil
}
