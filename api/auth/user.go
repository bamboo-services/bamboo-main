package apiAuth

import "github.com/bamboo-services/bamboo-main/internal/entity"

type UserInfoResponse struct {
	User entity.SystemUser `json:"user"`
}
