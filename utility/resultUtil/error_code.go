package resultUtil

type ErrorCode struct {
	output  string
	code    uint16
	message string
}

var (
	ServerInternalError = ErrorCode{output: "ServerInternalError", code: 50000, message: "服务器内部错误"}
)
