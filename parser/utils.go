package parser

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

// Read filename in "./temp/" folder
func ReadTempFile(filename string) (*os.File, error) {
	fullTempFilePath := fmt.Sprintf("./temp/%s", filename)
	f, err := os.Open(fullTempFilePath)
	if err != nil {
		return nil, fmt.Errorf("[parseXML] Error opening XML file: %w", err)
	}
	return f, nil
}

// Computes a SHA256 hash of a map[string]any data object.
func ComputeSHA256Hash(data map[string]interface{}) string {
	jsonString, _ := json.Marshal(data)
	hasher := sha256.New()
	hasher.Write(jsonString)
	return hex.EncodeToString(hasher.Sum(nil))
}
