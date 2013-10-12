package main

import (
	"fmt"
	"log"
	"os"
)

type LogEntry struct {
	term         int
	index        int
	response     int
	serialNumber int
}

var logFile *os.File
var LatestEntry *LogEntry

const (
	PersistLocation = "/tmp/persist/"
)

func getHighestTerm() int {
	if LatestEntry != nil {
		return LatestEntry.term
	}
	return -1
}

func getNextTerm() int {
	return getHighestTerm() + 1
}

func (l *LogEntry) persist() {
	fmt.Fprintf(logFile, "%d|%d|%d|%d\n", l.term, l.index, l.response, l.serialNumber)
	LatestEntry = l
}

func getLogEntry(serialNumber int) (*LogEntry, error) {
	file, err := os.Open(PersistLocation + getMyUniqueId())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	for {
		var l LogEntry
		_, err := fmt.Fscanf(file, "%d|%d|%d|%d", &l.term, &l.index, &l.response, &l.serialNumber)
		if err != nil {
			return nil, err
		} else if l.serialNumber == serialNumber {
			return &l, nil
		}
	}
}

func initLogger() error {
	file, err := os.OpenFile(PersistLocation+getMyUniqueId(),
		os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err == nil {
		logFile = file
	}
	return err
}
