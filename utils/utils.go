package utils

import (
	"encoding/json"
	"strconv"
	log "github.com/sirupsen/logrus"
)

// Print a struct cleanly.
func PrettyPrintln(message string, data interface{}) {
	obj, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	log.Info(message, string(obj))
}

// Make a struct look pretty
func Pretty(data interface{}) string {
	obj, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	return string(obj)
}

// Convert int64 to slice pf bytes
func Int64ToByte(i int64) []byte {
	return []byte(strconv.FormatInt(i, 10))
}
