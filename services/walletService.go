package services

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"github.com/brucetieu/blockchain/repository"
	reps "github.com/brucetieu/blockchain/representations"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ripemd160"

	"github.com/akamensky/base58"
	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"
)

var (
	ChecksumLen = 4
	Version = byte(0)
)
type WalletService interface {
	CreateWallet() (reps.Wallet, error)
	CreateKeyPair() (ecdsa.PrivateKey, []byte)
	CreatePubKeyHash(pubKey []byte) ([]byte , error)
	CreateChecksum(pubKeyHash []byte) []byte
	CreateAddress(pubKey []byte) ([]byte, error)
}

type walletService struct {
	blockchainRepo repository.BlockchainRepository
}

func NewWalletService(blockchainRepo repository.BlockchainRepository) WalletService {
	return &walletService{
		blockchainRepo: blockchainRepo,
	}
}

func (ws *walletService) CreateKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	privKey, _ := ecdsa.GenerateKey(curve, rand.Reader)

	// Public key is a combination of x and y coordinates on elliptic curve
	pubKey := append(privKey.X.Bytes(), privKey.Y.Bytes()...)

	return *privKey, pubKey
}
func (ws *walletService) CreateWallet() (reps.Wallet, error) {
	privKey, pubKey := ws.CreateKeyPair()

	walletAddress, err := ws.CreateAddress(pubKey)
	if err != nil {
		return reps.Wallet{}, err
	}

	log.Info("wallet address: ", string(walletAddress))

	wallet := reps.Wallet{uuid.Must(uuid.NewRandom()).String(), string(walletAddress), privKey.D.Bytes(), pubKey}

	// Persist
	err = ws.blockchainRepo.CreateWallet(wallet)
	if err != nil {
		return reps.Wallet{}, err
	}

	return wallet, nil
}

func (ws *walletService) CreatePubKeyHash(pubKey []byte) ([]byte, error) {
	pubHash := sha256.Sum256(pubKey)

	ripemdHasher := ripemd160.New()
	_, err := ripemdHasher.Write(pubHash[:])
	if err != nil {
		logrus.Error(err.Error())
	}

	pubKeyHash := ripemdHasher.Sum(nil)

	return pubKeyHash, err
}

func (ws *walletService) CreateChecksum(pubKeyHash []byte) []byte {
	pubKeyHashSum := sha256.Sum256(pubKeyHash)
	pubKeyHashSum2 := sha256.Sum256(pubKeyHashSum[:])

	return pubKeyHashSum2[:ChecksumLen] // checksum is first 4 bytes of second hash
}

func (ws *walletService) CreateAddress(pubKey []byte) ([]byte, error) {
	pubKeyHash, err := ws.CreatePubKeyHash(pubKey)
	if err != nil {
		return []byte{}, err
	}

	// Version + pubKeyHash
	versionedPubKeyHash := append([]byte{Version}, pubKeyHash...)

	checksum := ws.CreateChecksum(pubKeyHash)

	// version + pubKeyHash + checksum
	finalHash := append(versionedPubKeyHash, checksum...)

	address := base58Encode(finalHash)

	return address, nil
}

func base58Encode(hash []byte) []byte {
	base58encoded := base58.Encode(hash)

	return []byte(base58encoded)
}