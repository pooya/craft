package main

import (
	"fmt"
	"net/http"
)

const (
	FOLLOWER  = iota
	CANDIDATE = iota
	LEADER    = iota
)

var status int

func clientHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is: %d\n", status)
}

func startServer(port int) error {
	http.HandleFunc("/", clientHandler)
	status = LEADER

	strPort := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(strPort, nil)
}
