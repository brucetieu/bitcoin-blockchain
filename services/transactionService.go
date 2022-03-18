package services

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/brucetieu/blockchain/repository"
	reps "github.com/brucetieu/blockchain/representations"
	log "github.com/sirupsen/logrus"
)

var Reward = 50

type TransactionService interface {
	CreateCoinbaseTxn(to string, data string) *reps.Transaction
	SetID(txnRep reps.Transaction) []byte
	GetUnspentTransactions(address string) []reps.Transaction
	CanUnlock(input *reps.TxnInput, data string) bool
	CanBeUnlockedWith(output *reps.TxnOutput, data string) bool
	IsCoinbaseTransaction(txn *reps.Transaction) bool
	GetUnspentTxnOutputs(address string) []reps.TxnOutput
	
}

type transactionService struct {
	blockchainRepo repository.BlockchainRepository
	blockAssembler BlockAssemblerFac
}

func NewTransactionService(blockchainRepo repository.BlockchainRepository) TransactionService {
	return &transactionService{
		blockchainRepo: blockchainRepo,
		blockAssembler: BlockAssembler,
	}
}

// A coinbase transaction is a special type of transaction which doesn’t require previously existing outputs. It creates the output
func (ts *transactionService) CreateCoinbaseTxn(to string, data string) *reps.Transaction {
	log.WithFields(log.Fields{"to": to, "data": data}).Info("Creating coinbase transaction")
	if data == "" {
		data = fmt.Sprintf("Coins to: %s", to)
	}

	txnIn := reps.TxnInput{[]byte{}, -1, data}
	txnOut := reps.TxnOutput{Reward, to}

	txnRep := reps.Transaction{}
	txnRep.Inputs = []reps.TxnInput{txnIn}
	txnRep.Outputs = []reps.TxnOutput{txnOut}
	txnID := ts.SetID(txnRep)

	txnRep.ID = txnID

	return &txnRep
}

func (ts *transactionService) SetID(txnRep reps.Transaction) []byte {
	txnRepInBytes, err := json.Marshal(txnRep)
	if err != nil {
		log.Error("Error marshalling object: ", err.Error())
	}
	hashID := sha256.Sum256(txnRepInBytes)
	return hashID[:]
}

func (ts *transactionService) GetUnspentTransactions(address string) []reps.Transaction {
	var unspentTxns []reps.Transaction

	// key: transaction id, value: list of output indices 
	spentTxns := make(map[string][]int) 

	allBlocks, err := ts.blockchainRepo.GetBlockchain()
	if err != nil {
		log.WithField("error", err.Error()).Error("Error getting blocks in blockchain")
	}

	for _, byteBlock := range allBlocks {
		block := ts.blockAssembler.ToBlockStructure(byteBlock)

		for _, txn := range block.Transactions {
			txnId := hex.EncodeToString(txn.ID)

		Outputs:
			for outputIdx, output := range txn.Outputs {
				if spentTxns[txnId] != nil {
					for _, inputOutIdx := range spentTxns[txnId] {
						if outputIdx == inputOutIdx {
							continue Outputs  // Skip this spent transaction
						}
					}
				}

				// If we get to here no OutIdx is referenced in an input
				if ts.CanBeUnlockedWith(&output, address) {
					unspentTxns = append(unspentTxns, *txn)
				}
			}

			if !ts.IsCoinbaseTransaction(txn) {
				for _, input := range txn.Inputs {
					if ts.CanUnlock(&input, address) {
						txnId := hex.EncodeToString(input.TxnID)
						spentTxns[txnId] = append(spentTxns[txnId], input.OutIdx)
					}
				}
			}
		}
	}
	return unspentTxns
}

func (ts *transactionService) GetUnspentTxnOutputs(address string) []reps.TxnOutput {
	unspentTxnOutputs := make([]reps.TxnOutput, 0)
	unspentTxns := ts.GetUnspentTransactions(address)

	for _, unspentTxn := range unspentTxns {
		for _, output := range unspentTxn.Outputs {
			if ts.CanBeUnlockedWith(&output, address) {
				unspentTxnOutputs = append(unspentTxnOutputs, output)
			}
		}
	}

	return unspentTxnOutputs
}

func (ts *transactionService) CanUnlock(input *reps.TxnInput, data string) bool {
	return strings.EqualFold(input.ScriptSig, data)
}

func (ts *transactionService) CanBeUnlockedWith(output *reps.TxnOutput, data string) bool {
	return strings.EqualFold(output.ScriptPubKey, data)
}

func (ts *transactionService) IsCoinbaseTransaction(txn *reps.Transaction) bool {
	return len(txn.Inputs) == 1 && len(txn.Inputs[0].TxnID) == 0 && txn.Inputs[0].OutIdx == -1
}
