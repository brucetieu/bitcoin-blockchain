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
	GetGenesisBlock() (*reps.Block, error)
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

	// Attempt to get genesis block. If it doesn't exist, then create a block.
	genesisBlock, err := bc.GetGenesisBlock()
	if err != nil {
		log.Error("Error getting genesis block: ", err.Error())
		coinbaseTxn := bc.transactionService.CreateCoinbaseTxn(to, "First transaction in Blockchain")
		newBlock := bc.blockService.CreateBlock([]reps.Transaction{coinbaseTxn}, []byte{})
		// serializedBlock := bc.blockAssembler.ToBlockBytes(newBlock)

		utils.PrettyPrintln("newBlock", newBlock)
		// for _, txn := range newBlock.Transactions {

			// for i := 0; i < len(txn.Inputs); i++ {
			// 	utils.PrettyPrintln("txnInput: ", txn.Inputs[i])
			// 	txn.Inputs[i].TxnID = txn.TxnID

			// 	err := bc.blockchainRepo.CreateTxnInput(txn.Inputs[i])
			// 	if err != nil {
			// 		log.Error("Error creating txnInput: ", err.Error())
			// 		return &reps.Block{}, doesExist, err
			// 	}
			// }

			// for _, txnOutput := range txn.Outputs {
			// 	err := bc.blockchainRepo.CreateTxnOutput(txnOutput)
			// 	if err != nil {
			// 		log.Error("Error creating txnOutput: ", err.Error())
			// 		return &reps.Block{}, doesExist, err
			// 	}
			// }
			// temp := make([]reps.Transaction, 0)
			// t := reps.Transaction{
			// 	TxnID: newBlock.Transactions[0].TxnID,
			// 	BlockID: newBlock.Transactions[0].BlockID,
			// }
			// temp = append(temp, t)
			// err := bc.blockchainRepo.CreateTransaction(newBlock.Transactions)
			// if err != nil {
			// 	log.Error("Error creating transaction", err.Error())
			// 	return &reps.Block{}, doesExist, err
			// }
		// }
		// Need: <k,v> = <block.Hash, serialized(block)>
		//       <k,v> = <"lastBlock", block.Hash>
		err = bc.blockchainRepo.CreateBlock(*newBlock)
		if err != nil {
			log.Error("Error creating blockchain: ", err.Error())
			return &reps.Block{}, doesExist, err
		}

		return newBlock, doesExist, err
	}

	// Genesis block always has an empty previous hash
	if len(genesisBlock.PrevHash) == 0 {
		doesExist = true
	}

	return genesisBlock, doesExist, nil
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
func (bc *blockchainService) GetGenesisBlock() (*reps.Block, error) {
	log.Info("Getting Genesis Block...")
	genesis, err := bc.blockchainRepo.GetGenesisBlock()
	if err != nil {
		return &reps.Block{}, err
	}

	return &genesis, nil
	// genesis := &reps.Block{}
	// blocks, err := bc.blockchainRepo.GetBlockchain()
	// if err != nil {
	// 	log.WithField("error", err.Error()).Error("Error getting all blocks in blockchain")
	// 	return &reps.Block{}, err
	// }

	// for _, block := range blocks {
	// 	decodedB := bc.blockAssembler.ToBlockStructure(block)
	// 	if len(decodedB.PrevHash) == 0 {
	// 		genesis = decodedB
	// 		break
	// 	}
	// }

	// if len(genesis.Hash) == 0 && len(genesis.PrevHash) == 0 {
	// 	return genesis, fmt.Errorf("no genesis block exists, please create a blockchain first")
	// }

	// return genesis, nil
}
