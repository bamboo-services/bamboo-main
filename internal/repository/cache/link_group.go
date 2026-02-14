package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
	"github.com/redis/go-redis/v9"
)

type LinkGroupCache struct {
	RDB *redis.Client
	TTL time.Duration
}

func (c *LinkGroupCache) Get(ctx context.Context, id int64) (*entity.LinkGroup, error) {
	if c == nil || c.RDB == nil {
		return nil, nil
	}

	value, err := c.RDB.Get(ctx, constants.RedisLinkGroup.Get(id).String()).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var group entity.LinkGroup
	if err = json.Unmarshal([]byte(value), &group); err != nil {
		return nil, err
	}

	return &group, nil
}

func (c *LinkGroupCache) Set(ctx context.Context, group *entity.LinkGroup) error {
	if c == nil || c.RDB == nil || group == nil {
		return nil
	}

	payload, err := json.Marshal(group)
	if err != nil {
		return err
	}

	return c.RDB.Set(ctx, constants.RedisLinkGroup.Get(group.ID).String(), payload, c.TTL).Err()
}

func (c *LinkGroupCache) Delete(ctx context.Context, id int64) error {
	if c == nil || c.RDB == nil {
		return nil
	}

	return c.RDB.Del(ctx, constants.RedisLinkGroup.Get(id).String()).Err()
}
