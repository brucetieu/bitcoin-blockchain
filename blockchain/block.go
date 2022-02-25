package blockchain

import (
	"bytes"
	"crypto/sha256"
)

type Block struct {
	Data []byte
	PrevHash []byte
	Hash []byte
}

func (b *Block) CreateHash() {
	joined := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(joined)
	b.Hash = hash[:]
}
func CreateBlock(data string, prevHash []byte) *Block {
	newBlock := &Block{
		Data: []byte(data),
		PrevHash: prevHash,
	}
	newBlock.CreateHash()
	return newBlock
}