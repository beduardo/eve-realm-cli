package main

import (
	"fmt"
	"os"
)

var (
	Version   = "dev"
	GitHash   = "unknown"
	BuildDate = "unknown"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("eve-realm %s (git: %s, built: %s)\n", Version, GitHash, BuildDate)
		return
	}
	fmt.Println("eve-realm — thin client for the Eve Realm platform")
}
