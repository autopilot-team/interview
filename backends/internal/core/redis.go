package core

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// RedisOptions contains configuration options for creating a new Redis client.
type RedisOptions struct {
	// URL is the Redis connection URL in the format:
	// redis[s]://[[user][:password]@][host][:port][/db-number]
	URL string
}

// NewRedis creates and returns a new Redis client using the provided options.
// It validates the connection by performing a PING command.
//
// The context passed will be used for the initial connection test.
// If the connection fails or the PING command fails, an error is returned.
func NewRedis(ctx context.Context, opts RedisOptions) (*redis.Client, error) {
	opt, err := redis.ParseURL(opts.URL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
