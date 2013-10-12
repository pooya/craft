package main

import (
	"fmt"
	"log"
	"net/http"
)

func getStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is: %d\n", getMyState())
}

const (
	StatusPath    = "/status/"
	CommandPath   = "/command/"
	HeartbeatPath = "/heartbeat/"
	VoteForMePath = "/voteforme/"
)

const (
	lenStatusPath    = len(StatusPath)
	lenCommandPath   = len(CommandPath)
	lenHeartbeatPath = len(HeartbeatPath)
	lenVoteForMePath = len(CommandPath)
)

/*
Command is a number that we should set the status to.
    http://<>/command/cmd/serialNumber
*/
func handleCommand(w http.ResponseWriter, r *http.Request) {
	if !amILeader() {
		http.Redirect(w, r, getLeader(), http.StatusFound)
		return
	}
	idAndNumber := r.URL.Path[lenCommandPath:]
	var cmd, serialNumber int
	fmt.Sscanf(idAndNumber, "%d/%d", &cmd, &serialNumber)
	resp, err := processCommand(cmd, serialNumber)

	if err != nil {
		http.NotFound(w, r)
	} else {
		fmt.Fprintf(w, "Response: %d\n", resp)
	}
}

func handleHeartBeat(w http.ResponseWriter, r *http.Request) {
	sender := r.URL.Path[lenHeartbeatPath:]
	log.Print("received heartbeat from ", sender)
	heartbeatChan <- true
	fmt.Fprintf(w, "Got it %s.\n", sender)
}

func handleVoteRequest(w http.ResponseWriter, r *http.Request) {
	voteRequest := r.URL.Path[lenVoteForMePath:]

	var term int
	var sender string
	fmt.Sscanf(voteRequest, "%d/%s", &term, &sender)
	log.Print("received vote request from ", sender, " and term ", term)
	voteIfEligible(sender, term)
}

func startServer(port int) error {
	http.HandleFunc(StatusPath, getStatus)
	http.HandleFunc(CommandPath, handleCommand)
	http.HandleFunc(HeartbeatPath, handleHeartBeat)
	http.HandleFunc(VoteForMePath, handleVoteRequest)
	strPort := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(strPort, nil)
}
