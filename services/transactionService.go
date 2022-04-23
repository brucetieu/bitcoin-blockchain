package services

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	// "strings"

	"github.com/akamensky/base58"
	"github.com/brucetieu/blockchain/repository"
	reps "github.com/brucetieu/blockchain/representations"
	"github.com/brucetieu/blockchain/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/ripemd160"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	// "github.com/google/uuid"
)

var Reward = 500 // Initial reward miner gets for mining the first block

type TransactionService interface {
	NewTxnOutput(value int, address string) reps.TxnOutput

	// SetID(txnRep reps.Transaction) []byte
	CreateCoinbaseTxn(to string, data string) reps.Transaction
	CreateTransaction(from string, to string, amount int) (reps.Transaction, error)
	CreateTrimmedTxnCopy(txn reps.Transaction) reps.Transaction

	GetTransactions() ([]reps.Transaction, error)
	GetTransaction(txnId string) (reps.Transaction, error)
	GetUnspentTransactions(address []byte) []reps.Transaction
	GetUnspentTxnOutputs(address []byte) []reps.TxnOutput
	GetSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int)

	// CanUnlock(input reps.TxnInput, data string) bool
	// CanBeUnlockedWith(output reps.TxnOutput, data string) bool
	IsCoinbaseTransaction(txn reps.Transaction) bool

	VerifyTransaction(txn reps.Transaction) (bool, error)
	VerifySignature(currTxn reps.Transaction, prevTxns map[string]reps.Transaction) (bool, error)

	GetBalances() ([]reps.AddressBalance, error)
	GetBalance(address string) (int, error)
}

type transactionService struct {
	blockchainRepo  repository.BlockchainRepository
	walletService   WalletService
	blockAssembler  BlockAssemblerFac
	txnAssembler    TxnAssemblerFac
	walletAssembler WalletAssemblerFac
}

func NewTransactionService(blockchainRepo repository.BlockchainRepository, walletService WalletService) TransactionService {
	return &transactionService{
		blockchainRepo:  blockchainRepo,
		walletService:   walletService,
		blockAssembler:  BlockAssembler,
		txnAssembler:    TxnAssembler,
		walletAssembler: WalletAssembler,
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

	txnRep := ts.ToCoinbaseTxn(to, data)
	log.Info("txnRep in CreateCoinbaseTxn: ", utils.Pretty(txnRep))

	return txnRep
}

// Given an address, create a coinbase transaction representation
func (ts *transactionService) ToCoinbaseTxn(to string, data string) reps.Transaction {
	var txnOut reps.TxnOutput
	var txnIn reps.TxnInput
	var txnRep reps.Transaction

	txnInputId := uuid.Must(uuid.NewRandom()).String()
	// txnOutputId := uuid.Must(uuid.NewRandom()).String()

	txnOut = ts.NewTxnOutput(Reward, to)
	// txnOut.OutputID = txnOutputId
	// txnOut.Value = Reward
	// txnOut.PubKeyHash = to

	txnIn.InputID = txnInputId
	txnIn.PrevTxnID = []byte{}
	txnIn.OutIdx = -1
	txnIn.PubKey = []byte(data)
	txnIn.Signature = nil // don't sign coinbase txn

	txnRep.Outputs = []reps.TxnOutput{txnOut}
	txnRep.Inputs = []reps.TxnInput{txnIn}

	// Put this here to ensure we get a different hash each time
	currTxnID := ts.txnAssembler.SetID(txnRep)

	txnRep.ID = currTxnID
	txnOut.CurrTxnID = currTxnID
	txnIn.CurrTxnID = currTxnID

	return txnRep
}

// Create a transaction. This does the following:
// 1. Create locked outputs (populate PubKeyHash in the output)
// 2. Create new input referencing locked outputs
// 3. sign the transaction
func (ts *transactionService) CreateTransaction(from string, to string, amount int) (reps.Transaction, error) {
	log.WithFields(log.Fields{"from": from, "to": to, "amount": amount}).Info("Creating transaction...")

	var transaction reps.Transaction
	txnOutput := ts.NewTxnOutput(amount, to)
	txnInputs := make([]reps.TxnInput, 0)
	txnOutputs := make([]reps.TxnOutput, 0)

	// Check that a wallet exists to send coins from
	wallet, err := ts.walletService.GetWallet(from)
	if err != nil {
		return reps.Transaction{}, err
	}

	pubKeyBytes, _ := hex.DecodeString(wallet.PublicKey)
	pubKeyHash, _ := ts.walletService.CreatePubKeyHash(pubKeyBytes)
	privKey := ts.walletAssembler.ToECDSAPrivateKey(wallet.PrivateKey)

	totalUnspentAmount, validOutputs := ts.GetSpendableOutputs(pubKeyHash, amount)
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
			input.PubKey = pubKeyBytes
			txnInputs = append(txnInputs, input)
		}
	}

	transaction.Inputs = txnInputs

	// Amount sender gave to receiver
	txnOutputs = append(txnOutputs, txnOutput)

	// Any change associated with sender
	if totalUnspentAmount > amount {
		txnOutputChange := ts.NewTxnOutput(totalUnspentAmount-amount, from)
		txnOutputs = append(txnOutputs, txnOutputChange)
	}

	transaction.Outputs = txnOutputs

	// txnId := ts.txnAssembler.SetID(transaction)
	txnId := ts.txnAssembler.HashTransaction(transaction)

	for i := 0; i < len(txnInputs); i++ {
		txnInputs[i].CurrTxnID = txnId
	}

	for j := 0; j < len(txnOutputs); j++ {
		txnOutputs[j].CurrTxnID = txnId
	}

	transaction.ID = txnId

	// sign transaction
	transaction, err = ts.SignTransaction(transaction, privKey)
	if err != nil {
		return reps.Transaction{}, err
	}

	return transaction, nil
}

