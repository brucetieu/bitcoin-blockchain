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

	groupRoute.GET("/", blockchainHandler.BlockchainHome)

	groupRoute.POST("/blockchain", blockchainHandler.CreateBlockchain)
	groupRoute.GET("/blockchain", blockchainHandler.GetBlockchain)
	groupRoute.POST("/blockchain/block", blockchainHandler.AddToBlockchain)
	groupRoute.GET("/blockchain/block/genesis", blockchainHandler.GetGenesisBlock)
	groupRoute.GET("/blockchain/block/last", blockchainHandler.GetLastBlock)
	groupRoute.GET("/blockchain/block/:blockId", blockchainHandler.GetBlock)

	// groupRoute.GET("/blockchain/balance/:address", transactionHandler.GetBalance)
	// groupRoute.GET("/blockchain/balances", transactionHandler.GetBalances)
	// groupRoute.GET("/blockchain/addresses", transactionHandler.GetAddresses)
	groupRoute.GET("/blockchain/transactions", transactionHandler.GetTransactions)
	groupRoute.GET("/blockchain/transactions/:transactionId", transactionHandler.GetTransaction)

	groupRoute.POST("/blockchain/wallets", walletHandler.CreateWallet)
	// groupRoute.GET("/blockchain/wallets/addresses", transactionHandler.GetAddresses)
	groupRoute.GET("/blockchain/wallets", transactionHandler.GetWallets)
	groupRoute.GET("/blockchain/wallets/:address", transactionHandler.GetWallet)
	// groupRoute.GET("/blockchain/wallets/:address/balance", transactionHandler.GetAddresses)
}
