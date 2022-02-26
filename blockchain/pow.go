package blockchain

import (
	"bytes"
	"crypto/sha256"
	"github.com/brucetieu/blockchain/utils"
	"math/big"
)

// First TargetBits number of bits of sha256(block header + counter) need to all be 0.
const TargetBits = 12

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)

	// means the first TargetBits number of bits will be 0. e.g. 0000000000001...
	target.Lsh(target, uint(256-TargetBits))
	return &ProofOfWork{block, target}
}

func (pow *ProofOfWork) Solve() (int64, []byte) {
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
func (pow *ProofOfWork) HashData() []byte {
	joined := bytes.Join([][]byte{
		pow.Block.Data,
		pow.Block.PrevHash,
		utils.Int64ToByte(pow.Block.Timestamp),
		utils.Int64ToByte(pow.Block.Nounce),
	}, []byte{})
	hash := sha256.Sum256(joined)
	return hash[:]
}

func (pow *ProofOfWork) ValidateProof() bool {
	proposedHash := pow.HashData()
	proposedHashInt := new(big.Int)
	proposedHashInt.SetBytes(proposedHash)
	if proposedHashInt.Cmp(pow.Target) == -1 {
		return true
	}
	return false
}
