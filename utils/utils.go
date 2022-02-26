package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

// Print a struct cleanly.
func PrettyPrintln(message string, data interface{}) {
	obj, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	fmt.Println(message, string(obj))
}

func Int64ToByte(i int64) []byte {
	return []byte(strconv.FormatInt(i, 10))
}
