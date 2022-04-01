package handlers

import (
	"net/http"

	"github.com/brucetieu/blockchain/services"
	"github.com/brucetieu/blockchain/utils"
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

func (th *TransactionHandler) GetBalance(c *gin.Context) {
	address := c.Param("address")
	balance := 0

	unspentTxnOutputs := th.transactionService.GetUnspentTxnOutputs(address)
	log.Info("unspentTxnOutputs in GetBalance: ", utils.Pretty(unspentTxnOutputs))

	for _, unspentOutput := range unspentTxnOutputs {
		balance += unspentOutput.Value
	}

	c.JSON(http.StatusOK, gin.H{"address": address, "balance": balance})
}

func (th *TransactionHandler) GetTransactions(c *gin.Context) {
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

	txn, err := th.transactionService.GetTransaction(txnId)
	if err != nil {
		log.Error("error getting transaction: ", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"transaction": th.assemblerService.ToReadableTransaction(txn)})
	}
}