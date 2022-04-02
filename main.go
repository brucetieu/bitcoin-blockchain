package main

import (
	"fmt"

	"github.com/brucetieu/blockchain/db"
	"github.com/brucetieu/blockchain/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Blockchain")

	db.ConnectDatabase()
	// defer db.DB.

	// r := gin.Default()
	router := gin.Default()

	routes.InitRoutes(router)

	// port :8080 by default
	router.Run()

	// for _, block := range bc.Blocks {
	// 	fmt.Println("Timestamp: ", block.Timestamp)
	// 	fmt.Println("Block data: ", string(block.Data))
	// 	fmt.Printf("Block prev hash: %x\n", block.PrevHash)
	// 	fmt.Printf("Block hash: %x\n", block.Hash)
	// 	fmt.Println("Block nounce: ", block.Nounce)

	// 	pow := blockchain.NewProofOfWork(block)
	// 	fmt.Printf("valid proof? %t\n\n", pow.ValidateProof())
	// }

}
