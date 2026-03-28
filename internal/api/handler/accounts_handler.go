package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pragadeesh-c/pismo-tech-case/internal/api/models"
	"github.com/pragadeesh-c/pismo-tech-case/internal/service"
	"github.com/rs/zerolog/log"
)

type AccountsHandler struct {
	service *service.AccountsService
}

func NewAccountsHandler(service *service.AccountsService) *AccountsHandler {
	return &AccountsHandler{service: service}
}

// CreateAccount godoc
// @Summary 	Create account
// @Description Creates a new account with a document number
// @Tags 		accounts
// @Accept 		json
// @Produce 	json
// @Param 		request body models.CreateAccountRequest 					true "Account input"
// @Success 	201 {object} models.SuccessResponse{data=service.Account}
// @Failure 	400 {object} models.ErrResponse 							"BAD_REQUEST"
// @Failure 	409 {object} models.ErrResponse 							"ACCOUNT_ALREADY_EXISTS"
// @Failure 	500 {object} models.ErrResponse 							"INTERNAL_ERROR"
// @Router 		/accounts [post]
func (h *AccountsHandler) CreateAccount(c *gin.Context) {
	var request models.CreateAccountRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Err(err).
			Msg("invalid request body")
		ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
		return
	}

	account, err := h.service.Create(c.Request.Context(), service.CreateAccountInput{
		DocumentNumber: request.DocumentNumber,
	})
	if err != nil {
		if errors.Is(err, service.ErrAccountAlreadyExists) {
			log.Err(err).
				Msg("account already exists")
			ErrorResponse(c, http.StatusConflict, "ACCOUNT_ALREADY_EXISTS", "account already exists")
			return
		}
		if errors.Is(err, service.ErrDocNumEmpty) {
			log.Err(err).
				Msg("document number is empty")
			ErrorResponse(c, http.StatusBadRequest, "DOCUMENT_NUMBER_EMPTY", "document number is empty")
			return
		}
		log.Err(err).
			Msg("error creating account")
		ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "not possible to create the account")
		return
	}

	SuccessResponse(c, http.StatusCreated, "account created", account)
}

// GetAccount 		godoc
// @Summary 		Get Account Details
// @Description  	Retrieves account details for the given account ID
// @Tags 			accounts
// @Produce 		json
// @Param 			accountId path int true "Account ID"
// @Success 	200 {object} models.SuccessResponse{data=service.Account}
// @Failure 	400 {object} models.ErrResponse 							"INVALID_ACCOUNT_ID"
// @Failure 	404 {object} models.ErrResponse								"ACCOUNT_NOT_FOUND"
// @Failure 	500 {object} models.ErrResponse 							"INTERNAL_ERROR"
// @Router 		/accounts/{accountId} [get]
func (h *AccountsHandler) GetAccount(c *gin.Context) {
	accountIdStr := c.Param("accountId")
	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		log.Err(err).
			Msg("error converting account id from string to int")

		ErrorResponse(c, http.StatusBadRequest, "INVALID_ACCOUNT_ID", "account id is invalid")
		return
	}

	account, err := h.service.GetAccountByID(c.Request.Context(), accountId)
	if err != nil {
		if errors.Is(err, service.ErrAccountNotFound) {
			log.Err(err).
				Msg("account not found")
			ErrorResponse(c, http.StatusNotFound, "ACCOUNT_NOT_FOUND", "account not found")
			return
		}
		if errors.Is(err, service.ErrInvalidAccountID) {
			log.Err(err).
				Msg("invalid account id")

			ErrorResponse(c, http.StatusBadRequest, "INVALID_ACCOUNT_ID", "account id is invalid")
			return
		}

		log.Err(err).
			Msg("error fetching account")
		ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "not possible to get the account")
		return
	}

	SuccessResponse(c, http.StatusOK, "account fetched successfully", account)
}
