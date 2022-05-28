package main

import (
	"os"

	"github.com/brucetieu/blockchain/db"
	"github.com/brucetieu/blockchain/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	_ "github.com/brucetieu/blockchain/docs" // load swagger docs
)

// @title        Bitcoin Blockchain API documentation
// @version      1.0.0
// @description  This is a Bitcoin Blockchain API

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /bitcoin
func main() {
	log.Info("Bitcoin Blockchain App")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.ConnectDatabase()

	router := gin.Default()
	routes.InitRoutes(router)

	// port 5000 by default
	_ = router.Run(":" + os.Getenv("PORT"))
}

