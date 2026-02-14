package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
	"github.com/redis/go-redis/v9"
)

type LinkFriendCache struct {
	RDB *redis.Client
	TTL time.Duration
}

func (c *LinkFriendCache) Get(ctx context.Context, id int64) (*entity.LinkFriend, error) {
	if c == nil || c.RDB == nil {
		return nil, nil
	}

	value, err := c.RDB.Get(ctx, constants.RedisLinkFriend.Get(id).String()).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var link entity.LinkFriend
	if err = json.Unmarshal([]byte(value), &link); err != nil {
		return nil, err
	}

	return &link, nil
}

func (c *LinkFriendCache) Set(ctx context.Context, link *entity.LinkFriend) error {
	if c == nil || c.RDB == nil || link == nil {
		return nil
	}

	payload, err := json.Marshal(link)
	if err != nil {
		return err
	}

	return c.RDB.Set(ctx, constants.RedisLinkFriend.Get(link.ID).String(), payload, c.TTL).Err()
}

func (c *LinkFriendCache) Delete(ctx context.Context, id int64) error {
	if c == nil || c.RDB == nil {
		return nil
	}

	return c.RDB.Del(ctx, constants.RedisLinkFriend.Get(id).String()).Err()
}
