package representations

type Wallet struct {
	ID         string `json:"id,omitempty" gorm:"primary_key"`
	Address    string  `json:"address,omitempty"`
	PrivateKey []byte `json:"privateKey,omitempty"`
	PublicKey  string`json:"publicKey,omitempty"`
}

// This represents balance information for a wallet (address)
type AddressBalance struct {
	Address string `json:"address,omitempty"`
	PublicKey string `json:"publicKey,omitempty"`
	Balance int `json:"balance"`
}
