package main

import (
	"flag"
	"fmt"
)

var UniqueId string

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

func getMyUniqueId() string {
	return UniqueId
}

func main() {
	hostId, port, err := processCommandLineArguments()
	if err != nil {
		fmt.Println("Problem parsing arguments:", err)
		return
	}
	UniqueId = fmt.Sprintf("%d_%d", hostId, port)
	err = startServer(port)
	if err != nil {
		fmt.Println("Problem starting server", err)
		return
	}
}
