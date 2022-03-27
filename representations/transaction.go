package representations

// import "github.com/google/uuid"

type Transaction struct {
	ID     []byte   `json:"txnId" gorm:"primary_key"`
	BlockID string `json:"blockId" `
	// ID string `gorm:"primary_key;type:char(36);"` // references block id
	// TxnID      []byte `json:"txnId" gorm:"primary_key"`
	Inputs  []TxnInput `json:"txnInputs" gorm:"foreignKey:TxnID;association_foreignkey:ID;PRELOAD:false"`
	Outputs []TxnOutput `json:"txnOutputs" gorm:"foreignKey:TxnID;association_foreignkey:ID;PRELOAD:false"`
}

type ReadableTransaction struct {
	ID      string             `json:"id"`
	Inputs  []ReadableTxnInput `json:"txnInputs"`
	Outputs []TxnOutput        `json:"txnOutputs"`
}

// TxnID -> Id of transaction that the output is inside of
// OutIdx -> Index of an output in a transaction
// ScriptSig ->  Script which provides data to be used in an outputs ScriptPubKey
type TxnInput struct {
	InputID string `json:"inputId" gorm:"primary_key"`
	TxnID     []byte `json:"txnId" gorm:"column:txn_id"`
	OutIdx    int    `json:"outIdx"`
	ScriptSig string `json:"scriptSig"`
}

type ReadableTxnInput struct {
	TxnID     string `json:"txnId"`
	OutIdx    int    `json:"outIdx"`
	ScriptSig string `json:"scriptSig"`
}

// Value -> Stores coins
// ScriptPubKey -> Value needed to unlock a transaction
type TxnOutput struct {
	OutputID string `json:"outputId" gorm:"primary_key"`
	TxnID []byte `json:"txnId" gorm:"column:txn_id"` 
	Value        int    `json:"value"`
	ScriptPubKey string `json:"scriptPubKey"`
}
