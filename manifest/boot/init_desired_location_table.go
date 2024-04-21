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
	"github.com/gogf/gf/v2/os/glog"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
)

func InitialDesiredLocationTable(ctx context.Context) {
	// 记录日志，开始初始化数据表
	glog.Info(ctx, "[BOOT] 初始化期望位置表")

	// 初始化期望位置表
	insertLocationData(ctx, 1, "favorite", "最喜欢", "这是我最喜欢的东西，我当然要置顶啦", true)
	insertLocationData(ctx, 100, "fail", "失效的", "这是失效的友链，希望你快回来嗷", false)
}

// insertLocationData 插入数据，用于信息初始化进行的操作
func insertLocationData(ctx context.Context, sort uint, name string, displayName string, desc string, reveal bool) {
	if record, _ := dao.XfDesiredLocation.Ctx(ctx).Where(do.XfDesiredLocation{Name: name}).One(); record == nil {
		if _, err := dao.XfDesiredLocation.Ctx(ctx).Data(
			do.XfDesiredLocation{
				Sort:        sort,
				Name:        name,
				DisplayName: displayName,
				Description: desc,
				Reveal:      reveal,
			}).Insert(); err != nil {
			glog.Infof(ctx, "[SQL] 数据表 xf_desired_color 中插入键 %s 失败", name)
			glog.Errorf(ctx, "[SQL] 错误信息：%v", err.Error())
		} else {
			glog.Debugf(ctx, "[SQL] 数据表 xf_desired_color 中插入键 %s 成功", name)
		}
	}
}
