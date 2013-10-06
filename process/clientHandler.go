package main

import (
	"fmt"
	"net/http"
)

func clientHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is: %d\n", getMyState())
}

func startServer(port int) error {
	http.HandleFunc("/", clientHandler)
	strPort := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(strPort, nil)
}
