package main

import (
	"fmt"
	"net/http"
)

func getStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is: %d\n", getMyState())
}

const lenPath = len("/command/")

/*
Command is a number that we should set the status to.
    http://<>/command/cmd/serialNumber
*/
func handleCommand(w http.ResponseWriter, r *http.Request) {
	if !amILeader() {
		http.Redirect(w, r, getLeader(), http.StatusFound)
		return
	}
	idAndNumber := r.URL.Path[lenPath:]
	var cmd, serialNumber int
	fmt.Sscanf(idAndNumber, "%d/%d", &cmd, &serialNumber)
	fmt.Fprintf(w, "Response: %s", processCommand(cmd, serialNumber))
}

func startServer(port int) error {
	http.HandleFunc("/status", getStatus)
	http.HandleFunc("/command/", handleCommand)
	strPort := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(strPort, nil)
}
