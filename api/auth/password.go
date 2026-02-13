package apiAuth

type PasswordChangeRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=100" example:"oldpassword123"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=100" example:"newpassword123"`
}

type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email" example:"admin@example.com"`
}

type ConfirmResetPasswordRequest struct {
	Token       string `json:"token" binding:"required,min=32,max=64" example:"abc123..."`
	NewPassword string `json:"new_password" binding:"required,min=6,max=100" example:"newpassword123"`
}

type VerifyResetTokenRequest struct {
	Token string `form:"token" binding:"required,min=32,max=64" example:"abc123..."`
}

type PasswordChangeResponse = MessageResponse
type PasswordResetResponse = MessageResponse
