package main

import (
	"fmt"
)

var LatestLogEntry LogEntry

func processCommand(cmd int, serialNumber int) string {
	if serialNumber <= LatestLogEntry.serialNumber {
		/* TODO extract the response from the log and send it to the client. */
		return ""
	}
	response := LatestLogEntry.response + cmd
	fmt.Println("Response is: ", response)
	l := &LogEntry{LatestLogEntry.term, LatestLogEntry.index + 1, response, serialNumber}
	l.persist()
	LatestLogEntry = *l
	return fmt.Sprintf("%d", response)
}
