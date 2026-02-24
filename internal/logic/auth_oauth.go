package logic

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/pkg/constants"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xUtil "github.com/bamboo-services/bamboo-base-go/common/utility"
	"github.com/gin-gonic/gin"
	bSdkLogic "github.com/phalanx-labs/beacon-sso-sdk/logic"
	bSdkModels "github.com/phalanx-labs/beacon-sso-sdk/models"
)

func (a *AuthLogic) LoginByOAuth(ctx *gin.Context, userinfo *bSdkModels.OAuthUserinfo, accessToken string) (*entity.SystemUser, string, *time.Time, *time.Time, *xError.Error) {
	if accessToken == "" {
		return nil, "", nil, nil, xError.NewError(ctx, xError.ParameterEmpty, "访问令牌不能为空", false)
	}

	user, xErr := a.SyncOAuthUser(ctx, userinfo)
	if xErr != nil {
		return nil, "", nil, nil, xErr
	}

	now := time.Now()
	expiredAt := now.Add(24 * time.Hour)

	oauthLogic := bSdkLogic.NewBusiness(ctx)
	introspection, introspectionErr := oauthLogic.Introspection(ctx, "access_token", accessToken)
	if introspectionErr != nil {
		xLog.WithName(xLog.NamedLOGC, "AUTH").Warn(ctx, fmt.Sprintf("OAuth introspection 获取失败: %v", introspectionErr))
	} else if introspection != nil && introspection.Exp > 0 {
		expiredAt = time.Unix(introspection.Exp, 0)
	}

	return user, accessToken, &now, &expiredAt, nil
}

func (a *AuthLogic) SyncOAuthUser(ctx *gin.Context, userinfo *bSdkModels.OAuthUserinfo) (*entity.SystemUser, *xError.Error) {
	if userinfo == nil || strings.TrimSpace(userinfo.Sub) == "" {
		return nil, xError.NewError(ctx, xError.ParameterError, "OAuth 用户标识无效", false)
	}

	var user *entity.SystemUser
	now := time.Now()

	user, found, xErr := a.repo.user.GetByOAuthUserID(ctx, userinfo.Sub)
	if xErr != nil {
		return nil, xErr
	}

	if found {
		if user.Status == constants.StatusInactive {
			return nil, xError.NewError(ctx, xError.Forbidden, "用户已被禁用", false)
		}
		if xErr := a.updateOAuthProfile(ctx, user, userinfo, &now); xErr != nil {
			return nil, xErr
		}
		return user, nil
	}

	oauthEmail := strings.TrimSpace(userinfo.Email)
	if oauthEmail != "" {
		user, found, xErr = a.repo.user.GetByEmail(ctx, oauthEmail)
		if xErr != nil {
			return nil, xErr
		}
		if found {
			if user.Status == constants.StatusInactive {
				return nil, xError.NewError(ctx, xError.Forbidden, "用户已被禁用", false)
			}
			if user.OAuthUserID != nil && *user.OAuthUserID != "" && *user.OAuthUserID != userinfo.Sub {
				return nil, xError.NewError(ctx, xError.Unauthorized, "该本地账户已绑定其他 OAuth 用户", false)
			}
			if xErr := a.updateOAuthProfile(ctx, user, userinfo, &now); xErr != nil {
				return nil, xErr
			}
			return user, nil
		}
	}

	passwordRaw := fmt.Sprintf("oauth:%s:%d", userinfo.Sub, now.UnixNano())
	hashedPassword, hashErr := xUtil.Password().EncryptString(passwordRaw)
	if hashErr != nil {
		return nil, xError.NewError(ctx, xError.ServerInternalError, "初始化本地账户密码失败", false, hashErr)
	}

	username, xErr := a.generateUniqueOAuthUsername(ctx, userinfo)
	if xErr != nil {
		return nil, xErr
	}

	oauthID := strings.TrimSpace(userinfo.Sub)
	newUser := entity.SystemUser{
		OAuthUserID: &oauthID,
		Username:    username,
		Password:    hashedPassword,
		Email:       a.resolveOAuthEmail(userinfo),
		Nickname:    a.resolveOAuthNickname(userinfo),
		Avatar:      a.resolveOAuthAvatar(userinfo),
		Role:        constants.RoleUser,
		Status:      constants.StatusActive,
		EmailVerify: strings.TrimSpace(userinfo.Email) != "",
		LastLoginAt: &now,
	}

	user, xErr = a.repo.user.Create(ctx, &newUser)
	if xErr != nil {
		return nil, xErr
	}

	return user, nil
}

