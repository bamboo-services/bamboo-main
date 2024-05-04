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

package info

import (
	"context"
	"github.com/gogf/gf/v2/os/glog"
	"sync"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/entity"
	"xiaoMain/internal/model/vo"
	"xiaoMain/internal/service"
)

type sInfoLogic struct {
}

func init() {
	service.RegisterInfoLogic(New())
}

func New() *sInfoLogic {
	return &sInfoLogic{}
}

// GetMainInfo 获取主要信息
// 用于获取主要信息，如果返回成功则返回具体的信息，若某些情况下无法获取则获取的内容为空
// 接口的返回都会有结果，如果返回错误将会返回空值
//
// 参数：
// ctx: 请求的上下文，用于管理超时和取消信号。
//
// 返回：
// getMainInfo: 如果获取成功，返回具体的信息；否则返回空值。
func (s *sInfoLogic) GetMainInfo(ctx context.Context) *vo.MainVO {
	glog.Notice(ctx, "[LOGIC] 执行 InfoLogic:GetMainInfo 服务层")
	getMainInfo := new(vo.MainVO)
	wg := sync.WaitGroup{}
	// 获取名字信息
	wg.Add(5)
	go func() {
		defer wg.Done()
		getMainInfo.SiteName = s.getXfIndexTableData(ctx, "site_name")
	}()
	// 获取版本信息
	go func() {
		defer wg.Done()
		getMainInfo.Version = s.getXfIndexTableData(ctx, "version")
	}()
	// 获取作者信息
	go func() {
		defer wg.Done()
		getMainInfo.Author = s.getXfIndexTableData(ctx, "author")
	}()
	// 获取站点描述
	go func() {
		defer wg.Done()
		getMainInfo.Description = s.getXfIndexTableData(ctx, "description")
	}()
	// 获取站点关键字
	go func() {
		defer wg.Done()
		getMainInfo.Keywords = s.getXfIndexTableData(ctx, "keywords")
	}()

	// 等待所有协程执行完毕
	wg.Wait()
	// 处理返回错误
	return getMainInfo
}

// GetBloggerInfo 获取站长信息
// 主要用于获取站点作者的一些基本信息，如果返回成功则返回具体的信息，若某些情况下无法获取则获取的内容为空
// 接口的返回都会有结果，如果返回错误将会返回空值
//
// 参数：
// ctx: 请求的上下文，用于管理超时和取消信号。
//
// 返回：
// getBloggerInfo: 如果获取成功，返回具体的信息；否则返回空值。
func (s *sInfoLogic) GetBloggerInfo(ctx context.Context) *vo.BloggerVO {
	glog.Notice(ctx, "[LOGIC] 执行 InfoLogic:GetBloggerInfo 服务层")
	getBloggerInfo := new(vo.BloggerVO)
	wg := sync.WaitGroup{}
	// 获取站长名字
	wg.Add(4)
	go func() {
		defer wg.Done()
		getBloggerInfo.Name = s.getXfIndexTableData(ctx, "blogger_name")
	}()
	// 获取站长昵称
	go func() {
		defer wg.Done()
		getBloggerInfo.Nick = s.getXfIndexTableData(ctx, "blogger_nick")
	}()
	// 获取站长邮箱
	go func() {
		defer wg.Done()
		getBloggerInfo.Email = s.getXfIndexTableData(ctx, "email")
	}()
	// 获取站长的描述
	go func() {
		defer wg.Done()
		getBloggerInfo.Description = s.getXfIndexTableData(ctx, "blogger_description")
	}()

	// 等待所有的进程结束
	wg.Wait()
	return getBloggerInfo
}

// getXfIndexTableData 从数据库获取主要信息
// 用于获取主要信息，如果返回成功则返回具体的信息，若某些情况下无法获取则获取的内容为空
// 忽略了报错行为，所以获取不到数据都会返回空的内容，若成功获取会返回具体的内容
//
// 参数：
// ctx: 请求的上下文，用于管理超时和取消信号。
// wg: 等待组，用于等待协程执行完毕。
// name: 获取的信息名称。
//
// 返回：
// getData: 如果获取成功，返回具体的信息；否则返回空值。
func (s *sInfoLogic) getXfIndexTableData(ctx context.Context, name string) string {
	glog.Infof(ctx, "[LOGIC-PRIVATE] 获取 XfIndex 数据库中 %s 的信息", name)
	var getInfo *entity.XfIndex
	if name != "" {
		err := dao.XfIndex.Ctx(ctx).Where(do.XfIndex{Key: name}).Scan(&getInfo)
		if err == nil {
			glog.Debugf(ctx, "[SQL] 数据表 xf_index 中获取键 %s 成功, 值为 %s", name, getInfo.Value)
			return getInfo.Value
		}
	}
	return ""
}
