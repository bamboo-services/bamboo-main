/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2025 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package v1

import (
	"github.com/XiaoLFeng/bamboo-utils/bmodels"
	"github.com/gogf/gf/v2/frame/g"
	"go/types"
)

type AuthPasswordResetReq struct {
	g.Meta      `path:"/auth/pass/reset" method:"PUT" tags:"认证控制器" summary:"重置密码" dc:"重置密码接口，允许用户通过邮箱或手机号重置其登录密码"`
	Email       string `json:"email" v:"required|email#请输入邮箱地址|请输入有效的邮箱地址" dc:"邮箱地址，用于接收重置密码的链接或验证码"`
	VerifyKey   string `json:"verify_key" in:"query" v:"regex:^(|[0-9A-Za-z]+)$#验证码格式不正确" dc:"验证码，用于验证用户身份，通常通过邮箱或短信发送"`
	NewPassword string `json:"new_password" in:"query" dc:"新密码，用户希望设置的新登录密码"`
}

type AuthPasswordResetRes struct {
	g.Meta `mime:"application/json; charset=utf-8" dc:"重置密码响应"`
	*bmodels.ResponseDTO[types.Nil]
}
