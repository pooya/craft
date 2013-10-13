package handler

import (
	"fmt"

	"logger"
)

var LatestLogEntry logger.LogEntry

func digLogAndFindResponseFor(serialNumber int) (int, error) {
	l, err := logger.GetLogEntry(serialNumber)
	if err != nil {
		return 0, err
	}
	return l.Response, nil
}

func processCommand(cmd int, serialNumber int) (int, error) {
	if serialNumber <= LatestLogEntry.SerialNumber {
		return digLogAndFindResponseFor(serialNumber)
	}
	response := LatestLogEntry.Response + cmd
	fmt.Println("Response is: ", response)
	l := &logger.LogEntry{LatestLogEntry.Term, LatestLogEntry.Index + 1, response, serialNumber}
	l.Persist()
	LatestLogEntry = *l
	return response, nil
}
