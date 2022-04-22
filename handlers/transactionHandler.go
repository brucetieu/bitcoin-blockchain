package handlers

import (
	"net/http"

	"github.com/brucetieu/blockchain/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type TransactionHandler struct {
	transactionService services.TransactionService
	assemblerService  services.TxnAssemblerFac
}

func NewTransactionHandler(transactionService services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		assemblerService: services.TxnAssembler,
	}
}

// func (th *TransactionHandler) GetBalances(c *gin.Context) {
// 	balances, err := th.transactionService.GetBalances()
// 	if err != nil {
// 		log.Error("error getting balances: ", err.Error())
// 		c.JSON(http.StatusInternalServerError, err.Error())
// 	} else {
// 		c.JSON(http.StatusOK, gin.H{"balances": balances})
// 	}
// }

// func (th *TransactionHandler) GetBalance(c *gin.Context) {
// 	address := c.Param("address")

// 	balance, err := th.transactionService.GetBalance(address)
// 	if err != nil {
// 		log.Error("error getting balance: ", err.Error())
// 		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 	} else {
// 		c.JSON(http.StatusOK, gin.H{"address": address, "balance": balance})
// 	}

// }

func (th *TransactionHandler) GetTransactions(c *gin.Context) {
	log.Info("GetTransactions called")
	txns, err := th.transactionService.GetTransactions()
	if err != nil {
		log.Error("error getting transactions: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"transactions": th.assemblerService.ToReadableTransactions(txns)})
	}
}

func (th *TransactionHandler) GetTransaction(c *gin.Context) {
	txnId := c.Param("transactionId")
	log.Info("GetTransaction called with transactionId: " + txnId)

	txn, err := th.transactionService.GetTransaction(txnId)
	if err != nil {
		log.Error("error getting transaction: ", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"transaction": th.assemblerService.ToReadableTransaction(txn)})
	}
}

func (th *TransactionHandler) GetWallets(c *gin.Context) {
	wallets, err := th.transactionService.GetWallets()
	if err != nil {
		log.Error("error getting all wallets: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"wallets": wallets})
	}
}

func (th *TransactionHandler) GetWallet(c *gin.Context) {
	address := c.Param("address")

	wallet, err := th.transactionService.GetWallet(address)
	if err != nil {
		log.Errorf("error getting wallet with address: %s %s", address, err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"wallet": wallet})
	}
}


// func (th *TransactionHandler) GetAddresses(c *gin.Context) {
// 	addresses, err := th.transactionService.GetAddresses()
// 	if err != nil {
// 		log.Error("error getting all addresses in transactions: ", err.Error())
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 	} else {
// 		c.JSON(http.StatusOK, gin.H{"addresses": addresses})
// 	}
// }