package logger

import (
	"fmt"
	"log"
	"os"

	"config"
)

type LogEntry struct {
	Term         int
	Index        int
	Response     int
	SerialNumber int
}

var logFile *os.File
var LatestEntry *LogEntry
var latestTermInput, latestTermOutput chan int

const (
	PersistLocation = "/tmp/persist/"
)

func initLatestTerm(i chan int, o chan int) {
	latestTerm := 0
	for {
		select {
		case latestTerm = <-i:
		case o <- latestTerm:
		}
	}
}

func GetHighestTerm() int {
	return <-latestTermOutput
}

func SetHighestTerm(term int) {
	latestTermInput <- term
}

func IncrementNextTerm() {
	latestTermInput <- 1 + <-latestTermOutput
}

func GetNextTerm() int {
	return GetHighestTerm() + 1
}

func (l *LogEntry) Persist() {
	fmt.Fprintf(logFile, "%d|%d|%d|%d\n", l.Term, l.Index, l.Response, l.SerialNumber)
	LatestEntry = l
}

func GetLogEntry(serialNumber int) (*LogEntry, error) {
	if config.UniqueId == "" {
		panic("config not initialized yet.")
	}
	fmt.Println(config.UniqueId)
	file, err := os.Open(PersistLocation + config.UniqueId)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	for {
		var l LogEntry
		_, err := fmt.Fscanf(file, "%d|%d|%d|%d", &l.Term, &l.Index, &l.Response, &l.SerialNumber)
		if err != nil {
			return nil, err
		} else if l.SerialNumber == serialNumber {
			return &l, nil
		}
	}
}

func Init() {
	latestTermInput, latestTermOutput = make(chan int), make(chan int)
	go initLatestTerm(latestTermInput, latestTermOutput)
	file, err := os.OpenFile(PersistLocation+config.UniqueId,
		os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("Could not open log file")
	}
	logFile = file
}
