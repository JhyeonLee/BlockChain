package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
)

// one block
type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Heights  int    `json:"height"`
}

// one blockchain
type blockchain struct {
	blocks []*Block
}

var b *blockchain
var once sync.Once // initializing blockchain once

// calculate hash : hash(i) = fn( data(i) + hash(i-1) )
func (b *Block) calculateHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

// get last block's hash(previous hash)
func getLastHash() string {
	totalBlocks := len(GetBlockchain().blocks)
	if totalBlocks == 0 {
		return ""
	}
	return GetBlockchain().blocks[totalBlocks-1].Hash
}

// create a block
func createBlock(data string) *Block {
	newBlock := Block{data, "", getLastHash(), len(GetBlockchain().blocks) + 1}
	newBlock.calculateHash()
	return &newBlock
}

// add a block on blockchain
func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

// Singleton Pattern to initialize blockchain
func GetBlockchain() *blockchain {
	if b == nil {
		once.Do(func() { // initializing blockchain once
			b = &blockchain{}
			b.AddBlock("Genesis Block")
		})
	}
	return b
}

func (b *blockchain) AllBlocks() []*Block {
	return b.blocks
}

var ErrNotFound = errors.New("block not found")

func (b *blockchain) GetBlock(height int) (*Block, error) {
	if height > len(b.blocks) {
		return nil, ErrNotFound
	}
	return b.blocks[height-1], nil
}
