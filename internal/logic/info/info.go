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

func (s *sInfoLogic) GetMainInfo(ctx context.Context) (getMainInfo *vo.MainVO) {
	glog.Info(ctx, "[LOGIC] 执行 InfoLogic:GetMainInfo 服务层")
	getMainInfo = new(vo.MainVO)
	wg := sync.WaitGroup{}
	// 获取名字信息
	getMainInfo.Name = s.getMainInfoData(ctx, &wg, "name")
	// 获取版本信息
	getMainInfo.Version = s.getMainInfoData(ctx, &wg, "version")
	// 获取作者信息
	getMainInfo.Author = s.getMainInfoData(ctx, &wg, "author")
	// TODO-Info[2024050401] 还有其他信息等待获取

	// 等待所有协程执行完毕
	wg.Wait()
	// 处理返回错误
	return getMainInfo
}

func (s *sInfoLogic) getMainInfoData(ctx context.Context, wg *sync.WaitGroup, name string) (getData string) {
	var getInfo *entity.XfIndex
	wg.Add(1)
	go func() {
		defer wg.Done()
		if name != "" {
			_ = dao.XfIndex.Ctx(ctx).Where(do.XfIndex{Key: name}).Scan(&getInfo)
		}
	}()
	if getInfo != nil {
		return getInfo.Value
	} else {
		return ""
	}
}
