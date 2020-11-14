package utils

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"
)

// SplitCheckfile - split a checkfile into an array
func SplitCheckfile(checkfileData []byte) map[string]string {
	checkList := make(map[string]string)

	rawLineEntries := strings.Split(string(checkfileData), "\n")
	// Avoid getting last line EOF
	lineEntries := rawLineEntries[:len(rawLineEntries)-1]

	for _, lineEntry := range lineEntries {
		lineEntryData := strings.Split(lineEntry, " ")
		checkList[lineEntryData[0]] = lineEntryData[1]
	}

	return checkList
}

// Compare - check if the give checksum is equal to the hash
func Compare(checksum string, hash hash.Hash) error {
	decodedChecksum, err := hex.DecodeString(checksum)
	if err != nil {
		fmt.Println(err)
	}

	if !bytes.Equal(decodedChecksum, hash.Sum(nil)) {
		return ErrInvalid
	}

	return nil
}
