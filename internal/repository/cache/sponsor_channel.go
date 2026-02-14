package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
	"github.com/redis/go-redis/v9"
)

type SponsorChannelCache struct {
	RDB *redis.Client
	TTL time.Duration
}

func (c *SponsorChannelCache) Get(ctx context.Context, id int64) (*entity.SponsorChannel, error) {
	if c == nil || c.RDB == nil {
		return nil, nil
	}

	value, err := c.RDB.Get(ctx, constants.RedisSponsorChan.Get(id).String()).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var channel entity.SponsorChannel
	if err = json.Unmarshal([]byte(value), &channel); err != nil {
		return nil, err
	}

	return &channel, nil
}

func (c *SponsorChannelCache) Set(ctx context.Context, channel *entity.SponsorChannel) error {
	if c == nil || c.RDB == nil || channel == nil {
		return nil
	}

	payload, err := json.Marshal(channel)
	if err != nil {
		return err
	}

	return c.RDB.Set(ctx, constants.RedisSponsorChan.Get(channel.ID).String(), payload, c.TTL).Err()
}

func (c *SponsorChannelCache) Delete(ctx context.Context, id int64) error {
	if c == nil || c.RDB == nil {
		return nil
	}

	return c.RDB.Del(ctx, constants.RedisSponsorChan.Get(id).String()).Err()
}
