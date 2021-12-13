package blockchain

import (
	"sync"

	"github.com/JhyeonLee/BlockChain/db"
	"github.com/JhyeonLee/BlockChain/utils"
)

const (
	defaultDifficulty  int = 2 // difficulty
	difficultyInterval int = 5 // how many blocks for a interval
	blockInterval      int = 2 // time per mining a block ex. 2 minutes per a block mined
	allowedRange       int = 2 // rooom for interval ~> make some range for time, not just certain time
)

// one blockchain
type blockchain struct {
	// blocks []*Block // without db
	// with db
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDIffuculty int    `json:"currentDifficulty"`
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
	b.CurrentDIffuculty = block.Difficulty
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

func (b *blockchain) recalculateDifficulty() int {
	allBlocks := b.Blocks()
	newestBlock := allBlocks[0]
	lastRecalculatedBlock := allBlocks[difficultyInterval-1]
	// minutes between newest and lastRecalculated
	actualTime := (newestBlock.Timestamp / 60) - (lastRecalculatedBlock.Timestamp / 60) // Because timestamp is UNIX, each is divided by 60
	expectedTime := difficultyInterval * blockInterval
	// too fst to mine : make harder ~> difficulty up // some rooms(ramge) for time
	if actualTime <= (expectedTime - allowedRange) {
		return b.CurrentDIffuculty + 1
	} else if actualTime >= (expectedTime + allowedRange) { // too slow to mine : make easier ~> difficulty down // some rooms(range) for time
		return b.CurrentDIffuculty - 1
	}
	return b.CurrentDIffuculty
}

// update difficulty per interval
func (b *blockchain) difficulty() int {
	if b.Height == 0 { // first block
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 { // on the interval
		// recalculate the difficulty
		return b.recalculateDifficulty()
	} else {
		return b.CurrentDIffuculty
	}
}

// GetBlockchain version "with db"
func Blockchain() *blockchain {
	if b == nil { // first time creating db, create *.db file
		once.Do(func() { // initializing blockchain once
			b = &blockchain{
				Height: 0,
			}
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
