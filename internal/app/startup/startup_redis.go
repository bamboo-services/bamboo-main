package startup

import (
	"context"
	"fmt"
	"time"

	"github.com/bamboo-services/bamboo-main/internal/models/base"
	"github.com/bamboo-services/bamboo-main/pkg/constants"

	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/common/utility/context"
	"github.com/redis/go-redis/v9"
)

func (r *reg) nosqlInit(ctx context.Context) (any, error) {
	log := xLog.WithName(xLog.NamedINIT)
	log.Debug(ctx, "正在连接缓存")

	cfg, err := xCtxUtil.Get[*base.BambooConfig](ctx, constants.ContextCustomConfig)
	if err != nil {
		return nil, fmt.Errorf("获取配置失败: %w", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.NoSQL.Host, cfg.NoSQL.Port),
		Password:     cfg.NoSQL.Pass,
		DB:           cfg.NoSQL.Database,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := rdb.Ping(pingCtx).Err(); err != nil {
		return nil, fmt.Errorf("连接 Redis 失败: %w", err)
	}

	log.Info(ctx, "缓存连接成功")
	return rdb, nil
}