func (a *AuthLogic) updateOAuthProfile(ctx *gin.Context, user *entity.SystemUser, userinfo *bSdkModels.OAuthUserinfo, now *time.Time) *xError.Error {
	updates := map[string]any{
		"oauth_user_id": userinfo.Sub,
		"last_login_at": now,
	}

	if oauthEmail := strings.TrimSpace(userinfo.Email); oauthEmail != "" && oauthEmail != user.Email {
		exists, xErr := a.repo.user.ExistsByEmailExceptID(ctx, oauthEmail, user.ID)
		if xErr != nil {
			return xErr
		}
		if !exists {
			updates["email"] = oauthEmail
			updates["email_verify"] = true
		}
	}

	nickname := a.resolveOAuthNickname(userinfo)
	if nickname != nil && (user.Nickname == nil || *user.Nickname != *nickname) {
		updates["nickname"] = nickname
	}

	avatar := a.resolveOAuthAvatar(userinfo)
	if avatar != nil && (user.Avatar == nil || *user.Avatar != *avatar) {
		updates["avatar"] = avatar
	}

	updatedUser, xErr := a.repo.user.UpdateFieldsByID(ctx, user.ID, updates)
	if xErr != nil {
		return xErr
	}
	*user = *updatedUser

	return nil
}

func (a *AuthLogic) generateUniqueOAuthUsername(ctx *gin.Context, userinfo *bSdkModels.OAuthUserinfo) (string, *xError.Error) {
	base := sanitizeOAuthIdentifier(userinfo.PreferredUsername)
	if base == "" {
		base = sanitizeOAuthIdentifier(userinfo.Nickname)
	}
	if base == "" {
		base = sanitizeOAuthIdentifier(userinfo.Sub)
	}
	if base == "" {
		base = "oauth_user"
	}

	if len(base) > 42 {
		base = base[:42]
	}

	for i := 0; i < 100; i++ {
		candidate := base
		if i > 0 {
			suffix := fmt.Sprintf("_%d", i)
			if len(candidate)+len(suffix) > 50 {
				candidate = candidate[:50-len(suffix)]
			}
			candidate += suffix
		}

		exists, xErr := a.repo.user.ExistsByUsername(ctx, candidate)
		if xErr != nil {
			return "", xErr
		}
		if !exists {
			return candidate, nil
		}
	}

	return "", xError.NewError(ctx, xError.ServerInternalError, "生成唯一用户名失败", false)
}

func (a *AuthLogic) resolveOAuthEmail(userinfo *bSdkModels.OAuthUserinfo) string {
	oauthEmail := strings.TrimSpace(userinfo.Email)
	if oauthEmail != "" {
		return oauthEmail
	}

	identifier := sanitizeOAuthIdentifier(userinfo.Sub)
	if identifier == "" {
		identifier = fmt.Sprintf("oauth_%d", time.Now().UnixNano())
	}
	if len(identifier) > 60 {
		identifier = identifier[:60]
	}
	return fmt.Sprintf("%s@oauth.local", identifier)
}

func (a *AuthLogic) resolveOAuthNickname(userinfo *bSdkModels.OAuthUserinfo) *string {
	nickname := strings.TrimSpace(userinfo.Nickname)
	if nickname == "" {
		nickname = strings.TrimSpace(userinfo.PreferredUsername)
	}
	if nickname == "" {
		return nil
	}
	if len(nickname) > 100 {
		nickname = nickname[:100]
	}
	return &nickname
}

func (a *AuthLogic) resolveOAuthAvatar(userinfo *bSdkModels.OAuthUserinfo) *string {
	if userinfo == nil || userinfo.Raw == nil {
		return nil
	}

	candidates := []string{"avatar", "avatar_url", "picture"}
	for _, key := range candidates {
		if value, exists := userinfo.Raw[key]; exists {
			url, ok := value.(string)
			if ok {
				url = strings.TrimSpace(url)
				if url != "" {
					if len(url) > 500 {
						url = url[:500]
					}
					return &url
				}
			}
		}
	}

	return nil
}

func sanitizeOAuthIdentifier(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	var builder strings.Builder
	builder.Grow(len(value))

	for _, r := range value {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r):
			builder.WriteRune(unicode.ToLower(r))
		case r == '-', r == '_', r == '.':
			builder.WriteRune(r)
		case unicode.IsSpace(r):
			builder.WriteRune('_')
		default:
			builder.WriteRune('_')
		}
	}

	result := strings.Trim(builder.String(), "_.-")
	return result
}
