package services

import (
	// "fmt"
	"fmt"
	"sort"

	"github.com/brucetieu/blockchain/repository"
	reps "github.com/brucetieu/blockchain/representations"
	"github.com/brucetieu/blockchain/utils"

	log "github.com/sirupsen/logrus"
)

type BlockchainService interface {
	AddToBlockChain(from string, to string, amount int) (reps.Block, error)
	CreateBlockchain(to string) (*reps.Block, bool, error)
	GetBlockchain() ([]reps.Block, error)
	GetGenesisBlock() (*reps.Block, error)
	GetBlock(blockId string) (reps.Block, error)
}

type blockchainService struct {
	blockchainRepo     repository.BlockchainRepository
	blockService       BlockService
	transactionService TransactionService
	blockAssembler     BlockAssemblerFac
}

func NewBlockchainService(blockchainRepo repository.BlockchainRepository,
	blockService BlockService, transactionService TransactionService) BlockchainService {
	return &blockchainService{
		blockchainRepo:     blockchainRepo,
		blockService:       blockService,
		transactionService: transactionService,
		blockAssembler:     BlockAssembler,
	}
}

func (bc *blockchainService) CreateBlockchain(to string) (*reps.Block, bool, error) {
	// Try to get genesis block
	genesis, err := bc.GetGenesisBlock()
	if err != nil {
		log.Info("Genesis doesn't exist, so creating it now...")
		coinbaseTxn := bc.transactionService.CreateCoinbaseTxn(to, "First transaction in Blockchain")
		newBlock, err := bc.blockService.CreateBlock([]reps.Transaction{coinbaseTxn}, []byte{})

		// Persist
		// err = bc.blockchainRepo.CreateBlock(newBlock)
		if err != nil {
			log.Error("Error creating blockchain: ", err.Error())
			return &reps.Block{}, false, err
		}

		return &newBlock, false, nil
	}

	// // Genesis doesn't exist, so create it.
	// if count == 0 {
	// 	log.Info("Genesis doesn't exist, so creating it now...")
	// 	coinbaseTxn := bc.transactionService.CreateCoinbaseTxn(to, "First transaction in Blockchain")
	// 	newBlock, err := bc.blockService.CreateBlock([]reps.Transaction{coinbaseTxn}, []byte{})

	// 	// Persist
	// 	// err = bc.blockchainRepo.CreateBlock(newBlock)
	// 	if err != nil {
	// 		log.Error("Error creating blockchain: ", err.Error())
	// 		return &reps.Block{}, false, err
	// 	}

	// 	return &newBlock, false, nil
	// }

	// Genesis does exist, so return it
	return genesis, true, nil
}

func (bc *blockchainService) AddToBlockChain(from string, to string, amount int) (reps.Block, error) {
	// Check genesis exists first before adding block
	// _, err := bc.GetGenesisBlock()
	// if err != nil {
	// 	return reps.Block{}, err
	// }

	// if count == 0 {
	// 	return reps.Block{}, fmt.Errorf("Cannot add to blockchain without Genesis block")
	// }

	// Check if there is at least a genesis block in the blockchain
	lastBlock, err := bc.blockchainRepo.GetLastBlock()
	if err != nil {
		errMsg := fmt.Errorf("%s, cannot create a block without genesis", err.Error())
		return reps.Block{}, errMsg
	}

	// Create a new transaction
	newTxn, err := bc.transactionService.CreateTransaction(from, to, amount)

	// Create a new block and persist
	newBlock, err := bc.blockService.CreateBlock([]reps.Transaction{newTxn}, lastBlock.Hash)
	if err != nil {
		return reps.Block{}, err
	}

	return newBlock, nil
}

// Get all blocks in the blockchain
func (bc *blockchainService) GetBlockchain() ([]reps.Block, error) {
	blocks, err := bc.blockchainRepo.GetBlockchain()
	if err != nil {
		log.WithField("error", err.Error()).Error("Error getting all blocks in blockchain")
		return []reps.Block{}, err
	}

	// Ensure that genesis block is last
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].Timestamp > blocks[j].Timestamp
	})

	return blocks, nil
}

// Get the first block in the block chain.
func (bc *blockchainService) GetGenesisBlock() (*reps.Block, error) {
	log.Info("Getting Genesis Block...")

	// Get Genesis block from db
	genesis, err := bc.blockchainRepo.GetGenesisBlock()
	if err != nil {
		errMsg := fmt.Errorf("%s, genesis does not exist", err.Error())
		return &reps.Block{}, errMsg
	}

	log.Info("Returned genesis block: ", utils.Pretty(genesis))
	return &genesis, nil
}

func (bc *blockchainService) GetBlock(blockId string) (reps.Block, error) {
	block, err := bc.blockchainRepo.GetBlockById(blockId)
	if err != nil {
		log.Error("error getting block: ", err.Error())
		errMsg := fmt.Errorf("%s, id: %s", err.Error(), blockId)
		return reps.Block{}, errMsg
	}

	// if count == 0 {
	// 	log.Error("no block found wth id ", blockId)
	// 	return reps.Block{}, fmt.Errorf("no block found with id %s", blockId)
	// }

	return block, nil
}