package repository

import (
	// "strings"

	// "fmt"

	"github.com/brucetieu/blockchain/db"
	"github.com/brucetieu/blockchain/utils"

	// "github.com/brucetieu/blockchain/utils"
	// badger "github.com/dgraph-io/badger/v3"
	// log "github.com/sirupsen/logrus"

	reps "github.com/brucetieu/blockchain/representations"
)

var lastBlockHashKey = "lastBlockHash" // This gives us back block.Hash for the last block

type BlockchainRepository interface {
	// CreateBlock(hash []byte, blockByte []byte) ([]byte, error)
	CreateTransaction(txn []reps.Transaction) error
	GetTransactions(blockId string) ([]reps.Transaction, error)
	CreateBlock(block reps.Block) error
	GetGenesisBlock() (reps.Block, int, error)
	GetBlockchain() ([]reps.Block, error)

	CreateTxnInput(txnInput reps.TxnInput) error
	GetTxnInputs(txnId []byte) ([]reps.TxnInput, error)
	GetTxnOutputs(txnId []byte) ([]reps.TxnOutput, error)
	CreateTxnOutput(txnOutput reps.TxnOutput) error

	// GetBlockchain() ([][]byte, error)
}

type blockchainRepository struct {
}

func NewBlockchainRepository() BlockchainRepository {
	return &blockchainRepository{}
}

func (repo *blockchainRepository) CreateTxnInput(txnInput reps.TxnInput) error {
	if err := db.DB.Create(&txnInput).Error; err != nil {
		return err
	}
	return nil
}

func (repo *blockchainRepository) CreateTxnOutput(txnOutput reps.TxnOutput) error {
	if err := db.DB.Create(&txnOutput).Error; err != nil {
		return err
	}
	return nil
}

func (repo *blockchainRepository) CreateTransaction(txn []reps.Transaction) error {
	utils.PrettyPrintln("txn", txn)
	if err := db.DB.Create(&txn).Error; err != nil {
		return err
	}
	return nil
}

// Get all transactions in a given block.
func (repo *blockchainRepository) GetTransactions(blockId string) ([]reps.Transaction, error) {
	var transactions []reps.Transaction

	err := db.DB.Where("block_id = ?", blockId).
			Find(&transactions).
			Error
	
	if err != nil {
		return []reps.Transaction{}, err
	}

	return transactions, nil
}

// Get all transaction inputs for a given transaction
func (repo *blockchainRepository) GetTxnInputs(txnId []byte) ([]reps.TxnInput, error) {
	var txnInputs []reps.TxnInput

	err := db.DB.Where("curr_txn_id = ?", txnId).
			Find(&txnInputs).
			Error
	
	if err != nil {
		return []reps.TxnInput{}, err
	}

	return txnInputs, nil
}

// Get all transaction outputs for a given transaction
func (repo *blockchainRepository) GetTxnOutputs(txnId []byte) ([]reps.TxnOutput, error) {
	var txnOutputs []reps.TxnOutput

	err := db.DB.Where("curr_txn_id = ?", txnId).
			Find(&txnOutputs).
			Error
	
	if err != nil {
		return []reps.TxnOutput{}, err
	}

	return txnOutputs, nil
}

// Get the first block in the blockchain
func (repo *blockchainRepository) GetGenesisBlock() (reps.Block, int, error) {
	var genesisBlock reps.Block

	res := db.DB.Where("prev_hash = ''").Find(&genesisBlock)
	if res.Error != nil {
		return reps.Block{}, -1, res.Error
	}	

	return genesisBlock, int(res.RowsAffected), nil

}

func (repo *blockchainRepository) CreateBlock(block reps.Block) error {
	// res := db.DB.Create(&block)

	// utils.PrettyPrintln("res: ", res.Value)
	if err := db.DB.Create(&block).Error; err != nil {
		return err
	}

	// if err := db.DB.Save(&block).Error; err != nil {
	// 	return err
	// }

	return nil
}

func (repo *blockchainRepository) GetBlockchain() ([]reps.Block, error) {
	var blocks []reps.Block

	if err := db.DB.Find(&blocks).Error; err != nil {
		return []reps.Block{}, nil
	}

	return blocks, nil
}
// Get the last serialized block in the blockchain
// func (repo *blockchainRepository) GetBlock() ([]byte, error) {
// 	var blockStructByte []byte

// 	err := db.DB.View(func(txn *badger.Txn) error {
// 		item, err := txn.Get([]byte(lastBlockHashKey))
// 		if err != nil {
// 			log.WithFields(log.Fields{"warning": err.Error()}).Warn("Blockchain doesn't exist.")
// 			return err
// 		}

// 		// This should give us block.Hash
// 		blockStructByte, err = item.ValueCopy(nil)
// 		if err != nil {
// 			log.Error("Error retreiving value: ", err.Error())
// 			return err
// 		}

// 		// This should give us the serialized block structure
// 		item, err = txn.Get(blockStructByte)
// 		if err != nil {
// 			log.Error("Error retreiving value: ", err.Error())
// 			return err
// 		}

// 		// Get the actual value
// 		blockStructByte, err = item.ValueCopy(nil)
// 		if err != nil {
// 			log.Error("Error retreiving value: ", err.Error())
// 			return err
// 		}

// 		return err
// 	})

// 	if err != nil {
// 		return blockStructByte, err
// 	}
// 	return blockStructByte, nil
// }

// // Add block to the blockchain, returns the last serialized block which is the one just created
// func (repo *blockchainRepository) CreateBlock(hash []byte, blockByte []byte) ([]byte, error) {
// 	errUpdate := db.DB.Update(func(txn *badger.Txn) error {
// 		errSet := txn.Set(hash, blockByte)
// 		errSet = txn.Set([]byte(lastBlockHashKey), hash)
// 		if errSet != nil {
// 			log.WithField("error", errSet.Error()).Error("Error setting entry")
// 			return errSet
// 		}

// 		return nil
// 	})

// 	if errUpdate != nil {
// 		log.Error("Error creating new blockchain")
// 		return []byte{}, errUpdate
// 	}

// 	return repo.GetBlock()
// }

// // Get all blocks in the blockchain
// func (repo *blockchainRepository) GetBlockchain() ([][]byte, error) {
// 	log.Info("Printing blockchain")
// 	blocks := make([][]byte, 0)

// 	err := db.DB.View(func(txn *badger.Txn) error {
// 		opts := badger.DefaultIteratorOptions
// 		opts.PrefetchSize = 10
// 		it := txn.NewIterator(opts)
// 		defer it.Close()
// 		for it.Rewind(); it.Valid(); it.Next() {
// 			item := it.Item()
// 			k := item.Key()
// 			err := item.Value(func(v []byte) error {
// 				if !strings.EqualFold(string(k), lastBlockHashKey) { // only want encoded block for now
// 					blocks = append(blocks, v)
// 				}
// 				return nil
// 			})
// 			if err != nil {
// 				return nil
// 			}
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return [][]byte{}, err
// 	}

// 	return blocks, nil
// }
