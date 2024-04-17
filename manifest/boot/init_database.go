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
 */

package boot

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"
	uuid "github.com/satori/go.uuid"
	"xiaoMain/internal/consts"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/utility"
)

// InitialDatabase 是一个初始化数据库的函数。
// 它会检查数据库表是否完整，并插入必要的数据。
// 这个函数在应用程序的启动过程中被调用。
//
// 参数:
// ctx: context.Context 对象，用于管理 go 协程和其他与上下文相关的任务。
//
// 返回: 无
func InitialDatabase(ctx context.Context) {
	/*
	 * 检查数据表是否完善
	 */
	// 记录日志，开始初始化数据表
	glog.Info(ctx, "[BOOT] 数据表初始化中")
	// 初始化信息表
	initialSQL(ctx, "xf_index")
	// 初始化登录信息表
	initialSQL(ctx, "xf_token")
	// 初始化日志表
	initialSQL(ctx, "xf_logs")
	// 初始化链接表
	initialSQL(ctx, "xf_link_list")
	// 初始化验证码表
	initialSQL(ctx, "xf_verification_code")

	/**
	 * 检查数据表信息是否完整
	 */
	// 记录日志，开始初始化数据库表信息
	glog.Info(ctx, "[BOOT] 数据库表信息初始化中")
	// 插入软件版本信息
	insertIndexData(ctx, "version", consts.XiaoMainVersion)
	// 插入软件作者信息
	insertIndexData(ctx, "author", consts.XiaoMainAuthor)

	// 生成并插入用户的唯一 UUID
	insertIndexData(ctx, "uuid", uuid.NewV4().String())
	// 新建默认用户
	insertIndexData(ctx, "user", "admin")
	// 设置默认用户密码
	insertIndexData(ctx, "password", utility.PasswordEncode("admin-admin"))
	// 设置默认用户邮箱
	insertIndexData(ctx, "email", "admin@xiaoMain.com")

	// 设置允许登录的节点数（设备数）
	insertIndexData(ctx, "auth_limit", "3")

	// SMTP 邮件服务器配置
	insertIndexData(ctx, "smtp_host", "smtp.x-lf.com")
	// SMTP 邮件服务器端口(默认)
	insertIndexData(ctx, "smtp_port_tls", "25")
	// SMTP 邮件服务器端口(SSL)
	insertIndexData(ctx, "smtp_port_ssl", "465")
	// SMTP 邮件服务器用户名
	insertIndexData(ctx, "smtp_username", "noreplay@xiaoMain.com")
	// SMTP 邮件服务器密码
	insertIndexData(ctx, "smtp_password", "password")

	// 初始化邮件模板(user-change-password)
	insertIndexData(ctx, "mail_template_user_change_password", getMailTemplate("user-change-password"))
	// 初始化邮件模板(user-forget-password)
	insertIndexData(ctx, "mail_template_user_forget_password", getMailTemplate("user-forget-password"))

	/**
	 * 初始化完毕结束任务
	 */
	// 记录日志，数据表初始化完毕
	glog.Info(ctx, "[BOOT] 数据表初始化完毕")
}

func getMailTemplate(template string) string {
	return gfile.GetContents("resource/template/mail/" + template + ".html")
}

// insertIndexData 插入数据，用于信息初始化进行的操作
func insertIndexData(ctx context.Context, key string, value string) {
	var err error
	if record, _ := dao.XfIndex.Ctx(ctx).Where("key=?", key).One(); record == nil {
		if _, err = dao.XfIndex.Ctx(ctx).Data(do.XfIndex{Key: key, Value: value}).Insert(); err != nil {
			glog.Infof(ctx, "[SQL] 数据表 xf_index 中插入键 %s 失败", key)
			glog.Errorf(ctx, "[SQL] 错误信息：%v", err.Error())
		} else {
			glog.Debugf(ctx, "[SQL] 数据表 xf_index 中插入键 %s成功", key)
		}
	}
}

// initialSQL 初始化数据库
func initialSQL(ctx context.Context, databaseName string) {
	if _, err := g.DB().Exec(ctx, "SELECT * FROM "+databaseName); err != nil {
		// 创建数据表
		errTransaction := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// 读取文件
			getFileContent := gfile.GetContents("internal/sql/" + databaseName + ".sql")
			// 创建 xf_index.sql 表
			if _, err := tx.Exec(getFileContent); err != nil {
				return err
			}
			return nil
		})
		if errTransaction != nil {
			glog.Panicf(ctx, "[SQL] 数据表 %s 创建失败", databaseName)
		} else {
			glog.Debugf(ctx, "[SQL] 数据表 %s 创建成功", databaseName)
		}
	} else {
		glog.Debugf(ctx, "[SQL] 数据表 %s 已存在", databaseName)
	}
}
