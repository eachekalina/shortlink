package cache

import (
	"context"
	"errors"
	"github.com/eachekalina/shortlink/internal/errs"
	"github.com/redis/go-redis/v9"
	"log/slog"
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
	slog.Debug("saving link to cache", "path", path, "link", link)
	err := c.r.Set(ctx, path, link, c.ttl).Err()
	if err != nil {
		slog.Error("failed to save link", "path", path, "link", link, "err", err)
		return err
	}
	return nil
}

func (c *Cache) GetLink(ctx context.Context, path string) (string, error) {
	slog.Debug("trying to get link from cache", "path", path)
	res, err := c.r.Get(ctx, path).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			slog.Debug("link not found in cache", "path", path)
			return "", errs.ErrNotFound
		}
		slog.Error("error while accessing cache", "path", path, "err", err)
		return "", err
	}
	return res, nil
}
