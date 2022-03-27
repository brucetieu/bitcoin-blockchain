package representations

import (
	// "github.com/google/uuid"
)


// type Block struct {
// 	ID     string   `json:"id" gorm:"primary_key;column:block_id"`
// 	Timestamp    int64          `json:"timestamp"`
// 	Transactions []Transaction `json:"transactions" gorm:"foreignKey:TxnID"`
// 	PrevHash     []byte         `json:"prevHash"`
// 	Hash         []byte         `json:"hash"`
// 	Nounce       int64          `json:"nounce"`
// }

type CreateBlockInput struct {
	From   string `json:"from" binding:"required"`
	To     string `json:"to" binding:"required"`
	Amount int    `json:"amount" binding:"required"`
}

type Block struct {
	ID string `gorm:"primary_key;type:char(36);column:block_id"`
	Timestamp    int64          `json:"timestamp"`
	Transactions []Transaction `json:"transactions" gorm:"foreignKey:BlockID"`
	PrevHash     []byte         `json:"prevHash"`
	Hash         []byte         `json:"hash"`
	Nounce       int64          `json:"nounce"`
}
