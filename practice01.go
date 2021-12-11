package main

import (
	"fmt"

	"github.com/JhyeonLee/BlockChain/person"
)

func main() {

	x := person.Person{}
	x.SetDetails("jhyeon", 100)

	fmt.Println(x)

	fmt.Println(x.Name())
}
