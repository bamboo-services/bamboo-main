package constants

const (
	// 上下文键
	ContextKeyUser   = "user"
	ContextKeyUserID = "user_id"
	ContextKeyToken  = "token"
	
	// HTTP 头部
	HeaderAuthorization = "Authorization"
	HeaderContentType   = "Content-Type"
	
	// 认证相关
	TokenPrefix = "Bearer "
	
	// 系统状态
	StatusActive   = 1
	StatusInactive = 0
	
	// 链接状态
	LinkStatusPending  = 0 // 待审核
	LinkStatusApproved = 1 // 已通过
	LinkStatusRejected = 2 // 已拒绝
	
	// 链接失效状态
	LinkFailNormal = 0 // 正常
	LinkFailBroken = 1 // 失效
	
	// 邮件类型
	EmailTypeApply        = "apply"
	EmailTypeApproved     = "approved"
	EmailTypeRejected     = "rejected"
	EmailTypePasswordReset = "password_reset"
)