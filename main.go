package main

import (
	// "crypto/sha1"
	// "encoding/hex"
	"fmt"

	// "net/http"

	// "github.com/brucetieu/blockchain/blockchain"
	// "github.com/brucetieu/blockchain/controllers"
	"github.com/brucetieu/blockchain/db"
	"github.com/brucetieu/blockchain/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Blockchain")

	db.ConnectDatabase()
	defer db.DB.Close()

	// r := gin.Default()
	router := gin.Default()

	routes.InitRoutes(router)

	router.Run()

	// r.GET("/", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "Welcome to the Blockchain",
	// 	})
	// })

	// r.POST("/blockchain", controllers.CreateBlockchain)
	// r.POST("/blockchain/add", )
	// r.Run()

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
