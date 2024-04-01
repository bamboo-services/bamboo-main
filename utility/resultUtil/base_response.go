package resultUtil

type BaseResponse struct {
	Output  string      `json:"output"`
	Code    uint16      `json:"code"`
	Message string      `json:"message"`
	Time    string      `json:"time"`
	Data    interface{} `json:"data"`
}

type ErrorBaseResponse struct {
	Output  string    `json:"output"`
	Code    uint16    `json:"code"`
	Message string    `json:"message"`
	Time    string    `json:"time"`
	Data    errorData `json:"data"`
}

type errorData struct {
	ErrorMessage string      `json:"errorMessage"`
	ErrorData    interface{} `json:"errorData"`
}
