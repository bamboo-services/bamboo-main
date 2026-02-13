package apiAuth

type MessageResponse struct {
	Message string `json:"message"`
}

type LogoutResponse = MessageResponse
