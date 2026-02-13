package logic

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/models/dto"
	"github.com/bamboo-services/bamboo-main/pkg/constants"

	xError "github.com/bamboo-services/bamboo-base-go/error"
	xLog "github.com/bamboo-services/bamboo-base-go/log"
	xUtil "github.com/bamboo-services/bamboo-base-go/utility"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/utility/ctxutil"
	"github.com/gin-gonic/gin"
	bSdkLogic "github.com/phalanx/beacon-sso-sdk/logic"
	bSdkModels "github.com/phalanx/beacon-sso-sdk/models"
	"gorm.io/gorm"
)

func (a *AuthLogic) LoginByOAuth(ctx *gin.Context, userinfo *bSdkModels.OAuthUserinfo, accessToken string) (*dto.SystemUserDetailDTO, string, *time.Time, *time.Time, *xError.Error) {
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

	return buildSystemUserDetailDTO(user), accessToken, &now, &expiredAt, nil
}

func (a *AuthLogic) SyncOAuthUser(ctx *gin.Context, userinfo *bSdkModels.OAuthUserinfo) (*entity.SystemUser, *xError.Error) {
	if userinfo == nil || strings.TrimSpace(userinfo.Sub) == "" {
		return nil, xError.NewError(ctx, xError.ParameterError, "OAuth 用户标识无效", false)
	}

	db := xCtxUtil.MustGetDB(ctx)
	var user entity.SystemUser
	now := time.Now()

	err := db.Where("oauth_user_id = ?", userinfo.Sub).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询 OAuth 用户失败", false, err)
	}

	if err == nil {
		if user.Status == constants.StatusInactive {
			return nil, xError.NewError(ctx, xError.Forbidden, "用户已被禁用", false)
		}
		if xErr := a.updateOAuthProfile(ctx, db, &user, userinfo, &now); xErr != nil {
			return nil, xErr
		}
		return &user, nil
	}

	oauthEmail := strings.TrimSpace(userinfo.Email)
	if oauthEmail != "" {
		err = db.Where("email = ?", oauthEmail).First(&user).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xError.NewError(ctx, xError.DatabaseError, "按邮箱查询用户失败", false, err)
		}
		if err == nil {
			if user.Status == constants.StatusInactive {
				return nil, xError.NewError(ctx, xError.Forbidden, "用户已被禁用", false)
			}
			if user.OAuthUserID != nil && *user.OAuthUserID != "" && *user.OAuthUserID != userinfo.Sub {
				return nil, xError.NewError(ctx, xError.Unauthorized, "该本地账户已绑定其他 OAuth 用户", false)
			}
			if xErr := a.updateOAuthProfile(ctx, db, &user, userinfo, &now); xErr != nil {
				return nil, xErr
			}
			return &user, nil
		}
	}

	passwordRaw := fmt.Sprintf("oauth:%s:%d", userinfo.Sub, now.UnixNano())
	hashedPassword, hashErr := xUtil.EncryptPasswordString(passwordRaw)
	if hashErr != nil {
		return nil, xError.NewError(ctx, xError.ServerInternalError, "初始化本地账户密码失败", false, hashErr)
	}

	username, xErr := a.generateUniqueOAuthUsername(ctx, db, userinfo)
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

	if err = db.Create(&newUser).Error; err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "创建本地映射用户失败", false, err)
	}

	return &newUser, nil
}

func (a *AuthLogic) updateOAuthProfile(ctx *gin.Context, db *gorm.DB, user *entity.SystemUser, userinfo *bSdkModels.OAuthUserinfo, now *time.Time) *xError.Error {
	updates := map[string]any{
		"oauth_user_id": userinfo.Sub,
		"last_login_at": now,
	}

	if oauthEmail := strings.TrimSpace(userinfo.Email); oauthEmail != "" && oauthEmail != user.Email {
		exists, xErr := a.emailExists(ctx, db, oauthEmail, user.ID)
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

	if err := db.Model(user).Updates(updates).Error; err != nil {
		return xError.NewError(ctx, xError.DatabaseError, "更新本地映射用户失败", false, err)
	}

	if err := db.Where("id = ?", user.ID).First(user).Error; err != nil {
		return xError.NewError(ctx, xError.DatabaseError, "刷新本地映射用户失败", false, err)
	}

	return nil
}

func (a *AuthLogic) emailExists(ctx *gin.Context, db *gorm.DB, email string, exceptID int64) (bool, *xError.Error) {
	var count int64
	err := db.Model(&entity.SystemUser{}).Where("email = ? AND id <> ?", email, exceptID).Count(&count).Error
	if err != nil {
		return false, xError.NewError(ctx, xError.DatabaseError, "检查邮箱唯一性失败", false, err)
	}
	return count > 0, nil
}

func (a *AuthLogic) generateUniqueOAuthUsername(ctx *gin.Context, db *gorm.DB, userinfo *bSdkModels.OAuthUserinfo) (string, *xError.Error) {
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

		var count int64
		err := db.Model(&entity.SystemUser{}).Where("username = ?", candidate).Count(&count).Error
		if err != nil {
			return "", xError.NewError(ctx, xError.DatabaseError, "检查用户名唯一性失败", false, err)
		}
		if count == 0 {
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

func buildSystemUserDetailDTO(user *entity.SystemUser) *dto.SystemUserDetailDTO {
	return &dto.SystemUserDetailDTO{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Role:        user.Role,
		Status:      user.Status,
		EmailVerify: user.EmailVerify,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
