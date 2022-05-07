package representations

import (
	"crypto/sha256"

	log "github.com/sirupsen/logrus"
)

type MerkleTree struct {
	Root *MerkleNode
}

type MerkleNode struct {
	Data  []byte
	Left  *MerkleNode
	Right *MerkleNode
}

func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	merkleNode := MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		merkleNode.Data = hash[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		merkleNode.Data = hash[:]
	}

	merkleNode.Left = left
	merkleNode.Right = right

	return &merkleNode
}

// Data is a list of transactions
func NewMerkleTree(txns [][]byte) *MerkleTree {
	log.Info("Creating new merkle tree")
	merkleNodes := make([]*MerkleNode, 0)
	// If number of transactions is odd, make a copy of the last transaction and append it to the transactions to satisfy merkle tree structure
	if len(txns)%2 != 0 {
		txns = append(txns, txns[len(txns)-1])
	}

	// Create a leaf merkle tree node for each transaction
	for _, txn := range txns {
		merkleNode := NewMerkleNode(nil, nil, txn)
		merkleNodes = append(merkleNodes, merkleNode)
	}

	// Build merkle tree from bottom up
	// e.g. 4 leafs = 7 nodes = 3 levels
	for i := 0; i < len(txns)/2; i++ {
		treeLevel := make([]*MerkleNode, 0)

		for j := 0; j < len(merkleNodes); j += 2 {
			merkleNode := NewMerkleNode(merkleNodes[j], merkleNodes[j+1], nil)
			treeLevel = append(treeLevel, merkleNode)
		}

		merkleNodes = treeLevel
	}

	return &MerkleTree{merkleNodes[0]}
}
