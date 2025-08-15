package app

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config represents the complete application configuration
type Config struct {
	// Application holds core application settings
	App struct {
		// AssetsURL holds the assets URL
		AssetsURL string `env:"APP_ASSETS_URL" envDefault:"http://localhost:2998"`

		// BaseURL holds the base URL
		BaseURL string `env:"APP_BASE_URL" envDefault:"http://localhost:3001"`

		// Cache holds cache configuration
		Cache struct {
			ValkeyURLs string `env:"API_CACHE_VALKEY_URLS" envDefault:"redis://localhost:6379/0,redis://localhost:6380/0"`
		}

		// Cloudflare holds Cloudflare configuration
		Cloudflare struct {
			TurnstileSecretKey string `env:"CF_TURNSTILE_SECRET_KEY" envDefault:"1x0000000000000000000000000000000AA"`
		}

		// CompanyName holds the company name
		CompanyName string `env:"APP_COMPANY_NAME" envDefault:"Autopilot Technologies Pte. Ltd."`

		// CORS holds CORS configuration
		CORS struct {
			// AllowedOrigins is a list of origins a cross-domain request can be executed from.
			// In production, this should be your specific domain(s).
			// Example: https://app.autopilot.domain
			AllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS" envDefault:"http://localhost:2999,http://localhost:3000"`
		}

		// DashboardURL holds the dashboard URL
		DashboardURL string `env:"APP_DASHBOARD_URL" envDefault:"http://localhost:3000"`

		// Database holds database configuration
		Database struct {
			// Worker is the worker database URL
			Worker string `env:"WORKER_DB_URL" envDefault:"postgres://postgres:postgres@localhost:5432/worker?sslmode=disable&search_path=public&pool_max_conns=25&pool_min_conns=2&pool_max_conn_lifetime=1h&pool_max_conn_idle_time=30m&pool_health_check_period=1m"`
		}

		// Domain holds the application domain
		Domain string `env:"APP_DOMAIN" envDefault:"localhost:3001"`

		// Environment holds the application environment
		Environment string `env:"APP_ENV" envDefault:"development"`

		// Mailer holds mailer configuration
		Mailer struct {
			SMTPURL string `env:"SMTP_URL" envDefault:"smtp://localhost:1025"`
		}

		// Name holds the application name
		Name string `env:"APP_NAME" envDefault:"Autopilot"`

		// Observability holds monitoring and logging configuration
		Observability struct {
			APIEndpoint string `env:"OTEL_API_ENDPOINT" envDefault:"localhost:4317"`
		}

		// RateLimiter holds rate limiter configuration
		RateLimiter struct {
			ValkeyURLs string `env:"API_RATE_LIMITER_VALKEY_URLS" envDefault:"redis://localhost:6379/0,redis://localhost:6380/0"`
		}

		// Server holds HTTP server configuration
		Server struct {
			Host string `env:"HOST" envDefault:"0.0.0.0"`
			Port string `env:"PORT" envDefault:"3001"`
		}

		// Service holds the service name
		Service string `env:"APP_SERVICE" envDefault:"api"`

		// Support holds support configuration
		Support struct {
			Email string `env:"APP_SUPPORT_EMAIL" envDefault:"support@autopilot.is"`
			Name  string `env:"APP_SUPPORT_NAME" envDefault:"Autopilot Support"`
			URL   string `env:"APP_SUPPORT_URL" envDefault:"mailto:support@autopilot.is"`
		}

		// Version holds the application version
		Version string `env:"APP_VERSION" envDefault:"development"`
	}

	// Identity holds identity module configuration
	Identity struct {
		// Cache holds cache configuration
		Cache struct {
			ValkeyURLs string `env:"IDENTITY_CACHE_VALKEY_URLS" envDefault:"redis://localhost:6379/1,redis://localhost:6380/1"`
		}

		// Database holds database configuration
		Database struct {
			PrimaryWriter  string   `env:"IDENTITY_PRIMARY_WRITER_DB_URL" envDefault:"postgres://postgres:postgres@localhost:5432/identity?sslmode=disable&search_path=public&pool_max_conns=25&pool_min_conns=2&pool_max_conn_lifetime=1h&pool_max_conn_idle_time=30m&pool_health_check_period=1m"`
			PrimaryReaders []string `env:"IDENTITY_PRIMARY_READER_DB_URLS" envDefault:""`
		}

		// Storage holds S3 storage configuration
		Storage struct {
			Endpoint        string `env:"AWS_ENDPOINT" envDefault:"http://localhost:9000"`
			Region          string `env:"AWS_REGION" envDefault:"us-east-1"`
			AccessKeyID     string `env:"AWS_ACCESS_KEY_ID" envDefault:"minioadmin"`
			SecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY" envDefault:"minioadmin"`
			Bucket          string `env:"IDENTITY_S3_BUCKET" envDefault:"autopilot-development-identity"`
			UsePathStyle    bool   `env:"S3_USE_PATH_STYLE" envDefault:"true"`
		}
	}

	// Payment holds payment module configuration
	Payment struct {
		// Cache holds cache configuration
		Cache struct {
			ValkeyURLs string `env:"PAYMENT_CACHE_VALKEY_URLS" envDefault:"redis://localhost:6379/2,redis://localhost:6380/2"`
		}

		// Database holds database configuration
		Database struct {
			LivePrimaryWriter  string   `env:"PAYMENT_LIVE_PRIMARY_WRITER_DB_URL" envDefault:"postgres://postgres:postgres@localhost:5432/payment_live?sslmode=disable&search_path=public&pool_max_conns=25&pool_min_conns=2&pool_max_conn_lifetime=1h&pool_max_conn_idle_time=30m&pool_health_check_period=1m"`
			LivePrimaryReaders []string `env:"PAYMENT_LIVE_PRIMARY_READER_DB_URLS" envDefault:""`
			TestPrimaryWriter  string   `env:"PAYMENT_TEST_PRIMARY_WRITER_DB_URL" envDefault:"postgres://postgres:postgres@localhost:5432/payment_test?sslmode=disable&search_path=public&pool_max_conns=25&pool_min_conns=2&pool_max_conn_lifetime=1h&pool_max_conn_idle_time=30m&pool_health_check_period=1m"`
			TestPrimaryReaders []string `env:"PAYMENT_TEST_PRIMARY_READER_DB_URLS" envDefault:""`
		}

		// Storage holds S3 storage configuration
		Storage struct {
			Endpoint        string `env:"AWS_ENDPOINT" envDefault:"http://localhost:9000"`
			Region          string `env:"AWS_REGION" envDefault:"us-east-1"`
			AccessKeyID     string `env:"AWS_ACCESS_KEY_ID" envDefault:"minioadmin"`
			SecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY" envDefault:"minioadmin"`
			Bucket          string `env:"PAYMENT_S3_BUCKET" envDefault:"autopilot-development-payment"`
			UsePathStyle    bool   `env:"S3_USE_PATH_STYLE" envDefault:"true"`
		}
	}
}

// NewConfig creates a new Config instance with values from environment variables
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}
