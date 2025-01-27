package v1

import (
	"autopilot/backends/api/internal/app"
	"autopilot/backends/api/internal/service"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// V1 is the v1 API handler
type V1 struct {
	*app.Container
	paymentService *service.Payment
	// sessionService *service.Session
}

const (
	TagSessions = "Sessions"
)

func BasePath(path string) string {
	return fmt.Sprintf("/v1%s", path)
}

// AddRoutes adds the v1 API docs/routes to the http server
func AddRoutes(container *app.Container, api huma.API) error {
	paymentService, err := service.NewPayment(container)
	if err != nil {
		return err
	}

	v1 := &V1{
		Container:      container,
		paymentService: paymentService,
		// sessionService: service.NewSession(container),
	}

	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		OperationID: "sign-in",
		Path:        BasePath("/identity/sign-in"),
		Summary:     "Sign in to your account",
		Tags:        []string{TagSessions},
	}, v1.SignIn)

	return nil
}
