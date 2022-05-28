package routes

import (
	"github.com/brucetieu/blockchain/handlers"
	"github.com/brucetieu/blockchain/repository"
	"github.com/brucetieu/blockchain/services"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitRoutes(route *gin.Engine) {
	services.BlockAssembler = services.NewBlockAssemblerFac()
	services.TxnAssembler = services.NewTxnAssemblerFac()
	services.WalletAssembler = services.NewWalletAssemblerFac()

	blockchainRepo := repository.NewBlockchainRepository()
	blockService := services.NewBlockService(blockchainRepo)

	walletService := services.NewWalletService(blockchainRepo)
	transactionService := services.NewTransactionService(blockchainRepo, walletService)
	blockchainService := services.NewBlockchainService(blockchainRepo, blockService, transactionService, walletService)

	blockchainHandler := handlers.NewBlockchainHandler(blockchainService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	walletHandler := handlers.NewWalletHandler(walletService)

	groupRoute := route.Group("/")

	// Health check
	groupRoute.GET("/bitcoin", blockchainHandler.BlockchainHome)

	// Blockchain handlers
	groupRoute.POST("/bitcoin/blockchain", blockchainHandler.CreateBlockchain)
	groupRoute.GET("/bitcoin/blockchain", blockchainHandler.GetBlockchain)

	// Block handlers
	groupRoute.POST("/bitcoin/blockchain/block", blockchainHandler.AddToBlockchain)
	groupRoute.GET("/bitcoin/blockchain/block/genesis", blockchainHandler.GetGenesisBlock)
	groupRoute.GET("/bitcoin/blockchain/block/last", blockchainHandler.GetLastBlock)
	groupRoute.GET("/bitcoin/blockchain/block/:blockId", blockchainHandler.GetBlock)

	// Transaction handlers
	groupRoute.GET("/bitcoin/blockchain/transactions", transactionHandler.GetTransactions)
	groupRoute.GET("/bitcoin/blockchain/transactions/:transactionId", transactionHandler.GetTransaction)

	// Wallet handlers
	groupRoute.POST("/bitcoin/blockchain/wallets", walletHandler.CreateWallet)
	groupRoute.GET("/bitcoin/blockchain/wallets", walletHandler.GetWallets)
	groupRoute.GET("/bitcoin/blockchain/wallets/balances", transactionHandler.GetBalances)
	groupRoute.GET("/bitcoin/blockchain/wallets/:address", walletHandler.GetWallet)
	groupRoute.GET("/bitcoin/blockchain/wallets/:address/balance", transactionHandler.GetBalance)

	// swagger
	groupRoute.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
