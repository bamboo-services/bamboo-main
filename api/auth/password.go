/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package apiAuth

// PasswordChangeRequest 密码修改请求
type PasswordChangeRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=100" example:"oldpassword123"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=100" example:"newpassword123"`
}

// PasswordResetRequest 密码重置请求
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email" example:"admin@example.com"`
}

// ConfirmResetPasswordRequest 确认重置密码请求
type ConfirmResetPasswordRequest struct {
	Token       string `json:"token" binding:"required,min=32,max=64" example:"abc123..."`
	NewPassword string `json:"new_password" binding:"required,min=6,max=100" example:"newpassword123"`
}

// VerifyResetTokenRequest 重置令牌验证请求
type VerifyResetTokenRequest struct {
	Token string `form:"token" binding:"required,min=32,max=64" example:"abc123..."`
}
