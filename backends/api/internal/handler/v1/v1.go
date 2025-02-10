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
	service *service.Manager
}

const (
	TagIdentity = "Identity"
)

func BasePath(path string) string {
	return fmt.Sprintf("/v1%s", path)
}

// AddRoutes adds the v1 API docs/routes to the http server
func AddRoutes(container *app.Container, api huma.API, serviceManager *service.Manager) error {
	v1 := &V1{
		Container: container,
		service:   serviceManager,
	}

	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		OperationID: "get-session",
		Path:        BasePath("/identity/me"),
		Summary:     "Get current user session",
		Tags:        []string{TagIdentity},
	}, v1.Me)

	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		OperationID: "sign-in",
		Path:        BasePath("/identity/sign-in"),
		Summary:     "Authenticate and create a new session",
		Tags:        []string{TagIdentity},
	}, v1.SignIn)

	huma.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		OperationID: "sign-out",
		Path:        BasePath("/identity/sign-out"),
		Summary:     "Terminate current session",
		Tags:        []string{TagIdentity},
	}, v1.SignOut)

	return nil
}
