package services

import (
	"bytes"
	"crypto/sha256"
	"math/big"

	"github.com/brucetieu/blockchain/representations"
	"github.com/brucetieu/blockchain/utils"
)

var (
	TargetBits = 12
)

type PowService interface {
	Solve() (int64, []byte)
	HashData() []byte
	ValidateProof() bool
}

type powService struct {
	Block          *representations.Block
	Target         *big.Int
	txnAssembler TxnAssemblerFac
}

func NewProofOfWorkService(block *representations.Block) PowService {
	target := big.NewInt(1)

	// means the first TargetBits number of bits will be 0. e.g. 0000000000001...
	target.Lsh(target, uint(256-TargetBits))

	return &powService{
		Target:         target,
		Block:          block,
		txnAssembler: TxnAssembler,
	}
}

func (pow *powService) Solve() (int64, []byte) {
	nounce := 0
	var solvedHash []byte
	solvedHashInt := new(big.Int)

	for {
		pow.Block.Nounce = int64(nounce)
		solvedHash = pow.HashData()
		solvedHashInt.SetBytes(solvedHash)

		// Check if HASH(data + nounce) < target number
		if solvedHashInt.Cmp(pow.Target) == -1 {
			break
		} else {
			nounce++
		}
	}

	// miner is basically trying to solve for nounce.
	return int64(nounce), solvedHash
}

// sha256 hash the block data and nounce
func (pow *powService) HashData() []byte {
	joined := bytes.Join([][]byte{
		pow.txnAssembler.HashTransactions(pow.Block.Transactions),
		pow.Block.PrevHash,
		utils.Int64ToByte(pow.Block.Timestamp),
		utils.Int64ToByte(pow.Block.Nounce),
	}, []byte{})
	hash := sha256.Sum256(joined)
	return hash[:]
}

func (pow *powService) ValidateProof() bool {
	proposedHash := pow.HashData()
	proposedHashInt := new(big.Int)
	proposedHashInt.SetBytes(proposedHash)

	return proposedHashInt.Cmp(pow.Target) == -1
}
