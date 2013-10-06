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
	MIN_WAIT_BEFORE_CANDIDACY = 1000 // 1000ms
	MAX_WAIT_BEFORE_CANDIDACY = 5000 // 5000ms
)

var status int
var random *rand.Rand
var LatestEvent int64
var votes chan int

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

func transitionToLeader() {
	if status != CANDIDATE {
		panic("should be follower")
	}
	log.Print("I am the leader now.")
	status = LEADER
}

func captureVotes() {
	nVotes := 0
	for {
		x := <-votes
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

func transitionToCandidate() {
	if status != FOLLOWER {
		panic("should be follower")
	}
	log.Print("I am a candidate now.")
	status = CANDIDATE
	go captureVotes()
	votes <- 1
}

func transitionToFollower() {
	log.Print("I am a follower now.")
	status = FOLLOWER
}

func selectLeader() {
	heartbeat := true
	for {
		select {
		/*
		   case //TODO received heartbeat
		       heartbeat = true
		*/
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
	votes = make(chan int)
	go selectLeader()
}
