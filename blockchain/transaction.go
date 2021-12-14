package blockchain

import (
	"errors"
	"time"

	"github.com/JhyeonLee/BlockChain/utils"
)

const (
	minerReward int = 50
)

// mempool: unconfirmed transactions
type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

type Tx struct {
	ID        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

func (t *Tx) getId() {
	t.ID = utils.Hash(t)
}

type TxIn struct {
	TxID  string `json:"txID"`  // transaction that created output being used by input
	Index int    `json:"index"` // where TxID is
	Owner string `json:"owner"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

// Unspent Transaction Ouput: which TxOut has been spent or not
type UTxOut struct { // these are not user has not spent yet
	TxID   string
	Index  int
	Amount int
}

func isOnMempool(uTxOut *UTxOut) bool {
	// break loop no.1 : return in the loop
	// no.2 : break label in the loop ~> I choose it because this project is not just project but also a study
	exists := false
Outer:
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
				break Outer
			}
		}
	}
	return exists
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}

/* func makeTx(from, to string, amount int) (*Tx, error) {
	if Blockchain().BalanceByAddress(from) < amount { // from has not enough money for transaction
		return nil, errors.New("not enoguh money")
	}
	var txIns []*TxIn
	var txOuts []*TxOut
	// transaction input
	total := 0
	oldTxOuts := Blockchain().TxOutsByAddress(from)
	for _, txOut := range oldTxOuts {
		if total > amount { // enough transaction input
			break
		}
		txIn := &TxIn{txOut.Owner, txOut.Amount}
		txIns = append(txIns, txIn)
		total += txOut.Amount
	}
	// changes and transaction output
	change := total - amount
	if change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)
	// transaction
	tx := &Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return tx, nil
} */

func makeTx(from, to string, amount int) (*Tx, error) {
	// check if user has not enough money
	if BalanceByAddress(from, Blockchain()) < amount {
		return nil, errors.New("not enough money")
	}
	var txOuts []*TxOut
	var txIns []*TxIn

	// transaction input
	total := 0
	uTxOuts := UTxOutsByAddress(from, Blockchain())
	for _, uTxOut := range uTxOuts {
		if total >= amount {
			// enough transaction input that is more than amount needed
			break
		}
		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, from}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}

	// caculate change or not
	if change := total - amount; change != 0 {
		changetxOut := &TxOut{from, change}
		txOuts = append(txOuts, changetxOut)
	}
	// transactio output for to
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)
	// create a transaction
	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return tx, nil
}

// function that REST API will be calling
func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("jhyeon", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx("jhyeon")
	// all transactions
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil // make mempool empty
	return txs
}
