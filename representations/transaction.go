package representations

type Transaction struct {
	ID []byte 
	Inputs []TxnInput
	Outputs []TxnOutput
}

type ReadableTransaction struct {
	ID string `json:"id"`
	Inputs []TxnInput `json:"txnInputs"`
	Outputs []TxnOutput `json:"txnOutputs"`
}

// TxnID -> Id of transaction that the output is inside of
// OutIdx -> Index of an output in a transaction
// ScriptSig ->  Script which provides data to be used in an outputs ScriptPubKey 
type TxnInput struct {
	TxnID []byte `json:"txnId"`
	OutIdx int  `json:"outIdx"`
	ScriptSig string `json:"scriptSig"`
}

// Value -> Stores coins
// ScriptPubKey -> Value needed to unlock a transaction
type TxnOutput struct {
	Value int `json:"value"`
	ScriptPubKey string `json:"scriptPubKey"`
}