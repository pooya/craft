package main

import (
	"fmt"
	"log"
	"net/http"
)

type Node struct {
	ip   string
	port int
}

func (node *Node) sendRequest(req string) {
	url := fmt.Sprintf("http://%s:%d/%s", node.ip, node.port, req)
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
