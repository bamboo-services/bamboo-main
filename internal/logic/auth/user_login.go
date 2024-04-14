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
 *
 */

package auth

import (
	"context"
	"encoding/base64"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	v1 "xiaoMain/api/auth/v1"
	"xiaoMain/internal/dao"
	"xiaoMain/internal/model/entity"
	"xiaoMain/internal/service"
	"xiaoMain/utility"
	"xiaoMain/utility/result"
)

type sAuthLogic struct {
}

func init() {
	service.RegisterAuthLogic(New())
}

func New() *sAuthLogic {
	return &sAuthLogic{}
}

// IsUserLogin
// 检查用户是否已经登录，若用户已经完成登录的相关操作，若用户的此次登录有效则返回 true 否则返回 false
func (aL *sAuthLogic) IsUserLogin(ctx context.Context) bool {
	glog.Info(ctx, "[LOGIC] 执行 AuthLogic:IsUserLogin 服务层")
	// 根据 ctx 获取 Request 信息
	getRequest := ghttp.RequestFromCtx(ctx)
	// 获取用户的 UUID(UID) 以及 认证密钥
	getUserUUID := getRequest.Header.Get("X-User-Uid")
	getUserAuthorize, err := utility.TokenLeftBearer(getRequest.Header.Get("Authorization"))
	if err != nil {
		return false
	}
	// 对内容进行校验
	if getUserUUID != "" && getUserAuthorize != "" {
		var getTokenDO entity.XfToken
		err := dao.XfToken.Ctx(ctx).
			Where("user_uuid", getUserUUID).
			Where("user_token", getUserAuthorize).
			Limit(1).Scan(&getTokenDO)
		if err != nil {
			glog.Error(ctx, "[LOGIC] 获取数据库出错", err)
			return false
		}
		// 检查是否过期
		if gtime.Timestamp() < getTokenDO.ExpiredAt.Timestamp() {
			// 验证登录有效
			if getUserAuthorize == getTokenDO.UserToken {
				glog.Infof(ctx, "[LOGIC] 用户UID %s 任然登录状态", getTokenDO.UserUuid)
				return true
			} else {
				glog.Warning(ctx, "[LOGIC] 用户登录已失效")
				return false
			}
		}
	}
	glog.Warning(ctx, "[LOGIC] 用户未登录")
	return false
}

// CheckUserLogin
// 对用户的登录进行检查。主要用于对用户输入的信息与数据库的内容进行校验，当用户名与用户校验通过后 isCorrect 返回正确值，否则返回错误的内容
// 并且当用户正常登录后，将会返回用户的 UUID 作为下一步的登录操作
func (aL *sAuthLogic) CheckUserLogin(ctx context.Context, getData *v1.UserLoginReq) (userUUID *string, isCorrect bool) {
	glog.Info(ctx, "[LOGIC] 执行 AuthLogic:CheckUserLogin 服务层")
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
		return nil, false
	}
	getUserUsernameErr := dao.XfIndex.Ctx(ctx).Where("key = ?", "user").Scan(&userUsernameEntity)
	if getUserUsernameErr != nil {
		glog.Error(ctx, "[LOGIC] 获取数据库出错", getUserUsernameErr)
		result.DatabaseError.Response(getRequest)
		return nil, false
	}
	getUserPasswordErr := dao.XfIndex.Ctx(ctx).Where("key = ?", "password").Scan(&userPasswordEntity)
	if getUserPasswordErr != nil {
		glog.Error(ctx, "[LOGIC] 获取数据库出错", getUserPasswordErr)
		result.DatabaseError.Response(getRequest)
		return nil, false
	}
	// 对账号密码进行校验
	if userUsernameEntity.Value == getData.User {
		// 原密码处理
		handlingPasswords := base64.StdEncoding.EncodeToString([]byte(getData.Pass))
		// 密码校验
		if bcrypt.CompareHashAndPassword([]byte(userPasswordEntity.Value), []byte(handlingPasswords)) == nil {
			glog.Info(ctx, "[LOGIC] 用户校验通过")
			return &userUUIDEntity.Value, true
		} else {
			glog.Info(ctx, "[LOGIC] 密码错误")
			return nil, false
		}
	} else {
		glog.Info(ctx, "[LOGIC] 用户名未找到")
		return nil, false
	}
}

// RegisteredUserLogin
// 对用户的登录内容进行登记，将用户的 UUID 传入后存入 token 数据表中，作为用户登录的登录依据。在检查用户是否登录时候，此数据表的内容作为登录
// 依据。
//
// 依据 index 数据表字段 key 中的 auth_limit 所对应的 value 的大小作为允许登录节点数的限制
func (aL *sAuthLogic) RegisteredUserLogin(
	ctx context.Context,
	userUUID string,
	remember bool,
) (userToken *entity.XfToken, err error) {
	glog.Info(ctx, "[LOGIC] 执行 AuthLogic:RegisteredUserLogin 服务层")
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
	// 直接插入登录信息
	newToken := g.Map{
		"UserUuid":     userUUID,
		"UserToken":    uuid.NewV4().String(),
		"UserIp":       getRequest.GetClientIp(),
		"verification": uuid.NewV4(),
		"UserAgent":    getRequest.GetHeader("User-Agent"),
		"ExpiredAt":    gtime.NewFromTimeStamp(remPassword),
	}
	// 数据库操作
	insert, insertErr := dao.XfToken.Ctx(ctx).Insert(&newToken)
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
