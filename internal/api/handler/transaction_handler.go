package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pragadeesh-c/pismo-tech-case/internal/api/models"
	"github.com/pragadeesh-c/pismo-tech-case/internal/service"
	"github.com/rs/zerolog/log"
)

type TransactionHandler struct {
	service *service.TransactionService
}

func NewTransactionHandler(service *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// CreateTransaction 	godoc
// @Summary 			Create Transaction
// @Description 		Creates a new transaction for an account id
// @Tags 				transactions
// @Accept 				json
// @Produce 			json
// @Param 				request body models.CreateTransactionRequest 					true "Transaction input"
// @Success 			201 {object} models.SuccessResponse{data=service.Transaction}
// @Failure 			400 {object} models.ErrResponse 							"BAD_REQUEST"
// @Failure 			404 {object} models.ErrResponse 							"ACCOUNT_NOT_FOUND"
// @Failure 			500 {object} models.ErrResponse 							"INTERNAL_ERROR"
// @Router 				/transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var request models.CreateTransactionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Err(err).
			Msg("error binding request json to struct")
		ErrorResponse(c, http.StatusBadRequest, models.ErrCodeBadRequest, "invalid request body")
		return
	}

	transaction, err := h.service.Create(c.Request.Context(), service.CreateTransaction{
		AccountID:     request.AccountID,
		OperationType: request.OperationType,
		Amount:        request.Amount,
	})

	if err != nil {
		log.Err(err).
			Int("account_id", request.AccountID).
			Int("operation_type", request.OperationType).
			Msg("create transaction failed")
		switch {
		case errors.Is(err, service.ErrAccountNotFound):
			ErrorResponse(c, http.StatusNotFound, models.ErrCodeAccountNotFound, "account not found")
		case errors.Is(err, service.ErrInvalidOperationType):
			ErrorResponse(c, http.StatusBadRequest, models.ErrCodeInvalidOperationType, "invalid operation type")
		case errors.Is(err, service.ErrInvalidAmount):
			ErrorResponse(c, http.StatusBadRequest, models.ErrCodeInvalidAmount, "invalid amount")
		default:
			ErrorResponse(c, http.StatusInternalServerError, models.ErrCodeInternalError, "error occurred while creating transaction")
		}
		return
	}

	SuccessResponse(c, http.StatusCreated, "transaction created successfully", transaction)
}
