package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
	"wallet/internal/errs"
)

type Client struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) *Client {
	return &Client{
		rdb: rdb,
	}
}

func (c *Client) Get(ctx context.Context, key string, dest any) error {
	if dest == nil {
		return nil
	}

	value, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return errs.ErrNotFound
		}
		return err
	}

	if err = json.Unmarshal([]byte(value), dest); err != nil {
		return err
	}

	return nil
}

func (c *Client) Set(ctx context.Context, key string, data any, expiresAfter time.Duration) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if _, err = c.rdb.Set(ctx, key, raw, expiresAfter).Result(); err != nil {
		return err
	}

	return nil
}

func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	_, err := c.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return true, nil
}
