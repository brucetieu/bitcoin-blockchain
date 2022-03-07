package services

import (
	// "github.com/brucetieu/blockchain/db"
	"github.com/brucetieu/blockchain/repository"
	"github.com/brucetieu/blockchain/representations"
	// "github.com/brucetieu/blockchain/representations"
	// badger "github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
)

type BlockchainService interface {
	// AddToBlockChain(data string) (*representations.Block, error)
	CreateBlockchain() (*representations.Block, bool, error)
}

type blockchainService struct {
	blockChainRepo repository.BlockchainRepository
	blockService BlockService
	// DB *badger.DB
	// LastBlock []byte // -> the last block in serialized format
	// IsGenesis string
}

func NewBlockchainService(blockchainRepo repository.BlockchainRepository,
	blockService BlockService) BlockchainService {
	return &blockchainService{
		blockChainRepo: blockchainRepo,
		blockService: blockService,
	}
}

func (bc *blockchainService) CreateBlockchain() (*representations.Block, bool, error) {
	doesExist := false
	existingBlockchain, err := bc.blockChainRepo.GetBlockchain()
	decodedExistingBlock := bc.blockService.Deserialize(existingBlockchain)
	if err != nil {
		log.Error("Error getting blockchain", err.Error())
		return decodedExistingBlock, doesExist, err
	}

	if len(existingBlockchain) > 0 {
		doesExist = true
		return decodedExistingBlock, doesExist, err
	} 

	newBlock := bc.blockService.CreateBlock("Start of Blockchain", []byte{})
	serializedBlock := bc.blockService.Serialize(newBlock)

	newBlockchain, err := bc.blockChainRepo.CreateBlockchain(serializedBlock)
	if err != nil {
		log.Error("Error creating blockchain", err.Error())
		return nil, doesExist, err
	}

	return bc.blockService.Deserialize(newBlockchain), doesExist, err
	// return bc.blockChainRepo.CreateBlockchain(serializedBlock), doesExist
	// doesExist := false
	// errUpdate := db.DB.Update(func(txn *badger.Txn) error {
	// 	lastBlockHash, getErr := txn.Get([]byte("lastBlockSerial"))
	// 	if getErr != nil {
	// 		log.WithFields(log.Fields{"warning": getErr.Error()}).Warn("Blockchain doesn't exist... creating a new one.")
			
	// 		firstBlock := CreateBlock("Start of Blockchain", []byte{})
	// 		// err = txn.Set(firstBlock.Hash, firstBlock.StructToByte())
	// 		errSet := txn.Set([]byte("lastBlockSerial"), firstBlock.Serialize())
	// 		if errSet != nil {
	// 			log.Error("Error setting entry")
	// 		}
	// 		blockchain.LastBlock = firstBlock.Serialize()
	// 	} else {
	// 		log.Info("Blockchain already exists...")
	// 		doesExist = true
	// 		valErr := lastBlockHash.Value(func(val []byte) error {
	// 			blockchain.LastBlock = val

	// 			return nil
	// 		})
	// 		if valErr != nil {
	// 			log.Error("Error getting value", valErr.Error())
	// 		}
	// 	}
	// 	return nil
	//   })
	
	// if errUpdate != nil {
	// 	log.Error("Error creating new blockchain")
	// }
	// // blockchain.DB = db
	// blockchain.IsGenesis = "true"
	// return doesExist, &blockchain
}

// This assumes blockchain is already created
// 1. Get the last block hash in the blockchain
// 2. Create a new block using this last block hash as the previous hash
// 3. Serialize this new block, and set this block as the last block hash in the blockchain
// func (bc *blockchainService) AddToBlockChain(data string) (*representations.Block, error) {
// 	var err error
// 	err = db.DB.Update(func(txn *badger.Txn) error {
// 		lastBlockHash, getErr := txn.Get([]byte("lastBlockSerial"))
// 		if getErr != nil {
// 			log.WithFields(log.Fields{"error": getErr.Error()}).Error("Error getting key lastBlockHash. Blockchain probably has not been created")
// 			return getErr
// 		} else {
// 			valErr := lastBlockHash.Value(func(val []byte) error {
// 				newBlock := CreateBlock(data, val)
// 				// err = txn.Set(newBlock.Hash, newBlock.StructToByte())
// 				// err = txn.Set([]byte("lastBlockHash"), newBlock.Hash)
// 				err = txn.Set([]byte("lastBlockSerial"), newBlock.Serialize())
// 				bc.LastBlock = newBlock.Serialize()
// 				// bc.LastBlock = newBlock.Hash
// 				return nil
// 			})
// 			if valErr != nil {
// 				log.Error("Error getting value", valErr.Error())
// 				return valErr
// 			}
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		log.WithField("error", err.Error()).Error("Error adding block to block chain")
// 		return Deserialize(bc.LastBlock), err
// 	}
	
// 	return Deserialize(bc.LastBlock), nil
// }


