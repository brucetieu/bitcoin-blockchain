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
	CreateBlockchain(address string) (reps.Block, bool, error)
	GetBlockchain() ([]reps.Block, error)
	GetGenesisBlock() (reps.Block, error)
	GetBlock(blockId string) (reps.Block, error)
	GetLastBlock() (reps.Block, error)
}

type blockchainService struct {
	blockchainRepo     repository.BlockchainRepository
	blockService       BlockService
	transactionService TransactionService
	walletService WalletService
	blockAssembler     BlockAssemblerFac
}

func NewBlockchainService(blockchainRepo repository.BlockchainRepository,
	blockService BlockService, transactionService TransactionService, walletService WalletService) BlockchainService {
	return &blockchainService{
		blockchainRepo:     blockchainRepo,
		blockService:       blockService,
		transactionService: transactionService,
		walletService: walletService,
		blockAssembler:     BlockAssembler,
	}
}

// address is wallet address
func (bc *blockchainService) CreateBlockchain(address string) (reps.Block, bool, error) {
	// validate address
	addressValid, err := bc.walletService.ValidateAddress(address)
	if err != nil {
		log.Error(err.Error())
		return reps.Block{}, false, err
	}

	if !addressValid {
		log.Errorf("error: address of %s is not valid", address)
		return reps.Block{}, false, fmt.Errorf("error: address of %s is not valid", address)
	}
	// if !bc.walletService.ValidateAddress(address) {
	// 	log.Errorf("error: address of %s is not valid", address)
	// 	return reps.Block{}, false, fmt.Errorf("error: address of %s is not valid", address)
	// }

	// Try to get genesis block
	genesis, err := bc.GetGenesisBlock()
	if err != nil {
		log.Info("Genesis doesn't exist, so creating it now...")
		coinbaseTxn := bc.transactionService.CreateCoinbaseTxn(address, "First transaction in Blockchain")
		newBlock, err := bc.blockService.CreateBlock([]reps.Transaction{coinbaseTxn}, []byte{})

		// Persist
		if err != nil {
			log.Error("Error creating blockchain: ", err.Error())
			return reps.Block{}, false, err
		}

		return newBlock, false, nil
	}

	// Genesis does exist, so return it
	return genesis, true, nil
}

func (bc *blockchainService) AddToBlockChain(from string, to string, amount int) (reps.Block, error) {
	// Validate from and to are valid addresses
	addressValid, err := bc.walletService.ValidateAddress(from)
	if err != nil {
		log.Error(err.Error())
		return reps.Block{}, err
	}
	if !addressValid {
		log.Errorf("error: address of %s is not valid", from)
		return reps.Block{}, fmt.Errorf("error: address of %s is not valid", from)
	}

	addressValid, err = bc.walletService.ValidateAddress(to)
	if err != nil {
		log.Error(err.Error())
		return reps.Block{}, err
	}
	if !addressValid {
		log.Errorf("error: address of %s is not valid", to)
		return reps.Block{}, fmt.Errorf("error: address of %s is not valid", to)
	}

	// Check if there is at least a genesis block in the blockchain
	lastBlock, err := bc.blockchainRepo.GetLastBlock()
	if err != nil {
		errMsg := fmt.Errorf("%s, cannot create a block without genesis", err.Error())
		return reps.Block{}, errMsg
	}

	// Create a new transaction. 
	newTxn, err := bc.transactionService.CreateTransaction(from, to, amount)
	if err != nil {
		return reps.Block{}, err
	}

	// Verify the signatures on transaction inputs
	txns := []reps.Transaction{newTxn}
	for _, txn := range txns {
		if !bc.transactionService.VerifyTransaction(txn) {
			log.WithField("error", err.Error()).Error("error: invalid transaction")
			return reps.Block{}, err
		}
	}

	// Create a new block with new transaction and persist
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
func (bc *blockchainService) GetGenesisBlock() (reps.Block, error) {
	log.Info("Getting Genesis Block...")

	// Get Genesis block from db
	genesis, err := bc.blockchainRepo.GetGenesisBlock()
	if err != nil {
		errMsg := fmt.Errorf("%s, genesis does not exist", err.Error())
		return reps.Block{}, errMsg
	}

	log.Info("Returned genesis block: ", utils.Pretty(genesis))
	return genesis, nil
}

func (bc *blockchainService) GetBlock(blockId string) (reps.Block, error) {
	block, err := bc.blockchainRepo.GetBlockById(blockId)
	if err != nil {
		log.Error("error getting block: ", err.Error())
		errMsg := fmt.Errorf("%s, id: %s", err.Error(), blockId)
		return reps.Block{}, errMsg
	}

	return block, nil
}

func (bc *blockchainService) GetLastBlock() (reps.Block, error) {
	lastBlock, err := bc.blockchainRepo.GetLastBlock()
	if err != nil {
		errMsg := fmt.Errorf("%s, genesis does not exist", err.Error())
		return reps.Block{}, errMsg
	}

	return lastBlock, nil
}