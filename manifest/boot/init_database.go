package boot

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"xiaoMain/internal/dao"
)

// InitialDatabase 数据库初始化操作
func InitialDatabase(ctx context.Context) {
	/*
	 * 检查数据表是否完善
	 */
	glog.Info(ctx, "[BOOT] 数据表初始化中")
	// 初始化信息表
	if _, err := dao.XfIndex.DB().Exec(ctx, "SELECT * FROM xf_index"); err != nil {
		// 创建数据表
		err := dao.XfIndex.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// 创建 xf_index 表
			if _, err := tx.Exec("CREATE TABLE xf_index ( id serial CONSTRAINT xf_index_pk PRIMARY KEY, key varchar(40) NOT NULL, value text, created_at timestamp NOT NULL DEFAULT now(), updated_at timestamp )"); err != nil {
				return err
			}
			// 创建 COMMIT
			if _, err := tx.Exec("comment on table xf_index is '筱主页信息数据表';comment on column xf_index.id is '主键';comment on column xf_index.key is '键';comment on column xf_index.value is '值';comment on column xf_index.created_at is '创建时间';comment on column xf_index.updated_at is '修改时间'"); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			panic("\t> 数据表 xf_index 创建失败")
		} else {
			glog.Debug(ctx, "\t> 数据表 xf_index 创建成功")
		}
	} else {
		glog.Debug(ctx, "\t> 数据表 xf_index 已存在")
	}
	// 初始化登录信息表
	if _, err := dao.XfToken.DB().Exec(ctx, "SELECT * FROM xf_token"); err != nil {
		err := dao.XfIndex.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// 创建 xf_logs 表
			if _, err := tx.Exec("create table xf_token( id bigserial constraint xf_token_pk primary key, user_uuid uuid NOT null, user_token uuid NOT null, created_at timestamp default now() NOT null, expired_at timestamp NOT NULL)"); err != nil {
				return err
			}
			// 创建 COMMIT
			if _, err := tx.Exec("comment ON table xf_token is '登录信息表'; comment ON column xf_token.id is '主键'; comment ON column xf_token.user_uuid is '用户 UUID'; comment ON column xf_token.user_token is '用户 TOKEN'; comment ON column xf_token.created_at is '创建时间'; comment ON column xf_token.expired_at is '过期时间';"); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			panic("\t> 数据表 xf_token 创建失败")
		} else {
			glog.Debug(ctx, "\t> 数据表 xf_token 创建成功")
		}
	} else {
		glog.Debug(ctx, "\t> 数据表 xf_token 已存在")
	}
	// 初始化日志表
	if _, err := dao.XfLogs.DB().Exec(ctx, "SELECT * FROM xf_logs"); err != nil {
		err := dao.XfIndex.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// 创建 xf_logs 表
			if _, err := tx.Exec("CREATE TABLE xf_logs ( id bigserial CONSTRAINT xf_logs_pk PRIMARY KEY, type smallint NOT NULL DEFAULT 0, log varchar NOT NULL, created_at timestamp NOT NULL DEFAULT now() )"); err != nil {
				return err
			}
			// 创建 COMMIT
			if _, err := tx.Exec("comment on table xf_logs is '日志数据表'; comment on column xf_logs.id is '主键'; comment on column xf_logs.type is '日志类型'; comment on column xf_logs.log is '日志内容'; comment on column xf_logs.created_at is '日志时间'"); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			panic("\t> 数据表 xf_logs 创建失败")
		} else {
			glog.Debug(ctx, "\t> 数据表 xf_logs 创建成功")
		}
	} else {
		glog.Debug(ctx, "\t> 数据表 xf_logs 已存在")
	}

	/**
	 * 检查数据表信息是否完整
	 */
	glog.Info(ctx, "[BOOT] 数据库表信息初始化中")
	insertData(ctx, "version", "1.0.0")                      // 插入版本信息
	insertData(ctx, "uuid", fmt.Sprintf("%s", uuid.NewV4())) // 生成用户的唯一 UUID
	insertData(ctx, "user", "admin")                         // 新建默认用户
	// 设置初始化密码
	getBase64Password := base64.StdEncoding.EncodeToString([]byte("admin-admin"))
	getEncodePassword, err := bcrypt.GenerateFromPassword([]byte(getBase64Password), bcrypt.DefaultCost)
	if err == nil {
		insertData(ctx, "password", string(getEncodePassword)) // 默认用户密码
	} else {
		glog.Error(ctx, "[BOOT] 密码加密失败")
	}
}

// insertData 插入数据，用于信息初始化进行的操作
func insertData(ctx context.Context, key string, value string) {
	var err error
	if record, _ := dao.XfIndex.Ctx(ctx).Where("key='" + key + "'").One(); record == nil {
		if _, err = dao.XfIndex.Ctx(ctx).Data(g.Map{"Key": key, "Value": value}).Insert(); err != nil {
			glog.Errorf(ctx, "\t> 数据表 xf_index 中 %s 数据创建失败", key)
		} else {
			glog.Debugf(ctx, "\t> 数据表 xf_index 中 %s 数据创建成功", key)
		}
	} else {
		glog.Debugf(ctx, "\t> 数据表 xf_index 中 %s 数据已存在", key)
	}
}
