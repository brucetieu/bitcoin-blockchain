package db

import (
	// "github.com/jinzhu/gorm"
	// "github.com/brucetieu/blockchain/models"
	reps "github.com/brucetieu/blockchain/representations"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Connect to running postgres container. database is the container name
	// dbURL := "postgres://postgres:pass@database:5432/blockchain"
	dbURL := "postgres://postgres:pass@localhost:5432/blockchain"

	database, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	database.Logger.LogMode(logger.Info)

	database.AutoMigrate(&reps.Block{})
	database.AutoMigrate(&reps.Transaction{})
	database.AutoMigrate(&reps.TxnInput{})
	database.AutoMigrate(&reps.TxnOutput{})
	database.AutoMigrate(&reps.WalletGorm{})

	DB = database
}
