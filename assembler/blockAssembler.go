package assembler

import (
	"github.com/brucetieu/blockchain/representations"
	"fmt"
	"encoding/base64"
)

func ToBlockMap(block *representations.Block) map[string]interface{} {
	prevHash := fmt.Sprintf("%x", block.PrevHash)
	hash := fmt.Sprintf("%x", block.Hash)
	encodedContent := base64.StdEncoding.EncodeToString(block.Data)
	decodedContent, _ := base64.StdEncoding.DecodeString(encodedContent)
	data := make(map[string]interface{})
	data["timestamp"] = block.Timestamp
	data["prevHash"] = prevHash
	data["hash"] = hash
	data["nounce"] = block.Nounce
	data["content"] = string(decodedContent)
	return data
}