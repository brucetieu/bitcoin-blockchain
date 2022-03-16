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
	assemblerService services.BlockAssemblerFac
}

func NewBlockchainHandler(blockchainService services.BlockchainService) *BlockchainHandler {
	return &BlockchainHandler{
		blockchainService: blockchainService,
		assemblerService: services.BlockAssembler,
	}
}

func (bch *BlockchainHandler) CreateBlockchain(c *gin.Context) {
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
	data := bch.assemblerService.ToBlockMap(decodedGenesis)

	if exists {
		c.JSON(http.StatusOK, gin.H{"message": "Blockchain already exists."})
	} else {
		c.JSON(http.StatusCreated, gin.H{"data": data, "message": "Blockchain created."})
	}
}

// func (bch *BlockchainHandler) AddToBlockchain(c *gin.Context) {
// 	// Validate input
// 	var input reps.CreateBlockInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Create block and persist to db
// 	newBlock, err := bch.blockchainService.AddToBlockChain(input.Data)
// 	if err != nil {
// 		log.WithField("error", err.Error()).Error("Error adding block")
// 		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Format return data to be readable
// 	data := assembler.ToBlockMap(newBlock)

// 	c.JSON(http.StatusCreated, gin.H{"data": data})
// }

// Print out all blocks in blockchain
func (bch *BlockchainHandler) GetBlockchain(c *gin.Context) {
	blockchain, err := bch.blockchainService.GetBlockchain()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := make([]map[string]interface{}, 0)

	for _, block := range blockchain {
		data = append(data, bch.assemblerService.ToBlockMap(block))
	}

	c.JSON(http.StatusOK, gin.H{"data": data})

}

// Get the first block in block chain
func (bch *BlockchainHandler) GetGenesisBlock(c *gin.Context) {
	genesis, err := bch.blockchainService.GetGenesisBlock()
	if err != nil {
		log.WithField("error", err.Error()).Error("Error getting genesis block")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	formattedGenesis := bch.assemblerService.ToBlockMap(genesis)

	c.JSON(http.StatusOK, gin.H{"data": formattedGenesis})
}
