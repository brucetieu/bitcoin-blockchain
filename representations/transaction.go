package representations

// import "github.com/google/uuid"

// ID -> Unique id of this transaction
// BlockID -> Which block is this transaction in?
// Inputs and Outputs -> In both these tables, curr_txn_id is equal to id of transaction. This helps us to track which transaction did these inputs and outputs come from
type Transaction struct {
	ID      []byte      `json:"txnId" gorm:"primary_key"`
	BlockID string      `json:"blockId"`
	Inputs  []TxnInput  `json:"txnInputs" gorm:"foreignKey:CurrTxnID;association_foreignkey:ID"`
	Outputs []TxnOutput `json:"txnOutputs" gorm:"foreignKey:CurrTxnID;association_foreignkey:ID"`
}

type ReadableTransaction struct {
	ID      string              `json:"id"`
	BlockID string      `json:"blockId"`
	Inputs  []ReadableTxnInput  `json:"txnInputs"`
	Outputs []ReadableTxnOutput `json:"txnOutputs"`
}

type ReadableTxnInput struct {
	CurrTxnID string `json:"currTxnId"`
	PrevTxnID string `json:"prevTxnId"`
	OutIdx    int    `json:"outIdx"`
	PubKey string    `json:"pubKey"`
	Signature string `json:"signature"`
}

type ReadableTxnOutput struct {
	CurrTxnID    string `json:"currTxnId"`
	Value        int    `json:"value"`
	PubKeyHash string `json:"pubKeyHash"`
}

// InputID -> unique id of the TxnInput
// CurrTxnId -> What transaction is this input currently in?
// OutIdx -> From which output index was used to create this input?
// PrevTxnID -> From which previous transaction was and ouptut used to create this input?
// ScriptSig ->  Script which provides data to be used in an outputs ScriptPubKey
type TxnInput struct {
	InputID string `json:"inputId" gorm:"primary_key"`

	CurrTxnID []byte `json:"currTxnId" gorm:"column:curr_txn_id"`
	PrevTxnID []byte `json:"prevTxnId" gorm:"column:prev_txn_id"`
	OutIdx    int    `json:"outIdx"`
	// ScriptSig string `json:"scriptSig"`
	Signature []byte `json:"signature"`// signature of the entire transaction
	PubKey []byte `json:"pubKey"` // not hashed
}

// OutputID -> Unique id representing the output
// CurrTxnID -> What transaction is this output currently in?
// Value -> Stores coins
// ScriptPubKey -> Value needed to unlock a transaction
type TxnOutput struct {
	OutputID string `json:"outputId" gorm:"primary_key"`

	CurrTxnID    []byte `json:"currTxnId" gorm:"column:curr_txn_id"`
	Value        int    `json:"value"`
	PubKeyHash   []byte `json:"pubKeyHash"`// locks the output
	// ScriptPubKey string `json:"scriptPubKey"`
}
