package groups

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

func GenerateHash(data interface{}) [32]byte {
	jsonString, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error: %s", err)
		panic(err)
	}

	return sha256.Sum256([]byte(jsonString))
}
