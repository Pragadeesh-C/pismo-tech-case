package models

type SuccessResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message,omitempty" example:"account created"`
	Data    any    `json:"data,omitempty"`
}

type ErrResponse struct {
	Status bool       `json:"status" example:"false"`
	Error  *ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
