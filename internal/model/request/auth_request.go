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

package request

// AuthLoginReq 登录请求
type AuthLoginReq struct {
	Username string `json:"username" binding:"required,min=1,max=50" example:"admin"`
	Password string `json:"password" binding:"required,min=6,max=100" example:"password123"`
}

// AuthRegisterReq 注册请求
type AuthRegisterReq struct {
	Username string  `json:"username" binding:"required,min=1,max=50" example:"admin"`
	Email    string  `json:"email" binding:"required,email" example:"admin@example.com"`
	Nickname *string `json:"nickname" binding:"omitempty,min=1,max=50" example:"筱锋"`
	Password string  `json:"password" binding:"required,min=6,max=100" example:"password123"`
}

// AuthPasswordChangeReq 修改密码请求
type AuthPasswordChangeReq struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=100" example:"oldpassword123"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=100" example:"newpassword123"`
}

// AuthPasswordResetReq 重置密码请求
type AuthPasswordResetReq struct {
	Email string `json:"email" binding:"required,email" example:"admin@example.com"`
}

// AuthVerifyEmailReq 验证邮箱请求
type AuthVerifyEmailReq struct {
	Token string `form:"token" binding:"required,min=32,max=64" example:"abc123..."`
}

// AuthConfirmResetPasswordReq 确认重置密码请求
type AuthConfirmResetPasswordReq struct {
	Token       string `json:"token" binding:"required,min=32,max=64" example:"abc123..."`
	NewPassword string `json:"new_password" binding:"required,min=6,max=100" example:"newpassword123"`
}

// AuthVerifyResetTokenReq 验证重置Token请求
type AuthVerifyResetTokenReq struct {
	Token string `form:"token" binding:"required,min=32,max=64" example:"abc123..."`
}
