package services

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/brucetieu/blockchain/repository"
	reps "github.com/brucetieu/blockchain/representations"
	"github.com/brucetieu/blockchain/utils"
	"golang.org/x/crypto/ripemd160"

	"github.com/akamensky/base58"
	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"
)

var (
	ChecksumLen = 4
	Version     = byte(0)
)

type WalletService interface {
	CreateWallet() (reps.Wallet, error)
	GetWallet(address string) (reps.Wallet, error)
	// GetWallets() ([]reps.WalletGorm, error)
	CreateKeyPair() (ecdsa.PrivateKey, []byte)
	CreatePubKeyHash(pubKey []byte) ([]byte, error)
	CreateChecksum(pubKeyHash []byte) []byte
	CreateAddress(pubKey []byte) ([]byte, error)
	ValidateAddress(address string) (bool, error)
}

type walletService struct {
	blockchainRepo repository.BlockchainRepository
	walletAssember   WalletAssemblerFac
}

func NewWalletService(blockchainRepo repository.BlockchainRepository) WalletService {
	return &walletService{
		blockchainRepo: blockchainRepo,
		walletAssember: WalletAssembler,
	}
}

func (ws *walletService) CreateKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	privKey, _ := ecdsa.GenerateKey(curve, rand.Reader)

	// Public key is a combination of x and y coordinates on elliptic curve
	pubKey := append(privKey.X.Bytes(), privKey.Y.Bytes()...)

	log.Info(fmt.Sprintf("pubKey: %x\n", pubKey))
	return *privKey, pubKey
}

func (ws *walletService) GetWallet(address string) (reps.Wallet, error) {
	walletGorm, err := ws.blockchainRepo.GetWallet(address)
	if err != nil {
		return reps.Wallet{}, err
	}

	utils.PrettyPrintln("walletGorm in ws.GetWallet: ", walletGorm)
	return ws.walletAssember.ToWallet(&walletGorm), nil
}	

func (ws *walletService) CreateWallet() (reps.Wallet, error) {
	privKey, pubKey := ws.CreateKeyPair()

	walletAddress, err := ws.CreateAddress(pubKey)
	if err != nil {
		return reps.Wallet{}, err
	}

	log.Info("wallet address: ", string(walletAddress))

	wallet := reps.Wallet{uuid.Must(uuid.NewRandom()).String(), string(walletAddress), privKey, pubKey}

	utils.PrettyPrintln("wallet: ", wallet)
	// Persist
	err = ws.blockchainRepo.CreateWallet(ws.walletAssember.ToGormWallet(&wallet))
	if err != nil {
		return reps.Wallet{}, err
	}

	return wallet, nil
}

// pubKeyHash = ripemd160(sha256(pubKey))
func (ws *walletService) CreatePubKeyHash(pubKey []byte) ([]byte, error) {
	pubHash := sha256.Sum256(pubKey)

	ripemdHasher := ripemd160.New()
	_, err := ripemdHasher.Write(pubHash[:])
	if err != nil {
		log.Error(err.Error())
	}

	pubKeyHash := ripemdHasher.Sum(nil)

	log.Info(fmt.Sprintf("pubKeyHash: %x\n", pubKeyHash))
	return pubKeyHash, err
}

// checksum = sha256(sha256(pubKeyHash))
func (ws *walletService) CreateChecksum(pubKeyHash []byte) []byte {
	pubKeyHashSum := sha256.Sum256(pubKeyHash)
	pubKeyHashSum2 := sha256.Sum256(pubKeyHashSum[:])

	log.Info(fmt.Sprintf("checksum: %x\n", pubKeyHashSum2[:ChecksumLen]))
	return pubKeyHashSum2[:ChecksumLen] // checksum is first 4 bytes of second hash
}

func (ws *walletService) CreateAddress(pubKey []byte) ([]byte, error) {
	pubKeyHash, err := ws.CreatePubKeyHash(pubKey)
	if err != nil {
		return []byte{}, err
	}

	// Version + pubKeyHash
	versionedPubKeyHash := append([]byte{Version}, pubKeyHash...)
	log.Info(fmt.Sprintf("versionedPubKeyHash: %x", versionedPubKeyHash))

	checksum := ws.CreateChecksum(pubKeyHash)

	// version + pubKeyHash + checksum
	finalHash := append(versionedPubKeyHash, checksum...)
	log.Info(fmt.Sprintf("finalHash: %x", finalHash))

	address := base58Encode(finalHash)

	return address, nil
}

func (ws *walletService) ValidateAddress(address string) (bool, error) {
	// query address from db first
	_, err := ws.GetWallet(address)
	if err != nil {
		errMsg := fmt.Errorf("%s: wallet with address %s does not exist", err.Error(), address)
		return false, errMsg
	}

	pubKeyHash := base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-ChecksumLen:]
	pubKeyHash = pubKeyHash[1:len(pubKeyHash)-ChecksumLen]
	expectedChecksum := ws.CreateChecksum(pubKeyHash)

	if bytes.Compare(actualChecksum, expectedChecksum) == 0 {
		return true, nil
	}

	return false, nil
	// return bytes.Compare(actualChecksum, expectedChecksum) == 0
}

func base58Decode(address []byte) []byte {
	base58Decoded, _ := base58.Decode(string(address))
	return base58Decoded
}

func base58Encode(hash []byte) []byte {
	base58encoded := base58.Encode(hash)

	return []byte(base58encoded)
}
