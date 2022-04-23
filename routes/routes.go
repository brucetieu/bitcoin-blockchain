package routes

import (
	"github.com/brucetieu/blockchain/handlers"
	"github.com/brucetieu/blockchain/repository"
	"github.com/brucetieu/blockchain/services"
	"github.com/gin-gonic/gin"
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
	groupRoute.GET("/", blockchainHandler.BlockchainHome)

	// Blockchain handlers
	groupRoute.POST("/blockchain", blockchainHandler.CreateBlockchain)
	groupRoute.GET("/blockchain", blockchainHandler.GetBlockchain)

	// Block handlers
	groupRoute.POST("/blockchain/block", blockchainHandler.AddToBlockchain)
	groupRoute.GET("/blockchain/block/genesis", blockchainHandler.GetGenesisBlock)
	groupRoute.GET("/blockchain/block/last", blockchainHandler.GetLastBlock)
	groupRoute.GET("/blockchain/block/:blockId", blockchainHandler.GetBlock)

	// Transaction handlers
	groupRoute.GET("/blockchain/transactions", transactionHandler.GetTransactions)
	groupRoute.GET("/blockchain/transactions/:transactionId", transactionHandler.GetTransaction)

	// Wallet handlers
	groupRoute.POST("/blockchain/wallets", walletHandler.CreateWallet)
	groupRoute.GET("/blockchain/wallets", walletHandler.GetWallets)
	groupRoute.GET("/blockchain/wallets/balances", transactionHandler.GetBalances)
	groupRoute.GET("/blockchain/wallets/:address", walletHandler.GetWallet)
	groupRoute.GET("/blockchain/wallets/:address/balance", transactionHandler.GetBalance)
}
