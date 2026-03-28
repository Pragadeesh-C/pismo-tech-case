package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pragadeesh-c/pismo-tech-case/internal/api/models"
)

func SuccessResponse(c *gin.Context, status int, message string, data any) {
	c.JSON(status, models.SuccessResponse{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, status int, code string, message string) {
	c.JSON(status, models.ErrResponse{
		Status: false,
		Error: &models.ErrorBody{
			Code:    code,
			Message: message,
		},
	})
}
