package main

import (
	"fmt"
	"os"
)

type LogEntry struct {
	term  int
	index int
	blob  interface{}
}

const (
	PersistLocation = "/tmp/persist/"
)

func (l *LogEntry) persist() {
	file, err := os.Open(PersistLocation + getMyUniqueId())
	if err != nil {
		fmt.Println("Fatail", err)
	}
	fmt.Fprintf(file, "%d|%d", l.term, l.index)
}

func accept(l LogEntry) {
	l.persist()
}
