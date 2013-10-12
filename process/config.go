package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func parseLine(line string) {
	var id, port int
	var ip string

	_, err := fmt.Sscanf(line, "%d %d %s", &id, &port, &ip)
	if err != nil {
		log.Fatal("Could not read the line: ", err)
	}
	fmt.Print("Read ", id, port, " : ", ip, "\n")
	addNode(getNode(id, ip, port))
}

func parseConfig() {
	fi, err := os.Open("config")
	if err != nil {
		log.Fatal("Could not open the config file", err)
	}
	defer func() {
		if err := fi.Close(); err != nil {
			log.Fatal("Could not close the config file", err)
		}
	}()
	r := bufio.NewReader(fi)
	for {
		buf, err := r.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("Error reading config file", err)
		}
		parseLine(string(buf))
	}
}
