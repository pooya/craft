package main

import (
	"flag"
	"fmt"
)

var hostId = flag.Int("hostId", 0, "The unique id for the host.")
var port = flag.Int("port", 0, "The per-host unique port")

func printUsage() {
	fmt.Println("Usage: ./identity hostId port")
}

func main() {
	flag.Parse()

	if *hostId == 0 || *port == 0 {
		fmt.Println("Cannot proceed hostId ", *hostId, ", and port ", *port)
		printUsage()
		return
	}
	fmt.Println("Accepting ", *hostId, ":", *port)
}
