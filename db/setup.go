package db

import (
	reps "github.com/brucetieu/blockchain/representations"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Connect to running postgres container. database is the container name
	dbURL := "postgres://postgres:pass@database:5432/blockchain"

	// For debugging code
	// dbURL := "postgres://postgres:pass@localhost:5432/blockchain"

	database, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	database.Logger.LogMode(logger.Info)

	_ = database.AutoMigrate(&reps.Block{})
	_ = database.AutoMigrate(&reps.Transaction{})
	_ = database.AutoMigrate(&reps.TxnInput{})
	_ = database.AutoMigrate(&reps.TxnOutput{})
	_ = database.AutoMigrate(&reps.Wallet{})

	DB = database
}
