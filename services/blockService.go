package services

import (
	"time"

	reps "github.com/brucetieu/blockchain/representations"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type BlockService interface {
	CreateBlock(txns []reps.Transaction, prevHash []byte) *reps.Block
}

type blockService struct {
}

func NewBlockService() BlockService {
	return &blockService{}
}

// Create a single block in the block chain.
func (bs *blockService) CreateBlock(txns []reps.Transaction, prevHash []byte) *reps.Block {
	log.Info("Mining block...")
	id := uuid.Must(uuid.NewRandom()).String()

	// Set ID in transactions to be Id of block
	for i := 0; i < len(txns); i++ {
		txns[i].BlockID = id
	}

	newBlock := &reps.Block{
		ID: id,
		Timestamp:    time.Now().UnixMilli(),
		Transactions: txns,
		PrevHash:     prevHash,
	}
	// proof := bs.powService.Solve()
	proof := NewProofOfWorkService(newBlock)
	nounce, hash := proof.Solve()
	newBlock.Nounce = nounce
	newBlock.Hash = hash
	return newBlock
}
