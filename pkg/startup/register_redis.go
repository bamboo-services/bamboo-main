package startup

import (
	"context"
	"fmt"
	"time"

	xConsts "github.com/bamboo-services/bamboo-base-go/constants"
	"github.com/redis/go-redis/v9"
)

// RedisInit 初始化Redis连接池。
//
// 此方法根据配置文件创建Redis客户端连接，并测试连接有效性。
// 如果连接失败，程序将触发 panic 并终止运行。
//
// 功能包括:
//   - 根据配置创建Redis客户端
//   - 测试Redis连接可用性
//   - 设置连接池参数
//
// 注意: 使用此方法之前需确保配置文件已正确加载，并且Redis服务正常运行。
func (r *Reg) RedisInit() {
	r.Serv.Logger.Named(xConsts.LogINIT).Info("初始化Redis连接")

	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", r.Config.NoSQL.Host, r.Config.NoSQL.Port),
		Password:     r.Config.NoSQL.Pass,
		DB:           r.Config.NoSQL.Database,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// 测试Redis连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic("[Redis] Redis连接失败: " + err.Error())
	}

	r.Rdb = rdb
	r.Serv.Logger.Named(xConsts.LogINIT).Info("Redis连接初始化完成")
}