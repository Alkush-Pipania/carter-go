package redis

import (
	"context"
	"time"
)

func (c *Client) SetSession(ctx context.Context, sessionID, userID string, ttl time.Duration) error {
	return c.rdb.Set(ctx, "session:"+sessionID, userID, ttl).Err()
}

func (c *Client) GetSession(ctx context.Context, sessionID string) (string, error) {
	return c.rdb.Get(ctx, "session:"+sessionID).Result()
}

func (c *Client) DeleteSession(ctx context.Context, sessionID string) error {
	return c.rdb.Del(ctx, "session:"+sessionID).Err()
}
