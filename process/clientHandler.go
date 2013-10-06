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

func main() {
	var port int
	port = 8080
	http.HandleFunc("/", clientHandler)
	status = LEADER

	strPort := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(strPort, nil)
	if err != nil {
		fmt.Println("Problem registering http handler", err)
	}
}
