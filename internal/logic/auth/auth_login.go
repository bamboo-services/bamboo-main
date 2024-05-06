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
	"encoding/base64"
	"errors"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	v1 "xiaoMain/api/auth/v1"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/do"
	"xiaoMain/internal/model/entity"
	"xiaoMain/utility"
	"xiaoMain/utility/result"
)

const ServerInternalErrorString = "服务器内部错误"

// IsUserLogin 是一个用于检查用户是否已经登录的函数。
// 它主要用于获取用户的 UUID 和认证密钥，并对这些信息进行校验。如果用户已经完成登录的相关操作，并且此次登录有效，则返回 true 和空字符串。否则，返回 false 和错误信息。
//
// 参数:
// ctx: 上下文对象，用于传递和控制请求的生命周期。
//
// 返回值:
// hasLogin: 如果用户已经登录并且此次登录有效，返回 true。否则返回 false。
// message: 如果用户未登录或登录已失效，返回错误信息。否则返回空字符串。
func (s *sAuthLogic) IsUserLogin(ctx context.Context) (hasLogin bool, message string) {
	glog.Notice(ctx, "[LOGIC] 执行 AuthLogic:IsUserLogin 服务层")
	// 根据 ctx 获取 Request 信息
	getRequest := ghttp.RequestFromCtx(ctx)
	// 获取用户的 UUID(UID) 以及 认证密钥
	getUserUUID, err := utility.GetUUIDFromHeader(getRequest)
	if err != nil {
		return false, err.Error()
	}
	getUserAuthorize, err := utility.GetAuthorizationFromHeader(getRequest)
	if err != nil {
		return false, err.Error()
	}
	// 对内容进行校验
	if getUserUUID != nil && getUserAuthorize != nil && *getUserUUID != "" && *getUserAuthorize != "" {
		var getTokenDO entity.XfToken
		err := dao.XfToken.Ctx(ctx).
			Where(do.XfToken{UserUuid: getUserUUID, UserToken: getUserAuthorize}).
			Limit(1).Scan(&getTokenDO)
		if err != nil {
			glog.Warning(ctx, "[LOGIC] 用户登录过期 [无法从数据库取得数据]")
			return false, "用户登录已失效"
		}
		// 检查是否过期
		if gtime.Timestamp() < getTokenDO.ExpiredAt.Timestamp() {
			// 验证登录有效
			if *getUserAuthorize == getTokenDO.UserToken {
				glog.Noticef(ctx, "[LOGIC] 用户UID %s 任然登录状态", getTokenDO.UserUuid)
				return true, ""
			} else {
				glog.Warning(ctx, "[LOGIC] 用户登录已失效 [数据已取得，但已失效]")
				return false, "用户登录已失效"
			}
		}
	}
	glog.Warning(ctx, "[LOGIC] 用户未登录")
	return false, "用户未登录"
}

// CheckUserLogin 是一个用于检查用户登录的函数。
// 它主要用于对用户输入的信息与数据库的内容进行校验，当用户名与用户校验通过后 isCorrect 返回正确值，否则返回错误的内容。
// 并且当用户正常登录后，将会返回用户的 UUID 作为下一步的登录操作。
//
// 参数:
// ctx: 上下文对象，用于传递和控制请求的生命周期。
// getData: 用户的登录请求数据，包含了用户的用户名和密码。
//
// 返回值:
// userUUID: 如果用户登录成功，返回用户的 UUID 字符串。
// isCorrect: 如果用户登录成功，返回 true。否则返回 false。
// errMessage: 如果用户登录失败，返回错误信息。否则返回空字符串。
func (s *sAuthLogic) CheckUserLogin(
	ctx context.Context,
	getData *v1.AuthLoginReq,
) (userUUID *string, isCorrect bool, errMessage string) {
	glog.Notice(ctx, "[LOGIC] 执行 AuthLogic:CheckUserLogin 服务层")
	// 根据 ctx 获取 Request 信息
	getRequest := ghttp.RequestFromCtx(ctx)
	// 接收数据处理用户登录
	var userUUIDEntity entity.XfIndex
	var userUsernameEntity entity.XfIndex
	var userPasswordEntity entity.XfIndex
	getUserUUIDErr := dao.XfIndex.Ctx(ctx).Where("key = ?", "uuid").Scan(&userUUIDEntity)
	if getUserUUIDErr != nil {
		glog.Error(ctx, "[LOGIC] 获取数据库出错", getUserUUIDErr)
		result.DatabaseError.Response(getRequest)
		return nil, false, ServerInternalErrorString
	}
	getUserUsernameErr := dao.XfIndex.Ctx(ctx).Where("key = ?", "user").Scan(&userUsernameEntity)
	if getUserUsernameErr != nil {
		glog.Error(ctx, "[LOGIC] 获取数据库出错", getUserUsernameErr)
		result.DatabaseError.Response(getRequest)
		return nil, false, ServerInternalErrorString
	}
	getUserPasswordErr := dao.XfIndex.Ctx(ctx).Where("key = ?", "password").Scan(&userPasswordEntity)
	if getUserPasswordErr != nil {
		glog.Error(ctx, "[LOGIC] 获取数据库出错", getUserPasswordErr)
		result.DatabaseError.Response(getRequest)
		return nil, false, ServerInternalErrorString
	}
	// 对账号密码进行校验
	if userUsernameEntity.Value == getData.User {
		// 原密码处理
		handlingPasswords := base64.StdEncoding.EncodeToString([]byte(getData.Pass))
		// 密码校验
		if bcrypt.CompareHashAndPassword([]byte(userPasswordEntity.Value), []byte(handlingPasswords)) == nil {
			glog.Notice(ctx, "[LOGIC] 用户校验通过")
			return &userUUIDEntity.Value, true, ""
		} else {
			glog.Notice(ctx, "[LOGIC] 密码错误")
			return nil, false, "密码错误"
		}
	} else {
		glog.Notice(ctx, "[LOGIC] 用户名未找到")
		return nil, false, "用户名未找到"
	}
}

