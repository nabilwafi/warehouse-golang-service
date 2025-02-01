package web

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
	Paging  Paging `json:"paging,omitempty"`
}

type Paging struct {
	Size      int `json:"size,omitempty"`
	Page      int `json:"page,omitempty"`
	TotalData int `json:"total_data,omitempty"`
}

func SuccessResponse(status int, message string, payload any) *Response {
	return &Response{
		Status:  status,
		Message: message,
		Data:    payload,
	}
}

func ErrorResponse(status int, message string, error error) *Response {
	return &Response{
		Status:  status,
		Message: message,
		Error:   error.Error(),
	}
}
