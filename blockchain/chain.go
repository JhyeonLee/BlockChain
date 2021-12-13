package blockchain

import (
	"sync"

	"github.com/JhyeonLee/BlockChain/db"
	"github.com/JhyeonLee/BlockChain/utils"
)

// one blockchain
type blockchain struct {
	// blocks []*Block // without db
	// with db
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

var b *blockchain
var once sync.Once // initializing blockchain once

// When restoring, checkpoint was encoded and it should be docoded
func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func (b *blockchain) Blocks() []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

// GetBlockchain version "with db"
func Blockchain() *blockchain {
	if b == nil { // first time creating db, create *.db file
		once.Do(func() { // initializing blockchain once
			b = &blockchain{"", 0}
			checkpoint := db.Checkpoint() // search for checkpoint on db
			if checkpoint == nil {
				b.AddBlock("Genesis Block") //if there is no block on db, initialize db
			} else { // else, restore b from bytes
				b.restore(checkpoint)
			}
		})
	}
	return b
}

/*
// without db
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
*/
