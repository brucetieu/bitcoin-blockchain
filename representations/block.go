package representations

// Format of payload when mining a block
type CreateBlockInput struct {
	From   string `json:"from" binding:"required"`
	To     string `json:"to" binding:"required"`
	Amount int    `json:"amount" binding:"required"`
}

// Block representation in bitcoin blockchain
type Block struct {
	ID           string        `gorm:"primary_key;type:char(36);column:block_id"`
	Timestamp    int64         `json:"timestamp"`
	Transactions []Transaction `json:"transactions" gorm:"foreignKey:BlockID"`
	PrevHash     []byte        `json:"prevHash"`
	Hash         []byte        `json:"hash"`
	Nounce       int64         `json:"nounce"`
}


type ReadableBlock struct {
	ID           string                `gorm:"primary_key;type:char(36);column:block_id"`
	Timestamp    int64                 `json:"timestamp"`
	Transactions []ReadableTransaction `json:"transactions" gorm:"foreignKey:BlockID"`
	PrevHash     string                `json:"prevHash"`
	Hash         string                `json:"hash"`
	Nounce       int64                 `json:"nounce"`
}
