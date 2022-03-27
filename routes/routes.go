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

	transactionService := services.NewTransactionService(blockchainRepo)
	blockchainService := services.NewBlockchainService(blockchainRepo, blockService, transactionService)

	blockchainHandler := handlers.NewBlockchainHandler(blockchainService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	groupRoute := route.Group("/")
	groupRoute.POST("/blockchain", blockchainHandler.CreateBlockchain)
	groupRoute.GET("/blockchain", blockchainHandler.GetBlockchain)
	groupRoute.GET("/blockchain/genesis", blockchainHandler.GetGenesisBlock)
	groupRoute.GET("/blockchain/balance/:address", transactionHandler.GetBalance)
	// groupRoute.POST("/blockchain/block", blockchainHandler.AddToBlockchain)

	// groupRoute.GET("/blockchain/balances")
	// groupRoute.GET("/blockchain/addresses", blockchainHandler.GetAddresses)
	// groupRoute.GET("/blockchain/transactions")
	// groupRoute.GET("/blockchain/transactions/:transactionId")

}
