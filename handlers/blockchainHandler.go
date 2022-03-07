package handlers

import (
	// "fmt"
	"fmt"
	"net/http"

	// "github.com/brucetieu/blockchain/blockchain"
	"github.com/brucetieu/blockchain/services"
	"github.com/gin-gonic/gin"
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
	prevHash := fmt.Sprintf("%x", decodedGenesis.PrevHash)
	hash := fmt.Sprintf("%x", decodedGenesis.Hash)
	data := make(map[string]interface{})
	data["timestamp"] = decodedGenesis.Timestamp
	data["prevHash"] = prevHash
	data["hash"] = hash
	data["nounce"] = decodedGenesis.Nounce
	
	if exists {
		c.JSON(http.StatusOK, gin.H{"data": data, "message": "Blockchain already exists."})
	} else {
		c.JSON(http.StatusCreated, gin.H{"data": data, "message": "Blockchain created."})
	}
}
// type CreateBlockInput struct {
// 	Data string `json:"data" binding:"required"`
// }



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