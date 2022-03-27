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
	dbURL := "postgres://postgres:pass@localhost:5432/blockchain"
	// dsn := "host=localhost user=postgres password=postgres port=5432 sslmode=disable"

	database, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.Logger.LogMode(logger.Info)


	// database.AutoMigrate(&models.Transaction{}, &models.TxnOutput{}, &models.TxnInput{})
	database.AutoMigrate(&reps.Block{})
	database.AutoMigrate(&reps.Transaction{})
	database.AutoMigrate(&reps.TxnInput{})
	database.AutoMigrate(&reps.TxnOutput{})
	// database.AutoMigrate(&reps.Block{}, &reps.Transaction{})
	// database.AutoMigrate(&reps.TxnOutput{}, &reps.TxnInput{})


	DB = database
}
