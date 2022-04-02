package services

import (
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/brucetieu/blockchain/repository"
	reps "github.com/brucetieu/blockchain/representations"
	"github.com/brucetieu/blockchain/utils"
	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"
	// "github.com/google/uuid"
)

var Reward = 50

type TransactionService interface {
	// SetID(txnRep reps.Transaction) []byte
	CreateCoinbaseTxn(to string, data string) reps.Transaction
	CreateTransaction(from string, to string, amount int) (reps.Transaction, error)

	GetTransactions() ([]reps.Transaction, error)
	GetTransaction(txnId string) (reps.Transaction, error)
	GetUnspentTransactions(address string) []reps.Transaction
	GetUnspentTxnOutputs(address string) []reps.TxnOutput
	GetSpendableOutputs(from string, amount int) (int, map[string][]int)

	CanUnlock(input reps.TxnInput, data string) bool
	CanBeUnlockedWith(output reps.TxnOutput, data string) bool
	IsCoinbaseTransaction(txn reps.Transaction) bool

	GetAddresses() (map[string]bool, error)
	GetBalance(address string) (int, error)
	GetBalances() ([]reps.AddressBalance, error)
}

type transactionService struct {
	blockchainRepo repository.BlockchainRepository
	blockAssembler BlockAssemblerFac
	txnAssembler   TxnAssemblerFac
}

func NewTransactionService(blockchainRepo repository.BlockchainRepository) TransactionService {
	return &transactionService{
		blockchainRepo: blockchainRepo,
		blockAssembler: BlockAssembler,
		txnAssembler:   TxnAssembler,
	}
}

// func (ts *transactionService) SetID(txnRep reps.Transaction) []byte {
// 	txnRepInBytes, err := json.Marshal(txnRep)
// 	if err != nil {
// 		log.Error("Error marshalling object: ", err.Error())
// 	}
// 	hashID := sha256.Sum256(txnRepInBytes)
// 	return hashID[:]
// }

// A coinbase transaction is a special type of transaction which doesn’t require previously existing outputs. It creates the output
func (ts *transactionService) CreateCoinbaseTxn(to string, data string) reps.Transaction {
	log.WithFields(log.Fields{"to": to, "data": data}).Info("Creating coinbase transaction")
	if data == "" {
		data = fmt.Sprintf("Coins to: %s", to)
	}

	txnRep := ts.txnAssembler.ToCoinbaseTxn(to, data)
	log.Info("txnRep in CreateCoinbaseTxn: ", utils.Pretty(txnRep))

	return txnRep
}

func (ts *transactionService) CreateTransaction(from string, to string, amount int) (reps.Transaction, error) {
	log.WithFields(log.Fields{"from": from, "to": to, "amount": amount}).Info("Creating transaction...")

	var transaction reps.Transaction
	var txnOutput reps.TxnOutput
	txnInputs := make([]reps.TxnInput, 0)
	txnOutputs := make([]reps.TxnOutput, 0)

	totalUnspentAmount, validOutputs := ts.GetSpendableOutputs(from, amount)
	log.WithFields(log.Fields{"totalUnspentAmount": totalUnspentAmount, "validOutputs": utils.Pretty(validOutputs)}).Info("Got spendable outputs")

	// Not enough coins to send
	if amount > totalUnspentAmount {
		err := fmt.Errorf("%s only has %d coins to send to %s, not %d, Cancelling transaction", from, totalUnspentAmount, to, amount)
		log.Error(err)
		return reps.Transaction{}, err
	}

	// For each found unspent output an input referencing it is created
	for txnId, outputIndices := range validOutputs {
		decodedTxnId, err := hex.DecodeString(txnId)
		if err != nil {
			log.WithField("error", err.Error()).Error("Error decoding transactionId to bytes")
		}

		for _, outputIdx := range outputIndices {
			input := reps.TxnInput{}
			inputID := uuid.Must(uuid.NewRandom()).String()
			input.InputID = inputID
			input.PrevTxnID = decodedTxnId
			input.OutIdx = outputIdx
			input.ScriptSig = from
			txnInputs = append(txnInputs, input)
		}
	}

	transaction.Inputs = txnInputs

	txnOutput.OutputID = uuid.Must(uuid.NewRandom()).String()
	txnOutput.Value = amount
	txnOutput.ScriptPubKey = to

	// Amount sender gave to receiver
	txnOutputs = append(txnOutputs, txnOutput)

	// Any change associated with sender
	if totalUnspentAmount > amount {
		var txnOutputChange reps.TxnOutput
		txnOutputChange.OutputID = uuid.Must(uuid.NewRandom()).String()
		txnOutputChange.Value = totalUnspentAmount - amount
		txnOutputChange.ScriptPubKey = from

		txnOutputs = append(txnOutputs, txnOutputChange)
	}

	transaction.Outputs = txnOutputs

	txnId := ts.txnAssembler.SetID(transaction)

	for i := 0; i < len(txnInputs); i++ {
		txnInputs[i].CurrTxnID = txnId
	}

	for j := 0; j < len(txnOutputs); j++ {
		txnOutputs[j].CurrTxnID = txnId
	}

	transaction.ID = txnId

	return transaction, nil
}

