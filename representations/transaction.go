package representations

type Transaction struct {
	ID []byte 
	Inputs []TxnInput
	Outputs []TxnOutput
}

// ID -> Id of transaction that the output is inside of
// OutIdx -> Index of an output in a transaction
// ScriptSig ->  Script which provides data to be used in an outputs ScriptPubKey 
type TxnInput struct {
	ID []byte
	OutIdx int 
	ScriptSig string
}

// Value -> Stores coins
// ScriptPubKey -> Value needed to unlock a transaction
type TxnOutput struct {
	Value int
	ScriptPubKey string
}