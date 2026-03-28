package route

import (
	"github.com/gin-gonic/gin"
	docs "github.com/pragadeesh-c/pismo-tech-case/cmd/docs"
	"github.com/pragadeesh-c/pismo-tech-case/internal/api/handler"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handlers struct {
	Account     *handler.AccountsHandler
	Transaction *handler.TransactionHandler
}

func RegisterRoutes(r *gin.Engine, h *Handlers) {
	docs.SwaggerInfo.BasePath = "/api/v1"
	api := r.Group("/api/v1")
	{
		api.POST("/accounts", h.Account.CreateAccount)
		api.GET("/accounts/:accountId", h.Account.GetAccount)
		api.POST("/transactions", h.Transaction.CreateTransaction)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
