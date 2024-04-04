/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋(https://www.x-lf.com)
 *
 * 本文件包含 XiaoMain 的源代码，该项目的所有源代码均遵循MIT开源许可证协议。
 * --------------------------------------------------------------------------------
 * 许可证声明：
 *
 * 版权所有 (c) 2016-2024 筱锋。保留所有权利。
 *
 * 本软件是“按原样”提供的，没有任何形式的明示或暗示的保证，包括但不限于
 * 对适销性、特定用途的适用性和非侵权性的暗示保证。在任何情况下，
 * 作者或版权持有人均不承担因软件或软件的使用或其他交易而产生的、
 * 由此引起的或以任何方式与此软件有关的任何索赔、损害或其他责任。
 *
 * 使用本软件即表示您了解此声明并同意其条款。
 *
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 * 免责声明：
 *
 * 使用本软件的风险由用户自担。作者或版权持有人在法律允许的最大范围内，
 * 对因使用本软件内容而导致的任何直接或间接的损失不承担任何责任。
 * --------------------------------------------------------------------------------
 *
 */

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
			if _, err := tx.Exec("create table xf_token( id bigserial constraint xf_token_pk primary key, user_uuid uuid not null, user_token uuid not null, user_ip varchar not null, user_agent varchar not null, verification uuid not null, created_at timestamp default now() not null, expired_at timestamp)"); err != nil {
				return err
			}
			// 创建 COMMIT
			if _, err := tx.Exec("comment on table xf_token is '登录信息表'; comment on column xf_token.id is '主键'; comment on column xf_token.user_uuid is '用户 UUID'; comment on column xf_token.user_token is '用户 TOKEN'; comment on column xf_token.user_ip is '用户 IP 地址'; comment on column xf_token.user_agent is '用户 Agent'; comment on column xf_token.verification is '验证用户是否是唯一用户'; comment on column xf_token.created_at is '创建时间'; comment on column xf_token.expired_at is '修改时间'"); err != nil {
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
	// 初始化链接表
	if _, err := dao.XfLinkList.DB().Exec(ctx, "SELECT * FROM xf_link_list"); err != nil {
		err := dao.XfIndex.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// 创建 xf_link_list 表
			if _, err := tx.Exec("create table xf_link_list( id bigserial constraint xf_link_list_pk primary key, webmaster_email varchar, service_provider varchar not null, site_name varchar not null, site_url varchar not null, site_logo varchar not null, site_description varchar not null, site_rss_url varchar, has_adv boolean default false not null, desired_location integer default 0 not null, location integer default '-1'::integer not null, desired_color integer default 0 not null, color integer default '-1'::integer not null, webmaster_remark varchar, remark varchar, status smallint default 0 not null, created_at timestamp default now() not null, updated_at timestamp, deleted_at timestamp)"); err != nil {
				return err
			}
			// 创建 COMMIT
			if _, err := tx.Exec("comment on table xf_link_list is '链接列表'; comment on column xf_link_list.id is '主键'; comment on column xf_link_list.webmaster_email is '站长邮箱'; comment on column xf_link_list.service_provider is '服务提供商'; comment on column xf_link_list.site_name is '站点名字'; comment on column xf_link_list.site_url is '站点地址'; comment on column xf_link_list.site_logo is '站点 logo'; comment on column xf_link_list.site_description is '站点描述'; comment on column xf_link_list.site_rss_url is '站点订阅地址'; comment on column xf_link_list.has_adv is '是否有广告'; comment on column xf_link_list.desired_location is '理想位置'; comment on column xf_link_list.location is '所在位置'; comment on column xf_link_list.desired_color is '理想颜色'; comment on column xf_link_list.color is '颜色'; comment on column xf_link_list.webmaster_remark is '站长留言'; comment on column xf_link_list.remark is '我的留言'; comment on column xf_link_list.status is '0 待审核，1 通过，-1 审核拒绝'; comment on column xf_link_list.created_at is '创建时间'; comment on column xf_link_list.updated_at is '修改时间'; comment on column xf_link_list.deleted_at is '删除时间'"); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			panic("\t> 数据表 xf_link_list 创建失败")
		} else {
			glog.Debug(ctx, "\t> 数据表 xf_link_list 创建成功")
		}
	} else {
		glog.Debug(ctx, "\t> 数据表 xf_link_list 已存在")
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
	insertData(ctx, "auth_limit", "3") // 允许登录的节点数（设备数）
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
