package handlers

import (
	"net/http"

	reps "github.com/brucetieu/blockchain/representations"
	"github.com/brucetieu/blockchain/services"
	"github.com/brucetieu/blockchain/utils"
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

func (bch *BlockchainHandler) BlockchainHome(ctx *gin.Context) {
	log.Info("Checking if blockchain is up...")
	ctx.JSON(http.StatusOK, "Blockchain healthy")
}

// CreateBlockchain ... Create the blockchain
// @Summary      Create the blockchain
// @Description  Create a blockchain by mining the genesis block
// @Tags         Blocks
// @Param        BlockchainInput  body      representations.CreateBlockchainInput  true  "Create Blockchain"
// @Success      201              {object}  representations.ReadableBlock
// @Success      200              {object}  representations.ReadableBlock
// @Failure      404              {object}  HTTPError
// @Router       /blockchain [post]
func (bch *BlockchainHandler) CreateBlockchain(ctx *gin.Context) {
	log.Info("Creating Blockchain")

	// Validate input
	var input reps.CreateBlockchainInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		log.WithField("error", err.Error()).Error("Error validating input: ", utils.Pretty(input))
		NewError(ctx, http.StatusBadRequest, err)
		return
	}

	// Create the genesis if it doesn't exist. Otherwise return a message that blockchain already exists
	decodedGenesis, exists, err := bch.blockchainService.CreateBlockchain(input.To)
	if err != nil {
		log.WithField("error", err.Error()).Error("Error creating blockchain")
		NewError(ctx, http.StatusNotFound, err)
		return
	}

	// Format return data to be readable
	data := bch.assemblerService.ToReadableBlock(decodedGenesis)

	if exists {
		ctx.JSON(http.StatusOK, gin.H{"message": "Blockchain already exists."})
	} else {
		ctx.JSON(http.StatusCreated, gin.H{"block": data, "message": "Blockchain created."})
	}
}

// AddToBlockchain ... Mine or add a block to the blockchain
// @Summary      Add a block
// @Description  Add a block to the end of the blockchain
// @Tags         Blocks
// @Param        BlockInput  body      representations.CreateBlockInput  true  "Mine block"
// @Success      201         {object}  representations.ReadableBlock
// @Failure      400         {object}  HTTPError
// @Failure      500         {object}  HTTPError
// @Router       /blockchain/block [post]
func (bch *BlockchainHandler) AddToBlockchain(ctx *gin.Context) {
	// Validate input
	var input reps.CreateBlockInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		NewError(ctx, http.StatusBadRequest, err)
		return
	}

	log.Info("Adding Block to blockchain: ", utils.Pretty(input))

	// Create block and persist to db
	newBlock, err := bch.blockchainService.AddToBlockChain(input.From, input.To, input.Amount)
	if err != nil {
		log.WithField("error", err.Error()).Error("Error adding block")
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	// Format return data to be readable
	data := bch.assemblerService.ToReadableBlock(newBlock)

	ctx.JSON(http.StatusCreated, gin.H{"block": data})
}

// GetBlockchain ... Print out all blocks in blockchain
// @Summary      Get all blocks
// @Description  Get all blocks on the blockchain
// @Tags         Blocks
// @Success      200  {array}   representations.ReadableBlock
// @Failure      500  {object}  HTTPError
// @Router       /blockchain [get]
func (bch *BlockchainHandler) GetBlockchain(ctx *gin.Context) {
	log.Info("Printing out the Blockchain")
	blockchain, err := bch.blockchainService.GetBlockchain()
	if err != nil {
		log.WithField("error", err.Error()).Error("Error getting blockchain")
		NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	data := make([]reps.ReadableBlock, 0)

	for _, block := range blockchain {
		data = append(data, bch.assemblerService.ToReadableBlock(block))
	}

	ctx.JSON(http.StatusOK, gin.H{"blockchain": data})
}

// GetGenesisBlock ... Get the genesis block
// @Summary      Get the genesis block
// @Description  Get the genesis block on the blockchain
// @Tags         Blocks
// @Success      200  {object}  representations.ReadableBlock
// @Failure      404  {object}  HTTPError
// @Router       /blockchain/block/genesis [get]
func (bch *BlockchainHandler) GetGenesisBlock(ctx *gin.Context) {
	log.Info("Getting Genesis block")
	genesis, err := bch.blockchainService.GetGenesisBlock()
	if err != nil {
		log.WithField("error", err.Error()).Error("Error getting genesis block")
		NewError(ctx, http.StatusNotFound, err)
	} else {
		formattedGenesis := bch.assemblerService.ToReadableBlock(genesis)
		ctx.JSON(http.StatusOK, gin.H{"genesis": formattedGenesis})
	}
}

// GetBlock ... Get block by block ID
// @Summary      Get a block
// @Description  Get a block on the blockchain by block ID
// @Tags         Blocks
// @Param        blockId  path      string  true  "Block ID"
// @Success      200      {object}  representations.ReadableBlock
// @Failure      404      {object}  HTTPError
// @Router       /blockchain/block/{blockId} [get]
func (bch *BlockchainHandler) GetBlock(ctx *gin.Context) {
	blockId := ctx.Param("blockId")
	log.Info("Getting block with blockId: ", blockId)

	block, err := bch.blockchainService.GetBlock(blockId)
	if err != nil {
		log.WithField("error", err.Error()).Error("Error getting block")
		NewError(ctx, http.StatusNotFound, err)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"block": bch.assemblerService.ToReadableBlock(block)})
	}
}

// GetLastBlock ... Get last block in blockchain. If it's a genesis, it will return it.
// @Summary      Get the last block
// @Description  Get the last block on the blockchain
// @Tags         Blocks
// @Success      200  {object}  representations.ReadableBlock
// @Failure      404  {object}  HTTPError
// @Router       /blockchain/block/last [get]
func (bch *BlockchainHandler) GetLastBlock(ctx *gin.Context) {
	log.Info("Getting last block...")

	lastBlock, err := bch.blockchainService.GetLastBlock()
	if err != nil {
		log.WithField("error", err.Error()).Error("Error getting last block")
		NewError(ctx, http.StatusNotFound, err)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"block": bch.assemblerService.ToReadableBlock(lastBlock)})
	}
}