// Get transaction on a block by transactionId
func (tx *transactionService) GetTransaction(txnId string) (reps.Transaction, error) {
	log.Info("Attempting to get transaction with transaction id: ", txnId)
	txnIdByte, err := hex.DecodeString(txnId)
	if err != nil {
		log.Error("error decoding string to byte: ", err.Error())
	}

	txn, err := tx.blockchainRepo.GetTransaction(txnIdByte)
	if err != nil {
		errMsg := fmt.Errorf("%s, id: %s", err.Error(), txnId)
		return reps.Transaction{}, errMsg
	}

	return txn, nil
}

// Get all transactions that exist on blockchain
func (ts *transactionService) GetTransactions() ([]reps.Transaction, error) {
	log.Info("Attempting to get all transactions on the blockchain")
	txns, err := ts.blockchainRepo.GetTransactions()
	if err != nil {
		return []reps.Transaction{}, err
	}

	return txns, nil
}

// Get balances for each address / wallet
func (ts *transactionService) GetBalances() ([]reps.AddressBalance, error) {
	log.Info("Attempting to get the balance for each wallet / address")
	wallets, err := ts.walletService.GetWallets()
	if err != nil {
		return []reps.AddressBalance{}, err
	}

	addressBalances := make([]reps.AddressBalance, 0)

	for _, wallet := range wallets {
		balance := 0

		pubKeyHash := base58Decode([]byte(wallet.Address))
		pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-ChecksumLen]

		unspentTxnOutputs := ts.GetUnspentTxnOutputs(pubKeyHash)
		log.Info("unspentTxnOutputs in GetBalances for address: "+wallet.Address, utils.Pretty(unspentTxnOutputs))

		for _, unspentOutput := range unspentTxnOutputs {
			balance += unspentOutput.Value
		}

		addressBalances = append(addressBalances, reps.AddressBalance{Address: wallet.Address, Balance: balance})
	}

	return addressBalances, nil
}

