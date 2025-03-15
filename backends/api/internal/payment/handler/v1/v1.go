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
	}

	// Payments Endpoints
	httpx.Register(api, huma.Operation{
		Method:      http.MethodPost,
		OperationID: "create-payment",
		Path:        BasePath("/payments"),
		Summary:     "Create payment",
		Tags:        []string{TagPayment.Name},
	}, v1.CreatePayment, api.WithUnauthenticated())

	return nil
}
