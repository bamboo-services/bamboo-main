package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
	"github.com/redis/go-redis/v9"
)

type LinkColorCache struct {
	RDB *redis.Client
	TTL time.Duration
}

func (c *LinkColorCache) Get(ctx context.Context, id int64) (*entity.LinkColor, error) {
	if c == nil || c.RDB == nil {
		return nil, nil
	}

	value, err := c.RDB.Get(ctx, constants.RedisLinkColor.Get(id).String()).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var color entity.LinkColor
	if err = json.Unmarshal([]byte(value), &color); err != nil {
		return nil, err
	}

	return &color, nil
}

func (c *LinkColorCache) Set(ctx context.Context, color *entity.LinkColor) error {
	if c == nil || c.RDB == nil || color == nil {
		return nil
	}

	payload, err := json.Marshal(color)
	if err != nil {
		return err
	}

	return c.RDB.Set(ctx, constants.RedisLinkColor.Get(color.ID).String(), payload, c.TTL).Err()
}

func (c *LinkColorCache) Delete(ctx context.Context, id int64) error {
	if c == nil || c.RDB == nil {
		return nil
	}

	return c.RDB.Del(ctx, constants.RedisLinkColor.Get(id).String()).Err()
}
