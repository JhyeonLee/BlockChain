package main

import (
	"crypto/sha256"
	"fmt"
)

type block struct {
	data     string
	hash     string
	prevHash string
}

func main() {
	// the first block
	genesisBlock := block{"Genesis Block", "", ""}
	// sha256 : hash function algorithm
	// type of input and ouput : byte slice
	hash := sha256.Sum256([]byte(genesisBlock.data + genesisBlock.prevHash))
	// convert byte type to base 16
	hexHash := fmt.Sprintf("%x", hash)
	genesisBlock.hash = hexHash

	fmt.Println(genesisBlock)
}
