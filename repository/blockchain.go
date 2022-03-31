package repository

import (
	"github.com/brucetieu/blockchain/db"
	// "github.com/brucetieu/blockchain/utils"

	reps "github.com/brucetieu/blockchain/representations"
)

type BlockchainRepository interface {
	CreateTransaction(txn []reps.Transaction) error
	GetTransactionsByBlockId(blockId string) ([]reps.Transaction, error)
	GetTransactions() ([]reps.Transaction, error)
	// GetTransactionsByTxnId(txnId []byte) ([]reps.Transaction, error)

	CreateBlock(block reps.Block) error
	GetGenesisBlock() (reps.Block, int, error)
	GetBlockchain() ([]reps.Block, error)
	GetLastBlock() (reps.Block, error)
	GetBlockById(blockId string) (reps.Block, int, error)

	CreateTxnOutput(txnOutput reps.TxnOutput) error
	CreateTxnInput(txnInput reps.TxnInput) error
	GetTxnInputs(txnId []byte) ([]reps.TxnInput, error)
	GetTxnOutputs(txnId []byte) ([]reps.TxnOutput, error)
}

type blockchainRepository struct {
}

func NewBlockchainRepository() BlockchainRepository {
	return &blockchainRepository{}
}

func (repo *blockchainRepository) CreateTxnInput(txnInput reps.TxnInput) error {
	if err := db.DB.Create(&txnInput).Error; err != nil {
		return err
	}
	return nil
}

func (repo *blockchainRepository) CreateTxnOutput(txnOutput reps.TxnOutput) error {
	if err := db.DB.Create(&txnOutput).Error; err != nil {
		return err
	}
	return nil
}

func (repo *blockchainRepository) CreateTransaction(txn []reps.Transaction) error {
	if err := db.DB.Create(&txn).Error; err != nil {
		return err
	}
	return nil
}

// Get the last (non genesis) block in the blockchain
func (repo *blockchainRepository) GetLastBlock() (reps.Block, error) {
	var lastBlock reps.Block

	err := db.DB.Limit(1).Order("timestamp desc").Find(&lastBlock).Error
	if err != nil {
		return reps.Block{}, err
	}

	return lastBlock, nil
}

func (repo *blockchainRepository) GetBlockById(blockId string) (reps.Block, int, error) {
	var block reps.Block

	res := db.DB.
		Preload("Transactions").
		Where("block_id = ?", blockId).
		Find(&block)
	if res.Error != nil {
		return reps.Block{}, -1, res.Error
	}

	txns, err := repo.GetTransactionsByBlockId(block.ID)
	if err != nil {
		return reps.Block{}, -1, err
	}

	block.Transactions = txns

	return block, int(res.RowsAffected), nil
}

func (repo *blockchainRepository) GetTransactions() ([]reps.Transaction, error) {
	var transactions []reps.Transaction

	err := db.DB.
		Preload("Inputs").
		Preload("Outputs").
		Find(&transactions).
		Error

	if err != nil {
		return []reps.Transaction{}, err
	}

	return transactions, nil
}

// Get all transactions in a given block.
func (repo *blockchainRepository) GetTransactionsByBlockId(blockId string) ([]reps.Transaction, error) {
	var transactions []reps.Transaction

	err := db.DB.Where("block_id = ?", blockId).
		Preload("Inputs").
		Preload("Outputs").
		Find(&transactions).
		Error

	if err != nil {
		return []reps.Transaction{}, err
	}

	return transactions, nil
}

// func (repo *blockchainRepository) GetTransactionsByTxnId(txnId []byte) ([]reps.Transaction, error) {
// 	var transactions []reps.Transaction

// 	err := db.DB.Where("id = ?", txnId).
// 		Preload("Inputs").
// 		Preload("Outputs").
// 		Find(&transactions).
// 		Error

// 	if err != nil {
// 		return []reps.Transaction{}, err
// 	}

// 	return transactions, nil
// }

// Get all transaction inputs for a given transaction
func (repo *blockchainRepository) GetTxnInputs(txnId []byte) ([]reps.TxnInput, error) {
	var txnInputs []reps.TxnInput

	err := db.DB.Where("curr_txn_id = ?", txnId).
		Find(&txnInputs).
		Error

	if err != nil {
		return []reps.TxnInput{}, err
	}

	return txnInputs, nil
}

// Get all transaction outputs for a given transaction
func (repo *blockchainRepository) GetTxnOutputs(txnId []byte) ([]reps.TxnOutput, error) {
	var txnOutputs []reps.TxnOutput

	err := db.DB.Where("curr_txn_id = ?", txnId).
		Find(&txnOutputs).
		Error

	if err != nil {
		return []reps.TxnOutput{}, err
	}

	return txnOutputs, nil
}

// Get the first block in the blockchain
func (repo *blockchainRepository) GetGenesisBlock() (reps.Block, int, error) {
	var genesisBlock reps.Block

	res := db.DB.
		Preload("Transactions").
		Where("prev_hash = ''").
		Find(&genesisBlock)
	if res.Error != nil {
		return reps.Block{}, -1, res.Error
	}

	txns, err := repo.GetTransactionsByBlockId(genesisBlock.ID)
	if err != nil {
		return reps.Block{}, -1, err
	}

	genesisBlock.Transactions = txns

	return genesisBlock, int(res.RowsAffected), nil

}

func (repo *blockchainRepository) CreateBlock(block reps.Block) error {
	if err := db.DB.Create(&block).Error; err != nil {
		return err
	}

	return nil
}

func (repo *blockchainRepository) GetBlockchain() ([]reps.Block, error) {
	var blocks []reps.Block

	if err := db.DB.
		Preload("Transactions").
		Find(&blocks).Error; err != nil {
		return []reps.Block{}, nil
	}

	for i := 0; i < len(blocks); i++ {
		blockId := blocks[i].ID

		txns, err := repo.GetTransactionsByBlockId(blockId)
		if err != nil {
			return []reps.Block{}, err
		}

		blocks[i].Transactions = txns
	}

	return blocks, nil
}
