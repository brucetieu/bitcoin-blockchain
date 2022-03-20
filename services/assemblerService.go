package services

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	// "fmt"

	reps "github.com/brucetieu/blockchain/representations"
	log "github.com/sirupsen/logrus"
)

var BlockAssembler BlockAssemblerFac

type blockAssembler struct{}

func NewBlockAssemblerFac() BlockAssemblerFac {
	return &blockAssembler{}
}

type BlockAssemblerFac interface {
	ToBlockBytes(block *reps.Block) []byte
	ToBlockStructure(data []byte) *reps.Block
	HashTransactions(txns []*reps.Transaction) []byte
	ToBlockMap(block *reps.Block) map[string]interface{}
}

func (b *blockAssembler) ToBlockBytes(block *reps.Block) []byte {
	byteStruct, err := json.Marshal(block)
	if err != nil {
		log.Error("Unable to marshal", err.Error())
	}

	return byteStruct
}

func (b *blockAssembler) ToBlockStructure(data []byte) *reps.Block {
	var block reps.Block
	err := json.Unmarshal(data, &block)
	if err != nil {
		log.Error("Unable to unmarshal: ", err.Error())
	}
	return &block
}

// Hash all transaction ids
func (b *blockAssembler) HashTransactions(txns []*reps.Transaction) []byte {
	allTxns := make([][]byte, 0)

	for _, txn := range txns {
		allTxns = append(allTxns, txn.ID)
	}

	hashedTxns := sha256.Sum256(bytes.Join(allTxns, []byte{}))
	return hashedTxns[:]
}

func (b *blockAssembler) ToBlockMap(block *reps.Block) map[string]interface{} {
	data := make(map[string]interface{})
	data["timestamp"] = block.Timestamp
	data["prevHash"] = hex.EncodeToString(block.PrevHash)
	data["hash"] = hex.EncodeToString(block.Hash)
	data["nounce"] = block.Nounce

	var transactions []*reps.ReadableTransaction
	for _, txn := range block.Transactions {
		var inputs []reps.ReadableTxnInput
		for _, in := range txn.Inputs {
			input := reps.ReadableTxnInput{
				TxnID:     hex.EncodeToString(in.TxnID),
				OutIdx:    in.OutIdx,
				ScriptSig: in.ScriptSig,
			}
			inputs = append(inputs, input)
		}
		transaction := &reps.ReadableTransaction{
			ID:      hex.EncodeToString(txn.ID),
			Inputs:  inputs,
			Outputs: txn.Outputs,
		}
		transactions = append(transactions, transaction)
	}
	data["transactions"] = transactions

	return data
}
