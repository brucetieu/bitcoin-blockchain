package services

import (
	"fmt"

	"github.com/brucetieu/blockchain/repository"
	reps "github.com/brucetieu/blockchain/representations"

	log "github.com/sirupsen/logrus"
)

type BlockchainService interface {
	// AddToBlockChain(data string) (*reps.Block, error)
	CreateBlockchain(to string) (*reps.Block, bool, error)
	GetBlockchain() ([]*reps.Block, error)
	GetGenesisBlock() (*reps.Block, error)
}

type blockchainService struct {
	blockChainRepo repository.BlockchainRepository
	blockService   BlockService
	transactionService TransactionService
	blockAssembler BlockAssemblerFac
}

func NewBlockchainService(blockchainRepo repository.BlockchainRepository,
	blockService BlockService, transactionService TransactionService) BlockchainService {
	return &blockchainService{
		blockChainRepo: blockchainRepo,
		blockService:   blockService,
		transactionService: transactionService,
		blockAssembler: BlockAssembler,
	}
}

func (bc *blockchainService) CreateBlockchain(to string) (*reps.Block, bool, error) {
	doesExist := false

	// Attempt to get genesis block. If it doesn't exist, then create a block.
	genesisBlock, err := bc.GetGenesisBlock()
	if err != nil {
		log.Error("Error getting genesis block: ", err.Error())
		coinbaseTxn := bc.transactionService.CreateCoinbaseTxn(to, "First transaction in Blockchain")
		newBlock := bc.blockService.CreateBlock([]*reps.Transaction{coinbaseTxn}, []byte{})
		serializedBlock := bc.blockAssembler.ToBlockBytes(newBlock)

		// Need: <k,v> = <block.Hash, serialized(block)>
		//       <k,v> = <"lastBlock", block.Hash>
		newBlockchain, err := bc.blockChainRepo.CreateBlock(newBlock.Hash, serializedBlock)
		if err != nil {
			log.Error("Error creating blockchain: ", err.Error())
			return nil, doesExist, err
		}

		return bc.blockAssembler.ToBlockStructure(newBlockchain), doesExist, err
	}

	// Genesis block always has an empty previous hash
	if len(genesisBlock.PrevHash) == 0 {
		doesExist = true
	}

	return genesisBlock, doesExist, nil
}

// 1. Get the last block hash in the blockchain
// 2. Create a new block using this last block hash as the previous hash
// 3. Serialize this new block, and set this block as the last block hash in the blockchain
// func (bc *blockchainService) AddToBlockChain(data string) (*reps.Block, error) {
// 	lastBlock, err := bc.blockChainRepo.GetBlock()
// 	if err != nil {
// 		log.WithField("error", err.Error()).Error("Error getting last block. Blockchain probably has not been created")
// 		return &reps.Block{}, fmt.Errorf("error getting last block. Blockchain probably has not been created: %s", err.Error())
// 	}

// 	decodedLastBlock := bc.blockService.Deserialize(lastBlock)
// 	newBlock := bc.blockService.CreateBlock(data, decodedLastBlock.Hash)

// 	// Prsist to db
// 	_, err = bc.blockChainRepo.CreateBlock(newBlock.Hash, bc.blockService.Serialize(newBlock))
// 	if err != nil {
// 		return &reps.Block{}, err
// 	}

// 	return newBlock, err

// }

// Get all blocks in the blockchain
func (bc *blockchainService) GetBlockchain() ([]*reps.Block, error) {
	var allBlocks []*reps.Block
	blocks, err := bc.blockChainRepo.GetBlockchain()
	if err != nil {
		log.WithField("error", err.Error()).Error("Error getting all blocks in blockchain")
		return allBlocks, err
	}
	for _, block := range blocks {
		decodedB := bc.blockAssembler.ToBlockStructure(block)
		allBlocks = append(allBlocks, decodedB)
	}

	return allBlocks, nil
}

// Get the first block in the block chain.
func (bc *blockchainService) GetGenesisBlock() (*reps.Block, error) {
	genesis := &reps.Block{}
	blocks, err := bc.blockChainRepo.GetBlockchain()
	if err != nil {
		log.WithField("error", err.Error()).Error("Error getting all blocks in blockchain")
		return &reps.Block{}, err
	}

	for _, block := range blocks {
		decodedB := bc.blockAssembler.ToBlockStructure(block)
		if len(decodedB.PrevHash) == 0 {
			genesis = decodedB
			break
		}
	}

	if len(genesis.Hash) == 0 && len(genesis.PrevHash) == 0 {
		return genesis, fmt.Errorf("no genesis block exists, please create a blockchain first")
	}

	return genesis, nil
}
