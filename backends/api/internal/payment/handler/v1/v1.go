package v1

import (
	"autopilot/backends/api/internal/payment/service"
	"autopilot/backends/api/pkg/app"
	"autopilot/backends/api/pkg/httpx"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// V1 is the v1 API handler
type V1 struct {
	*app.Container
	payment *service.Manager
}

var TagPayment = huma.Tag{
	Name:        "Payment",
	Description: `Payment Service`,
}

func BasePath(path string) string {
	return fmt.Sprintf("/v1%s", path)
}

// AddRoutes adds the v1 API docs/routes to the http server
func AddRoutes(container *app.Container, humaAPI huma.API, service *service.Manager, auth httpx.Authenticator) error {
	api := httpx.InitHandler(humaAPI, container.Mode, auth)
	api.AddTags(&TagPayment)

	v1 := &V1{
		Container: container,
		payment:   service,
	}

	// Payment Endpoints
	httpx.Register(api, huma.Operation{
		Method:      http.MethodPost,
		OperationID: "create-payment",
		Path:        BasePath("/payments"),
		Summary:     "Create payment",
		Description: "Create a new payment with idempotency support",
		Tags:        []string{TagPayment.Name},
	}, v1.CreatePayment)

	httpx.Register(api, huma.Operation{
		Method:      http.MethodGet,
		OperationID: "get-payment",
		Path:        BasePath("/payments/{payment_id}"),
		Summary:     "Get payment details",
		Description: "Get details of a specific payment",
		Tags:        []string{TagPayment.Name},
	}, v1.GetPayment)

	httpx.Register(api, huma.Operation{
		Method:      http.MethodGet,
		OperationID: "list-payments",
		Path:        BasePath("/payments"),
		Summary:     "List payments",
		Description: "List payments with filtering options",
		Tags:        []string{TagPayment.Name},
	}, v1.ListPayments)

	httpx.Register(api, huma.Operation{
		Method:      http.MethodPatch,
		OperationID: "update-payment-status",
		Path:        BasePath("/payments/{payment_id}/status"),
		Summary:     "Update payment status",
		Description: "Update payment status (webhook endpoint)",
		Tags:        []string{TagPayment.Name},
	}, v1.UpdatePaymentStatus, api.WithUnauthenticated())

	// Refund Endpoints
	httpx.Register(api, huma.Operation{
		Method:      http.MethodPost,
		OperationID: "initiate-refund",
		Path:        BasePath("/refunds"),
		Summary:     "Initiate refund",
		Description: "Create a new refund with validation and idempotency support",
		Tags:        []string{TagPayment.Name},
	}, v1.InitiateRefund)

	httpx.Register(api, huma.Operation{
		Method:      http.MethodGet,
		OperationID: "get-refund",
		Path:        BasePath("/refunds/{refund_id}"),
		Summary:     "Get refund details",
		Description: "Get details of a specific refund",
		Tags:        []string{TagPayment.Name},
	}, v1.GetRefund)

	httpx.Register(api, huma.Operation{
		Method:      http.MethodGet,
		OperationID: "list-refunds",
		Path:        BasePath("/refunds"),
		Summary:     "List refunds",
		Description: "List refunds with filtering options",
		Tags:        []string{TagPayment.Name},
	}, v1.ListRefunds)

	httpx.Register(api, huma.Operation{
		Method:      http.MethodGet,
		OperationID: "list-merchant-refunds",
		Path:        BasePath("/merchants/{merchant_id}/refunds"),
		Summary:     "List merchant refunds",
		Description: "List all refunds associated with a particular merchant",
		Tags:        []string{TagPayment.Name},
	}, v1.ListMerchantRefunds)

	httpx.Register(api, huma.Operation{
		Method:      http.MethodPost,
		OperationID: "cancel-refund",
		Path:        BasePath("/refunds/{refund_id}/cancel"),
		Summary:     "Cancel refund",
		Description: "Cancel a pending or processing refund",
		Tags:        []string{TagPayment.Name},
	}, v1.CancelRefund)

	httpx.Register(api, huma.Operation{
		Method:      http.MethodPatch,
		OperationID: "update-refund-status",
		Path:        BasePath("/refunds/{refund_id}/status"),
		Summary:     "Update refund status",
		Description: "Update refund status (webhook endpoint for async processing)",
		Tags:        []string{TagPayment.Name},
	}, v1.UpdateRefundStatus, api.WithUnauthenticated())

	return nil
}
