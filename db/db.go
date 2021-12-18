package db

import (
	"fmt"
	"os"

	"github.com/JhyeonLee/BlockChain/utils"
	bolt "go.etcd.io/bbolt"
)

const (
	dbName       = "blockchain"
	dataBucket   = "data"   // saving a data(last hash)
	blocksBucket = "blocks" // saving all blocks

	checkpoint = "checkpoint" // restoring db
)

var db *bolt.DB

func getDbName() string {
	/* for i, a := range os.Args {
		fmt.Println(i, a) // a[1][6:] : port
	} */
	// fmt.Println(os.Args[1][6:]) //port

	port := os.Args[2][6:]
	return fmt.Sprintf("%s_%s.db", dbName, port) // blockchain_{port}.db

}

func DB() *bolt.DB {
	if db == nil {
		// if db is not exist, Initialize db
		// Open db : path, permission, oprion
		dbPointer, err := bolt.Open(getDbName(), 0600, nil)
		utils.HandleErr(err)
		db = dbPointer

		// 2 Buckets : saving a data(last hash), and saving all blocks
		err = db.Update(func(t *bolt.Tx) error { // Read-Write transactions
			_, err := t.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleErr(err)
			_, err = t.CreateBucketIfNotExists([]byte(blocksBucket))
			return err
		})
		utils.HandleErr(err)
	}
	return db
}

func Close() {
	DB().Close()
}

func SaveBlock(hash string, data []byte) {
	// fmt.Printf("Saving Block %s\nData: %b\n", hash, data)
	err := DB().Update(func(t *bolt.Tx) error { // read-write transaction
		bucket := t.Bucket([]byte(blocksBucket)) // call blocksBucket
		err := bucket.Put([]byte(hash), data)    // save hash(key) : data(value) pair
		return err
	})
	utils.HandleErr(err)
}

func SaveBlockchain(data []byte) {
	// first time creating db, create *.db file
	err := DB().Update(func(t *bolt.Tx) error { // read-write transaction
		bucket := t.Bucket([]byte(dataBucket))      // call dtaBucket
		err := bucket.Put([]byte(checkpoint), data) // save last block
		return err
	})
	utils.HandleErr(err)
}

func Checkpoint() []byte { // search from dataBucket
	var data []byte
	DB().View(func(t *bolt.Tx) error { // read-only transaction
		bucket := t.Bucket([]byte(dataBucket)) // call dataBucket
		data = bucket.Get([]byte(checkpoint))  // get data of bucket checkpoint
		return nil
	})
	return data
}

func Block(hash string) []byte { // search from blocksBucket
	var data []byte
	DB().View(func(t *bolt.Tx) error { // read-only transaction
		bucket := t.Bucket([]byte(blocksBucket)) // call blocksBucket
		data = bucket.Get([]byte(hash))          // get data of bucket checkpoint
		return nil
	})
	return data
}

func EmptyBlocks() {
	DB().Update(func(t *bolt.Tx) error {
		utils.HandleErr(t.DeleteBucket([]byte(blocksBucket)))
		_, err := t.CreateBucket([]byte(blocksBucket))
		utils.HandleErr(err)
		return nil
	})
}
