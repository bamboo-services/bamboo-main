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
	"github.com/gogf/gf/v2/frame/g"
	"sync"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/dto/flow"
	"xiaoMain/internal/model/entity"
	"xiaoMain/internal/service"
)

type sInfo struct {
}

func init() {
	service.RegisterInfo(New())
}

func New() *sInfo {
	return &sInfo{}
}

// GetMainInfo
//
// # 获取主要信息
//
// 用于获取主要信息，如果返回成功则返回具体的信息，若某些情况下无法获取则获取的内容为空
//
// # 参数:
//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
//
// # 返回:
//   - *vo.SiteMainDTO: 如果获取成功，返回具体的信息；否则返回空值。
func (s *sInfo) GetMainInfo(ctx context.Context) *flow.SiteMainDTO {
	g.Log().Notice(ctx, "[LOGIC] 执行 InfoLogic:GetMainInfo 服务层")
	getMainInfo := new(flow.SiteMainDTO)
	wg := sync.WaitGroup{}
	// 获取名字信息
	wg.Add(5)
	go func() {
		defer wg.Done()
		getMainInfo.SiteName = s.GetIndexTableData(ctx, "site_name")
	}()
	// 获取版本信息
	go func() {
		defer wg.Done()
		getMainInfo.Version = s.GetIndexTableData(ctx, "version")
	}()
	// 获取作者信息
	go func() {
		defer wg.Done()
		getMainInfo.Author = s.GetIndexTableData(ctx, "author")
	}()
	// 获取站点描述
	go func() {
		defer wg.Done()
		getMainInfo.Description = s.GetIndexTableData(ctx, "description")
	}()
	// 获取站点关键字
	go func() {
		defer wg.Done()
		getMainInfo.Keywords = s.GetIndexTableData(ctx, "keywords")
	}()

	// 等待所有协程执行完毕
	wg.Wait()
	// 处理返回错误
	return getMainInfo
}

// GetBloggerInfo
//
// # 获取站长信息
//
// 用于获取站长信息，如果返回成功则返回具体的信息，若某些情况下无法获取则获取的内容为空
//
// # 参数:
//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
//
// # 返回:
//   - *vo.SiteBloggerDTO: 如果获取成功，返回具体的信息；否则返回空值。
func (s *sInfo) GetBloggerInfo(ctx context.Context) *flow.SiteBloggerDTO {
	g.Log().Notice(ctx, "[LOGIC] 执行 InfoLogic:GetBloggerInfo 服务层")
	getBloggerInfo := new(flow.SiteBloggerDTO)
	wg := sync.WaitGroup{}
	// 获取站长名字
	wg.Add(4)
	go func() {
		defer wg.Done()
		getBloggerInfo.Name = s.GetIndexTableData(ctx, "blogger_name")
	}()
	// 获取站长昵称
	go func() {
		defer wg.Done()
		getBloggerInfo.Nick = s.GetIndexTableData(ctx, "blogger_nick")
	}()
	// 获取站长邮箱
	go func() {
		defer wg.Done()
		getBloggerInfo.Email = s.GetIndexTableData(ctx, "email")
	}()
	// 获取站长的描述
	go func() {
		defer wg.Done()
		getBloggerInfo.Description = s.GetIndexTableData(ctx, "blogger_description")
	}()

	// 等待所有的进程结束
	wg.Wait()
	return getBloggerInfo
}

// GetIndexTableData
//
// # 获取 Index 数据库中的信息
//
// 用于获取 Index 数据库中的信息，如果成功则返回具体的信息，否则返回空值
//
// # 参数:
//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
//   - name: 需要获取的信息名称(string)
//
// # 返回:
//   - string: 如果获取成功，返回具体的信息；否则返回空值。
func (s *sInfo) GetIndexTableData(ctx context.Context, name string) string {
	g.Log().Infof(ctx, "[LOGIC-PRIVATE] 获取 Index 数据库中 %s 的信息", name)
	var getInfo *entity.Index
	if name != "" {
		err := dao.Index.Ctx(ctx).Where(do.Index{Key: name}).Scan(&getInfo)
		if err == nil {
			g.Log().Debugf(ctx, "[SQL] 数据表 xf_index 中获取键 %s 成功, 值为 %s", name, getInfo.Value)
			return getInfo.Value
		}
	}
	return ""
}
