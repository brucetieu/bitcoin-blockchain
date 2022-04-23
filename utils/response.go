package utils

type NewBlockchainResponse struct {
	Timestamp int64  `json:"timestamp,omitempty"`
	IsGenesis string `json:"isGenesis,omitempty"`
}
