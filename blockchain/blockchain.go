package blockchain

import (
	"github.com/brucetieu/blockchain/db"
	badger "github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
)

type Blockchain interface {
	AddToBlockChain(data string) (*Block, error)
	CreateBlockchain() (bool, *Blockchain)
}
type blockchain struct {
	// DB *badger.DB
	LastBlock []byte // -> the last block in serialized format
	IsGenesis string
}

func NewBlockchain() *blockchain {
	return &blockchain{}
}

func (bc *blockchain) CreateBlockchain() (bool, *Blockchain) {
	// log.Info("Creating new blockchain")
	// db, err := badger.Open(badger.DefaultOptions("/tmp/blockchain"))
	// if err != nil {
	// 	log.WithFields(log.Fields{"error": err.Error()}).Fatal("Error opening badgerdb")
	// }
	
	var blockchain Blockchain
	doesExist := false
	errUpdate := db.DB.Update(func(txn *badger.Txn) error {
		lastBlockHash, getErr := txn.Get([]byte("lastBlockSerial"))
		if getErr != nil {
			log.WithFields(log.Fields{"warning": getErr.Error()}).Warn("Blockchain doesn't exist... creating a new one.")
			
			firstBlock := CreateBlock("Start of Blockchain", []byte{})
			// err = txn.Set(firstBlock.Hash, firstBlock.StructToByte())
			errSet := txn.Set([]byte("lastBlockSerial"), firstBlock.Serialize())
			if errSet != nil {
				log.Error("Error setting entry")
			}
			blockchain.LastBlock = firstBlock.Serialize()
		} else {
			log.Info("Blockchain already exists...")
			doesExist = true
			valErr := lastBlockHash.Value(func(val []byte) error {
				blockchain.LastBlock = val

				return nil
			})
			if valErr != nil {
				log.Error("Error getting value", valErr.Error())
			}
		}
		return nil
	  })
	
	if errUpdate != nil {
		log.Error("Error creating new blockchain")
	}
	// blockchain.DB = db
	blockchain.IsGenesis = "true"
	return doesExist, &blockchain
}

// This assumes blockchain is already created
// 1. Get the last block hash in the blockchain
// 2. Create a new block using this last block hash as the previous hash
// 3. Serialize this new block, and set this block as the last block hash in the blockchain
func (bc *blockchain) AddToBlockChain(data string) (*Block, error) {
	var err error
	err = db.DB.Update(func(txn *badger.Txn) error {
		lastBlockHash, getErr := txn.Get([]byte("lastBlockSerial"))
		if getErr != nil {
			log.WithFields(log.Fields{"error": getErr.Error()}).Error("Error getting key lastBlockHash. Blockchain probably has not been created")
			return getErr
		} else {
			valErr := lastBlockHash.Value(func(val []byte) error {
				newBlock := CreateBlock(data, val)
				// err = txn.Set(newBlock.Hash, newBlock.StructToByte())
				// err = txn.Set([]byte("lastBlockHash"), newBlock.Hash)
				err = txn.Set([]byte("lastBlockSerial"), newBlock.Serialize())
				bc.LastBlock = newBlock.Serialize()
				// bc.LastBlock = newBlock.Hash
				return nil
			})
			if valErr != nil {
				log.Error("Error getting value", valErr.Error())
				return valErr
			}
		}
		return nil
	})
	if err != nil {
		log.WithField("error", err.Error()).Error("Error adding block to block chain")
		return Deserialize(bc.LastBlock), err
	}
	
	return Deserialize(bc.LastBlock), nil
}


