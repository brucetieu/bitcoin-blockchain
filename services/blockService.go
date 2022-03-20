package services

import (
	"time"

	reps "github.com/brucetieu/blockchain/representations"
	log "github.com/sirupsen/logrus"
)

type BlockService interface {
	CreateBlock(txns []*reps.Transaction, prevHash []byte) *reps.Block
}

type blockService struct {
}

func NewBlockService() BlockService {
	return &blockService{}
}

// Create a single block in the block chain.
func (bs *blockService) CreateBlock(txns []*reps.Transaction, prevHash []byte) *reps.Block {
	log.Info("Mining block...")
	newBlock := &reps.Block{
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
