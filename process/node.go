package main

import (
	"fmt"
)

type Node struct {
	ip      string
	id      int
	port    int
	uniqeId string
}

var Nodes map[string]*Node
var Leader *Node

func getNode(id int, ip string, port int) *Node {
	var str string
	str = fmt.Sprintf("%d_%d", id, port)
	node := Node{ip, id, port, str}
	return &node
}

func addNode(node *Node) {
	Nodes[node.uniqeId] = node
}

func nodeInit() {
	Nodes = make(map[string]*Node)
}
