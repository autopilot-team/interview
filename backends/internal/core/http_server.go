package core

import (
	"autopilot/backends/internal/types"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
	chimdw "github.com/go-chi/chi/v5/middleware"
)

// HttpHeaders wraps HTTP headers to control their log representation
type HttpHeaders map[string]string

var sensitiveHeaders = map[string]bool{
	"authorization":   true,
	"cookie":          true,
	"x-csrf-token":    true,
	"x-forwarded-for": true,
}

// LogValue implements slog.LogValuer interface
func (h HttpHeaders) LogValue() slog.Value {
	if h == nil {
		return slog.GroupValue()
	}

	attrs := make([]slog.Attr, 0, len(h))
	for k, v := range h {
		if sensitiveHeaders[strings.ToLower(k)] || strings.Contains(strings.ToLower(k), "token") || strings.Contains(strings.ToLower(k), "key") {
			attrs = append(attrs, slog.String(k, "[REDACTED]"))
		} else {
			attrs = append(attrs, slog.String(k, v))
		}
	}

	return slog.GroupValue(attrs...)
}

// HttpServerOptions contains configuration options for the HTTP server
type HttpServerOptions struct {
	// ApiDocs is a list of API documentation to be served
	ApiDocs []huma.API

	// Mode specifies the application mode (debug/release)
	Mode types.Mode

	// I18nBundle is used for internationalization
	I18nBundle *I18nBundle

	// Logger is used for server-related logging
	Logger *slog.Logger

	// Host specifies the server binding address
	Host string

	// Port specifies the server listening port
	Port string

	// SpaFS is the embedded filesystem for SPA files
	SpaFS FS

	// SpaDir is the path to the SPA files in the embedded filesystem
	SpaDir string

	// Mailer is the optional mailer instance for email previews
	Mailer Mailer

	// Middlewares is a list of middlewares to be applied to the server
	Middlewares []func(http.Handler) http.Handler
}

// HttpServer represents a web server instance with static file serving capabilities.
// It embeds chi.Mux for routing and http.Server for the underlying server.
type HttpServer struct {
	*chi.Mux
	*http.Server
	APIDocs    []huma.API
	i18nBundle *I18nBundle
	logger     *slog.Logger
	spaFS      FS
	mailer     Mailer
}

// NewHttpServer creates and initializes a new HttpServer instance.
func NewHttpServer(opts HttpServerOptions) (*HttpServer, error) {
	if opts.Logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	router := chi.NewRouter()
	router.Use(chimdw.Recoverer)
	router.Use(opts.Middlewares...)

	server := &HttpServer{
		router,
		&http.Server{
			Addr:    opts.Host + ":" + opts.Port,
			Handler: router,
		},
		opts.ApiDocs,
		opts.I18nBundle,
		opts.Logger,
		opts.SpaFS,
		opts.Mailer,
	}

	if (opts.SpaFS != nil || opts.SpaFS != ZeroFS) && opts.SpaDir != "" {
		server.ServeStaticFiles(opts.SpaDir)
	}

	if opts.Mailer != nil {
		opts.Mailer.SetupPreviewRoutes(router)
	}

	return server, nil
}

// ServeStaticFiles configures the server to serve static files from the embedded filesystem.
// It serves files from the "dist" subdirectory and handles special cases for index.html and 404.html.
func (s *HttpServer) ServeStaticFiles(path string) {
	subFs, err := fs.Sub(s.spaFS, path)
	if err != nil {
		panic(fmt.Errorf("Failed to get the sub tree for the static files: %w", err))
	}

	s.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filePath := strings.Trim(r.URL.Path, "/")
		if filePath == "/" || filePath == "" {
			filePath = "index"
		}

		filePath += ".html"
		content, err := fs.ReadFile(subFs, filePath)
		if err == nil {
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write(content)
			return
		}

		if os.IsNotExist(err) && !strings.Contains(r.URL.Path, ".") {
			content, _ := fs.ReadFile(subFs, "404.html")
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write(content)
			return
		}

		http.FileServer(http.FS(subFs)).ServeHTTP(w, r)
	}))
}
