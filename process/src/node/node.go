package node

import (
	"fmt"

	"config"
)

type Node struct {
	Ip      string
	Id      int
	Port    int
	UniqeId string
}

type NodeVisitor func(node *Node)

var Nodes map[string]*Node

func ForAll(visitor NodeVisitor) {
	for _, node := range Nodes {
		visitor(node)
	}
}

func getNode(Id int, Ip string, Port int) *Node {
	var str string
	str = fmt.Sprintf("%d_%d", Id, Port)
	node := Node{Ip, Id, Port, str}
	return &node
}

func FindNode(uniqueId string) *Node {
	if node, ok := Nodes[uniqueId]; ok {
		return node
	}
	return nil
}

func addNode(node *Node) {
	Nodes[node.UniqeId] = node
}

func Init() {
	Nodes = make(map[string]*Node)
	handler := func(Id int, Ip string, Port int) {
		addNode(getNode(Id, Ip, Port))
	}
	config.RegisterNodeHandler(handler)
}
