package handlers

import (
	"net/http"

	reps "github.com/brucetieu/blockchain/representations"
	"github.com/brucetieu/blockchain/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type BlockchainHandler struct {
	blockchainService services.BlockchainService
	assemblerService  services.BlockAssemblerFac
}

func NewBlockchainHandler(blockchainService services.BlockchainService) *BlockchainHandler {
	return &BlockchainHandler{
		blockchainService: blockchainService,
		assemblerService:  services.BlockAssembler,
	}
}

func (bch *BlockchainHandler) CreateBlockchain(c *gin.Context) {
	log.Info("Creating Blockchain")
	// Validate input
	var input reps.CreateBlockchainInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	decodedGenesis, exists, err := bch.blockchainService.CreateBlockchain(input.To)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Format return data to be readable
	data := bch.assemblerService.ToReadableBlock(decodedGenesis)

	if exists {
		c.JSON(http.StatusOK, gin.H{"message": "Blockchain already exists."})
	} else {
		c.JSON(http.StatusCreated, gin.H{"block": data, "message": "Blockchain created."})
	}
}

func (bch *BlockchainHandler) AddToBlockchain(c *gin.Context) {
	// Validate input
	var input reps.CreateBlockInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create block and persist to db
	newBlock, err := bch.blockchainService.AddToBlockChain(input.From, input.To, input.Amount)
	if err != nil {
		log.WithField("error", err.Error()).Error("Error adding block")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Format return data to be readable
	data := bch.assemblerService.ToReadableBlock(newBlock)

	c.JSON(http.StatusCreated, gin.H{"block": data})
}

// Print out all blocks in blockchain
func (bch *BlockchainHandler) GetBlockchain(c *gin.Context) {
	blockchain, err := bch.blockchainService.GetBlockchain()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := make([]reps.ReadableBlock, 0)

	for _, block := range blockchain {
		data = append(data, bch.assemblerService.ToReadableBlock(block))
	}

	c.JSON(http.StatusOK, gin.H{"blockchain": data})

}

// Get the first block in block chain
func (bch *BlockchainHandler) GetGenesisBlock(c *gin.Context) {
	genesis, err := bch.blockchainService.GetGenesisBlock()
	if err != nil {
		log.WithField("error", err.Error()).Error("Error getting genesis block")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		formattedGenesis := bch.assemblerService.ToReadableBlock(genesis)
		c.JSON(http.StatusOK, gin.H{"genesis": formattedGenesis})
	}
}

// Get a block given a blockId
func (bch *BlockchainHandler) GetBlock(c *gin.Context) {
	blockId := c.Param("blockId")
	log.Info("Getting block with blockId: ", blockId)

	block, err := bch.blockchainService.GetBlock(blockId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"block": bch.assemblerService.ToReadableBlock(block)})
	}
}

// Get last block in blockchain. If it's a genesis, it will return it. 
func (bch *BlockchainHandler) GetLastBlock(c *gin.Context) {
	log.Info("Getting last block...")

	lastBlock, err := bch.blockchainService.GetLastBlock()
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"block": bch.assemblerService.ToReadableBlock(lastBlock)})
	}
}

