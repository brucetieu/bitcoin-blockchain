package representations

type Block struct {
	Timestamp int64  `json:"timestamp,omitempty"`
	Data      []byte  `json:"data,omitempty"`
	PrevHash  []byte   `json:"prevHash,omitempty"`
	Hash      []byte   `json:"hash,omitempty"`
	Nounce    int64    `json:"nounce,omitempty"`
}