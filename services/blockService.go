package services

import (
	// "bytes"
	// "crypto/sha256"
	"encoding/json"
	"time"

	"github.com/brucetieu/blockchain/representations"
	// "github.com/brucetieu/blockchain/utils"
	log "github.com/sirupsen/logrus"
)

// type Block struct {
// 	Timestamp int64  `json:"timestamp,omitempty"`
// 	Data      []byte  `json:"data,omitempty"`
// 	PrevHash  []byte   `json:"prevHash,omitempty"`
// 	Hash      []byte   `json:"hash,omitempty"`
// 	Nounce    int64    `json:"nounce,omitempty"`
// }

type BlockService interface {
	CreateBlock(data string, prevHash []byte) *representations.Block
	Serialize(*representations.Block) []byte
	Deserialize(data []byte) *representations.Block
}

type blockService struct {
	// powService PowService
}

// func NewBlockService(powService PowService) BlockService {
// 	return &blockService{
// 		powService:  powService,
// 	}
// }
func NewBlockService() BlockService {
	return &blockService{}
}
// Take block fields, add them together, then sha256 on the joined result.
// func (b *Block) CreateHash() {
// 	joined := bytes.Join([][]byte{
// 		b.Data,
// 		b.PrevHash,
// 		utils.Int64ToByte(b.Timestamp),
// 	}, []byte{})
// 	hash := sha256.Sum256(joined)
// 	b.Hash = hash[:]
// }

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

func (bs *blockService) Serialize(block *representations.Block) []byte{
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
