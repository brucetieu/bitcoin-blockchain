package representations

type Block struct {
	Timestamp    int64          `json:"timestamp,omitempty"`
	Transactions []*Transaction `json:"transactions,omitempty"`
	PrevHash     []byte         `json:"prevHash,omitempty"`
	Hash         []byte         `json:"hash,omitempty"`
	Nounce       int64          `json:"nounce,omitempty"`
}

type CreateBlockInput struct {
	From   string `json:"from" binding:"required"`
	To     string `json:"to" binding:"required"`
	Amount int    `json:"amount" binding:"required"`
}
