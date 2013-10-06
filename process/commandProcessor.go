package main

import (
	"fmt"
)

var LatestLogEntry LogEntry

func digLogAndFindResponseFor(serialNumber int) (int, error) {
	l, err := getLogEntry(serialNumber)
	if err != nil {
		return 0, err
	}
	return l.response, nil
}

func processCommand(cmd int, serialNumber int) (int, error) {
	if serialNumber <= LatestLogEntry.serialNumber {
		return digLogAndFindResponseFor(serialNumber)
	}
	response := LatestLogEntry.response + cmd
	fmt.Println("Response is: ", response)
	l := &LogEntry{LatestLogEntry.term, LatestLogEntry.index + 1, response, serialNumber}
	l.persist()
	LatestLogEntry = *l
	return response, nil
}
