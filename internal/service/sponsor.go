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

// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	v1 "xiaoMain/api/sponsor/v1"
	"xiaoMain/internal/model/entity"

	"github.com/google/uuid"
)

type (
	ISponsor interface {
		// AddSponsor
		//
		// # 赞助添加
		//
		// 用于添加赞助，添加一个赞助信息
		//
		// # 接口
		//   - ctx: 上下文
		//   - req: 请求参数
		//
		// # 返回
		//   - err: 错误信息
		AddSponsor(ctx context.Context, req *v1.SponsorAddReq) (err error)
		// GetSponsorByUUID
		//
		// # 获取赞助信息
		//
		// 用于获取赞助信息，从数据库中获取赞助的数据；
		// 获取的数据直接输出出来；
		//
		// # 参数
		//   - ctx: 上下文
		//   - uuid: 赞助UUID
		//
		// # 返回
		//   - err: 错误信息
		GetSponsorByUUID(ctx context.Context, uuid uuid.UUID) (*entity.Sponsor, error)
		// EditSponsor
		//
		// # 赞助编辑
		//
		// 用于编辑赞助，编辑一个赞助信息
		//
		// # 接口
		//   - ctx: 上下文
		//   - uuid: 赞助UUID
		//   - req: 请求参数
		//
		// # 返回
		//   - err: 错误信息
		EditSponsor(ctx context.Context, uuid uuid.UUID, req *v1.SponsorEditReq) (err error)
		// DelSponsor
		//
		// # 赞助删除
		//
		// 用于删除赞助，删除一个赞助信息
		//
		// # 接口
		//   - ctx: 上下文
		//   - uuid: 赞助UUID
		//
		// # 返回
		//   - err: 错误信息
		DelSponsor(ctx context.Context, uuid uuid.UUID) (err error)
		// GetSponsorType
		//
		// # 获取赞助类型
		//
		// 用于获取赞助类型，从数据库中获取赞助类型的数据；
		// 获取的数据直接输出出来；
		// 该接口主要用于前端选择框使用。
		//
		// # 参数
		//   - ctx: 上下文
		//
		// # 返回
		//   - sponsorTypeList: 赞助类型列表
		//   - err: 错误信息
		GetSponsorType(ctx context.Context) ([]*entity.SponsorType, error)
		// GetSingleSponsorTypeByName
		//
		// # 获取单个赞助类型
		//
		// 用于获取单个赞助类型，从数据库中获取赞助类型的数据；
		// 获取的数据直接输出出来；
		//
		// # 参数
		//   - ctx: 上下文
		//   - name: 赞助类型名称
		//
		// # 返回
		//   - sponsorType: 赞助类型
		//   - err: 错误信息
		GetSingleSponsorTypeByName(ctx context.Context, name string) (*entity.SponsorType, error)
		// GetSingleSponsorTypeById
		//
		// # 获取单个赞助类型
		//
		// 用于获取单个赞助类型，从数据库中获取赞助类型的数据；
		// 获取的数据直接输出出来；
		//
		// # 参数
		//   - ctx: 上下文
		//   - id: 赞助类型ID
		//
		// # 返回
		//   - sponsorType: 赞助类型
		//   - err: 错误信息
		GetSingleSponsorTypeById(ctx context.Context, id int) (*entity.SponsorType, error)
		// AddSponsorType
		//
		// # 添加赞助类型
		//
		// 用于添加赞助类型，将赞助类型添加到数据库中；
		// 添加成功返回 nil，否则返回错误。
		//
		// # 参数
		//   - ctx: 上下文
		//   - req: 添加赞助类型的请求
		//
		// # 返回
		//   - err: 错误信息
		AddSponsorType(ctx context.Context, req *v1.SponsorTypeAddReq) (err error)
		// EditSponsorType
		//
		// # 编辑赞助类型
		//
		// 用于编辑赞助类型，将赞助类型编辑到数据库中；
		// 编辑成功返回 nil，否则返回错误。
		//
		// # 参数
		//   - ctx: 上下文
		//   - req: 编辑赞助类型的请求
		//
		// # 返回
		//   - err: 错误信息
		EditSponsorType(ctx context.Context, req *v1.SponsorTypeEditReq) (err error)
		// DelSponsorType
		//
		// # 删除赞助类型
		//
		// 用于删除赞助类型，将赞助类型从数据库中删除；
		// 删除成功返回 nil，否则返回错误。
		//
		// # 参数
		//   - ctx: 上下文
		//   - req: 删除赞助类型的请求
		//
		// # 返回
		//   - err: 错误信息
		DelSponsorType(ctx context.Context, req int) (err error)
	}
)

var (
	localSponsor ISponsor
)

func Sponsor() ISponsor {
	if localSponsor == nil {
		panic("implement not found for interface ISponsor, forgot register?")
	}
	return localSponsor
}

func RegisterSponsor(i ISponsor) {
	localSponsor = i
}
