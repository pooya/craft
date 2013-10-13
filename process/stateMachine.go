package main

import (
	"log"
	"math/rand"
	"time"
)

const (
	FOLLOWER  = iota
	CANDIDATE = iota
	LEADER    = iota
)

const (
	CONNECTION_TIMEOUT        = 100   // 100ms
	HEARTBEAT_INTERVAL        = 1000  // 1s
	MIN_WAIT_BEFORE_CANDIDACY = 5000  // 1000ms
	MAX_WAIT_BEFORE_CANDIDACY = 10000 // 5000ms
)

var status int
var random *rand.Rand
var LatestEvent int64
var voteChan chan string
var Leader *Node

func getMyState() int {
	return status
}

func amILeader() bool {
	return status == LEADER
}

func getCandidacyTimeout() int {
	return random.Intn(MAX_WAIT_BEFORE_CANDIDACY-MIN_WAIT_BEFORE_CANDIDACY) +
		MIN_WAIT_BEFORE_CANDIDACY
}

func sendHeartBeats() {
	//TODO make the following nonblocking with a timeout
	requestSender := func(node *Node) {
		node.sendRequest(HeartbeatPath + getMyUniqueId())
	}
	for {
		if status != LEADER {
			return
		}
		ForAllNodes(requestSender)
		time.Sleep(HEARTBEAT_INTERVAL * time.Millisecond)
	}
}

func transitionToLeader() {
	if status != CANDIDATE {
		panic("should be follower")
	}
	log.Print("I am the leader now.")
	status = LEADER
	go sendHeartBeats()
}

func captureVotes() {
	nVotes := 0
	voters := make(map[string]bool)
	for {
		sender := <-voteChan
		if sender == "" {
			log.Print("Someone asked us not to be the leader. Stepping down.")
			transitionToFollower()
			return
		} else {
			if contains, ok := voters[sender]; !ok {
				if contains {
					log.Fatal("ok must match contains.")
				}
				nVotes++
				if nVotes > nProcesses/2 {
					transitionToLeader()
					return
				}
			}
		}
	}
}

func voteIfEligible(sender string, term int) {
	highestTerm := getHighestTerm()
	if term <= highestTerm {
		log.Print("Ignoring vote request for term: ", term,
			" since we are at ", highestTerm)
	} else {
		voteFor(sender)
		setHighestTerm(term)
	}
}

func transitionToCandidate() {
	if status == CANDIDATE {
		// restart the vote requests with a new term.
		voteChan <- ""
	} else if status == LEADER {
		log.Fatal("A leader should not be getting votes.")
	}
	log.Print("I am a candidate now.")
	status = CANDIDATE
	incrementNextTerm()
	go captureVotes()
	sendVoteRequests(getNextTerm())
	voteChan <- UniqueId
}

func transitionToFollower() {
	log.Print("I am a follower now.")
	status = FOLLOWER
}

func selectLeader() {
	heartbeat := true
	for {
		select {
		case <-heartbeatChan:
			log.Print("Got heartbeat.")
			voteChan <- ""
			heartbeat = true
		case <-time.After(time.Duration(getCandidacyTimeout()) * time.Millisecond):
			if !heartbeat {
				if status == LEADER {
					log.Print("No heartbeat from myself")
				} else {
					log.Print("Timer Expired, claim the throne")
					transitionToCandidate()
				}
			}
			heartbeat = false
		}
	}
}

func stateMachineInit() {
	transitionToFollower()
	random = rand.New(rand.NewSource(1))
	voteChan = make(chan string)
	go selectLeader()
}
