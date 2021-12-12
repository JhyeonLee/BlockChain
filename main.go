package main

import (
	"github.com/JhyeonLee/BlockChain/rest"
)

func main() {
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

	// go explorer.Start(3000)
	rest.Start(4000)
}

// When download a dependecy, using "sudo env "PATH=$PATH" go get -u github.com/..."
// ex> sudo env "PATH=$PATH" go get -u github.com/gorilla/mux
// if something error about go: ...: ...: permission denied -> using sudo env PATH=$PATH ...
// ex sudo env PATH=$PATH go run main.go
