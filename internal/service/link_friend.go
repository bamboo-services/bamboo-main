package service

import (
	"bamboo-main/internal/logic"
)

// LinkFriendService 友情链接服务类型别名
type LinkFriendService = logic.LinkFriendLogic

// NewLinkFriendService 创建友情链接服务实例
func NewLinkFriendService() *LinkFriendService {
	return logic.NewLinkFriendLogic()
}