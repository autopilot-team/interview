package v1

import (
	"autopilot/backends/api/internal/payment/model"
	"autopilot/backends/api/internal/payment/service/mocks"
	"autopilot/backends/api/pkg/app"
	"autopilot/backends/api/pkg/httpx"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInitiateRefund(t *testing.T) {
	// Setup
	_, api := humatest.New(t)
	mockPaymenter := mocks.NewPaymenter(t)
	mockContainer := &app.Container{}
	
	v1Handler := &V1{
		Container: mockContainer,
		payment: &service.Manager{
			Payment: mockPaymenter,
		},
	}

	// Register the route
	httpx.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/v1/refunds",
		OperationID: "initiate-refund",
	}, v1Handler.InitiateRefund)

	// Test data
	paymentID := uuid.New()
	refundID := uuid.New()
	merchantID := uuid.New()
	idempotencyKey := "test-idempotency-key"
	refundAmount := int64(5000) // $50.00

	payment := &model.Payment{
		ID:               paymentID,
		MerchantID:       merchantID,
		Amount:           10000, // $100.00
		Currency:         "USD",
		Status:           model.PaymentStatusSucceeded,
		RefundableAmount: 10000,
		RefundedAmount:   0,
	}

	refund := &model.Refund{
		ID:             refundID,
		PaymentID:      paymentID,
		IdempotencyKey: idempotencyKey,
		Amount:         refundAmount,
		Currency:       "USD",
		Status:         model.RefundStatusPending,
		Reason:         model.RefundReasonCustomerRequest,
		InitiatedAt:    time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	tests := []struct {
		name               string
		idempotencyKey     string
		requestBody        map[string]interface{}
		setupMocks         func()
		expectedStatus     int
		expectedBodyCheck  func(t *testing.T, body []byte)
	}{
		{
			name:           "Successful refund initiation",
			idempotencyKey: idempotencyKey,
			requestBody: map[string]interface{}{
				"payment_id": paymentID.String(),
				"amount":     refundAmount,
				"reason":     "customer_request",
				"metadata": map[string]interface{}{
					"order_id": "12345",
				},
			},
			setupMocks: func() {
				// Mock getting payment
				mockPaymenter.On("GetPaymentByID", mock.Anything, paymentID).
					Return(payment, nil).Once()

				// Mock checking for existing refund with idempotency key
				mockPaymenter.On("GetRefundByIdempotencyKey", mock.Anything, paymentID, idempotencyKey).
					Return(nil, nil).Once()

				// Mock listing existing refunds
				mockPaymenter.On("ListRefunds", mock.Anything, mock.MatchedBy(func(filter *model.RefundFilter) bool {
					return filter.PaymentID != nil && *filter.PaymentID == paymentID
				})).Return([]model.Refund{}, int64(0), nil).Once()

				// Mock creating refund
				mockPaymenter.On("CreateRefund", mock.Anything, mock.MatchedBy(func(r *model.Refund) bool {
					return r.PaymentID == paymentID &&
						r.Amount == refundAmount &&
						r.IdempotencyKey == idempotencyKey
				})).Return(refund, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBodyCheck: func(t *testing.T, body []byte) {
				var response model.Refund
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, refundID, response.ID)
				assert.Equal(t, paymentID, response.PaymentID)
				assert.Equal(t, refundAmount, response.Amount)
				assert.Equal(t, model.RefundStatusPending, response.Status)
			},
		},
		{
			name:           "Idempotent refund request",
			idempotencyKey: idempotencyKey,
			requestBody: map[string]interface{}{
				"payment_id": paymentID.String(),
				"amount":     refundAmount,
				"reason":     "customer_request",
			},
			setupMocks: func() {
				// Mock getting payment
				mockPaymenter.On("GetPaymentByID", mock.Anything, paymentID).
					Return(payment, nil).Once()

				// Mock finding existing refund with same idempotency key
				mockPaymenter.On("GetRefundByIdempotencyKey", mock.Anything, paymentID, idempotencyKey).
					Return(refund, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBodyCheck: func(t *testing.T, body []byte) {
				var response model.Refund
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, refundID, response.ID)
			},
		},
		{
			name:           "Payment not found",
			idempotencyKey: idempotencyKey,
			requestBody: map[string]interface{}{
				"payment_id": uuid.New().String(),
				"amount":     refundAmount,
				"reason":     "customer_request",
			},
			setupMocks: func() {
				// Mock payment not found
				mockPaymenter.On("GetPaymentByID", mock.Anything, mock.Anything).
					Return(nil, httpx.ErrPaymentNotFound).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedBodyCheck: func(t *testing.T, body []byte) {
				assert.Contains(t, string(body), "Payment not found")
			},
		},
		{
			name:           "Invalid payment status",
			idempotencyKey: idempotencyKey,
			requestBody: map[string]interface{}{
				"payment_id": paymentID.String(),
				"amount":     refundAmount,
				"reason":     "customer_request",
			},
			setupMocks: func() {
				// Mock getting payment with invalid status
				invalidPayment := *payment
				invalidPayment.Status = model.PaymentStatusPending
				mockPaymenter.On("GetPaymentByID", mock.Anything, paymentID).
					Return(&invalidPayment, nil).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedBodyCheck: func(t *testing.T, body []byte) {
				assert.Contains(t, string(body), "Payment must be succeeded or partially refunded")
			},
		},
		{
			name:           "Refund amount exceeds refundable amount",
			idempotencyKey: idempotencyKey,
			requestBody: map[string]interface{}{
				"payment_id": paymentID.String(),
				"amount":     15000, // $150.00 - exceeds payment amount
				"reason":     "customer_request",
			},
			setupMocks: func() {
				// Mock getting payment
				mockPaymenter.On("GetPaymentByID", mock.Anything, paymentID).
					Return(payment, nil).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedBodyCheck: func(t *testing.T, body []byte) {
				assert.Contains(t, string(body), "Refund amount exceeds refundable amount")
			},
		},
		{
			name:           "Duplicate refund attempt",
			idempotencyKey: "different-idempotency-key",
			requestBody: map[string]interface{}{
				"payment_id": paymentID.String(),
				"amount":     refundAmount,
				"reason":     "customer_request",
			},
			setupMocks: func() {
				// Mock getting payment
				mockPaymenter.On("GetPaymentByID", mock.Anything, paymentID).
					Return(payment, nil).Once()

				// Mock no existing refund with this idempotency key
				mockPaymenter.On("GetRefundByIdempotencyKey", mock.Anything, paymentID, "different-idempotency-key").
					Return(nil, nil).Once()

				// Mock existing pending refund with same amount
				existingRefund := model.Refund{
					ID:        uuid.New(),
					PaymentID: paymentID,
					Amount:    refundAmount,
					Status:    model.RefundStatusPending,
				}
				mockPaymenter.On("ListRefunds", mock.Anything, mock.MatchedBy(func(filter *model.RefundFilter) bool {
					return filter.PaymentID != nil && *filter.PaymentID == paymentID
				})).Return([]model.Refund{existingRefund}, int64(1), nil).Once()
			},
			expectedStatus: http.StatusConflict,
			expectedBodyCheck: func(t *testing.T, body []byte) {
				assert.Contains(t, string(body), "A refund for this amount is already being processed")
			},
		},
		{
			name:           "Partial refund",
			idempotencyKey: idempotencyKey,
			requestBody: map[string]interface{}{
				"payment_id": paymentID.String(),
				"amount":     3000, // $30.00 - partial refund
				"reason":     "product_issue",
				"reason_description": "Item was damaged",
			},
			setupMocks: func() {
				// Mock getting payment that's already partially refunded
				partiallyRefundedPayment := *payment
				partiallyRefundedPayment.Status = model.PaymentStatusPartiallyRefunded
				partiallyRefundedPayment.RefundedAmount = 2000 // $20.00 already refunded
				partiallyRefundedPayment.RefundableAmount = 8000 // $80.00 remaining
				
				mockPaymenter.On("GetPaymentByID", mock.Anything, paymentID).
					Return(&partiallyRefundedPayment, nil).Once()

				// Mock checking for existing refund
				mockPaymenter.On("GetRefundByIdempotencyKey", mock.Anything, paymentID, idempotencyKey).
					Return(nil, nil).Once()

				// Mock listing existing refunds (one completed refund)
				completedRefund := model.Refund{
					ID:        uuid.New(),
					PaymentID: paymentID,
					Amount:    2000,
					Status:    model.RefundStatusSucceeded,
				}
				mockPaymenter.On("ListRefunds", mock.Anything, mock.MatchedBy(func(filter *model.RefundFilter) bool {
					return filter.PaymentID != nil && *filter.PaymentID == paymentID
				})).Return([]model.Refund{completedRefund}, int64(1), nil).Once()

				// Mock creating partial refund
				partialRefund := *refund
				partialRefund.Amount = 3000
				partialRefund.Reason = model.RefundReasonProductIssue
				mockPaymenter.On("CreateRefund", mock.Anything, mock.MatchedBy(func(r *model.Refund) bool {
					return r.Amount == 3000
				})).Return(&partialRefund, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBodyCheck: func(t *testing.T, body []byte) {
				var response model.Refund
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, int64(3000), response.Amount)
				assert.Equal(t, model.RefundReasonProductIssue, response.Reason)
			},
		},
		{
			name:           "Missing idempotency key",
			idempotencyKey: "",
			requestBody: map[string]interface{}{
				"payment_id": paymentID.String(),
				"amount":     refundAmount,
				"reason":     "customer_request",
			},
			setupMocks:     func() {},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBodyCheck: func(t *testing.T, body []byte) {
				assert.Contains(t, string(body), "validation failed")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockPaymenter.ExpectedCalls = nil
			mockPaymenter.Calls = nil

			// Setup mocks for this test
			tt.setupMocks()

			// Create request
			bodyBytes, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/v1/refunds", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			if tt.idempotencyKey != "" {
				req.Header.Set("Idempotency-Key", tt.idempotencyKey)
			}

			// Execute request
			resp := httptest.NewRecorder()
			api.Adapter().ServeHTTP(resp, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, resp.Code, "Response body: %s", resp.Body.String())
			if tt.expectedBodyCheck != nil {
				tt.expectedBodyCheck(t, resp.Body.Bytes())
			}

			// Verify all mocks were called
			mockPaymenter.AssertExpectations(t)
		})
	}
}

func TestListMerchantRefunds(t *testing.T) {
	// Setup
	_, api := humatest.New(t)
	mockPaymenter := mocks.NewPaymenter(t)
	mockContainer := &app.Container{}
	
	v1Handler := &V1{
		Container: mockContainer,
		payment: &service.Manager{
			Payment: mockPaymenter,
		},
	}

	// Register the route
	httpx.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/v1/merchants/{merchant_id}/refunds",
		OperationID: "list-merchant-refunds",
	}, v1Handler.ListMerchantRefunds)

	// Test data
	merchantID := uuid.New()
	refunds := []model.Refund{
		{
			ID:        uuid.New(),
			PaymentID: uuid.New(),
			Amount:    5000,
			Currency:  "USD",
			Status:    model.RefundStatusSucceeded,
			Reason:    model.RefundReasonCustomerRequest,
			CreatedAt: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:        uuid.New(),
			PaymentID: uuid.New(),
			Amount:    3000,
			Currency:  "USD",
			Status:    model.RefundStatusPending,
			Reason:    model.RefundReasonProductIssue,
			CreatedAt: time.Now(),
		},
	}

	tests := []struct {
		name              string
		merchantID        string
		queryParams       string
		setupMocks        func()
		expectedStatus    int
		expectedBodyCheck func(t *testing.T, body []byte)
	}{
		{
			name:        "List all merchant refunds",
			merchantID:  merchantID.String(),
			queryParams: "",
			setupMocks: func() {
				mockPaymenter.On("ListRefunds", mock.Anything, mock.MatchedBy(func(filter *model.RefundFilter) bool {
					return filter.MerchantID != nil && 
						*filter.MerchantID == merchantID &&
						filter.Limit == 20 &&
						filter.Offset == 0
				})).Return(refunds, int64(2), nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBodyCheck: func(t *testing.T, body []byte) {
				var response struct {
					Refunds []model.Refund `json:"refunds"`
					Total   int64          `json:"total"`
					Limit   int            `json:"limit"`
					Offset  int            `json:"offset"`
				}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Len(t, response.Refunds, 2)
				assert.Equal(t, int64(2), response.Total)
				assert.Equal(t, 20, response.Limit)
				assert.Equal(t, 0, response.Offset)
			},
		},
		{
			name:        "Filter by status",
			merchantID:  merchantID.String(),
			queryParams: "?status=pending",
			setupMocks: func() {
				pendingStatus := model.RefundStatusPending
				mockPaymenter.On("ListRefunds", mock.Anything, mock.MatchedBy(func(filter *model.RefundFilter) bool {
					return filter.MerchantID != nil && 
						*filter.MerchantID == merchantID &&
						filter.Status != nil &&
						*filter.Status == pendingStatus
				})).Return([]model.Refund{refunds[1]}, int64(1), nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBodyCheck: func(t *testing.T, body []byte) {
				var response struct {
					Refunds []model.Refund `json:"refunds"`
					Total   int64          `json:"total"`
				}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Len(t, response.Refunds, 1)
				assert.Equal(t, model.RefundStatusPending, response.Refunds[0].Status)
			},
		},
		{
			name:        "Pagination",
			merchantID:  merchantID.String(),
			queryParams: "?limit=1&offset=1",
			setupMocks: func() {
				mockPaymenter.On("ListRefunds", mock.Anything, mock.MatchedBy(func(filter *model.RefundFilter) bool {
					return filter.Limit == 1 && filter.Offset == 1
				})).Return([]model.Refund{refunds[1]}, int64(2), nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBodyCheck: func(t *testing.T, body []byte) {
				var response struct {
					Refunds []model.Refund `json:"refunds"`
					Total   int64          `json:"total"`
					Limit   int            `json:"limit"`
					Offset  int            `json:"offset"`
				}
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Len(t, response.Refunds, 1)
				assert.Equal(t, int64(2), response.Total)
				assert.Equal(t, 1, response.Limit)
				assert.Equal(t, 1, response.Offset)
			},
		},
		{
			name:        "Invalid merchant ID",
			merchantID:  "invalid-uuid",
			queryParams: "",
			setupMocks:  func() {},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBodyCheck: func(t *testing.T, body []byte) {
				assert.Contains(t, string(body), "validation failed")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockPaymenter.ExpectedCalls = nil
			mockPaymenter.Calls = nil

			// Setup mocks for this test
			tt.setupMocks()

			// Create request
			url := fmt.Sprintf("/v1/merchants/%s/refunds%s", tt.merchantID, tt.queryParams)
			req := httptest.NewRequest(http.MethodGet, url, nil)

			// Execute request
			resp := httptest.NewRecorder()
			api.Adapter().ServeHTTP(resp, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, resp.Code, "Response body: %s", resp.Body.String())
			if tt.expectedBodyCheck != nil {
				tt.expectedBodyCheck(t, resp.Body.Bytes())
			}

			// Verify all mocks were called
			mockPaymenter.AssertExpectations(t)
		})
	}
}

func TestCancelRefund(t *testing.T) {
	// Setup
	_, api := humatest.New(t)
	mockPaymenter := mocks.NewPaymenter(t)
	mockContainer := &app.Container{}
	
	v1Handler := &V1{
		Container: mockContainer,
		payment: &service.Manager{
			Payment: mockPaymenter,
		},
	}

	// Register the route
	httpx.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/v1/refunds/{refund_id}/cancel",
		OperationID: "cancel-refund",
	}, v1Handler.CancelRefund)

	// Test data
	refundID := uuid.New()
	paymentID := uuid.New()
	
	pendingRefund := &model.Refund{
		ID:        refundID,
		PaymentID: paymentID,
		Amount:    5000,
		Currency:  "USD",
		Status:    model.RefundStatusPending,
		Reason:    model.RefundReasonCustomerRequest,
	}

	cancelledRefund := *pendingRefund
	cancelledRefund.Status = model.RefundStatusCancelled
	cancelledRefund.CancelledAt = &time.Time{}
	*cancelledRefund.CancelledAt = time.Now()

	tests := []struct {
		name              string
		refundID          string
		requestBody       map[string]interface{}
		setupMocks        func()
		expectedStatus    int
		expectedBodyCheck func(t *testing.T, body []byte)
	}{
		{
			name:     "Successfully cancel pending refund",
			refundID: refundID.String(),
			requestBody: map[string]interface{}{
				"reason": "Customer changed their mind",
			},
			setupMocks: func() {
				// Mock getting refund
				mockPaymenter.On("GetRefundByID", mock.Anything, refundID).
					Return(pendingRefund, nil).Once()

				// Mock cancelling refund
				reason := "Customer changed their mind"
				mockPaymenter.On("CancelRefund", mock.Anything, refundID, &reason).
					Return(&cancelledRefund, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBodyCheck: func(t *testing.T, body []byte) {
				var response model.Refund
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, model.RefundStatusCancelled, response.Status)
				assert.NotNil(t, response.CancelledAt)
			},
		},
		{
			name:        "Cancel without reason",
			refundID:    refundID.String(),
			requestBody: map[string]interface{}{},
			setupMocks: func() {
				// Mock getting refund
				mockPaymenter.On("GetRefundByID", mock.Anything, refundID).
					Return(pendingRefund, nil).Once()

				// Mock cancelling refund
				mockPaymenter.On("CancelRefund", mock.Anything, refundID, (*string)(nil)).
					Return(&cancelledRefund, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBodyCheck: func(t *testing.T, body []byte) {
				var response model.Refund
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, model.RefundStatusCancelled, response.Status)
			},
		},
		{
			name:        "Refund not found",
			refundID:    uuid.New().String(),
			requestBody: map[string]interface{}{},
			setupMocks: func() {
				mockPaymenter.On("GetRefundByID", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("refund not found")).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedBodyCheck: func(t *testing.T, body []byte) {
				assert.Contains(t, string(body), "Refund not found")
			},
		},
		{
			name:        "Cannot cancel succeeded refund",
			refundID:    refundID.String(),
			requestBody: map[string]interface{}{},
			setupMocks: func() {
				succeededRefund := *pendingRefund
				succeededRefund.Status = model.RefundStatusSucceeded
				mockPaymenter.On("GetRefundByID", mock.Anything, refundID).
					Return(&succeededRefund, nil).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedBodyCheck: func(t *testing.T, body []byte) {
				assert.Contains(t, string(body), "Only pending or processing refunds can be cancelled")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockPaymenter.ExpectedCalls = nil
			mockPaymenter.Calls = nil

			// Setup mocks for this test
			tt.setupMocks()

			// Create request
			bodyBytes, _ := json.Marshal(tt.requestBody)
			url := fmt.Sprintf("/v1/refunds/%s/cancel", tt.refundID)
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			resp := httptest.NewRecorder()
			api.Adapter().ServeHTTP(resp, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, resp.Code, "Response body: %s", resp.Body.String())
			if tt.expectedBodyCheck != nil {
				tt.expectedBodyCheck(t, resp.Body.Bytes())
			}

			// Verify all mocks were called
			mockPaymenter.AssertExpectations(t)
		})
	}
}