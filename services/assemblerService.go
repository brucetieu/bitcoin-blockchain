package services

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"

	// "fmt"

	reps "github.com/brucetieu/blockchain/representations"

	// "github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var (
	BlockAssembler  BlockAssemblerFac
	TxnAssembler    TxnAssemblerFac
	WalletAssembler WalletAssemblerFac
)

type (
	blockAssembler  struct{}
	txnAssembler    struct{}
	walletAssembler struct{}
)

func NewBlockAssemblerFac() BlockAssemblerFac {
	return &blockAssembler{}
}

type BlockAssemblerFac interface {
	ToBlockBytes(block *reps.Block) []byte
	ToBlockStructure(data []byte) *reps.Block
	ToReadableBlock(block reps.Block) reps.ReadableBlock
}

func NewTxnAssemblerFac() TxnAssemblerFac {
	return &txnAssembler{}
}

type TxnAssemblerFac interface {
	HashTransactions(txns []reps.Transaction) []byte
	HashTransaction(txn reps.Transaction) []byte
	ToReadableTransactions(txns []reps.Transaction) []reps.ReadableTransaction
	ToReadableTransaction(txn reps.Transaction) reps.ReadableTransaction
	ToTxnBytes(txn reps.Transaction) []byte
	// ToCoinbaseTxn(to string, data string) reps.Transaction
	SetID(txnRep reps.Transaction) []byte
}

func NewWalletAssemblerFac() WalletAssemblerFac {
	return &walletAssembler{}
}