func (tx *transactionService) GetTransaction(txnId string) (reps.Transaction, error) {
	txnIdByte, err := hex.DecodeString(txnId)
	if err != nil {
		log.Error("error decoding string to byte: ", err.Error())
		// return reps.Transaction{}, err
	}

	txn, err := tx.blockchainRepo.GetTransaction(txnIdByte)
	if err != nil {
		errMsg := fmt.Errorf("%s, id: %s", err.Error(), txnId)
		return reps.Transaction{}, errMsg
	}

	return txn, nil
}

func (ts *transactionService) GetTransactions() ([]reps.Transaction, error) {
	txns, err := ts.blockchainRepo.GetTransactions()
	if err != nil {
		return []reps.Transaction{}, err
	}

	return txns, nil
}

func (ts *transactionService) GetAddresses() (map[string]bool, error) {
	addresses, err := ts.blockchainRepo.GetAddresses()
	if err != nil {
		log.Error("error getting addresses: ", err.Error())
		return addresses, err
	}

	return addresses, nil
}

func (ts *transactionService) GetBalance(address string) (int, error) {
	balance := 0

	addresses, err := ts.GetAddresses()
	if err != nil {
		return -1, err
	}

	if _, exists := addresses[address]; exists {
		unspentTxnOutputs := ts.GetUnspentTxnOutputs(address)
		log.Info("unspentTxnOutputs in GetBalance: ", utils.Pretty(unspentTxnOutputs))
	
		for _, unspentOutput := range unspentTxnOutputs {
			balance += unspentOutput.Value
		}
	
	} else {
		errMsg := fmt.Errorf("could not get balance: address %s was not found", address)
		return -1, errMsg
	}

	return balance, nil
}

func (ts *transactionService) GetBalances() ([]reps.AddressBalance, error) {
	addresses, err := ts.GetAddresses()
	if err != nil {
		return []reps.AddressBalance{}, err
	}

	balances := make([]reps.AddressBalance, 0)

	for key := range addresses {
		balance := 0
		unspentTxnOutputs := ts.GetUnspentTxnOutputs(key)
		log.Info("unspentTxnOutputs in GetBalance: ", utils.Pretty(unspentTxnOutputs))
	
		for _, unspentOutput := range unspentTxnOutputs {
			balance += unspentOutput.Value
		}

		balances = append(balances, reps.AddressBalance{key, balance})
	}

	return balances, nil
}

