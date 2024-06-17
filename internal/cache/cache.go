package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Cache struct {
	r   *redis.Client
	ttl time.Duration
}

func New(r *redis.Client, ttl time.Duration) *Cache {
	return &Cache{
		r:   r,
		ttl: ttl,
	}
}

func (c *Cache) PutLink(ctx context.Context, path string, link string) error {
	return c.r.Set(ctx, path, link, c.ttl).Err()
}

func (c *Cache) GetLink(ctx context.Context, path string) (string, error) {
	return c.r.Get(ctx, path).Result()
}
