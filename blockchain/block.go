package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"github.com/brucetieu/blockchain/utils"
	log "github.com/sirupsen/logrus"
	"time"
)

type Block struct {
	Timestamp int64  `json:"timestamp,omitempty"`
	Data      []byte `json:"data,omitempty"`
	PrevHash  []byte `json:"prevHash,omitempty"`
	Hash      []byte `json:"hash,omitempty"`
	Nounce    int64  `json:"nounce,omitempty"`
}

type customByte []byte

func (cb *customByte) UnmarshalJSON(input []byte) error {
	*cb = customByte(input)
	return nil
}

// Take block fields, add them together, then sha256 on the joined result.
func (b *Block) CreateHash() {
	joined := bytes.Join([][]byte{
		b.Data,
		b.PrevHash,
		utils.Int64ToByte(b.Timestamp),
	}, []byte{})
	hash := sha256.Sum256(joined)
	b.Hash = hash[:]
}

// Create a single block in the block chain.
func CreateBlock(data string, prevHash []byte) *Block {
	newBlock := &Block{
		Timestamp: time.Now().UnixMilli(),
		Data:      []byte(data),
		PrevHash:  prevHash,
	}
	proof := NewProofOfWork(newBlock)
	nounce, hash := proof.Solve()
	newBlock.Nounce = nounce
	newBlock.Hash = hash
	return newBlock
}

func (b *Block) Serialize() []byte {
	byteStruct, err := json.Marshal(b)
	if err != nil {
		log.Error("Unable to marshal", err.Error())
	}

	return byteStruct
}

func Deserialize(data []byte) *Block {
	var block Block
	err := json.Unmarshal(data, &block)
	if err != nil {
		log.Error("Unable to unmarshal: ", err.Error())
	}
	return &block
}
