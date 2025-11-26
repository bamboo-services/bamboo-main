/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2025 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package startup

import (
	"bamboo-main/internal/task"

	xConsts "github.com/bamboo-services/bamboo-base-go/constants"
)

// mailWorker 全局邮件工作协程实例
var mailWorker *task.MailWorker

// MailInit 初始化邮件系统
//
// 创建并启动邮件守护协程，从 Redis 队列消费邮件任务并发送
func (r *Reg) MailInit() {
	r.Serv.Logger.Named(xConsts.LogINIT).Info("初始化邮件系统")

	// 创建邮件工作协程
	mailWorker = task.NewMailWorker(
		r.Rdb,
		&r.Config.Email,
		r.Serv.Logger.Sugar().Named("MAIL"),
	)

	// 启动守护协程
	mailWorker.Start()

	r.Serv.Logger.Named(xConsts.LogINIT).Info("邮件系统初始化完成")
}

// GetMailWorker 获取邮件工作协程实例（用于优雅关闭）
func GetMailWorker() *task.MailWorker {
	return mailWorker
}
