package handler

import (
	"fmt"
	"log"
	"net/http"

	"config"
	"logger"
	"node"
	"state"
)

func getStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is: %d\n", state.GetMyState())
}

func handleCommand(w http.ResponseWriter, r *http.Request) {

	idAndNumber := r.URL.Path[config.LenCommandPath:]
	var cmd, serialNumber int
	fmt.Sscanf(idAndNumber, "%d/%d", &cmd, &serialNumber)

	if !state.AmILeader() {
		if state.GetLeader() != nil {
			str := fmt.Sprintf("http://%s:%d%s%d/%d",
				state.GetLeader().Ip, state.GetLeader().Port, config.CommandPath, cmd, serialNumber)
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
	sender := r.URL.Path[config.LenHeartbeatPath:]
	log.Print("received heartbeat from ", sender)
	node := node.FindNode(sender)
	if node == nil {
		log.Fatal("sender is not part of config: ", sender)
	} else if state.AmILeader() && config.UniqueId != sender {
		panic("Got heartbeat from " + sender + ", but I am the leader")
	}
	state.SetLeader(node)
	state.HeartbeatChan <- sender
	fmt.Fprintf(w, "Got it %s.\n", sender)
}

func handleVoteRequest(w http.ResponseWriter, r *http.Request) {
	voteRequest := r.URL.Path[config.LenVoteForMePath:]

	var term int
	var sender string
	log.Print("received vote request: ", voteRequest)
	fmt.Sscanf(voteRequest, "%d/%s", &term, &sender)
	log.Print("received vote request from ", sender, " and term ", term)
	state.VoteIfEligible(sender, term)
}

func handleVote(w http.ResponseWriter, r *http.Request) {
	vote := r.URL.Path[config.LenVotePath:]
	var term int
	var sender string

	fmt.Sscanf(vote, "%d/%s", &term, &sender)
	log.Print("got vote from ", sender, " with term ", term)
	if term != logger.GetHighestTerm() {
		log.Print("Vote is stale, since latestTerm is: ", logger.GetHighestTerm())
		return
	}
	if state.GetMyState() == state.CANDIDATE {
		state.VoteChan <- sender
	} else {
		log.Printf("We are in state %d, and did not request votes\n",
			state.GetMyState())
	}
}

func Init(port int) {
	http.HandleFunc(config.StatusPath, getStatus)
	http.HandleFunc(config.CommandPath, handleCommand)
	http.HandleFunc(config.HeartbeatPath, handleHeartBeat)
	http.HandleFunc(config.VoteForMePath, handleVoteRequest)
	http.HandleFunc(config.VotePath, handleVote)
	strPort := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(strPort, nil)
	if err != nil {
		log.Fatal("Could not start server: ", err)
	}
}
