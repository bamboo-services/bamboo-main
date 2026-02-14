package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
	"github.com/redis/go-redis/v9"
)

type SponsorRecordCache struct {
	RDB *redis.Client
	TTL time.Duration
}

func (c *SponsorRecordCache) Get(ctx context.Context, id int64) (*entity.SponsorRecord, error) {
	if c == nil || c.RDB == nil {
		return nil, nil
	}

	value, err := c.RDB.Get(ctx, constants.RedisSponsorRecord.Get(id).String()).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var record entity.SponsorRecord
	if err = json.Unmarshal([]byte(value), &record); err != nil {
		return nil, err
	}

	return &record, nil
}

func (c *SponsorRecordCache) Set(ctx context.Context, record *entity.SponsorRecord) error {
	if c == nil || c.RDB == nil || record == nil {
		return nil
	}

	payload, err := json.Marshal(record)
	if err != nil {
		return err
	}

	return c.RDB.Set(ctx, constants.RedisSponsorRecord.Get(record.ID).String(), payload, c.TTL).Err()
}

func (c *SponsorRecordCache) Delete(ctx context.Context, id int64) error {
	if c == nil || c.RDB == nil {
		return nil
	}

	return c.RDB.Del(ctx, constants.RedisSponsorRecord.Get(id).String()).Err()
}
