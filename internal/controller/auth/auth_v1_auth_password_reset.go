package auth

import (
	"bamboo-main/api/auth/v1"
	"bamboo-main/internal/service"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/XiaoLFeng/bamboo-utils/bresult"
)

// AuthPasswordReset
//
// 用户重置密码，需要用户提供新密码和验证码（可选）。
// 如果没有提供验证码，则发送验证码到用户邮箱或手机；
// 如果提供了验证码，则验证验证码并重置密码。
// 如果重置成功，则返回成功响应。
func (c *ControllerV1) AuthPasswordReset(ctx context.Context, req *v1.AuthPasswordResetReq) (res *v1.AuthPasswordResetRes, err error) {
	blog.ControllerInfo(ctx, "AuthPasswordReset", "用户重置密码")

	// 检查是否是验证码验证
	if req.VerifyKey == "" {
		/* 发送验证码 */
		iEmail := service.Email()
		errorCode := iEmail.SendMailByPasswordReset(ctx, req.Email)
		if errorCode != nil {
			return nil, errorCode
		}
		return &v1.AuthPasswordResetRes{
			ResponseDTO: bresult.Success(ctx, "验证码已发送到您的邮箱，请查收"),
		}, nil
	} else {
		/* 验证验证码并重置密码 */
		iAuth := service.Auth()
		iCode := service.Code()

		// 验证验证码
		errorCode := iCode.VerifyEmailCode(ctx, req.Email, req.VerifyKey)
		if errorCode != nil {
			return nil, errorCode
		}

		// 重置密码
		errorCode = iAuth.ChangePassword(ctx, req.NewPassword)
		if errorCode != nil {
			return nil, errorCode
		}

		// 删除所有设备的登录令牌
		iToken := service.Token()
		errorCode = iToken.RemoveUserAllToken(ctx)
		if errorCode != nil {
			return nil, errorCode
		}

		return &v1.AuthPasswordResetRes{
			ResponseDTO: bresult.Success(ctx, "用户密码重置成功"),
		}, nil
	}
}
