package models

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number" binding:"required"`
}

type CreateTransactionRequest struct {
	AccountID     int     `json:"account_id" binding:"required" example:"1"`
	OperationType int     `json:"operation_type_id" binding:"required" example:"2"`
	Amount        float64 `json:"amount" binding:"required" example:"120.34"`
}
