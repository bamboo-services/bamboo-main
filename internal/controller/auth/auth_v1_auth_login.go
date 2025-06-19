package auth

import (
	"bamboo-main/internal/model/response"
	"bamboo-main/internal/service"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/berror"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/XiaoLFeng/bamboo-utils/bresult"
	"regexp"

	"bamboo-main/api/auth/v1"
)

// AuthLogin
//
// 用户登录，需要用户提供用户名和密码。
func (c *ControllerV1) AuthLogin(ctx context.Context, req *v1.AuthLoginReq) (res *v1.AuthLoginRes, err error) {
	blog.ControllerInfo(ctx, "AuthLogin", "用户 %s 登录", req.User)

	// 是否通过验证
	var isVerify = false

	// 验证用户信息
	iAuth := service.Auth()

	// 尝试手机号登录
	matched, regexpErr := regexp.MatchString("^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\\d{8}$", req.User)
	if regexpErr != nil {
		blog.ControllerError(ctx, "AuthLogin", "手机号正则表达式格式错误: %v", regexpErr)
		return nil, berror.ErrorAddData(&berror.ErrInternalServer, "正则表达式格式错误")
	}
	if matched {
		errorCode := iAuth.VerifyUserByPhone(ctx, req.User, req.Password)
		if errorCode != nil {
			return nil, errorCode
		}
		isVerify = true
	}

	// 尝试邮箱登录
	if !isVerify {
		matched, regexpErr = regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, req.User)
		if regexpErr != nil {
			blog.ControllerError(ctx, "AuthLogin", "邮箱正则表达式格式错误: %v", regexpErr)
			return nil, berror.ErrorAddData(&berror.ErrInternalServer, "正则表达式格式错误")
		}
		if matched {
			errorCode := iAuth.VerifyUserByEmail(ctx, req.User, req.Password)
			if errorCode != nil {
				return nil, errorCode
			}
			isVerify = true
		}
	}

	// 尝试用户名登录
	if !isVerify {
		errorCode := iAuth.VerifyUserByUsername(ctx, req.User, req.Password)
		if errorCode != nil {
			return nil, errorCode
		}
	}

	// 获取用户信息
	iUser := service.User()
	getUserDTO, errorCode := iUser.GetUserSimple(ctx)
	if errorCode != nil {
		blog.ControllerError(ctx, "AuthLogin", "获取用户信息失败: %v", errorCode)
		return nil, errorCode
	}
	// 登录成功，生成用户令牌
	iToken := service.Token()
	tokenDAO, errorCode := iToken.GenerateUserToken(ctx)
	if errorCode != nil {
		return nil, errorCode
	}

	return &v1.AuthLoginRes{
		ResponseDTO: bresult.SuccessHasData(ctx, "用户登录成功", &response.AuthLoginResponse{
			User:  getUserDTO,
			Token: tokenDAO,
		}),
	}, nil
}
