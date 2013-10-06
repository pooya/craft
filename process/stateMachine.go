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

func amILeader() bool {
	return status == LEADER
}

func getLeader() string {
	return "the other guy"
}

func init() {
	status = LEADER
}
