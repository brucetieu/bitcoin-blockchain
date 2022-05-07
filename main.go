package main

import (

	"github.com/brucetieu/blockchain/db"
	"github.com/brucetieu/blockchain/routes"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	
)

func main() {
	log.Info("Bitcoin Blockchain App")

	db.ConnectDatabase()

	router := gin.Default()
	routes.InitRoutes(router)

	// port :8080 by default
	_ = router.Run()

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