// Find out how much of the unspendable outputs from the sender can be spent given an amount
func (ts *transactionService) GetSpendableOutputs(from string, amount int) (int, map[string][]int) {
	log.WithFields(log.Fields{"from": from, "amount": amount}).Info("Calling GetSpendableOutputs")
	totalUnspentAmount := 0

	// <key>: transactionIds associated with spender
	// <value> list of all unspent output indices associated with sender for each transaction
	unspentOutIdxs := make(map[string][]int)
	unspentTxns := ts.GetUnspentTransactions(from)

	// Outer:
	for _, unspentTxn := range unspentTxns {
		txnId := hex.EncodeToString(unspentTxn.ID)

		for outputIdx, output := range unspentTxn.Outputs {
			if ts.CanBeUnlockedWith(output, from) /*&& totalUnspentAmount < amount*/ {
				unspentOutIdxs[txnId] = append(unspentOutIdxs[txnId], outputIdx)
				totalUnspentAmount += output.Value

				// if totalUnspentAmount >= amount {
				// 	break Outer
				// }
			}
		}
	}
	return totalUnspentAmount, unspentOutIdxs
}

// Get all transactions whose outputs aren't referenced in inputs
func (ts *transactionService) GetUnspentTransactions(address string) []reps.Transaction {
	var unspentTxns []reps.Transaction

	// key: transaction id, value: list of output indices
	spentTxns := make(map[string][]int)

	// TODO: Need to put the sorting in blockchainrepo perhaps
	blocks, err := ts.blockchainRepo.GetBlockchain()
	if err != nil {
		log.WithField("error", err.Error()).Error("Error getting blocks in blockchain")
	}

	// Need to process genesis block last
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].Timestamp > blocks[j].Timestamp
	})

	for _, block := range blocks {

		for _, txn := range block.Transactions {
			txnId := hex.EncodeToString(txn.ID)

		Outputs:
			for outputIdx, output := range txn.Outputs {
				// if spentTxns[txnId] != nil {
				if _, ok := spentTxns[txnId]; ok {
					for _, inputOutIdx := range spentTxns[txnId] {
						if outputIdx == inputOutIdx {
							continue Outputs // Skip this spent transaction
						}
					}
				}

				// If we get to here no OutIdx is referenced in an input; ie unspent
				if ts.CanBeUnlockedWith(output, address) {
					unspentTxns = append(unspentTxns, txn)
				}
			}

			// Non coinbase transaction will have an input with non negative OutIdx. Find all of them for a given transaction, these are the spent txns
			if !ts.IsCoinbaseTransaction(txn) {
				for _, input := range txn.Inputs {
					if ts.CanUnlock(input, address) {
						txnId := hex.EncodeToString(input.PrevTxnID)
						spentTxns[txnId] = append(spentTxns[txnId], input.OutIdx)
					}
				}
			}
		}

		// Don't process the Transactions in the Genesis block.
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return unspentTxns
}

// Get the outputs from unspent transactions
func (ts *transactionService) GetUnspentTxnOutputs(address string) []reps.TxnOutput {
	unspentTxnOutputs := make([]reps.TxnOutput, 0)
	unspentTxns := ts.GetUnspentTransactions(address)
	log.Info("Got unspentTxns in GetUnspentTxnOutputs", utils.Pretty(unspentTxns))

	for _, unspentTxn := range unspentTxns {
		for _, output := range unspentTxn.Outputs {
			if ts.CanBeUnlockedWith(output, address) {
				unspentTxnOutputs = append(unspentTxnOutputs, output)
			}
		}
	}

	return unspentTxnOutputs
}

func (ts *transactionService) CanUnlock(input reps.TxnInput, data string) bool {
	return strings.EqualFold(input.ScriptSig, data)
}

func (ts *transactionService) CanBeUnlockedWith(output reps.TxnOutput, data string) bool {
	return strings.EqualFold(output.ScriptPubKey, data)
}

func (ts *transactionService) IsCoinbaseTransaction(txn reps.Transaction) bool {
	return len(txn.Inputs) == 1 && len(txn.Inputs[0].PrevTxnID) == 0 && txn.Inputs[0].OutIdx == -1
}