type WalletAssemblerFac interface {
	ToPrivateKeyBytes(privateKey ecdsa.PrivateKey) []byte
	ToECDSAPrivateKey(privateKey []byte) ecdsa.PrivateKey
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

func (t *txnAssembler) ToTxnBytes(txn reps.Transaction) []byte {
	txnBytes, err := json.Marshal(txn)
	if err != nil {
		log.Error("Unable to marshal", err.Error())
	}

	return txnBytes
}

func (t *txnAssembler) HashTransaction(txn reps.Transaction) []byte {
	var hash [32]byte

	txnCopy := &txn
	txnCopy.ID = nil

	hash = sha256.Sum256(t.ToTxnBytes(*txnCopy))

	return hash[:]
}

// Hash all transaction ids
func (t *txnAssembler) HashTransactions(txns []reps.Transaction) []byte {
	allTxns := make([][]byte, 0)

	for _, txn := range txns {
		allTxns = append(allTxns, txn.ID)
	}

	hashedTxns := sha256.Sum256(bytes.Join(allTxns, []byte{}))
	return hashedTxns[:]
}

// Create txn id
func (t *txnAssembler) SetID(txnRep reps.Transaction) []byte {
	txnRepInBytes, err := json.Marshal(txnRep)
	if err != nil {
		log.Error("Error marshalling object: ", err.Error())
	}
	hashID := sha256.Sum256(txnRepInBytes)
	return hashID[:]
}

// func (t *txnAssembler) ToCoinbaseTxn(to string, data string) reps.Transaction {
// 	var txnOut reps.TxnOutput
// 	var txnIn reps.TxnInput
// 	var txnRep reps.Transaction

// 	txnInputId := uuid.Must(uuid.NewRandom()).String()
// 	txnOutputId := uuid.Must(uuid.NewRandom()).String()

// 	txnOut.OutputID = txnOutputId
// 	txnOut.Value = Reward
// 	txnOut.PubKeyHash = to

// 	txnIn.InputID = txnInputId
// 	txnIn.PrevTxnID = []byte{}
// 	txnIn.OutIdx = -1
// 	txnIn.ScriptSig = data

// 	txnRep.Outputs = []reps.TxnOutput{txnOut}
// 	txnRep.Inputs = []reps.TxnInput{txnIn}

// 	// Put this here to ensure we get a different hash each time
// 	currTxnID := t.SetID(txnRep)

// 	txnRep.ID = currTxnID
// 	txnOut.CurrTxnID = currTxnID
// 	txnIn.CurrTxnID = currTxnID

// 	return txnRep
// }

func (b *blockAssembler) ToReadableBlock(block reps.Block) reps.ReadableBlock {
	var readableBlock reps.ReadableBlock

	readableBlock.ID = block.ID
	readableBlock.Timestamp = block.Timestamp
	readableBlock.PrevHash = hex.EncodeToString(block.PrevHash)
	readableBlock.Hash = hex.EncodeToString(block.Hash)
	readableBlock.Nounce = block.Nounce

	var transactions []reps.ReadableTransaction
	for _, txn := range block.Transactions {

		var inputs []reps.ReadableTxnInput
		for _, in := range txn.Inputs {
			input := reps.ReadableTxnInput{
				CurrTxnID: hex.EncodeToString(txn.ID),
				PrevTxnID: hex.EncodeToString(in.PrevTxnID),
				OutIdx:    in.OutIdx,
				PubKey:    hex.EncodeToString(in.PubKey),
				Signature: hex.EncodeToString(in.Signature),
			}
			inputs = append(inputs, input)
		}

		var outputs []reps.ReadableTxnOutput
		for _, out := range txn.Outputs {
			output := reps.ReadableTxnOutput{
				CurrTxnID:  hex.EncodeToString(txn.ID),
				Value:      out.Value,
				PubKeyHash: hex.EncodeToString(out.PubKeyHash),
			}
			outputs = append(outputs, output)
		}

		transaction := reps.ReadableTransaction{
			BlockID: block.ID,
			ID:      hex.EncodeToString(txn.ID),
			Inputs:  inputs,
			Outputs: outputs,
		}

		transactions = append(transactions, transaction)
	}

	readableBlock.Transactions = transactions

	return readableBlock
}

func (t *txnAssembler) ToReadableTransactions(txns []reps.Transaction) []reps.ReadableTransaction {
	var transactions []reps.ReadableTransaction

	for _, txn := range txns {

		var inputs []reps.ReadableTxnInput
		for _, in := range txn.Inputs {
			input := reps.ReadableTxnInput{
				CurrTxnID: hex.EncodeToString(txn.ID),
				PrevTxnID: hex.EncodeToString(in.PrevTxnID),
				PubKey:    hex.EncodeToString(in.PubKey),
				Signature: hex.EncodeToString(in.Signature),
			}
			inputs = append(inputs, input)
		}

		var outputs []reps.ReadableTxnOutput
		for _, out := range txn.Outputs {
			output := reps.ReadableTxnOutput{
				CurrTxnID:  hex.EncodeToString(txn.ID),
				Value:      out.Value,
				PubKeyHash: hex.EncodeToString(out.PubKeyHash),
			}
			outputs = append(outputs, output)
		}

		transaction := reps.ReadableTransaction{
			BlockID: txn.BlockID,
			ID:      hex.EncodeToString(txn.ID),
			Inputs:  inputs,
			Outputs: outputs,
		}

		transactions = append(transactions, transaction)
	}

	return transactions
}

func (t *txnAssembler) ToReadableTransaction(txn reps.Transaction) reps.ReadableTransaction {
	readableTxn := reps.ReadableTransaction{
		ID:      hex.EncodeToString(txn.ID),
		BlockID: txn.BlockID,
	}

	var inputs []reps.ReadableTxnInput
	for _, in := range txn.Inputs {
		input := reps.ReadableTxnInput{
			CurrTxnID: hex.EncodeToString(txn.ID),
			PrevTxnID: hex.EncodeToString(in.PrevTxnID),
			OutIdx:    in.OutIdx,
			PubKey:    hex.EncodeToString(in.PubKey),
			Signature: hex.EncodeToString(in.Signature),
		}
		inputs = append(inputs, input)
	}

	var outputs []reps.ReadableTxnOutput
	for _, out := range txn.Outputs {
		output := reps.ReadableTxnOutput{
			CurrTxnID:  hex.EncodeToString(txn.ID),
			Value:      out.Value,
			PubKeyHash: hex.EncodeToString(out.PubKeyHash),
		}
		outputs = append(outputs, output)
	}

	readableTxn.Inputs = inputs
	readableTxn.Outputs = outputs

	return readableTxn
}

// Convert ecdsa.PrivateKey to slice of bytes
func (w *walletAssembler) ToPrivateKeyBytes(privateKey ecdsa.PrivateKey) []byte {
	gob.Register(elliptic.P256())
	var content bytes.Buffer

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(privateKey)
	if err != nil {
		log.Error("unable to encode", err.Error())
	}

	return content.Bytes()
}

// Convert byte representation of the private key to a ecdsa.PrivateKey
func (w *walletAssembler) ToECDSAPrivateKey(privKeyBytes []byte) ecdsa.PrivateKey {
	var privKey ecdsa.PrivateKey
	gob.Register(elliptic.P256())

	buf := bytes.NewBuffer(privKeyBytes)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&privKey)
	if err != nil {
		log.Error("Unable to decode: ", err.Error())
	}

	return privKey
}
