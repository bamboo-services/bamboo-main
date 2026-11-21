// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

package logcHelper

import (
	"context"
	"time"

	"github.com/bamboo-services/bamboo-base-go/foundation/v2/apperr"
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/logging"
	"github.com/bamboo-services/bamboo-base-go/plugins/cache/v2"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	redisModel "github.com/bamboo-services/bamboo-main/internal/models/cache/redis"
	"github.com/bamboo-services/bamboo-main/internal/repository"
)

// SessionTTL 默认登录会话有效期
const SessionTTL = 24 * time.Hour

// SessionMeta 会话元数据
type SessionMeta struct {
	ClientIP  string
	UserAgent string
}

// SessionLogic 会话管理服务实现
type SessionLogic struct {
	repo *repository.SessionRepo
}

// NewSessionLogic 创建 SessionLogic 实例，基于缓存管理器初始化会话仓储依赖。
func NewSessionLogic(m *cache.Manager, logger ...*logging.Logger) *SessionLogic {
	return &SessionLogic{
		repo: repository.NewSessionRepo(m, logger...),
	}
}

// CreateUserSession 创建用户会话
func (s *SessionLogic) CreateUserSession(ctx context.Context, user *entity.SystemUser, token string, meta SessionMeta, ttl time.Duration) *apperr.Error {
	now := time.Now()
	tokenSession := &redisModel.TokenSession{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		LoginIP:   meta.ClientIP,
		UserAgent: meta.UserAgent,
		CreatedAt: now,
		ExpiredAt: now.Add(ttl),
	}

	return s.repo.SaveSession(ctx, token, tokenSession, ttl)
}

// GetUserSession 读取用户会话；会话不存在时 found=false（非错误）
func (s *SessionLogic) GetUserSession(ctx context.Context, token string) (*redisModel.TokenSession, bool, *apperr.Error) {
	return s.repo.GetSession(ctx, token)
}

// DeleteUserSession 删除用户会话
func (s *SessionLogic) DeleteUserSession(ctx context.Context, token string) *apperr.Error {
	return s.repo.DeleteSession(ctx, token)
}
