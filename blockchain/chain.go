package blockchain

import (
	"encoding/json"
	"net/http"
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
	m                 sync.Mutex
}

var b *blockchain
var once sync.Once // initializing blockchain once

// method : should mutate struct ~>ex. func (b *blockchain) AddBlock()
// if not, it is function ~> ex. func Blocks(b *blockchain) []*Block

// When restoring, checkpoint was encoded and it should be docoded
func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func persistBlockchain(b *blockchain) {
	db.SaveBlockchain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock() *Block {
	block := createBlock(b.NewestHash, b.Height+1, getDifficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDIffuculty = block.Difficulty
	persistBlockchain(b)
	return block
}

func Blocks(b *blockchain) []*Block {
	b.m.Lock()
	defer b.m.Unlock()
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

func Txs(b *blockchain) []*Tx {
	var txs []*Tx
	for _, block := range Blocks(b) {
		txs = append(txs, block.Transactions...)
	}
	return txs
}

func FindTx(b *blockchain, targetID string) *Tx {
	for _, tx := range Txs(b) {
		if tx.ID == targetID {
			return tx
		}
	}
	return nil
}

func recalculateDifficulty(b *blockchain) int {
	allBlocks := Blocks(b)
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
func getDifficulty(b *blockchain) int {
	if b.Height == 0 { // first block
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 { // on the interval
		// recalculate the difficulty
		return recalculateDifficulty(b)
	} else {
		return b.CurrentDIffuculty
	}
}

/* // return all transaction outputs
func (b *blockchain) txOuts() []*TxOut {
	var txOuts []*TxOut
	blocks := b.Blocks()
	for _, block := range blocks {
		for _, tx := range block.Transactions {
			txOuts = append(txOuts, tx.TxOuts...)
		}
	}
	return txOuts
} */

/* // return transanction outputs filtered by address(owner)
// called in router
func (b *blockchain) TxOutsByAddress(address string) []*TxOut {
	var ownedTxOuts []*TxOut
	txOuts := b.txOuts()
	for _, txOut := range txOuts {
		if txOut.Owner == address {
			ownedTxOuts = append(ownedTxOuts, txOut)
		}
	}
	return ownedTxOuts
} */

// Unspent Transaction Ouputs by Address
// return outputs that have not been used by inputs yet
func UTxOutsByAddress(address string, b *blockchain) []*UTxOut {
	var uTxOuts []*UTxOut
	creatorTxs := make(map[string]bool) // spent Transaction Outputs

	for _, block := range Blocks(b) {
		for _, tx := range block.Transactions {
			// checking on Transaction Inputs
			for _, input := range tx.TxIns {
				if input.Signature == "COINBASE" {
					break
				}
				if FindTx(b, input.TxID).TxOuts[input.Index].Address == address {
					// which tranaction ID create output being used as input
					creatorTxs[input.TxID] = true
				}
			}
			// checking on Transaction Outputs
			for index, output := range tx.TxOuts {
				if output.Address == address {
					if _, ok := creatorTxs[tx.ID]; !ok { // if tx.ID is not spented: unspented
						uTxOut := &UTxOut{tx.ID, index, output.Amount}
						if !isOnMempool(uTxOut) { // check it whether is on mempool
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}
	return uTxOuts
}

func BalanceByAddress(address string, b *blockchain) int {
	txOuts := UTxOutsByAddress(address, b)
	var amount int
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}

// GetBlockchain version "with db"
func Blockchain() *blockchain {
	// first time creating db, create *.db file
	// b == nil occurs only one
	// and this function Blockchain() is also called only one
	// so delete if b == nil {}
	// but it will break code literally, problem like loading forever when request /blocks
	// , called "deadlock" which is application cannot continue
	// it is because createBlock() call Blockchain() and AddBlock() also call Blockchain and repeat it
	// it is fixed by not call Blockchain().difficulty() on createBlock()
	once.Do(func() { // initializing blockchain once
		b = &blockchain{
			Height: 0,
		}
		checkpoint := db.Checkpoint() // search for checkpoint on db
		if checkpoint == nil {
			b.AddBlock() //if there is no block on db, initialize db
		} else { // else, restore b from bytes
			b.restore(checkpoint)
		}
	})
	return b
}

func Status(b *blockchain, rw http.ResponseWriter) {
	b.m.Lock()
	defer b.m.Unlock()

	utils.HandleErr(json.NewEncoder(rw).Encode(b))
}

func (b *blockchain) Replace(newBlocks []*Block) {
	b.m.Lock()
	defer b.m.Unlock()
	// replacing to new blokcchian: mutate and persist blockchain
	b.CurrentDIffuculty = newBlocks[0].Difficulty
	b.Height = len(newBlocks)
	b.NewestHash = newBlocks[0].Hash
	persistBlockchain(b)

	db.EmptyBlocks()
	for _, block := range newBlocks {
		persistBlock(block)
	}
}

func (b *blockchain) AddpeerBlock(newBlock *Block) {
	b.m.Lock()
	m.m.Lock()
	defer b.m.Unlock()
	defer m.m.Unlock()

	b.Height += 1
	b.CurrentDIffuculty = newBlock.Difficulty
	b.NewestHash = newBlock.Hash
	persistBlockchain(b)
	persistBlock(newBlock)

	// mempool
	for _, tx := range newBlock.Transactions {
		_, ok := m.Txs[tx.ID]
		if ok {
			delete(m.Txs, tx.ID)
		}
	}
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
