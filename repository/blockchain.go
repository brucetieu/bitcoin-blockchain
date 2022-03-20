package repository

import (
	"strings"

	"github.com/brucetieu/blockchain/db"
	badger "github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
)

var lastBlockHashKey = "lastBlockHash" // This gives us back block.Hash for the last block

type BlockchainRepository interface {
	CreateBlock(hash []byte, blockByte []byte) ([]byte, error)
	GetBlock() ([]byte, error)
	GetBlockchain() ([][]byte, error)
}

type blockchainRepository struct {
}

func NewBlockchainRepository() BlockchainRepository {
	return &blockchainRepository{}
}

// Get the last serialized block in the blockchain
func (repo *blockchainRepository) GetBlock() ([]byte, error) {
	var blockStructByte []byte

	err := db.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(lastBlockHashKey))
		if err != nil {
			log.WithFields(log.Fields{"warning": err.Error()}).Warn("Blockchain doesn't exist.")
			return err
		}

		// This should give us block.Hash
		blockStructByte, err = item.ValueCopy(nil)
		if err != nil {
			log.Error("Error retreiving value: ", err.Error())
			return err
		}

		// This should give us the serialized block structure
		item, err = txn.Get(blockStructByte)
		if err != nil {
			log.Error("Error retreiving value: ", err.Error())
			return err
		}

		// Get the actual value
		blockStructByte, err = item.ValueCopy(nil)
		if err != nil {
			log.Error("Error retreiving value: ", err.Error())
			return err
		}

		return err
	})

	if err != nil {
		return blockStructByte, err
	}
	return blockStructByte, nil
}

// Add block to the blockchain, returns the last serialized block which is the one just created
func (repo *blockchainRepository) CreateBlock(hash []byte, blockByte []byte) ([]byte, error) {
	errUpdate := db.DB.Update(func(txn *badger.Txn) error {
		errSet := txn.Set(hash, blockByte)
		errSet = txn.Set([]byte(lastBlockHashKey), hash)
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
				if !strings.EqualFold(string(k), lastBlockHashKey) { // only want encoded block for now
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
