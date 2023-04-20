package main

import (
	"fmt"
	"os"
)

var (
	version = "dev"
	commit  = ""
	date    = ""
)

func main() {
	fmt.Printf("Hello 👋 \n")
	fmt.Printf(" - version: %s\n", version)
	fmt.Printf(" - commit: %s\n", commit)
	fmt.Printf(" - date: %s\n", date)
	os.Exit(0)
}
