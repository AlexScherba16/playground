package main

import (
	"log"
	"playground/internal/cli"
)

func main() {

	// Get parsed user input flags
	flags, err := cli.NewFlags()
	if err != nil {
		log.Fatalln(err.Error())
	}

	_ = flags
}
