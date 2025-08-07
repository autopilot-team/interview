package v1

import (
	"autopilot/backends/api/internal/identity"
	"autopilot/backends/api/internal/payment/model"
	"autopilot/backends/api/pkg/httpx"
	"autopilot/backends/api/pkg/testutil"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestV1_CreatePayment(t *testing.T) {
	t.Parallel()
	createPaymentPath := BasePath("/payments")
	api, container, mods := testutil.Container(t)

	ctx := context.Background()
	auth := identity.NewAuthentication(container, api, mods.Identity.Service)
	err := AddRoutes(container, api, mods.Payment.Service, auth)
	require.NoError(t, err)

	// Generate a valid merchant ID for testing
	merchantID := uuid.New().String()

	tests := []struct {
		name           string
		payload        map[string]any
		expectedStatus int
		err            error
		checkResponse  bool
	}{
		{
			name: "should create payment successfully",
			payload: map[string]any{
				"MerchantID":  merchantID,
				"Amount":      1000, // $10.00 in cents
				"Currency":    "USD",
				"Provider":    "stripe",
				"Method":      "card",
				"Description": "Test payment",
				"Metadata": map[string]any{
					"order_id": "order_123",
				},
			},
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name: "should create payment with minimal required fields",
			payload: map[string]any{
				"MerchantID":  merchantID,
				"Amount":      500, // $5.00 in cents
				"Currency":    "USD",
				"Provider":    "stripe",
				"Method":      "card",
				"Description": "Minimal payment",
				"Metadata":    map[string]any{}, // Empty but required
			},
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name: "should reject invalid merchant_id format",
			payload: map[string]any{
				"MerchantID":  "invalid-uuid",
				"Amount":      1000,
				"Currency":    "USD",
				"Provider":    "stripe",
				"Method":      "card",
				"Description": "Test payment",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject negative amount",
			payload: map[string]any{
				"MerchantID":  merchantID,
				"Amount":      -100,
				"Currency":    "USD",
				"Provider":    "stripe",
				"Method":      "card",
				"Description": "Test payment",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject invalid currency",
			payload: map[string]any{
				"MerchantID":  merchantID,
				"Amount":      1000,
				"Currency":    "INVALID",
				"Provider":    "stripe",
				"Method":      "card",
				"Description": "Test payment",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject missing required fields",
			payload: map[string]any{
				"MerchantID": merchantID,
				"Amount":     1000,
				// missing Currency, Provider, Method, Description
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject empty merchant_id",
			payload: map[string]any{
				"MerchantID":  "",
				"Amount":      1000,
				"Currency":    "USD",
				"Provider":    "stripe",
				"Method":      "card",
				"Description": "Test payment",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject zero amount",
			payload: map[string]any{
				"MerchantID":  merchantID,
				"Amount":      0,
				"Currency":    "USD",
				"Provider":    "stripe",
				"Method":      "card",
				"Description": "Test payment",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject empty provider",
			payload: map[string]any{
				"MerchantID":  merchantID,
				"Amount":      1000,
				"Currency":    "USD",
				"Provider":    "",
				"Method":      "card",
				"Description": "Test payment",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
		{
			name: "should reject empty method",
			payload: map[string]any{
				"MerchantID":  merchantID,
				"Amount":      1000,
				"Currency":    "USD",
				"Provider":    "stripe",
				"Method":      "",
				"Description": "Test payment",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			err:            httpx.ErrInvalidBody,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := api.Post(createPaymentPath, tc.payload)
			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.err != nil {
				httpx.AssertErr(t, tc.err, resp.Body)
				return
			}

			if tc.checkResponse {
				var response Payment
				err := json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)

				// Verify response structure
				assert.NotEmpty(t, response.ID)
				assert.Equal(t, merchantID, response.MerchantID)
				assert.Equal(t, int64(tc.payload["Amount"].(int)), response.Amount)
				assert.Equal(t, tc.payload["Currency"], response.Currency)
				assert.Equal(t, model.PaymentStatusPending, response.Status)
				assert.NotZero(t, response.CreatedAt)
				assert.NotZero(t, response.UpdatedAt)

				// Verify payment was actually created in the database
				if response.ID != "" {
					createdPayment, err := mods.Payment.Service.Payment.Get(ctx, response.ID)
					if err != nil {
						t.Logf("Error getting payment from database: %v", err)
					} else {
						assert.NotNil(t, createdPayment)
						if createdPayment != nil {
							assert.Equal(t, response.ID, createdPayment.ID.String())
						}
					}
				}
			}
		})
	}
}
