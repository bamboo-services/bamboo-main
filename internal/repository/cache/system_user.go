package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
	"github.com/redis/go-redis/v9"
)

type SystemUserCache struct {
	RDB *redis.Client
	TTL time.Duration
}

func (c *SystemUserCache) Get(ctx context.Context, id int64) (*entity.SystemUser, error) {
	if c == nil || c.RDB == nil {
		return nil, nil
	}

	value, err := c.RDB.Get(ctx, c.buildKey(id)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var user entity.SystemUser
	if err = json.Unmarshal([]byte(value), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *SystemUserCache) Set(ctx context.Context, user *entity.SystemUser) error {
	if c == nil || c.RDB == nil || user == nil {
		return nil
	}

	payload, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return c.RDB.Set(ctx, c.buildKey(user.ID), payload, c.TTL).Err()
}

func (c *SystemUserCache) Delete(ctx context.Context, id int64) error {
	if c == nil || c.RDB == nil {
		return nil
	}
	return c.RDB.Del(ctx, c.buildKey(id)).Err()
}

func (c *SystemUserCache) buildKey(id int64) string {
	return constants.RedisSystemUser.Get(id).String()
}
