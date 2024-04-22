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

func InitialDesiredColorTable(ctx context.Context) {
	glog.Info(ctx, "[BOOT] 初始化期望颜色表")

	// 初始化期望颜色表
	insertColorData(ctx, "red", "红色", "FF0000")
	insertColorData(ctx, "orange", "橙色", "FFA500")
	insertColorData(ctx, "yellow", "黄色", "FFFF00")
	insertColorData(ctx, "green", "绿色", "008000")
	insertColorData(ctx, "cyan", "青色", "00FFFF")
	insertColorData(ctx, "blue", "蓝色", "0000FF")
	insertColorData(ctx, "purple", "紫色", "800080")
	insertColorData(ctx, "pink", "粉色", "FFC0CB")
	insertColorData(ctx, "black", "黑色", "000000")
	insertColorData(ctx, "white", "白色", "FFFFFF")
	insertColorData(ctx, "gray", "灰色", "808080")
	insertColorData(ctx, "brown", "棕色", "A52A2A")
	insertColorData(ctx, "gold", "金色", "FFD700")
	insertColorData(ctx, "silver", "银色", "C0C0C0")
	insertColorData(ctx, "bronze", "铜色", "CD7F32")
	insertColorData(ctx, "rose", "玫瑰金", "FFC0CB")
	insertColorData(ctx, "champagne", "香槟金", "F7E7CE")
	insertColorData(ctx, "peach", "桃红", "FFDAB9")
	insertColorData(ctx, "apricot", "杏色", "FBCEB1")
	insertColorData(ctx, "coral", "珊瑚红", "FF7F50")
	insertColorData(ctx, "salmon", "鲑鱼红", "FA8072")
	insertColorData(ctx, "tomato", "番茄红", "FF6347")
	insertColorData(ctx, "maroon", "栗色", "800000")
	insertColorData(ctx, "burgundy", "酒红", "800020")
	insertColorData(ctx, "ruby", "红宝石", "E0115F")
	insertColorData(ctx, "sapphire", "蓝宝石", "0F52BA")
	insertColorData(ctx, "emerald", "翡翠绿", "50C878")
	insertColorData(ctx, "amethyst", "紫水晶", "9966CC")
	insertColorData(ctx, "topaz", "黄玉", "FFC87C")
	insertColorData(ctx, "turquoise", "绿松石", "40E0D0")
	insertColorData(ctx, "aquamarine", "海蓝宝石", "7FFFD4")
	insertColorData(ctx, "peridot", "橄榄石", "E6E200")
	insertColorData(ctx, "opal", "蛋白石", "A8C3BC")
	insertColorData(ctx, "pearl", "珍珠", "F0EAD6")
	insertColorData(ctx, "moonstone", "月光石", "E3E4FA")
	insertColorData(ctx, "diamond", "钻石", "B9F2FF")
}

// insertColorData 插入数据，用于信息初始化进行的操作
func insertColorData(ctx context.Context, name string, displayName string, color string) {
	if record, _ := dao.XfColor.Ctx(ctx).Where(do.XfColor{Name: name}).One(); record == nil {
		if _, err := dao.XfColor.Ctx(ctx).Data(
			do.XfColor{
				Name:        name,
				DisplayName: displayName,
				Color:       color,
			}).Insert(); err != nil {
			glog.Infof(ctx, "[SQL] 数据表 xf_desired_color 中插入键 %s 失败", name)
			glog.Errorf(ctx, "[SQL] 错误信息：%v", err.Error())
		} else {
			glog.Debugf(ctx, "[SQL] 数据表 xf_desired_color 中插入键 %s 成功", name)
		}
	}
}
