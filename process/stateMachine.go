package main

const (
	FOLLOWER  = iota
	CANDIDATE = iota
	LEADER    = iota
)

var status int

func getMyState() int {
	return status
}

func init() {
	status = LEADER
}
