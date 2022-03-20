package services

import (
	"encoding/json"
	"time"

	reps "github.com/brucetieu/blockchain/representations"
	log "github.com/sirupsen/logrus"
)

type BlockService interface {
	CreateBlock(txns []*reps.Transaction, prevHash []byte) *reps.Block
	Serialize(*reps.Block) []byte
	Deserialize(data []byte) *reps.Block
}

type blockService struct {
}

func NewBlockService() BlockService {
	return &blockService{}
}

// Create a single block in the block chain.
func (bs *blockService) CreateBlock(txns []*reps.Transaction, prevHash []byte) *reps.Block {
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

func (bs *blockService) Serialize(block *reps.Block) []byte {
	byteStruct, err := json.Marshal(block)
	if err != nil {
		log.Error("Unable to marshal", err.Error())
	}

	return byteStruct
}

func (bs *blockService) Deserialize(data []byte) *reps.Block {
	var block reps.Block
	err := json.Unmarshal(data, &block)
	if err != nil {
		log.Error("Unable to unmarshal: ", err.Error())
	}
	return &block
}
