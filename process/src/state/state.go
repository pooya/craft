package state

import (
	"log"
	"math/rand"
	"time"

	"config"
	"logger"
	"node"
)

const (
	FOLLOWER  = iota
	CANDIDATE = iota
	LEADER    = iota
)

var HeartbeatChan chan string
var statusInput, statusOutput chan int
var leaderInput, leaderOutput chan *node.Node

const (
	CONNECTION_TIMEOUT        = 100   // 100ms
	HEARTBEAT_INTERVAL        = 1000  // 1s
	MIN_WAIT_BEFORE_CANDIDACY = 5000  // 1000ms
	MAX_WAIT_BEFORE_CANDIDACY = 10000 // 5000ms
)

var random *rand.Rand
var VoteChan chan string

func initStatus(i chan int, o chan int) {
	status := FOLLOWER
	for {
		select {
		case status = <-i:
		case o <- status:
		}
	}
}

func initLeader(i chan *node.Node, o chan *node.Node) {
	var leader *node.Node
	leader = nil
	for {
		select {
		case leader = <-i:
		case o <- leader:
		}
	}
}

func GetMyState() int {
	return <-statusOutput
}

func AmILeader() bool {
	return GetMyState() == LEADER
}

func GetLeader() *node.Node {
	return <-leaderOutput
}

func getCandidacyTimeout() int {
	return random.Intn(MAX_WAIT_BEFORE_CANDIDACY-MIN_WAIT_BEFORE_CANDIDACY) +
		MIN_WAIT_BEFORE_CANDIDACY
}

func sendHeartBeats() {
	//TODO make the following nonblocking with a timeout
	requestSender := func(node *node.Node) {
		node.SendRequest(config.HeartbeatPath + config.UniqueId)
	}
	for {
		if GetMyState() != LEADER {
			return
		}
		node.ForAll(requestSender)
		time.Sleep(HEARTBEAT_INTERVAL * time.Millisecond)
	}
}

func transitionToLeader() {
	if GetMyState() != CANDIDATE {
		panic("should be follower")
	}
	log.Print("I am the leader now.")
	statusInput <- LEADER
	go sendHeartBeats()
}

func captureVotes() {
	nVotes := 0
	voters := make(map[string]bool)
	for {
		sender := <-VoteChan
		if sender == "" {
			log.Print("Someone asked us not to be the leader. Stepping down.")
			transitionToFollower()
			return
		} else {
			if _, ok := voters[sender]; ok {
				log.Print("vote from ", sender, " is already processed")
				continue
			}
			if node.FindNode(sender) == nil {
				log.Print("Received vote from unknown sender: ", sender)
				continue
			}
			nVotes++
			if nVotes > config.NProcesses/2 {
				transitionToLeader()
				return
			}
		}
	}
}

func VoteIfEligible(sender string, term int) {
	highestTerm := logger.GetHighestTerm()
	if term <= highestTerm {
		log.Print("Ignoring vote request for term: ", term,
			" since we are at ", highestTerm)
	} else {
		node.VoteFor(sender)
		logger.SetHighestTerm(term)
	}
}

func transitionToCandidate() {
	if GetMyState() == CANDIDATE {
		// restart the vote requests with a new term.
		VoteChan <- ""
	} else if GetMyState() == LEADER {
		log.Fatal("A leader should not be getting votes.")
	}
	log.Print("I am a candidate now.")
	statusInput <- CANDIDATE
	logger.IncrementNextTerm()
	go captureVotes()
	node.SendVoteRequests()
	VoteChan <- config.UniqueId
}

func transitionToFollower() {
	log.Print("I am a follower now.")
	statusInput <- FOLLOWER
}

func selectLeader() {
	heartbeat := true
	for {
		select {
		case sender := <-HeartbeatChan:
			if config.UniqueId == sender {
				log.Print("Got heartbeat from myself")
			} else {
				log.Print("I am: " + config.UniqueId + ", got heartbeat from " + sender)
				VoteChan <- ""
				node := node.FindNode(sender)
				if node == nil {
					log.Fatal("sender is not part of config: ", sender)
				}
				leaderInput <- node
			}
			heartbeat = true
		case <-time.After(time.Duration(getCandidacyTimeout()) * time.Millisecond):
			if !heartbeat {
				if GetMyState() == LEADER {
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

func Init() {
	statusInput, statusOutput = make(chan int), make(chan int)
	leaderInput, leaderOutput = make(chan *node.Node), make(chan *node.Node)
	go initStatus(statusInput, statusOutput)
	go initLeader(leaderInput, leaderOutput)
	HeartbeatChan = make(chan string)
	transitionToFollower()
	random = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	VoteChan = make(chan string)
	go selectLeader()
}
