package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

func PrettyPrintln(message string, data interface{}) {
	obj, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	fmt.Println(message, string(obj))
}
