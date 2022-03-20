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
	log.Info("unspentTxnOutputs in GetBalance: ", utils.Pretty(unspentTxnOutputs))

	for _, unspentOutput := range unspentTxnOutputs {
		balance += unspentOutput.Value
	}

	c.JSON(http.StatusOK, gin.H{"address": address, "balance": balance})
}
