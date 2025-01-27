package app

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Application holds core application settings
type Application struct {
	Environment string `env:"APP_ENV" envDefault:"development"`
	Name        string `env:"APP_NAME" envDefault:"Autopilot"`
	Service     string `env:"APP_SERVICE" envDefault:"payment"`
	Version     string `env:"APP_VERSION" envDefault:"development"`
	CompanyName string `env:"COMPANY_NAME" envDefault:"Autopilot Technologies Pte. Ltd."`
}

// Server holds HTTP server configuration
type Server struct {
	Host string `env:"HOST" envDefault:"0.0.0.0"`
	Port string `env:"PORT" envDefault:"3002"`
}

// Database holds database connection settings
type Database struct {
	LivePrimaryWriter  string   `env:"LIVE_PRIMARY_WRITER_DB_URL" envDefault:"postgres://postgres:postgres@localhost:5432/payment_live?sslmode=disable&search_path=public&pool_max_conns=25&pool_min_conns=2&pool_max_conn_lifetime=1h&pool_max_conn_idle_time=30m&pool_health_check_period=1m"`
	LivePrimaryReaders []string `env:"LIVE_PRIMARY_READER_DB_URLS" envDefault:""`
	TestPrimaryWriter  string   `env:"TEST_PRIMARY_WRITER_DB_URL" envDefault:"postgres://postgres:postgres@localhost:5432/payment_test?sslmode=disable&search_path=public&pool_max_conns=25&pool_min_conns=2&pool_max_conn_lifetime=1h&pool_max_conn_idle_time=30m&pool_health_check_period=1m"`
	TestPrimaryReaders []string `env:"TEST_PRIMARY_READER_DB_URLS" envDefault:""`
	LiveWorker         string   `env:"LIVE_WORKER_DB_URL" envDefault:"postgres://postgres:postgres@localhost:5432/payment_live?sslmode=disable&search_path=public&pool_max_conns=25&pool_min_conns=2&pool_max_conn_lifetime=1h&pool_max_conn_idle_time=30m&pool_health_check_period=1m"`
	TestWorker         string   `env:"TEST_WORKER_DB_URL" envDefault:"postgres://postgres:postgres@localhost:5432/payment_test?sslmode=disable&search_path=public&pool_max_conns=25&pool_min_conns=2&pool_max_conn_lifetime=1h&pool_max_conn_idle_time=30m&pool_health_check_period=1m"`
}

// Observability holds monitoring and logging configuration
type Observability struct {
	AxiomApiToken string `env:"AXIOM_API_TOKEN" envDefault:"xaat-6adbb824-ccc6-4267-8462-4fab18335352"`
}

// Config represents the complete application configuration
type Config struct {
	App           Application
	Server        Server
	Database      Database
	Observability Observability
}

// NewConfig creates a new Config instance with values from environment variables
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}
