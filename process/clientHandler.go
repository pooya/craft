package main

import (
	"fmt"
	"net/http"
)

func getStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is: %d\n", getMyState())
}

func getLeader(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "not me")
}

const lenPath = len("/command/")

/*
Command is a number that we should set the status to.
    http://<>/command/id/number
*/
func processCommand(w http.ResponseWriter, r *http.Request) {
	idAndNumber := r.URL.Path[lenPath:]
	var id, number int
	fmt.Sscanf(idAndNumber, "%d/%d", &id, &number)
	fmt.Fprintf(w, "Command is: %d, id is: %d\n", id, number)
}

func startServer(port int) error {
	http.HandleFunc("/status", getStatus)
	http.HandleFunc("/leader", getLeader)
	http.HandleFunc("/command/", processCommand)
	strPort := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(strPort, nil)
}
