package db

import (
	"fmt"
	"os"

	reps "github.com/brucetieu/blockchain/representations"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/jinzhu/gorm"
)

var (
	DB *gorm.DB

	pgUser = "postgres"
	pgPass = "pass"
	pgDbName = "blockchain"
	pgHost = "database"
	pgPort = "5432"
)

type PgEnvVars struct {
	PostgresUser string
	PostgresPass string
	PostgresDB string
	PostgresHost string
}

func ConnectDatabase() {

	dbURL := getPgConnectionString()
	log.Info("dbURL: ", dbURL)
	database, err := gorm.Open("postgres", dbURL)
	if err != nil {
		panic("Failed to connect to database! " + err.Error())
	}

	database.LogMode(true)

	_ = database.AutoMigrate(&reps.Block{})
	_ = database.AutoMigrate(&reps.Transaction{})
	_ = database.AutoMigrate(&reps.TxnInput{})
	_ = database.AutoMigrate(&reps.TxnOutput{})
	_ = database.AutoMigrate(&reps.Wallet{})

	DB = database
}

func getPgConnectionString () string {
	envVars := PgEnvVars{
		PostgresUser: pgUser,
		PostgresPass: pgPass,
		PostgresDB: pgDbName,
		PostgresHost: pgHost,
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	envPgUser := os.Getenv("POSTGRES_USER")
	if envPgUser != "" {
		envVars.PostgresUser = envPgUser
	}

	envPgPass := os.Getenv("POSTGRES_PASSWORD")
	if envPgPass != "" {
		envVars.PostgresPass = envPgPass
	}

	envPgDbName := os.Getenv("POSTGRES_DB");
	if envPgDbName != "" {
		envVars.PostgresDB = envPgDbName
	}

	envPgHostName := os.Getenv("POSTGRES_HOST_NAME")
	if envPgHostName != "" {
		envVars.PostgresHost = envPgHostName
	}
	
	dbURL := ""
	debugMode := os.Getenv("DEBUG")
	if debugMode == "false" {
	// Connect to running postgres container. database is the container name
		dbURL = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", envVars.PostgresHost, pgPort, envVars.PostgresUser, envVars.PostgresDB, envVars.PostgresPass)
	} else {
		dbURL = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", "localhost", pgPort, envVars.PostgresUser, envVars.PostgresDB, envVars.PostgresPass)
	}

	return dbURL
}
