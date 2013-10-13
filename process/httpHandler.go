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
	StatusPath    = "/status/"    // sample 0:8080/status/
	CommandPath   = "/command/"   // sample 0:8080/command/11/20
	HeartbeatPath = "/heartbeat/" // sample 0:8080/heartbeat/2_8080
	VoteForMePath = "/voteforme/" // sample 0:8080/voteforme/1/3_8080
	VotePath      = "/vote/"      // sample 0:8080/vote/1/3_8080
)

const (
	lenStatusPath    = len(StatusPath)
	lenCommandPath   = len(CommandPath)
	lenHeartbeatPath = len(HeartbeatPath)
	lenVoteForMePath = len(VoteForMePath)
	lenVotePath      = len(VotePath)
)

func handleCommand(w http.ResponseWriter, r *http.Request) {

	idAndNumber := r.URL.Path[lenCommandPath:]
	var cmd, serialNumber int
	fmt.Sscanf(idAndNumber, "%d/%d", &cmd, &serialNumber)

	if !amILeader() {
		if Leader != nil {
			str := fmt.Sprintf("http://%s:%d%s%d/%d",
				Leader.ip, Leader.port, CommandPath, cmd, serialNumber)
			log.Print(str)
			http.Redirect(w, r, str, http.StatusSeeOther)
		} else {
			http.NotFound(w, r)
		}
		return
	}
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
	node := findNode(sender)
	if node == nil {
		log.Print("sender is not part of config: ", sender)
	}
	Leader = node
	heartbeatChan <- true
	fmt.Fprintf(w, "Got it %s.\n", sender)
}

func handleVoteRequest(w http.ResponseWriter, r *http.Request) {
	voteRequest := r.URL.Path[lenVoteForMePath:]

	var term int
	var sender string
	log.Print("received vote request: ", voteRequest)
	fmt.Sscanf(voteRequest, "%d/%s", &term, &sender)
	log.Print("received vote request from ", sender, " and term ", term)
	voteIfEligible(sender, term)
}

func handleVote(w http.ResponseWriter, r *http.Request) {
	vote := r.URL.Path[lenVotePath:]
	var term int
	var sender string

	fmt.Sscanf(vote, "%d/%s", &term, &sender)
	log.Print("got vote from ", sender, " with term ", term)
    if term != LatestTerm {
        log.Print("Vote is stale, since latestTerm is: ", LatestTerm)
        return
    }
    voteChan <- sender
}

func startServer(port int) error {
	http.HandleFunc(StatusPath, getStatus)
	http.HandleFunc(CommandPath, handleCommand)
	http.HandleFunc(HeartbeatPath, handleHeartBeat)
	http.HandleFunc(VoteForMePath, handleVoteRequest)
	http.HandleFunc(VotePath, handleVote)
	strPort := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(strPort, nil)
}
