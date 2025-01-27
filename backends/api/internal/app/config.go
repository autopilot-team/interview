package app

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Application holds core application settings
type Application struct {
	BaseURL     string `env:"APP_BASE_URL" envDefault:"http://localhost:3001"`
	CompanyName string `env:"APP_COMPANY_NAME" envDefault:"Autopilot Technologies Pte. Ltd."`
	Domain      string `env:"APP_DOMAIN" envDefault:"localhost:3001"`
	Environment string `env:"APP_ENV" envDefault:"development"`
	Name        string `env:"APP_NAME" envDefault:"Autopilot"`
	Service     string `env:"APP_SERVICE" envDefault:"api"`
	Support     struct {
		Email string `env:"APP_SUPPORT_EMAIL" envDefault:"support@autopilot.com"`
		Name  string `env:"APP_SUPPORT_NAME" envDefault:"Autopilot Support"`
		URL   string `env:"APP_SUPPORT_URL" envDefault:"mailto:support@autopilot.com"`
	}
	Version string `env:"APP_VERSION" envDefault:"development"`
}

// CORS holds CORS configuration
type CORS struct {
	// AllowedOrigins is a list of origins a cross-domain request can be executed from.
	// In production, this should be your specific domain(s).
	// Example: https://app.autopilot.domain
	AllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS" envDefault:"http://localhost:2999,http://localhost:3000"`
}

// Cloudflare holds Cloudflare configuration
type Cloudflare struct {
	TurnstileSecretKey string `env:"CF_TURNSTILE_SECRET_KEY" envDefault:"1x0000000000000000000000000000000AA"`
}

// Database holds database connection settings
type Database struct {
	PrimaryWriter  string   `env:"PRIMARY_WRITER_DB_URL" envDefault:"postgres://postgres:postgres@localhost:5432/api?sslmode=disable&search_path=public&pool_max_conns=25&pool_min_conns=2&pool_max_conn_lifetime=1h&pool_max_conn_idle_time=30m&pool_health_check_period=1m"`
	PrimaryReaders []string `env:"PRIMARY_READER_DB_URLS" envDefault:""`
	Worker         string   `env:"WORKER_DB_URL" envDefault:"postgres://postgres:postgres@localhost:5432/api?sslmode=disable&search_path=public&pool_max_conns=25&pool_min_conns=2&pool_max_conn_lifetime=1h&pool_max_conn_idle_time=30m&pool_health_check_period=1m"`
}

// Mailer holds mailer configuration
type Mailer struct {
	SmtpUrl string `env:"SMTP_URL" envDefault:"smtp://localhost:1025"`
}

// Observability holds monitoring and logging configuration
type Observability struct {
	AxiomApiToken string `env:"AXIOM_API_TOKEN" envDefault:"xaat-6adbb824-ccc6-4267-8462-4fab18335352"`
}

// Server holds HTTP server configuration
type Server struct {
	Host string `env:"HOST" envDefault:"0.0.0.0"`
	Port string `env:"PORT" envDefault:"3001"`
	CORS CORS
}

// Services holds all services configuration
type Services struct {
	PaymentAddr string `env:"PAYMENT_SERVICE_ADDR" envDefault:"localhost:3002"`
}

// Config represents the complete application configuration
type Config struct {
	App           Application
	Cloudflare    Cloudflare
	Database      Database
	Mailer        Mailer
	Observability Observability
	Server        Server
	Services      Services
}

// NewConfig creates a new Config instance with values from environment variables
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}
