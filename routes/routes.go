package routes

import (
	"github.com/brucetieu/blockchain/handlers"
	"github.com/brucetieu/blockchain/repository"
	"github.com/brucetieu/blockchain/services"
	"github.com/gin-gonic/gin"
)

func InitRoutes(route *gin.Engine) {
	blockchainRepo := repository.NewBlockchainRepository()
	blockService := services.NewBlockService()
	blockchainService := services.NewBlockchainService(blockchainRepo, blockService)
	blockchainHandler := handlers.NewBlockchainHandler(blockchainService)

	groupRoute := route.Group("/")
	groupRoute.POST("/blockchain", blockchainHandler.CreateBlockchain)

}