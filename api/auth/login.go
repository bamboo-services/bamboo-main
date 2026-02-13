package apiAuth

import (
	"time"

	"github.com/bamboo-services/bamboo-main/internal/models/dto"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=1,max=50" example:"admin"`
	Password string `json:"password" binding:"required,min=6,max=100" example:"password123"`
}

type LoginResponse struct {
	User      dto.SystemUserDetailDTO `json:"user"`
	Token     string                  `json:"token"`
	CreatedAt time.Time               `json:"created_at"`
	ExpiredAt time.Time               `json:"expired_at"`
}
