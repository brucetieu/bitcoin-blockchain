package main

import (
	"fmt"

	"github.com/brucetieu/blockchain/blockchain"
	// "github.com/brucetieu/blockchain/utils"
)

func main() {
	fmt.Println("Blockchain")
	
	bc := blockchain.NewBlockchain()
	bc.AddToBlockChain("My first block")
	bc.AddToBlockChain("My second block")
	bc.AddToBlockChain("My third block")

	for _, block := range bc.Blocks {
		fmt.Println("Block data: ", string(block.Data))
		fmt.Printf("Block prev hash: %x\n", block.PrevHash)
		fmt.Printf("Block hash: %x\n", block.Hash)
	}
}