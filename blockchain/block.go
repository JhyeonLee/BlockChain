package blockchain

import (
	"errors"
	"strings"
	"time"

	"github.com/JhyeonLee/BlockChain/utils"
)

// one block
type Block struct {
	// Data       string `json:"data"`
	Hash       string `json:"hash"` // same ipnut same output, one-way, determinant
	PrevHash   string `json:"prevHash,omitempty"`
	Height     int    `json:"height"`
	Difficulty int    `json:"difiiculty"` // for POW, ex. how many 0s in front part of hash
	Nounce     int    `json:"nounce"`     // for POW, only can be changed
	Timestamp  int    `json:"timestamp"`

	Transactions []*Tx `json:"transactions"` // coinbase transaction: transaction for miner, created by blockchain
}

func persistBlock(b *Block) {
	dbStorage.SaveBlock(b.Hash, utils.ToBytes(b)) // save block on db
}

var ErrNotFound = errors.New("blocks not found")

// When restoring, checkpoint is encoded and it should be docoded
func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

func FindBlock(hash string) (*Block, error) {
	blockBytes := dbStorage.FindBlock(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		b.Timestamp = int(time.Now().Unix())
		hash := utils.Hash(b)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nounce++
		}
	}
}

func createBlock(preHash string, height, diff int) *Block {
	block := &Block{
		//Data:       data,
		Hash:       "",
		PrevHash:   preHash,
		Height:     height,
		Difficulty: diff,
		Nounce:     0,
	}
	// payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)
	// block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	block.Transactions = Mempool().TxToConfirm()
	block.mine()
	persistBlock(block) // save block on db
	return block
}
