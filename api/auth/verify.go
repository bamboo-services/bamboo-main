package apiAuth

type VerifyEmailRequest struct {
	Token string `form:"token" binding:"required,min=32,max=64" example:"abc123..."`
}
