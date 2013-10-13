package main

import (
	"flag"
	"fmt"
	"log"
    "logger"
    "state"
    "config"
    "node"
    "handler"
)

func processCommandLineArguments() (int, int, error) {
	var hostId = flag.Int("hostId", 0, "The unique id for the host.")
	var port = flag.Int("port", 0, "The per-host unique port")

	flag.Parse()

	if *hostId == 0 || *port == 0 {
		return 0, 0,
			fmt.Errorf("Cannot proceed with hostId %d and port %d\n"+
				"Usage: ./identity -hostId hostId -port port", *hostId, *port)
	}
	return *hostId, *port, nil
}

func handleCommandLine() (int, int) {
	hostId, port, err := processCommandLineArguments()
	if err != nil {
        log.Fatal("Problem parsing arguments:", err)
	}
    return hostId, port
}

func main() {
    hostId, port := handleCommandLine()
	node.Init()
	config.Init(fmt.Sprintf("%d_%d", hostId, port))
	logger.Init()
	state.Init()
	handler.Init(port)
}
