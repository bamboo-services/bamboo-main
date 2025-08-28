package service

import (
	"bamboo-main/internal/logic"
	"bamboo-main/pkg/startup"
)

// LinkFriendService 友情链接服务类型别名
type LinkFriendService = logic.LinkFriendLogic

// NewLinkFriendService 创建友情链接服务实例
func NewLinkFriendService(reg *startup.Reg) *LinkFriendService {
	return logic.NewLinkFriendLogic(reg)
}