package blockchain

type Blockchain struct {
	Blocks []*Block
}

func NewBlockchain() *Blockchain {
	b := CreateBlock("Start of Blockchain", []byte{})
	return &Blockchain{
		Blocks: []*Block{
			b,
		},
	}
}

// Add a single block to the blockchain. It always goes at the end.
func (bc *Blockchain) AddToBlockChain(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}
