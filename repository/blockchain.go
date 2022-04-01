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
	GetTransaction(txnId []byte) (reps.Transaction, error)
	GetAddresses() (map[string]bool, error)
	// GetTransactionsByTxnId(txnId []byte) ([]reps.Transaction, error)

	CreateBlock(block reps.Block) error
	GetGenesisBlock() (reps.Block, error)
	GetBlockchain() ([]reps.Block, error)
	GetLastBlock() (reps.Block, error)
	GetBlockById(blockId string) (reps.Block, error)

	CreateTxnOutput(txnOutput reps.TxnOutput) error
	CreateTxnInput(txnInput reps.TxnInput) error
	// GetTxnInputs(txnId []byte) ([]reps.TxnInput, error)
	// GetTxnOutputs(txnId []byte) ([]reps.TxnOutput, error)
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

func (repo *blockchainRepository) GetAddresses() (map[string]bool, error) {
	// Raw SQL
	addresses := make(map[string]bool)
	var address string

	rows, err := db.DB.Raw("SELECT DISTINCT script_pub_key FROM txn_outputs").Rows()
	if err != nil {
		return addresses, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&address)
		addresses[address] = true
	}
	// var addresses []string

	// err := db.DB.
	// 	Raw("SELECT DISTINCT script_pub_key FROM txn_outputs").
	// 	Scan(&addresses).
	// 	Error
	
	// if err != nil {
	// 	return []string{}, err
	// }

	return addresses, nil
}

// Get the last block in the blockchain
func (repo *blockchainRepository) GetLastBlock() (reps.Block, error) {
	var lastBlock reps.Block

	err := db.DB.Limit(1).Order("timestamp desc").First(&lastBlock).Error
	if err != nil {
		return reps.Block{}, err
	}

	return lastBlock, nil
}

// Get a block in the block chain by blockId
func (repo *blockchainRepository) GetBlockById(blockId string) (reps.Block, error) {
	var block reps.Block

	res := db.DB.
		Where("block_id = ?", blockId).
		First(&block)
	if res.Error != nil {
		return reps.Block{}, res.Error
	}

	txns, err := repo.GetTransactionsByBlockId(block.ID)
	if err != nil {
		return reps.Block{}, err
	}

	block.Transactions = txns

	return block, nil
}

// Get all transactions
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

// Get a single transaction
func (repo *blockchainRepository) GetTransaction(txnId []byte) (reps.Transaction, error) {
	var transaction reps.Transaction

	res := db.DB.Where("id = ?", txnId).
		Preload("Inputs").
		Preload("Outputs").
		First(&transaction)

	if res.Error != nil {
		return reps.Transaction{}, res.Error
	}

	return transaction, nil
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
func (repo *blockchainRepository) GetGenesisBlock() (reps.Block, error) {
	var genesisBlock reps.Block

	res := db.DB.
		Where("prev_hash = ''").
		First(&genesisBlock)
	if res.Error != nil {
		return reps.Block{}, res.Error
	}

	txns, err := repo.GetTransactionsByBlockId(genesisBlock.ID)
	if err != nil {
		return reps.Block{}, err
	}

	genesisBlock.Transactions = txns

	return genesisBlock, nil

}

func (repo *blockchainRepository) CreateBlock(block reps.Block) error {
	if err := db.DB.Create(&block).Error; err != nil {
		return err
	}

	return nil
}

// Get all blocks in blockchain
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
