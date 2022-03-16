package routes

import (
	"github.com/brucetieu/blockchain/handlers"
	"github.com/brucetieu/blockchain/repository"
	"github.com/brucetieu/blockchain/services"
	"github.com/gin-gonic/gin"
)

func InitRoutes(route *gin.Engine) {
	services.BlockAssembler = services.NewBlockAssemblerFac()
	blockchainRepo := repository.NewBlockchainRepository()
	blockService := services.NewBlockService()
	transactionService := services.NewTransactionService()
	blockchainService := services.NewBlockchainService(blockchainRepo, blockService, transactionService)
	blockchainHandler := handlers.NewBlockchainHandler(blockchainService)

	groupRoute := route.Group("/")
	groupRoute.POST("/blockchain", blockchainHandler.CreateBlockchain)
	groupRoute.GET("/blockchain", blockchainHandler.GetBlockchain)
	groupRoute.GET("/blockchain/genesis", blockchainHandler.GetGenesisBlock)
	// groupRoute.POST("/blockchain/block", blockchainHandler.AddToBlockchain)

}
