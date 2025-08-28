package startup

import (
	"fmt"

	"bamboo-main/internal/model/entity"

	xConsts "github.com/bamboo-services/bamboo-base-go/constants"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var mergeEntities = []interface{}{
	&entity.SystemUser{},
	&entity.LinkGroup{},
	&entity.LinkColor{},
	&entity.LinkFriend{},
	&entity.SystemLog{},
}

// DatabaseInit 初始化数据库连接并执行自动迁移。
//
// 此方法基于配置文件中的数据库设置建立PostgreSQL连接，并自动创建或更新表结构。
// 如果连接失败或自动迁移失败，程序将触发 panic 并终止运行。
//
// 功能包括:
//   - 根据配置创建PostgreSQL数据库连接
//   - 设置表前缀和命名策略
//   - 执行自动数据库迁移
//   - 配置日志级别（调试模式下显示详细日志）
//
// 注意: 使用此方法之前需确保配置文件已正确加载，并且数据库服务正常运行。
func (r *Reg) DatabaseInit() {
	r.Serv.Logger.Named(xConsts.LogINIT).Info("初始化数据库连接")

	// 构建PostgreSQL连接字符串
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		r.Config.Database.Host,
		r.Config.Database.Port,
		r.Config.Database.User,
		r.Config.Database.Pass,
		r.Config.Database.Name,
		r.Config.Database.SSLMode,
		r.Config.Database.TimeZone,
	)

	// 建立数据库连接
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   r.Config.Database.Prefix,
			SingularTable: true,
		},
	})
	if err != nil {
		panic("[Database] 数据库连接失败: " + err.Error())
	}

	// 自动迁移数据库表
	err = db.AutoMigrate(mergeEntities...)
	if err != nil {
		panic("[Database] 数据库迁移失败: " + err.Error())
	}

	// 检查是否启用 Debug 模式
	if r.Config.Xlf.Debug {
		r.Serv.Logger.Named(xConsts.LogINIT).Debug("数据库连接开启 Debug 模式")
		db = db.Debug()
	}

	r.DB = db
	r.Serv.Logger.Named(xConsts.LogINIT).Info("数据库连接初始化完成")
}