// RegisteredUserLogin 用于登记用户的登录信息。当用户完成登录操作后，该方法会将用户的 UUID 存入 token 数据表中，作为用户登录的依据。
// 在检查用户是否登录时，此数据表的内容作为登录依据。依据 index 数据表字段 key 中的 auth_limit 所对应的 value 的大小作为允许登录节点数的限制。
//
// 参数:
// ctx: 上下文对象，用于传递和控制请求的生命周期。
// userUUID: 用户的 UUID 字符串。
// remember: 用户是否选择记住登录状态的布尔值。
//
// 返回值:
// userToken: 用户的 token 信息，包含了用户的 UUID、token、IP、验证信息、User-Agent 和过期时间等信息。
// err: 如果登录登记成功，返回 nil。否则返回错误信息。
func (s *sAuthLogic) RegisteredUserLogin(
	ctx context.Context,
	userUUID string,
	remember bool,
) (userToken *entity.XfToken, err error) {
	glog.Notice(ctx, "[LOGIC] 执行 AuthLogic:RegisteredUserLogin 服务层")
	// 根据 ctx 获取 Request 信息
	getRequest := ghttp.RequestFromCtx(ctx)
	// 获取允许用户允许登录的节点数
	var getAuthLimit entity.XfIndex
	getAuthLimitErr := dao.XfIndex.Ctx(ctx).Where("key", "auth_limit").Scan(&getAuthLimit)
	if getAuthLimitErr != nil {
		glog.Error(ctx, "[LOGIC] 获取数据库出错", getAuthLimitErr)
		return nil, gerror.NewCode(gcode.CodeDbOperationError, getAuthLimitErr.Error())
	}
	// 获取用户的一登录信息
	var getAuthTokenList []entity.XfToken
	getAuthTokenErr := dao.XfToken.Ctx(ctx).Where("user_uuid", userUUID).OrderAsc("created_at").Scan(&getAuthTokenList)
	if getAuthTokenErr != nil {
		glog.Error(ctx, "[LOGIC] 获取数据库出错", getAuthTokenErr)
		return nil, gerror.NewCode(gcode.CodeDbOperationError, getAuthTokenErr.Error())
	}
	// 统计个数
	if len(getAuthTokenList) >= gconv.Int(getAuthLimit.Value) {
		glog.Warning(ctx, "[LOGIC] 用户登录节点已满")
		// 删除最早登录的用户内容
		_, deleteError := dao.XfToken.Ctx(ctx).Delete("id", getAuthTokenList[0].Id)
		if deleteError != nil {
			return nil, deleteError
		}
	}
	var remPassword int64
	if remember {
		remPassword = gtime.Timestamp() + 604800
	} else {
		remPassword = gtime.Timestamp() + 86400
	}
	// 数据库操作
	insert, insertErr := dao.XfToken.Ctx(ctx).Data(do.XfToken{
		UserUuid:     userUUID,
		UserToken:    uuid.NewV4().String(),
		UserIp:       getRequest.GetClientIp(),
		Verification: uuid.NewV4().String(),
		UserAgent:    getRequest.GetHeader("User-Agent"),
		ExpiredAt:    gtime.NewFromTimeStamp(remPassword),
	}).Insert()
	if insertErr != nil {
		glog.Error(ctx, "[LOGIC] 数据表插入失败", insertErr)
		return nil, gerror.NewCode(gcode.CodeDbOperationError, insertErr.Error())
	}
	// 获取数据
	getLastInsertID, getLastInsertIDErr := insert.LastInsertId()
	if getLastInsertIDErr != nil {
		glog.Error(ctx, "[LOGIC] 数据表查询出错", insertErr)
		return nil, gerror.NewCode(gcode.CodeDbOperationError, getLastInsertIDErr.Error())
	}
	var authToken entity.XfToken
	selectError := dao.XfToken.Ctx(ctx).Where("id", getLastInsertID).Scan(&authToken)
	if selectError != nil {
		return nil, gerror.NewCode(gcode.CodeDbOperationError, selectError.Error())
	}
	return &authToken, nil
}

// CheckUserHasConsoleUser 是一个用于检查用户是否存在于控制台用户列表中的函数。
// 它主要用于获取用户的用户名，并与数据库中的内容进行比对。如果用户名存在于数据库中，则返回 nil。否则，返回错误信息。
//
// 参数:
// ctx: 上下文对象，用于传递和控制请求的生命周期。
// username: 需要检查的用户名字符串。
//
// 返回值:
// err: 如果用户名存在于数据库中，返回 nil。否则返回错误信息。
func (s *sAuthLogic) CheckUserHasConsoleUser(ctx context.Context, username string) (err error) {
	glog.Notice(ctx, "[LOGIC] 执行 AuthLogic:CheckUserHasConsoleUser 服务层")
	// 根据 ctx 获取 Request 信息
	getRequest := ghttp.RequestFromCtx(ctx)
	// 获取用户的一登录信息
	var getUser entity.XfIndex
	getUserErr := dao.XfIndex.Ctx(ctx).Where(do.XfIndex{Key: "user"}).Scan(&getUser)
	if getUserErr != nil {
		glog.Error(ctx, "[LOGIC] 获取数据库出错", getUserErr)
		result.DatabaseError.Response(getRequest)
		return errors.New(ServerInternalErrorString)
	}
	if getUser.Value != username {
		return errors.New("用户名不匹配")
	}
	return nil
}
