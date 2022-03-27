package services

import (
	// "fmt"
	"sort"

	"github.com/brucetieu/blockchain/repository"
	reps "github.com/brucetieu/blockchain/representations"
	"github.com/brucetieu/blockchain/utils"

	log "github.com/sirupsen/logrus"
)

type BlockchainService interface {
	// AddToBlockChain(from string, to string, amount int) (*reps.Block, error)
	CreateBlockchain(to string) (*reps.Block, bool, error)
	GetBlockchain() ([]reps.Block, error)
	GetGenesisBlock() (*reps.Block, int, error)
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
	doesExist := false

	genesis, count, err := bc.GetGenesisBlock()
	if err != nil {
		return &reps.Block{}, doesExist, err
	}

	// Genesis doesn't exist, so create it.
	if count == 0 {
		coinbaseTxn := bc.transactionService.CreateCoinbaseTxn(to, "First transaction in Blockchain")
		newBlock := bc.blockService.CreateBlock([]reps.Transaction{coinbaseTxn}, []byte{})

		utils.PrettyPrintln("newBlock: ", newBlock)
		err = bc.blockchainRepo.CreateBlock(*newBlock)
		if err != nil {
			log.Error("Error creating blockchain: ", err.Error())
			return &reps.Block{}, doesExist, err
		}

		return newBlock, doesExist, nil
	}

	// Genesis does exist, so return it
	return genesis, true, nil
}

// 1. Get the last block hash in the blockchain
// 2. Create a new transaction
// 3. Create a new block with the transaction and use this last block hash as the previous hash
// 4. Set (persist) this block to the db
// func (bc *blockchainService) AddToBlockChain(from string, to string, amount int) (*reps.Block, error) {
// 	lastBlock, err := bc.blockchainRepo.GetBlock()
// 	if err != nil {
// 		log.WithField("error", err.Error()).Error("Error getting last block. Blockchain probably has not been created")
// 		return &reps.Block{}, fmt.Errorf("error getting last block. Blockchain probably has not been created: %s", err.Error())
// 	}

// 	decodedLastBlock := bc.blockAssembler.ToBlockStructure(lastBlock)

// 	// create new transaction
// 	newTxn, err := bc.transactionService.CreateTransaction(from, to, amount)
// 	if err != nil {
// 		return &reps.Block{}, err
// 	}
// 	newBlock := bc.blockService.CreateBlock([]*reps.Transaction{newTxn}, decodedLastBlock.Hash)

// 	// Prsist to db
// 	_, err = bc.blockchainRepo.CreateBlock(newBlock.Hash, bc.blockAssembler.ToBlockBytes(newBlock))
// 	if err != nil {
// 		return &reps.Block{}, err
// 	}

// 	return newBlock, err

// }

// Get all blocks in the blockchain
func (bc *blockchainService) GetBlockchain() ([]reps.Block, error) {
	// var allBlocks []*reps.Block

	blocks, err := bc.blockchainRepo.GetBlockchain()
	if err != nil {
		log.WithField("error", err.Error()).Error("Error getting all blocks in blockchain")
		return []reps.Block{}, err
	}

	utils.PrettyPrintln("blocks", blocks)

	// for _, block := range blocks {
	// 	decodedB := bc.blockAssembler.ToBlockStructure(block)
	// 	allBlocks = append(allBlocks, decodedB)
	// }

	// Ensure that genesis block is last
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].Timestamp > blocks[j].Timestamp
	})

	return blocks, nil
}

// Get the first block in the block chain.
func (bc *blockchainService) GetGenesisBlock() (*reps.Block, int, error) {
	log.Info("Getting Genesis Block...")

	// Get Genesis block from db
	genesis, count, err := bc.blockchainRepo.GetGenesisBlock()
	if err != nil {
		return &reps.Block{}, -1, err
	}

	// This means the block chain hasn't been created yet
	if count == 0 {
		return &reps.Block{}, count, nil
	}

	// Get transactions inside the block we got above
	transactions, err := bc.blockchainRepo.GetTransactions(genesis.ID)
	if err != nil {
		return &reps.Block{}, -1, err
	}

	// Get the inputs and outputs for each transaction in the block
	for i := 0; i < len(transactions); i++ {
		inputs, err := bc.blockchainRepo.GetTxnInputs(transactions[i].ID)
		if err != nil {
			return &reps.Block{}, -1, err
		}
		transactions[i].Inputs = inputs
	}
	
	for i := 0; i < len(transactions); i++ {
		outputs, err := bc.blockchainRepo.GetTxnOutputs(transactions[i].ID)
		if err != nil {
			return &reps.Block{}, -1, err
		}
		transactions[i].Outputs = outputs
	}

	// Update the genesis block with the transaction
	genesis.Transactions = transactions

	log.Info("Returned genesis block: ", utils.Pretty(genesis))
	return &genesis, count, nil
}
