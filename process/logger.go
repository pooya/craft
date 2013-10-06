package main

import (
	"fmt"
	"os"
)

type LogEntry struct {
	term         int
	index        int
	response     int
	serialNumber int
}

var logFile *os.File

const (
	PersistLocation = "/tmp/persist/"
)

func (l *LogEntry) persist() {
	fmt.Fprintf(logFile, "%d|%d|%d|%d\n", l.term, l.index, l.response, l.serialNumber)
}

func accept(l LogEntry) {
	l.persist()
}

func initLogger() error {
	file, err := os.OpenFile(PersistLocation+getMyUniqueId(),
		os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err == nil {
		logFile = file
	}
	return err
}
