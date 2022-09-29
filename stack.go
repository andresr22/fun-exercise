package main

import "fmt"

type Node struct {
	value map[string]interface{}
}

func (n *Node) String() string {
	return fmt.Sprint(n.value)
}

type Stack struct {
	nodes []*Node
	size  int
}

func NewStack() *Stack {
	return &Stack{}
}

func (s *Stack) Push(n *Node) {
	s.nodes = append(s.nodes[:s.size], n)
	s.size++
}

func (s *Stack) Pop() *Node {
	if s.size == 0 {
		return nil
	}
	s.size--
	return s.nodes[s.size]
}

func (s *Stack) Print() {
	for _, n := range s.nodes {
		fmt.Print(n.value, " ")
	}
}
