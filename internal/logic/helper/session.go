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

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xCache "github.com/bamboo-services/bamboo-base-go/major/cache"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	redisModel "github.com/bamboo-services/bamboo-main/internal/models/cache/redis"
	"github.com/bamboo-services/bamboo-main/internal/repository"
)

// SessionTTL 默认登录会话有效期
const SessionTTL = 24 * time.Hour

// SessionMeta 会话元数据
//
// 承载来自传输层（HTTP 请求）的客户端信息。由 handler 提取后以纯值下传，
// 避免 logic/repository 层耦合 gin.Context。
type SessionMeta struct {
	ClientIP  string
	UserAgent string
}

// SessionLogic 会话管理服务实现
//
// 负责组装领域会话对象（TokenSession）并交由 SessionRepo 持久化，
// 自身不再直接接触 Redis 命令。
type SessionLogic struct {
	repo *repository.SessionRepo
}

// NewSessionLogic 创建 SessionLogic 实例，基于缓存管理器初始化会话仓储依赖。
func NewSessionLogic(m *xCache.Manager) *SessionLogic {
	return &SessionLogic{
		repo: repository.NewSessionRepo(m),
	}
}

// CreateUserSession 创建用户会话
//
// ttl 为会话有效期；会话对象的 ExpiredAt 与 Redis TTL 均以此为准。
func (s *SessionLogic) CreateUserSession(ctx context.Context, user *entity.SystemUser, token string, meta SessionMeta, ttl time.Duration) *xError.Error {
	now := time.Now()
	tokenSession := &redisModel.TokenSession{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		LoginIP:   meta.ClientIP,
		UserAgent: meta.UserAgent,
		CreatedAt: now,
		ExpiredAt: now.Add(ttl),
	}

	return s.repo.SaveSession(ctx, token, tokenSession, ttl)
}

// GetUserSession 读取用户会话；会话不存在时 found=false（非错误）
func (s *SessionLogic) GetUserSession(ctx context.Context, token string) (*redisModel.TokenSession, bool, *xError.Error) {
	return s.repo.GetSession(ctx, token)
}

// DeleteUserSession 删除用户会话
func (s *SessionLogic) DeleteUserSession(ctx context.Context, token string) *xError.Error {
	return s.repo.DeleteSession(ctx, token)
}
