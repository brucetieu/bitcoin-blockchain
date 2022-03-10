package repository

import (
	"fmt"
	"strings"

	"github.com/brucetieu/blockchain/db"
	// "github.com/brucetieu/blockchain/representations"
	badger "github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
)

type BlockchainRepository interface {
	CreateBlock(hash []byte, blockByte []byte) ([]byte, error)
	// GetBlockchain() ([]byte, error)
	GetBlock() ([]byte, error)
	// AddBlock(hash []byte, blockByte []byte) (error)
	GetBlockchain() ([][]byte, error)
}

type blockchainRepository struct {
	// LastBlock []byte // -> the last block in serialized format
	// IsGenesis string
}

func NewBlockchainRepository() BlockchainRepository {
	return &blockchainRepository{}
}

// Get the last serialized block in the blockchain 
func (repo *blockchainRepository) GetBlock() ([]byte, error) {
	var valCopy []byte

	err := db.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lastBlockHash")) // This gives us back block.Hash
		if err != nil {
			log.WithFields(log.Fields{"warning": err.Error()}).Warn("Blockchain doesn't exist.")
			return err
		} 

		// block.Hash is a key itself in the db. Use it to get the value, which is the serialized block
		valCopy, err = item.ValueCopy(nil)
		if err != nil {
			log.Error("Error retreiving value: ", err.Error())
			return err
		}

		item, err = txn.Get(valCopy)
		if err != nil {
			log.Error("Error retreiving value: ", err.Error())
			return err
		}

		valCopy, err = item.ValueCopy(nil)
		if err != nil {
			log.Error("Error retreiving value: ", err.Error())
			return err
		}

		// err = item.Value(func(val []byte) error {
		// 	valCopy = append([]byte{}, val...)
		// 	return nil
		// })

		return err
	})
	
	if err != nil {
		return valCopy, err
	}
	return valCopy, nil
}

// Add block to the blockchain, returns the last serialized block which is the one just created
func (repo *blockchainRepository) CreateBlock(hash []byte, blockByte []byte) ([]byte, error) {
	errUpdate := db.DB.Update(func(txn *badger.Txn) error {
		errSet := txn.Set(hash, blockByte)
		errSet = txn.Set([]byte("lastBlockHash"), hash)
		if errSet != nil {
			log.WithField("error", errSet.Error()).Error("Error setting entry")
			return errSet
		}

		return nil
	  })
	
	if errUpdate != nil {
		log.Error("Error creating new blockchain")
		return []byte{}, errUpdate
	}

	return repo.GetBlock()
}

// func (repo *blockchainRepository) AddBlock(hash []byte, block []byte) (error) {
// 	errUpdate := db.DB.Update(func(txn *badger.Txn) error {
// 		errSet := txn.Set(hash, block)
// 		errSet = txn.Set([]byte("lastBlockHash"), hash)
// 		if errSet != nil {
// 			log.WithField("error", errSet.Error()).Error("Error setting entry")
// 			return errSet
// 		}
// 		return nil
// 	})

// 	if errUpdate != nil {
// 		log.WithField("error", errUpdate.Error()).Error("Error updating db")
// 		return errUpdate
// 	}
// 	return nil
// }

// Get all blocks in the blockchain
func (repo *blockchainRepository) GetBlockchain() ([][]byte, error) {
	log.Info("Printing blockchain")
	blocks := make([][]byte, 0)

	err := db.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				fmt.Printf("key=%s, value=%s\n", k, v)
				if !strings.EqualFold(string(k), "lastBlockHash") {
					blocks = append(blocks, v)
				}
				return nil
			})
			if err != nil {
				return nil
			}
		}
		return nil
	})
	if err != nil {
		return [][]byte{}, err
	}
	return blocks, nil
}

func (repo *blockchainRepository) GetGenesisBlock() {
	log.Info("Getting genesis block...")
	
}