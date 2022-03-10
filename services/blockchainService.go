package services

import (
	// "github.com/brucetieu/blockchain/db"
	"fmt"

	"github.com/brucetieu/blockchain/repository"
	"github.com/brucetieu/blockchain/representations"

	// "github.com/brucetieu/blockchain/representations"
	// badger "github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
)

type BlockchainService interface {
	AddToBlockChain(data string) (*representations.Block, error)
	CreateBlockchain() (*representations.Block, bool, error)
	GetBlockchain() ([]*representations.Block, error)
	GetGenesisBlock() (*representations.Block, error)
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

	// Attempt to get genesis block. If it doesn't exist, then create a block.
	genesisBlock, err := bc.GetGenesisBlock()
	if err != nil {
		log.Error("Error getting genesis block: ", err.Error())
		newBlock := bc.blockService.CreateBlock("Start of Blockchain", []byte{})
		serializedBlock := bc.blockService.Serialize(newBlock)
	
		// Need: <k,v> = <block.Hash, serialized(block)>
		//       <k,v> = <"lastBlock", block.Hash>
		newBlockchain, err := bc.blockChainRepo.CreateBlock(newBlock.Hash, serializedBlock)
		if err != nil {
			log.Error("Error creating blockchain: ", err.Error())
			return nil, doesExist, err
		}
	
		return bc.blockService.Deserialize(newBlockchain), doesExist, err
		// return &representations.Block{}, doesExist, err
	}

	// Genesis block always has an empty previous hash
	if len(genesisBlock.PrevHash) == 0 {
		doesExist = true
	} 

	return genesisBlock, doesExist, nil
	// newBlock := bc.blockService.CreateBlock("Start of Blockchain", []byte{})
	// serializedBlock := bc.blockService.Serialize(newBlock)

	// newBlockchain, err := bc.blockChainRepo.CreateBlockchain(serializedBlock)
	// if err != nil {
	// 	log.Error("Error creating blockchain", err.Error())
	// 	return nil, doesExist, err
	// }

	// return bc.blockService.Deserialize(newBlockchain), doesExist, err
}

// 1. Get the last block hash in the blockchain
// 2. Create a new block using this last block hash as the previous hash
// 3. Serialize this new block, and set this block as the last block hash in the blockchain
func (bc *blockchainService) AddToBlockChain(data string) (*representations.Block, error) {
	lastBlock, err := bc.blockChainRepo.GetBlock()
	if err != nil {
		log.WithField("error", err.Error()).Error("Error getting last block. Blockchain probably has not been created")
		return &representations.Block{}, fmt.Errorf("Error getting last block. Blockchain probably has not been created: %s", err.Error())
	}

	decodedLastBlock := bc.blockService.Deserialize(lastBlock)
	newBlock := bc.blockService.CreateBlock(data, decodedLastBlock.Hash)

	// Prsist to db
	_, err = bc.blockChainRepo.CreateBlock(newBlock.Hash, bc.blockService.Serialize(newBlock))
	if err != nil {
		return &representations.Block{}, err
	}

	return newBlock, err
	
	

	// var err error
	// err = db.DB.Update(func(txn *badger.Txn) error {
	// 	lastBlockHash, getErr := txn.Get([]byte("lastBlockSerial"))
	// 	if getErr != nil {
	// 		log.WithFields(log.Fields{"error": getErr.Error()}).Error("Error getting key lastBlockHash. Blockchain probably has not been created")
	// 		return getErr
	// 	} else {
	// 		valErr := lastBlockHash.Value(func(val []byte) error {
	// 			newBlock := CreateBlock(data, val)
	// 			// err = txn.Set(newBlock.Hash, newBlock.StructToByte())
	// 			// err = txn.Set([]byte("lastBlockHash"), newBlock.Hash)
	// 			err = txn.Set([]byte("lastBlockSerial"), newBlock.Serialize())
	// 			bc.LastBlock = newBlock.Serialize()
	// 			// bc.LastBlock = newBlock.Hash
	// 			return nil
	// 		})
	// 		if valErr != nil {
	// 			log.Error("Error getting value", valErr.Error())
	// 			return valErr
	// 		}
	// 	}
	// 	return nil
	// })
	// if err != nil {
	// 	log.WithField("error", err.Error()).Error("Error adding block to block chain")
	// 	return Deserialize(bc.LastBlock), err
	// }
	
	// return Deserialize(bc.LastBlock), nil
}

func (bc *blockchainService) GetBlockchain() ([]*representations.Block, error) {
	var allBlocks []*representations.Block
	blocks, err := bc.blockChainRepo.GetBlockchain()
	if err != nil {
		log.WithField("error", err.Error()).Error("Error getting all blocks in blockchain")
		return allBlocks, err
	}
	for _, block := range blocks {
		decodedB := bc.blockService.Deserialize(block)
		allBlocks = append(allBlocks, decodedB)
	}

	return allBlocks, nil
}

// Get the first block in the block chain.
func (bc *blockchainService) GetGenesisBlock() (*representations.Block, error) {
	genesis := &representations.Block{}
	blocks, err := bc.blockChainRepo.GetBlockchain()
	if err != nil {
		log.WithField("error", err.Error()).Error("Error getting all blocks in blockchain")
		return &representations.Block{}, err
	}

	for _, block := range blocks {
		decodedB := bc.blockService.Deserialize(block)
		if len(decodedB.PrevHash) == 0 {
			genesis = decodedB
			break
		}
	}

	if len(genesis.Hash) == 0 && len(genesis.PrevHash) == 0 {
		return genesis, fmt.Errorf("No genesis block exists, please create a blockchain first")
	}

	return genesis, nil
}

