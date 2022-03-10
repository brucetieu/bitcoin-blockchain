package handlers

import (
	// "encoding/base64"
	// "fmt"
	"net/http"

	// "github.com/brucetieu/blockchain/blockchain"
	"github.com/brucetieu/blockchain/assembler"
	reps "github.com/brucetieu/blockchain/representations"
	"github.com/brucetieu/blockchain/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
) 

type BlockchainHandler struct {
	blockchainService services.BlockchainService
}

func NewBlockchainHandler(blockchainService services.BlockchainService) *BlockchainHandler{
	return &BlockchainHandler{
		blockchainService: blockchainService,
	}
}

func (bch *BlockchainHandler) CreateBlockchain(c *gin.Context) {
	decodedGenesis, exists, err := bch.blockchainService.CreateBlockchain()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Format return data to be readable
	data := assembler.ToBlockMap(decodedGenesis)
	
	if exists {
		c.JSON(http.StatusOK, gin.H{"message": "Blockchain already exists."})
	} else {
		c.JSON(http.StatusCreated, gin.H{"data": data, "message": "Blockchain created."})
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
	newBlock, err := bch.blockchainService.AddToBlockChain(input.Data)
	if err != nil {
		log.WithField("error", err.Error()).Error("Error adding block")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Format return data to be readable
	data := assembler.ToBlockMap(newBlock)

	c.JSON(http.StatusCreated, gin.H{"data": data})
}

// Print out all blocks in blockchain
func (bch *BlockchainHandler) GetBlockchain(c *gin.Context) {
	blockchain, err := bch.blockchainService.GetBlockchain()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := make([]map[string]interface{}, 0)

	for _, block := range blockchain {
		data = append(data, assembler.ToBlockMap(block))
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

	formattedGenesis := assembler.ToBlockMap(genesis)

	c.JSON(http.StatusOK, gin.H{"data": formattedGenesis})
}


// func CreateBlockchain(c *gin.Context) {
// 	doesExist, bc := blockchain.NewBlockchain()
// 	block := blockchain.Deserialize(bc.LastBlock)
// 	baseHash := fmt.Sprintf("%x", block.Hash)
// 	basePrevHash := fmt.Sprintf("%x", block.PrevHash)
	
// 	data := make(map[string]interface{})
// 	data["timestamp"] = block.Timestamp
// 	data["baseHash"] = baseHash
// 	data["basePrevHash"] = basePrevHash
// 	data["isGenesis"] = bc.IsGenesis

// 	defer bc.DB.Close()
// 	if doesExist {
// 		c.JSON(http.StatusOK, gin.H{"data": data, "message": "Blockchain already exists."})
// 	} else {
// 		c.JSON(http.StatusCreated, gin.H{"data": data, "message": "Blockchain has been created."})
// 	}

// }

// func AddBlock(c *gin.Context) {
// 	// Validate input
// 	var input CreateBlockInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Create block
// 	block := 
// }