package representations

type Block struct {
	Timestamp int64  `json:"timestamp,omitempty"`
	Transactions      []*Transaction `json:"transactions,omitempty"`
	PrevHash  []byte `json:"prevHash,omitempty"`
	Hash      []byte `json:"hash,omitempty"`
	Nounce    int64  `json:"nounce,omitempty"`
}

type CreateBlockInput struct {
	Data string `json:"data" binding:"required"`
}


