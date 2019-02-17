package node

import (
	"encoding/json"
	"fmt"
)

type Data struct {
	X int
	Y int
}

type Node struct {
	Data Data
	Next *Node
}

type LinkNode struct {
	Node *Node
	Tail *Node
}

func CreatNode() *LinkNode {
	node := &Node{
		Next: nil,
	}
	return &LinkNode{
		Node: node,
		Tail: node,
	}
}

func (ln *LinkNode) HeadAddNode(node *Node) *LinkNode {
	if ln.Node.Next == nil && ln.Node == ln.Tail {
		ln.Node.Next = node
		ln.Tail = node
		return ln
	}
	node.Next = ln.Node.Next
	ln.Node.Next = node
	return ln
}

func (ln *LinkNode) TailDeleteNode() *LinkNode {
	if ln.Node.Next == nil {
		return ln
	}
	node := ln.Node.Next
	for {
		if node == ln.Tail {
			node.Next = nil
			ln.Tail = node
			break
		}
		if node.Next == ln.Tail {
			node.Next = nil
			ln.Tail = node
			break
		}
		node = node.Next
	}
	return ln
}

func (ln *LinkNode) PrintLinkNode() {
	bytes, _ := json.Marshal(ln)
	fmt.Println(string(bytes))
}
