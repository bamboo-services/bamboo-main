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

package startup

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
)

// initialSQL
//
// # 初始化 SQL
//
// 是一个初始化 SQL 的函数。
// 它会检查数据库表是否完整，并插入必要的数据。
// 这个函数在应用程序的启动过程中被调用。
//
// 参数:
//   - CTX: 上下文对象，用于传递和控制请求的生命周期。
//   - databaseName: 数据库名称
func (is *InitStruct) initialSQL(ctx context.Context, databaseName string) {
	_, err := g.DB().Exec(ctx, "SELECT * FROM information_schema.tables WHERE table_name = ?", databaseName)
	if err != nil {
		// 创建数据表
		errTransaction := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// 读取文件
			getFileContent := gfile.GetContents("resource/sql/" + databaseName + ".sql")
			// 创建 xf_index.sql 表
			if _, err := tx.Exec(getFileContent); err != nil {
				return err
			}
			return nil
		})
		if errTransaction != nil {
			g.Log().Panicf(ctx, "[BOOT] 数据表 %s 创建失败", databaseName)
		} else {
			g.Log().Debugf(ctx, "[BOOT] 数据表 %s 创建成功", databaseName)
		}
	}
}

// getMailTemplate
//
// # 获取邮件模板
//
// 是一个获取邮件模板的函数。
// 它会读取邮件模板文件，并返回文件内容。
//
// 参数:
//   - template: 模板名称
func (is *InitStruct) getMailTemplate(template string) string {
	return gfile.GetContents("resource/template/mail/" + template + ".html")
}

// insertIndexData
//
// # 插入索引数据
//
// 是一个插入索引数据的函数。
// 它会检查数据库表是否完整，并插入必要的数据。
// 这个函数在应用程序的启动过程中被调用。
//
// 参数:
//   - CTX: 上下文对象，用于传递和控制请求的生命周期。
//   - key: 索引键
//   - value: 索引值
func (is *InitStruct) insertIndexData(ctx context.Context, key string, value string) {
	var err error
	if record, _ := dao.Index.Ctx(ctx).Where(do.Index{Key: key}).One(); record == nil {
		if _, err = dao.Index.Ctx(ctx).Data(do.Index{Key: key, Value: value}).Insert(); err != nil {
			g.Log().Noticef(ctx, "[BOOT] 数据表 xf_index 中插入键 %s 失败", key)
			g.Log().Errorf(ctx, "[BOOT] 错误信息：%v", err.Error())
		} else {
			g.Log().Debugf(ctx, "[BOOT] 数据表 xf_index 中插入键 %s 成功", key)
		}
	}
}

// insertLocationData
//
// # 插入位置数据
//
// 插入数据，用于信息初始化进行的操作。
//
// 参数:
//   - CTX: 上下文对象，用于传递和控制请求的生命周期。
//   - sort: 排序
//   - name: 名称
//   - displayName: 显示名称
//   - desc: 描述
//   - reveal: 是否显示
func (is *InitStruct) insertLocationData(
	ctx context.Context,
	sort uint,
	name string,
	displayName string,
	desc string,
	reveal bool,
) {
	if record, _ := dao.Location.Ctx(ctx).Where(do.Location{Name: name}).One(); record == nil {
		if _, err := dao.Location.Ctx(ctx).Data(
			do.Location{
				Sort:        sort,
				Name:        name,
				DisplayName: displayName,
				Description: desc,
				Reveal:      reveal,
			}).Insert(); err != nil {
			g.Log().Noticef(ctx, "[BOOT] 数据表 xf_desired_color 中插入键 %s 失败", name)
			g.Log().Errorf(ctx, "[BOOT] 错误信息：%v", err.Error())
		} else {
			g.Log().Debugf(ctx, "[BOOT] 数据表 xf_desired_color 中插入键 %s 成功", name)
		}
	}
}

// insertColorData
//
// # 插入颜色数据
//
// 插入数据，用于信息初始化进行的操作。
//
// 参数:
//   - CTX: 上下文对象，用于传递和控制请求的生命周期。
//   - name: 名称
//   - displayName: 显示名称
//   - color: 颜色
func (is *InitStruct) insertColorData(ctx context.Context, name string, displayName string, color string) {
	if record, _ := dao.Color.Ctx(ctx).Where(do.Color{Name: name}).One(); record == nil {
		if _, err := dao.Color.Ctx(ctx).Data(
			do.Color{
				Name:        name,
				DisplayName: displayName,
				Color:       color,
			}).Insert(); err != nil {
			g.Log().Noticef(ctx, "[BOOT] 数据表 xf_desired_color 中插入键 %s 失败", name)
			g.Log().Errorf(ctx, "[BOOT] 错误信息：%v", err.Error())
		} else {
			g.Log().Debugf(ctx, "[BOOT] 数据表 xf_desired_color 中插入键 %s 成功", name)
		}
	}
}

// prepareCommonData
//
// # 准备常量数据
//
// 是一个准备常量数据的函数。
// 将会从数据库中读取数据，并赋值给指定的变量。
//
// 参数:
//   - key: 键
//
// 返回:
//   - interface{}: 值
func (is *InitStruct) prepareCommonData(key string) interface{} {
	record, err := dao.Index.Ctx(is.CTX).Where(do.Index{Key: key}).One()
	if err != nil {
		g.Log().Panicf(is.CTX, "[BOOT] 准备常量 %s, 错误: %v", key, err)
	}
	if !record.IsEmpty() {
		return record.GMap().Get("value")
	} else {
		return nil
	}
}
