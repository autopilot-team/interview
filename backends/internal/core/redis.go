package core

import (
	"context"
	"crypto/tls"
	"strings"

	"github.com/redis/go-redis/v9"
)

// RedisOptions contains configuration options for creating a new Redis client.
type RedisOptions struct {
	// URL is the Redis connection URL in the format:
	// redis[s]://[[user][:password]@][host][:port][/db-number]
	// For cluster mode, provide multiple URLs separated by commas
	URL string

	// IsCluster determines whether to use Redis Cluster client
	IsCluster bool

	// EnableTLS enables TLS for the Redis connection
	EnableTLS bool

	// TLSConfig is the TLS configuration to use for the Redis connection
	// If nil, a default configuration will be used when EnableTLS is true
	TLSConfig *tls.Config

	// SkipTLSVerify skips TLS certificate verification
	// This is useful for AWS ElastiCache which may have certificates that
	// don't fully comply with standards that Go enforces
	// CAUTION: Only use this in trusted environments
	SkipTLSVerify bool
}

// NewRedis creates and returns a new Redis client using the provided options.
// It supports both standalone and cluster modes.
// It validates the connection by performing a PING command.
//
// The context passed will be used for the initial connection test.
// If the connection fails or the PING command fails, an error is returned.
func NewRedis(ctx context.Context, opts RedisOptions) (redis.UniversalClient, error) {
	// Auto-detect AWS ElastiCache endpoints
	if strings.Contains(opts.URL, ".cache.amazonaws.com") {
		opts.EnableTLS = true
		opts.SkipTLSVerify = true
	}

	if opts.IsCluster {
		// For cluster mode, parse multiple URLs
		urls := strings.Split(opts.URL, ",")
		addrs := make([]string, len(urls))

		// Use the first URL to extract username/password
		var username, password string
		if len(urls) > 0 {
			firstOpt, err := redis.ParseURL(urls[0])
			if err != nil {
				return nil, err
			}
			username = firstOpt.Username
			password = firstOpt.Password
		}

		// Extract host:port from each URL
		for i, url := range urls {
			opt, err := redis.ParseURL(url)
			if err != nil {
				return nil, err
			}
			addrs[i] = opt.Addr
		}

		// Create a new cluster client
		clusterOpt := &redis.ClusterOptions{
			Addrs:    addrs,
			Username: username,
			Password: password,
		}

		// Configure TLS if enabled
		if opts.EnableTLS {
			tlsConfig := opts.TLSConfig
			if tlsConfig == nil {
				tlsConfig = &tls.Config{}
			}

			// For AWS ElastiCache, we might need to skip verification
			if opts.SkipTLSVerify {
				tlsConfig = &tls.Config{
					InsecureSkipVerify: true,
				}
			}

			clusterOpt.TLSConfig = tlsConfig
		}

		client := redis.NewClusterClient(clusterOpt)
		if err := client.Ping(ctx).Err(); err != nil {
			return nil, err
		}

		return client, nil
	}

	// For standalone mode (backward compatibility)
	opt, err := redis.ParseURL(opts.URL)
	if err != nil {
		return nil, err
	}

	// Configure TLS if enabled
	if opts.EnableTLS {
		tlsConfig := opts.TLSConfig
		if tlsConfig == nil {
			tlsConfig = &tls.Config{}
		}

		// For AWS ElastiCache, we might need to skip verification
		if opts.SkipTLSVerify {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		opt.TLSConfig = tlsConfig
	}

	client := redis.NewClient(opt)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
