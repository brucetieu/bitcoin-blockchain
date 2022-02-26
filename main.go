package main

import (
	"fmt"
	"github.com/brucetieu/blockchain/blockchain"
)

func main() {
	fmt.Println("Blockchain")

	bc := blockchain.NewBlockchain()
	bc.AddToBlockChain("My first block")
	bc.AddToBlockChain("My second block")
	bc.AddToBlockChain("My third block")

	for _, block := range bc.Blocks {
		fmt.Println("Timestamp: ", block.Timestamp)
		fmt.Println("Block data: ", string(block.Data))
		fmt.Printf("Block prev hash: %x\n", block.PrevHash)
		fmt.Printf("Block hash: %x\n", block.Hash)
		fmt.Println("Block nounce: ", block.Nounce)

		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("valid proof? %t\n\n", pow.ValidateProof())
	}
}
