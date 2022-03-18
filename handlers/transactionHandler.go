package handlers

import (
	"net/http"

	"github.com/brucetieu/blockchain/services"
	"github.com/gin-gonic/gin"
)


type TransactionHandler struct {
	transactionService services.TransactionService
}

func NewTransactionHandler(transactionService services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

func (th *TransactionHandler) GetBalance(c *gin.Context) {
	address := c.Param("address")
	balance := 0

	unspentTxnOutputs := th.transactionService.GetUnspentTxnOutputs(address)


	for _, unspentOutput := range unspentTxnOutputs {
		balance += unspentOutput.Value
	}

	c.JSON(http.StatusOK, gin.H{"address": address, "balance": balance})
}