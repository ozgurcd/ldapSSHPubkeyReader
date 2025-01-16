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

	var config Config

	err := readConfig(&config)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(-1)
	}

	err = doSearch(&config, args[1])
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(-2)
	}
}
