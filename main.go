package main

import (
	"github.com/JhyeonLee/BlockChain/cli"
	"github.com/JhyeonLee/BlockChain/db"
)

func main() {
	defer db.Close()
	db.InitDB()
	cli.Start()

	// 1. BLOCKCHAIN CONCEPT
	/*
		chain := blockchain.GetBlockchain()
		chain.AddBlock("Second Block")
		chain.AddBlock("Third Block")
		chain.AddBlock("Fourth Block")
		for _, block := range chain.AllBlocks() {
			fmt.Printf("Data: %s\n", block.Data)
			fmt.Printf("Hash: %s\n", block.Hash)
			fmt.Printf("Prev Hash: %s\n\n", block.PrevHash)
		}
	*/

	// 2. EXPLORER WITH *.gohtml AND REST API
	// go explorer.Start(3000)
	// rest.Start(4000)

	// 3. CLI
	// cli.Start()

	// 4. Database
	// No need to call .Start() for db, like exploer.Start(), rest.Start(), and cli.Start()
	// because db has only connection with blockchain
	//blockchain.Blockchain()
	/* make block to restore db and to check it
	blockchain.Blockchain().AddBlock("Second Block")
	blockchain.Blockchain().AddBlock("Third Block")
	blockchain.Blockchain().AddBlock("Fourth Block")
	*/
	// defer db.Close()
	// cli.Start()

	// 5. PoW(Proof of Work) about Mining
	// defer db.Close()
	// cli.Start()

	// 6. Wallet
	// wallet.Start()
	// wallet.Wallet()
	// cli.Start()

	// 7. Wallet
	// defer db.Close()
	// cli.Start()

	// goRoutine, blocking precess, channel(unbuffered), buffered channel
	// c := make(chan int, 10) // buffer 5 like a queue
	// go send(c)
	// receive(c)

	// 8. P2P
	// defer db.Close()
	// cli.Start()
}

// When download a dependecy, using "sudo env "PATH=$PATH" go get -u github.com/..."
// ex> sudo env "PATH=$PATH" go get -u github.com/gorilla/mux
// if something error about go: ...: ...: permission denied -> using sudo env PATH=$PATH ...
// ex sudo env PATH=$PATH go run main.go
// or give user permit
// sudo chown -R <username> <folder path>

// both sending and receiving are blocking
// func send(c chan<- int) { // send only
// 	for i := range [10]int{} {
// 		fmt.Printf(">> sending %d <<\n", i)
// 		c <- i // blocking
// 		fmt.Printf(">> sent %d <<\n", i)
// 	}
// 	close(c)
// }

// func receive(c <-chan int) { // receive only(read only)
// 	for {
// 		time.Sleep(10 * time.Second)
// 		a, ok := <-c // blocking
// 		if !ok {
// 			fmt.Println("We are done")
// 			break
// 		}
// 		fmt.Printf("|| receiving %d ||\n", a)
// 	}
// }
