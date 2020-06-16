package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args

	if len(args) == 1 {
		fmt.Printf("Missing argument username\n")
		os.Exit(-1)
	}

	// username is args[1]

	var config Config

	readConfig(&config)
	doSearch(&config, args[1])

}
