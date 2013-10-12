package main

import (
	"log"
	"math/rand"
	"time"
)

const (
	nProcesses = 1
)

const (
	FOLLOWER  = iota
	CANDIDATE = iota
	LEADER    = iota
)

const (
	CONNECTION_TIMEOUT        = 100  // 100ms
	HEARTBEAT_INTERVAL        = 1000 // 1s
	MIN_WAIT_BEFORE_CANDIDACY = 2000 // 1000ms
	MAX_WAIT_BEFORE_CANDIDACY = 5000 // 5000ms
)

var status int
var random *rand.Rand
var LatestEvent int64
var voteChan chan int

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

func getLeader() string {
	return "the other guy"
}

func sendHeartBeats() {
	for {
		if status != LEADER {
			return
		}
		log.Print("sending heartbeat")
		node := getNode(0, "0", 8080)
		addNode(node)

		//TODO make the following nonblocking with a timeout
		node.sendRequest(HeartbeatPath + getMyUniqueId())
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
	for {
		x := <-voteChan
		if x == 0 {
			log.Print("Someone asked us not to be the leader. Stepping down.")
			transitionToFollower()
			return
		} else {
			nVotes++
			if nVotes > nProcesses/2 {
				transitionToLeader()
				return
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
		//voteFor(sender)
	}
}

func transitionToCandidate() {
	if status != FOLLOWER {
		panic("should be follower")
	}
	log.Print("I am a candidate now.")
	status = CANDIDATE
	go captureVotes()
	sendVoteRequests(getNextTerm())
	voteChan <- 1
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
	voteChan = make(chan int)
	go selectLeader()
}
