package apiAuth

import (
	"time"

	"github.com/bamboo-services/bamboo-main/internal/entity"
)

type RegisterRequest struct {
	Username string  `json:"username" binding:"required,min=1,max=50" example:"admin"`
	Email    string  `json:"email" binding:"required,email" example:"admin@example.com"`
	Nickname *string `json:"nickname" binding:"omitempty,min=1,max=50" example:"筱锋"`
	Password string  `json:"password" binding:"required,min=6,max=100" example:"password123"`
}

type RegisterResponse struct {
	User      entity.SystemUser `json:"user"`
	Token     string            `json:"token"`
	CreatedAt time.Time         `json:"created_at"`
	ExpiredAt time.Time         `json:"expired_at"`
}
