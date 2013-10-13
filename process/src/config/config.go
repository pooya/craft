package config

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	StatusPath    = "/status/"    // sample 0:8080/status/
	CommandPath   = "/command/"   // sample 0:8080/command/11/20
	HeartbeatPath = "/heartbeat/" // sample 0:8080/heartbeat/2_8080
	VoteForMePath = "/voteforme/" // sample 0:8080/voteforme/1/3_8080
	VotePath      = "/vote/"      // sample 0:8080/vote/1/3_8080
)

const (
	LenStatusPath    = len(StatusPath)
	LenCommandPath   = len(CommandPath)
	LenHeartbeatPath = len(HeartbeatPath)
	LenVoteForMePath = len(VoteForMePath)
	LenVotePath      = len(VotePath)
)

var NProcesses int
var UniqueId string

type NodeHandler func(id int, ip string, port int)

var nodeHandler NodeHandler

func RegisterNodeHandler(handler NodeHandler) {
	nodeHandler = handler
}

func parseLine(line string) {
	var id, port int
	var ip string

	_, err := fmt.Sscanf(line, "%d %d %s", &id, &port, &ip)
	if err != nil {
		log.Fatal("Could not read the line: ", err)
	}
	fmt.Print("Read ", id, port, " : ", ip, "\n")
	nodeHandler(id, ip, port)
	NProcesses++
}

func Init(myId string) {
	UniqueId = myId
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
	NProcesses = 0
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
