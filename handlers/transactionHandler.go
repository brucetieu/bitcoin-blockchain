package handlers

import (
	"net/http"

	"github.com/brucetieu/blockchain/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type TransactionHandler struct {
	transactionService services.TransactionService
	assemblerService   services.TxnAssemblerFac
}

func NewTransactionHandler(transactionService services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		assemblerService:   services.TxnAssembler,
	}
}

// GetTransactions ... Get all transactions on the blockchain
// @Summary      Get all transactions
// @Description  Get all transactions that exist on the blockchain
// @Tags         Transactions
// @Success      200  {array}   representations.ReadableTransaction
// @Failure      500  {object}  HTTPError
// @Router       /blockchain/transactions [get]
func (th *TransactionHandler) GetTransactions(ctx *gin.Context) {
	log.Info("GetTransactions called")
	txns, err := th.transactionService.GetTransactions()
	if err != nil {
		log.Error("error getting transactions: ", err.Error())
		NewError(ctx, http.StatusInternalServerError, err)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"transactions": th.assemblerService.ToReadableTransactions(txns)})
	}
}

// GetTransactions ... Get a single transaction
// @Summary      Get a transaction
// @Description  Get a transaction on the blockchain
// @Tags         Transactions
// @Param        transactionId  path      string  true  "Transaction ID"
// @Success      200            {object}  representations.ReadableTransaction
// @Failure      404            {object}  HTTPError
// @Router       /blockchain/transactions/{transactionId} [get]
func (th *TransactionHandler) GetTransaction(ctx *gin.Context) {
	txnId := ctx.Param("transactionId")
	log.Info("GetTransaction called with transactionId: " + txnId)

	txn, err := th.transactionService.GetTransaction(txnId)
	if err != nil {
		log.Error("error getting transaction: ", err.Error())
		NewError(ctx, http.StatusNotFound, err)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"transaction": th.assemblerService.ToReadableTransaction(txn)})
	}
}

// GetBalances ... Get the coin balance for each address on the blockchain
// @Summary      Get coin balances
// @Description  Get the coin balances for each address on the blockchain
// @Tags         Wallets
// @Success      200  {array}   representations.AddressBalance
// @Failure      404  {object}  HTTPError
// @Router       /blockchain/wallets/balances [get]
func (th *TransactionHandler) GetBalances(ctx *gin.Context) {
	log.Info("GetBalances called")

	balances, err := th.transactionService.GetBalances()
	if err != nil {
		log.Error("error getting transaction: ", err.Error())
		NewError(ctx, http.StatusNotFound, err)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"balances": balances})
	}
}

// GetBalances ... Get the coin balance for a single address on the blockchain
// @Summary      Get coin balance
// @Description  Get the coin balance for an address on the blockchain
// @Tags         Wallets
// @Success      200  {integer}  integer
// @Failure      404  {object}   HTTPError
// @Router       /blockchain/wallets/{address}/balance [get]
func (th *TransactionHandler) GetBalance(ctx *gin.Context) {
	log.Info("GetBalances called")
	address := ctx.Param("address")

	balance, err := th.transactionService.GetBalance(address)
	if err != nil {
		log.Error("error getting transaction: ", err.Error())
		NewError(ctx, http.StatusNotFound, err)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"balance": balance})
	}
}