// Get balance for a single address
func (ts *transactionService) GetBalance(address string) (int, error) {
	log.Info("Attempting to get the balance for the address: ", address)
	wallet, err := ts.walletService.GetWallet(address)
	if err != nil {
		return 0, err
	}

	balance := 0

	pubKeyHash := base58Decode([]byte(wallet.Address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-ChecksumLen]

	unspentTxnOutputs := ts.GetUnspentTxnOutputs(pubKeyHash)
	log.Info("unspentTxnOutputs in GetBalances for address: "+wallet.Address, utils.Pretty(unspentTxnOutputs))

	for _, unspentOutput := range unspentTxnOutputs {
		balance += unspentOutput.Value
	}

	return balance, nil
}

// Find out how much of the unspendable outputs from the sender can be spent given an amount
func (ts *transactionService) GetSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	log.WithFields(log.Fields{"from": hex.EncodeToString(pubKeyHash), "amount": amount}).Info("Calling GetSpendableOutputs")
	totalUnspentAmount := 0

	// <key>: transactionIds associated with spender
	// <value> list of all unspent output indices associated with sender for each transaction
	unspentOutIdxs := make(map[string][]int)
	unspentTxns := ts.GetUnspentTransactions(pubKeyHash)

	// Outer:
	for _, unspentTxn := range unspentTxns {
		txnId := hex.EncodeToString(unspentTxn.ID)

		for outputIdx, output := range unspentTxn.Outputs {
			if ts.IsLockedWithKey(output, pubKeyHash) /*&& totalUnspentAmount < amount*/ {
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
func (ts *transactionService) GetUnspentTransactions(pubKeyHash []byte) []reps.Transaction {
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
				if ts.IsLockedWithKey(output, pubKeyHash) {
					unspentTxns = append(unspentTxns, txn)
				}
			}

			// Non coinbase transaction will have an input with non negative OutIdx. Find all of them for a given transaction, these are the spent txns
			if !ts.IsCoinbaseTransaction(txn) {
				for _, input := range txn.Inputs {
					if ts.UsesKey(input, pubKeyHash) {
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
func (ts *transactionService) GetUnspentTxnOutputs(address []byte) []reps.TxnOutput {
	unspentTxnOutputs := make([]reps.TxnOutput, 0)
	unspentTxns := ts.GetUnspentTransactions(address)
	// log.Info("Got unspentTxns in GetUnspentTxnOutputs", utils.Pretty(unspentTxns))

	for _, unspentTxn := range unspentTxns {
		for _, output := range unspentTxn.Outputs {
			if ts.IsLockedWithKey(output, address) {
				unspentTxnOutputs = append(unspentTxnOutputs, output)
			}
		}
	}

	return unspentTxnOutputs
}

// Create a new transaction output. Sending of tokens "locks" the output
func (ts *transactionService) NewTxnOutput(value int, address string) reps.TxnOutput {
	txnOutput := reps.TxnOutput{
		OutputID:   uuid.Must(uuid.NewRandom()).String(),
		Value:      value,
		PubKeyHash: nil,
	}
	ts.Lock(&txnOutput, address)
	return txnOutput
}

func (ts *transactionService) SignTransaction(txn reps.Transaction, privKey ecdsa.PrivateKey) (reps.Transaction, error) {
	log.Info("Attempting to sign transaction: ", hex.EncodeToString(txn.ID))
	prevTxns := make(map[string]reps.Transaction)

	for _, input := range txn.Inputs {
		prevTxn, err := ts.blockchainRepo.GetTransaction(input.PrevTxnID)
		if err != nil {
			log.Error("error finding previous transaction with id: ", input.PrevTxnID)
			return reps.Transaction{}, err
		}
		prevTxns[hex.EncodeToString(prevTxn.ID)] = prevTxn
	}

	return ts.Sign(privKey, txn, prevTxns)
}

func (ts *transactionService) VerifyTransaction(txn reps.Transaction) (bool, error) {
	log.Info("Attempting to verify transaction: ", hex.EncodeToString(txn.ID))
	prevTxns := make(map[string]reps.Transaction)

	for _, input := range txn.Inputs {
		prevTxn, err := ts.blockchainRepo.GetTransaction(input.PrevTxnID)
		if err != nil {
			log.Error("error finding previous transaction with id: ", input.PrevTxnID)
			return false, err
		}
		prevTxns[hex.EncodeToString(prevTxn.ID)] = prevTxn
	}

	return ts.VerifySignature(txn, prevTxns)
}

func (ts *transactionService) Sign(privKey ecdsa.PrivateKey, txn reps.Transaction, prevTxns map[string]reps.Transaction) (reps.Transaction, error) {
	log.Info("Attempting to sign: ", hex.EncodeToString(txn.ID))
	if ts.IsCoinbaseTransaction(txn) {
		return reps.Transaction{}, nil
	}

	for _, in := range txn.Inputs {
		if prevTxns[hex.EncodeToString(in.PrevTxnID)].ID == nil {
			log.WithField("input currTxnID", in.PrevTxnID).Error("error: previous transaction does not exist")
			// return reps.Transaction{}, fmt.Errorf("error: previous transaction does not exist with id: %s", in.PrevTxnID)
		}

	}

	// trimmed txn copy is signed, not a full one
	txnCopy := ts.CreateTrimmedTxnCopy(txn)

	for inIdx, input := range txnCopy.Inputs {
		prevTxn := prevTxns[hex.EncodeToString(input.PrevTxnID)]
		txnCopy.Inputs[inIdx].Signature = nil

		// This txns input's pubKey is set to the outputs pubKeyHash from previous transaction
		txnCopy.Inputs[inIdx].PubKey = prevTxn.Outputs[input.OutIdx].PubKeyHash

		// Sign the Public key hashes stored in unlocked outputs. This identifies “sender” of a transaction.
		txnCopy.ID = ts.txnAssembler.HashTransaction(txnCopy)
		txnCopy.Inputs[inIdx].PubKey = nil // don't affect further iterations

		// sign txnCopy.Id with privKey
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txnCopy.ID)
		if err != nil {
			log.Error("error signing transaction: ", err.Error())
		}

		signature := append(r.Bytes(), s.Bytes()...)
		txn.Inputs[inIdx].Signature = signature
	}

	return txn, nil
	// Update db with signature
}

func (ts *transactionService) VerifySignature(currTxn reps.Transaction, prevTxns map[string]reps.Transaction) (bool, error) {
	log.Info("Attempting to verify signature of transaction: "+hex.EncodeToString(currTxn.ID)+" with inputs: ", utils.Pretty(currTxn.Inputs))
	txnCopy := ts.CreateTrimmedTxnCopy(currTxn)
	curve := elliptic.P256()

	for inIdx, in := range currTxn.Inputs {

		if in.Signature == nil {
			log.Error("Signature cannot be null, cannot verify.")
			return false, fmt.Errorf("cannot verify a null signature. Invalid transaction")
		}

		// need same data that was signed
		prevTxn := prevTxns[hex.EncodeToString(in.PrevTxnID)]
		txnCopy.Inputs[inIdx].Signature = nil
		txnCopy.Inputs[inIdx].PubKey = prevTxn.Outputs[in.OutIdx].PubKeyHash
		txnCopy.ID = ts.txnAssembler.HashTransaction(txnCopy)
		txnCopy.Inputs[inIdx].PubKey = nil

		// Unpack signature, signature is a pair of numbers
		r := big.Int{}
		s := big.Int{}
		sigLen := len(in.Signature)
		r.SetBytes(in.Signature[:(sigLen / 2)])
		s.SetBytes(in.Signature[(sigLen / 2):])

		// Unpack pubKey, pubKey is a pair of points
		x := big.Int{}
		y := big.Int{}
		pubKeyLen := len(in.PubKey)
		x.SetBytes(in.PubKey[:(pubKeyLen / 2)])
		y.SetBytes(in.PubKey[(pubKeyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}

		// verifies the signature in r, s of hash (txnCopy.ID) using the public key.
		if ecdsa.Verify(&rawPubKey, txnCopy.ID, &r, &s) == false {
			return false, fmt.Errorf("Signature: %x could not be verified", in.Signature)
		}
	}

	log.Info("Signature is valid")
	return true, nil
}

// Create copy of a transaction, but remove the signature and pub key from the input inside the transaction
func (ts *transactionService) CreateTrimmedTxnCopy(txn reps.Transaction) reps.Transaction {
	var inputs []reps.TxnInput
	var outputs []reps.TxnOutput

	for _, in := range txn.Inputs {
		inputs = append(inputs, reps.TxnInput{in.InputID, in.CurrTxnID, in.PrevTxnID, in.OutIdx, nil, nil})
	}

	for _, out := range txn.Outputs {
		outputs = append(outputs, reps.TxnOutput{out.OutputID, out.CurrTxnID, out.Value, out.PubKeyHash})
	}

	txnCopy := reps.Transaction{
		ID:      txn.ID,
		Inputs:  inputs,
		Outputs: outputs,
	}

	return txnCopy
}

// checks that an input uses a specific key to unlock an output.
// input.PubKey is not hashed, so hash it to compare it with hashed pub key of output
func (ts *transactionService) UsesKey(input reps.TxnInput, pubKeyHash []byte) bool {
	lockingHash, err := createPubKeyHash(input.PubKey)
	if err != nil {
		log.Warn("error creating pubKeyHash: ", err.Error())
	}

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

// Lock an output. When we send coins to someone, we know only their address
func (ts *transactionService) Lock(output *reps.TxnOutput, address string) {
	pubKeyHash, err := base58.Decode(address)
	if err != nil {
		log.WithField("error", err.Error()).Warn("error decoding address " + address)
	}

	output.PubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4] // remove version and checksum

	log.Infof("Locking output with address: %s with PubKeyHash of: ", address, hex.EncodeToString(output.PubKeyHash))
}

// checks if provided public key hash was used to lock the output
func (ts *transactionService) IsLockedWithKey(output reps.TxnOutput, pubKeyHash []byte) bool {
	return bytes.Compare(output.PubKeyHash, pubKeyHash) == 0
}

func (ts *transactionService) IsCoinbaseTransaction(txn reps.Transaction) bool {
	return len(txn.Inputs) == 1 && len(txn.Inputs[0].PrevTxnID) == 0 && txn.Inputs[0].OutIdx == -1
}

func createPubKeyHash(pubKey []byte) ([]byte, error) {
	pubHash := sha256.Sum256(pubKey)

	ripemdHasher := ripemd160.New()
	_, err := ripemdHasher.Write(pubHash[:])
	if err != nil {
		logrus.Error(err.Error())
	}

	pubKeyHash := ripemdHasher.Sum(nil)

	// log.Info(fmt.Sprintf("pubKeyHash: %x\n", pubKeyHash))
	return pubKeyHash, err
}
