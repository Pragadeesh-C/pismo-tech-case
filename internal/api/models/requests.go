package models

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number" binding:"required"`
}
