package repository

import (
	"github.com/brucetieu/blockchain/db"
	// "github.com/brucetieu/blockchain/representations"
	badger "github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
	
)

type BlockchainRepository interface {
	CreateBlockchain(blockByte []byte) ([]byte, error)
	GetBlockchain() ([]byte, error)
}

type blockchainRepository struct {
	// LastBlock []byte // -> the last block in serialized format
	// IsGenesis string
}

func NewBlockchainRepository() BlockchainRepository {
	return &blockchainRepository{}
}

// Get the genesis block to a block chain, if it exists
func (repo *blockchainRepository) GetBlockchain() ([]byte, error) {
	var valCopy []byte

	err := db.DB.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte("lastSerializedBlock"))
			if err != nil {
				log.WithFields(log.Fields{"warning": err.Error()}).Warn("Blockchain doesn't exist.")
				return err
			} 

			valCopy, err = item.ValueCopy(nil)
			if err != nil {
				log.Error("Error retreiving value", err.Error())
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

func (repo *blockchainRepository) CreateBlockchain(blockByte []byte) ([]byte, error) {
	var genesisBlockByte []byte

	errUpdate := db.DB.Update(func(txn *badger.Txn) error {
		errSet := txn.Set([]byte("lastSerializedBlock"), blockByte)
		if errSet != nil {
			log.Error("Error setting entry")
			return errSet
		}

		return nil
	  })
	
	if errUpdate != nil {
		log.Error("Error creating new blockchain")
		return genesisBlockByte, errUpdate
	}

	return repo.GetBlockchain()
}