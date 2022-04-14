package representations

import (
	// "crypto/ecdsa"
)

type Wallet struct {
	// PrivateKey ecdsa.PrivateKey
	ID string `json:"id" gorm:"primary_key"`
	Address string
	PrivateKey []byte
	PublicKey []byte
}
