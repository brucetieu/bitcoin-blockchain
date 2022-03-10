package services

import (
	"encoding/json"
	"time"

	"github.com/brucetieu/blockchain/representations"
	log "github.com/sirupsen/logrus"
)

type BlockService interface {
	CreateBlock(data string, prevHash []byte) *representations.Block
	Serialize(*representations.Block) []byte
	Deserialize(data []byte) *representations.Block
}

type blockService struct {
}

func NewBlockService() BlockService {
	return &blockService{}
}

// Create a single block in the block chain.
func (bs *blockService) CreateBlock(data string, prevHash []byte) *representations.Block {
	newBlock := &representations.Block{
		Timestamp: time.Now().UnixMilli(),
		Data:      []byte(data),
		PrevHash:  prevHash,
	}
	proof := NewProofOfWorkService(newBlock)
	nounce, hash := proof.Solve()
	newBlock.Nounce = nounce
	newBlock.Hash = hash
	return newBlock
}

func (bs *blockService) Serialize(block *representations.Block) []byte {
	byteStruct, err := json.Marshal(block)
	if err != nil {
		log.Error("Unable to marshal", err.Error())
	}

	return byteStruct
}

func (bs *blockService) Deserialize(data []byte) *representations.Block {
	var block representations.Block
	err := json.Unmarshal(data, &block)
	if err != nil {
		log.Error("Unable to unmarshal: ", err.Error())
	}
	return &block
}
