package db

import "github.com/dgraph-io/badger/v3"

var DB *badger.DB

func ConnectDatabase() {
	badgerDb, err := badger.Open(badger.DefaultOptions("/tmp/blockchain"))
	if err != nil {
		panic("Failed to connect to badger db.")
	}
	
	DB = badgerDb
}