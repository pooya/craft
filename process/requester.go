package main

import (
	"fmt"
	"log"
	"net/http"
)

func voteFor(sender string) {
	node := findNode(sender)
	if node == nil {
		/* TODO: we do not know this node, request for the nodes info. */
		/* Or, parse all of the nodes from a static file, and if node is nil,
		   we will not ignore it. */
	}
	node.sendRequest(VotePath + getMyUniqueId())
}

func (node *Node) sendVoteRequest(term int) {
	node.sendRequest(VoteForMePath + fmt.Sprint(term))
}

func (node *Node) sendRequest(req string) {
	url := fmt.Sprintf("http://%s:%d%s", node.ip, node.port, req)
	resp, err := http.Get(url)
	if err != nil {
		log.Print("Error sending request", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Print("sending request did not succeed: ", resp.Status)
	}
}

func sendVoteRequests(term int) {
	for _, node := range Nodes {
		node.sendVoteRequest(term)
	}
}
