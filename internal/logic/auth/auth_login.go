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

package auth

import (
	"context"
	"github.com/bamboo-services/bamboo-utils/bcode"
	"github.com/bamboo-services/bamboo-utils/berror"
	"github.com/bamboo-services/bamboo-utils/butil"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	v1 "xiaoMain/api/auth/v1"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/entity"
	"xiaoMain/utility"
)

// IsUserLogin
//
// # 用户是否已登录
//
// 用于检查用户是否登录，如果登录则返回 nil, 否则返回错误信息。
//
// # 参数:
//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
//
// # 返回:
//   - err: 如果验证过程中发生错误，返回错误信息。否则返回 nil.
func (s *sAuth) IsUserLogin(ctx context.Context) (err error) {
	g.Log().Notice(ctx, "[LOGIC] AuthLogic:IsUserLogin | 用户是否已登录")
	getUserAuthorize, err := utility.GetAuthorizationFromHeader(ctx)
	if err != nil {
		return err
	}
	// 对内容进行校验
	if getUserAuthorize != "" {
		var getTokenDO *entity.Token
		err := dao.Token.Ctx(ctx).Where(do.Token{Token: getUserAuthorize}).Limit(1).Scan(&getTokenDO)
		if err != nil {
			g.Log().Error(ctx, "[LOGIC] 获取数据库出错", err)
			return berror.NewErrorHasError(bcode.ServerInternalError, err)
		}
		if getTokenDO != nil {
			if gtime.Timestamp() < getTokenDO.ExpiredAt.Timestamp() {
				g.Log().Noticef(ctx, "[LOGIC] 用户UID %s 任然登录状态", getTokenDO.UserUuid)
				return nil
			} else {
				g.Log().Warning(ctx, "[LOGIC] 用户登录过期")
				return berror.NewError(bcode.Expired, "用户登录过期")
			}
		} else {
			g.Log().Warning(ctx, "[LOGIC] 用户登录过期 [无法从数据库取得数据]")
			return berror.NewError(bcode.NotExist, "用户登录过期")
		}
	}
	g.Log().Warning(ctx, "[LOGIC] 用户未登录")
	return berror.NewError(bcode.OperationFailed, "用户未登录")
}

// UserLogin
//
// # 进行用户登录检查
//
// 用于检查用户的登录信息是否正确，如果正确则返回 nil, 否则返回错误信息。
//
// # 参数:
//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
//   - getData: 用户登录信息(v1.AuthLoginReq)
//
// # 返回:
//   - userUUID: 如果用户登录成功，返回用户的 UUID。否则返回 nil.
//   - isCorrect: 如果用户登录成功，返回 true。否则返回 false.
func (s *sAuth) UserLogin(
	ctx context.Context,
	getData *v1.AuthLoginReq,
) (userUUID string, err error) {
	g.Log().Notice(ctx, "[LOGIC] AuthLogic:UserLogin | 进行用户登录检查")
	// 接收数据处理用户登录
	var getUserUUID *entity.Index
	var getUsername *entity.Index
	var getUserPassword *entity.Index
	err = dao.Index.Ctx(ctx).Where(do.Index{Key: "uuid"}).Scan(&getUserUUID)
	if err != nil {
		g.Log().Error(ctx, "[LOGIC] 获取数据库出错", err)
		return "", berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	err = dao.Index.Ctx(ctx).Where(do.Index{Key: "user"}).Scan(&getUsername)
	if err != nil {
		g.Log().Error(ctx, "[LOGIC] 获取数据库出错", err)
		return "", berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	err = dao.Index.Ctx(ctx).Where(do.Index{Key: "password"}).Scan(&getUserPassword)
	if err != nil {
		g.Log().Error(ctx, "[LOGIC] 获取数据库出错", err)
		return "", berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	// 三者不能为空
	if getUserUUID == nil || getUsername == nil || getUserPassword == nil {
		g.Log().Warning(ctx, "[LOGIC] 数据库中未找到用户信息")
		return "", berror.NewError(bcode.ServerInternalError, "数据不存在")
	}
	// 对账号密码进行校验
	if getUsername.Value == getData.User {
		if butil.PasswordVerify(getData.Pass, getUserPassword.Value) {
			g.Log().Notice(ctx, "[LOGIC] 用户校验通过")
			return getUserUUID.Value, nil
		} else {
			g.Log().Notice(ctx, "[LOGIC] 密码错误")
			return "", berror.NewError(bcode.OperationFailed, "密码错误")
		}
	} else {
		g.Log().Notice(ctx, "[LOGIC] 用户名未找到")
		return "", berror.NewError(bcode.NotExist, "用户名不存在")
	}
}

// RegisteredUserLogin
//
// # 注册用户登录
//
// 用于注册用户登录，如果注册成功则返回用户的 Token 信息，否则返回错误信息。
//
// # 参数:
//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
//   - userUUID: 用户的 UUID。
//   - remember: 是否记住密码。
//
// # 返回:
//   - userToken: 如果注册成功，返回用户的 Token 信息。否则返回 nil.
//   - err: 如果注册过程中发生错误，返回错误信息。否则返回 nil.
func (s *sAuth) RegisteredUserLogin(
	ctx context.Context,
	userUUID string,
	remember bool,
) (userToken string, err error) {
	g.Log().Notice(ctx, "[LOGIC] AuthLogic:RegisteredUserLogin | 注册用户登录")
	// 获取允许用户允许登录的节点数
	var getAuthLimit *entity.Index
	err = dao.Index.Ctx(ctx).Where(do.Index{Key: "auth_limit"}).Scan(&getAuthLimit)
	if err != nil {
		g.Log().Error(ctx, "[LOGIC] 获取数据库出错", err)
		return "", berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	// 获取用户的一登录信息
	var getAuthTokenList []entity.Token
	getAuthTokenErr := dao.Token.Ctx(ctx).Where("user_uuid", userUUID).OrderAsc("created_at").Scan(&getAuthTokenList)
	if getAuthTokenErr != nil {
		g.Log().Error(ctx, "[LOGIC] 获取数据库出错", getAuthTokenErr)
		return "", gerror.NewCode(gcode.CodeDbOperationError, getAuthTokenErr.Error())
	}
	// 统计个数
	if len(getAuthTokenList) >= gconv.Int(getAuthLimit.Value) {
		g.Log().Warning(ctx, "[LOGIC] 用户登录节点已满")
		// 删除最早登录的用户内容
		_, deleteError := dao.Token.Ctx(ctx).Delete("id", getAuthTokenList[0].Id)
		if deleteError != nil {
			return "", deleteError
		}
	}
	var remPassword int64
	if remember {
		remPassword = gtime.Timestamp() + 604800
	} else {
		remPassword = gtime.Timestamp() + 86400
	}
	// 数据库操作
	getVerificationCode := butil.GenerateRandUUID().String()
	_, err = dao.Token.Ctx(ctx).Data(do.Token{
		UserUuid:  userUUID,
		UserIp:    g.RequestFromCtx(ctx).GetClientIp(),
		Token:     getVerificationCode,
		UserAgent: g.RequestFromCtx(ctx).GetHeader("User-Agent"),
		ExpiredAt: gtime.NewFromTimeStamp(remPassword),
	}).OnConflict("id").Save()
	if err != nil {
		g.Log().Error(ctx, "[LOGIC] 数据表插入失败", err)
		return "", berror.NewErrorHasError(bcode.ServerInternalError, err)
	}
	return getVerificationCode, nil
}
