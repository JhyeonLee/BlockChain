package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/JhyeonLee/BlockChain/explorer"
	"github.com/JhyeonLee/BlockChain/rest"
)

func usage() {
	fmt.Printf("Welcome to BlockChain Coin\n\n")
	fmt.Printf("Please use the following flags:\n\n")
	fmt.Printf("-port:     Set the PORT of the server\n")
	fmt.Printf("-mode:     Choose between 'html' and 'rest'\n\n")
	runtime.Goexit() // Goexit makes everything finished but leave defer to run
}

func Start() {
	if len(os.Args) < 2 {
		usage()
	}
	/*
		// Version 01: FlagSet
		rest := flag.NewFlagSet("rest", flag.ExitOnError)

		portFlag := rest.Int("port", 4000, "Sets the port of the server")

		switch os.Args[1] {
		case "explorer":
			fmt.Println("Start Explorer")
		case "rest":
			rest.Parse(os.Args[2:])
			// fmt.Println("Start REST API")
		default:
			usage()
		}

		if rest.Parsed() {
			fmt.Println(portFlag)
			fmt.Println("Start Server")
		}
	*/

	// Version 02
	port := flag.Int("port", 4000, "Set port of the server")
	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'")

	flag.Parse()
	switch *mode {
	// my code challenge
	case "both":
		go explorer.Start(*port)
		rest.Start(*port + 1)
	// end
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	default:
		usage()
	}
}
