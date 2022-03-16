package services

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	reps "github.com/brucetieu/blockchain/representations"
	log "github.com/sirupsen/logrus"
)

var Reward = 50

type TransactionService interface {
	CreateCoinbaseTxn(to string, data string) *reps.Transaction
	SetID(txnRep reps.Transaction) []byte
}

type transactionService struct {}

func NewTransactionService() TransactionService {
	return &transactionService{}
}

// A coinbase transaction is a special type of transaction which doesn’t require previously existing outputs. It creates the output
func (t *transactionService) CreateCoinbaseTxn(to string, data string) *reps.Transaction {
	log.WithFields(log.Fields{"to": to, "data": data}).Info("Creating coinbase transaction")
	if data == "" {
		data = fmt.Sprintf("Coins to: %s", to)
	}

	txnIn := reps.TxnInput{[]byte{}, -1, data}
	txnOut := reps.TxnOutput{50, to}

	txnRep := reps.Transaction{}
	txnRep.Inputs = []reps.TxnInput{txnIn}
	txnRep.Outputs = []reps.TxnOutput{txnOut}
	txnID := t.SetID(txnRep)

	txnRep.ID = txnID

	return &txnRep
}

func (t *transactionService) SetID(txnRep reps.Transaction) []byte {
	txnRepInBytes, err := json.Marshal(txnRep)
	if err != nil {
		log.Error("Error marshalling object: ", err.Error())
	}
	hashID := sha256.Sum256(txnRepInBytes)
	return hashID[:]
}

