package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/JhyeonLee/BlockChain/db"
	"github.com/JhyeonLee/BlockChain/utils"
)

// one block
type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height   int    `json:"height"`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b)) // save block on db
}

var ErrNotFound = errors.New("blocks not found")

// When restoring, checkpoint is encoded and it should be docoded
func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}

func createBlock(data string, preHash string, height int) *Block {
	block := &Block{
		Data:     data,
		Hash:     "",
		PrevHash: preHash,
		Height:   height,
	}
	payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	block.persist() // save block on db
	return block
}
