package utils

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"hash"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
)

// SplitCheckfile split a checkfile into an array
func SplitCheckfile(f *os.File) []string {
	d, err := ioutil.ReadAll(f)
	if err != nil {
		color.RedString(err.Error())
	}

	s := strings.Split(string(d), "\n")
	trim := s[:len(s)-1]

	return trim
}

// Compare check if the give checksum (c) is equal to the hash (h)
func Compare(c string, h hash.Hash) bool {
	d, err := hex.DecodeString(c)
	if err != nil {
		fmt.Println(err)
	}

	return bytes.Equal(d, h.Sum(nil))
}