package main

import (
	"fmt"
	"log"
	"net/http"
)

const port string = ":4000"

// rw : data to user
// r : request
func home(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprint(rw, "Hello from home!")
}

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

	http.HandleFunc("/", home)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
