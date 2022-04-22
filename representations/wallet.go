package representations

import (
"crypto/ecdsa"
)

type Wallet struct {
	ID         string `json:"id" gorm:"primary_key"`
	Address    string  `json:"address"`
	PrivateKey ecdsa.PrivateKey `json:"privateKey"`
	PublicKey  []byte`json:"publicKey"`
}

// This type of wallet can be saved to gormdb
// WalletByte is the contents of the wallet serialized (privateKey, publicKey)
type WalletGorm struct {
	ID         string `json:"id" gorm:"primary_key"`
	Address string `json:"address"`
	Balance int `json:"balance,omitempty"`
	WalletByte []byte `json:"walletByte,omitempty"`// The entire contents of Wallet in byte form
}

