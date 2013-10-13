package main

import (
	"fmt"
	"log"
	"net/http"
)

func getVoteTail() string {
	return fmt.Sprintf("%d/%s", LatestTerm, getMyUniqueId())
}

func voteFor(sender string) {
	node := findNode(sender)
	if node == nil {
		log.Fatal("Could not find the node we wanted to vote for.")
	}
	node.sendRequest(VotePath + getVoteTail())
}

func (node *Node) sendVoteRequest(term int) {
	node.sendRequest(VoteForMePath + getVoteTail())
}

func (node *Node) sendRequest(req string) {
	url := fmt.Sprintf("http://%s:%d%s", node.ip, node.port, req)
	fmt.Println("url is: ", url)
	resp, err := http.Get(url)
	if err != nil {
		//log.Print("Error sending request", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Print("sending request did not succeed: ", resp.Status)
	}
}

func sendVoteRequests(term int) {
	f := func(node *Node) {
		node.sendVoteRequest(term)
	}
	ForAllNodes(f)
}
