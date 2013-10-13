package node

import (
	"fmt"
	"log"
	"net/http"

	"config"
	"logger"
)

func getVoteTail() string {
	return fmt.Sprintf("%d/%s", logger.GetHighestTerm, config.UniqueId)
}

func VoteFor(sender string) {
	node := FindNode(sender)
	if node == nil {
		log.Fatal("Could not find the node we wanted to vote for.")
	}
	node.SendRequest(config.VotePath + getVoteTail())
}

func (node *Node) sendVoteRequest(term int) {
	node.SendRequest(config.VoteForMePath + getVoteTail())
}

func (node *Node) SendRequest(req string) {
	url := fmt.Sprintf("http://%s:%d%s", node.Ip, node.Port, req)
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

func SendVoteRequests() {
	f := func(node *Node) {
		node.sendVoteRequest(logger.GetNextTerm())
	}
	ForAll(f)
}
