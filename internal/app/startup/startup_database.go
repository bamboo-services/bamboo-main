package startup

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/model/base"
	"github.com/bamboo-services/bamboo-main/pkg/constants"

	xLog "github.com/bamboo-services/bamboo-base-go/log"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/utility/ctxutil"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var migrateTables = []interface{}{
	&entity.SystemUser{},
	&entity.LinkGroup{},
	&entity.LinkColor{},
	&entity.LinkFriend{},
	&entity.SystemLog{},
	&entity.System{},
	&entity.SponsorChannel{},
	&entity.SponsorRecord{},
}

func (r *reg) databaseInit(ctx context.Context) (any, error) {
	log := xLog.WithName(xLog.NamedINIT)
	log.Debug(ctx, "正在连接数据库")

	cfg, cfgErr := xCtxUtil.Get[*base.BambooConfig](ctx, constants.ContextCustomConfig)
	if cfgErr != nil {
		return nil, fmt.Errorf("获取配置失败: %w", cfgErr)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Pass,
		cfg.Database.Name,
		cfg.Database.SSLMode,
		cfg.Database.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.Database.Prefix,
			SingularTable: true,
		},
		Logger: xLog.NewSlogLogger(slog.Default().WithGroup(xLog.NamedREPO), xLog.GormLoggerConfig{
			SlowThreshold:             200,
			LogLevel:                  xLog.LevelInfo,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
		}),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	if err := db.WithContext(ctx).AutoMigrate(migrateTables...); err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %w", err)
	}

	log.Info(ctx, "数据库连接成功")
	return db.WithContext(ctx), nil
}
