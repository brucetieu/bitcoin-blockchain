package blockchain

import (
	"bytes"
	"crypto/sha256"
	"github.com/brucetieu/blockchain/utils"
	"time"
)

type Block struct {
	Timestamp int64
	Data      []byte
	PrevHash  []byte
	Hash      []byte
	Nounce    int64
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
