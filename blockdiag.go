package blockdiag

import (
	"fmt"

	linq "github.com/ahmetalpbalkan/go-linq"
)

type Diag struct {
	Name  string
	Nodes map[string]*Node
	Edges map[string]*Edge
}

func (diag *Diag) FindCircular() bool {
	for _, n := range diag.Nodes {
		visitedNodes := &nodes{}

		fmt.Println("\nStart from Node:", n.Name)
		if !visitedNodes.exists(n.Name) {
			visitedNodes.keys = append(visitedNodes.keys, n.Name)
		}
		for _, c := range n.getChildNodes() {
			fmt.Println("Child:", c.Name)
			subFindCircular(c, visitedNodes)
		}
	}

	// Wrong
	return false
}

func subFindCircular(n *Node, visitedNodes *nodes) {
	fmt.Println("In Node:", n.Name)
	if visitedNodes.exists(n.Name) {
		fmt.Println("Found already visited Node:", n.Name)
		fmt.Print("Path: ")
		for _, p := range visitedNodes.keys {
			fmt.Print(p, " -> ")
		}
		fmt.Println(n.Name)
		return
	}
	visitedNodes.keys = append(visitedNodes.keys, n.Name)

	for _, c := range n.getChildNodes() {
		fmt.Println("Child:", c.Name)
		subFindCircular(c, visitedNodes)
	}
	visitedNodes.keys = visitedNodes.keys[:len(visitedNodes.keys)-1]
}

type Node struct {
	Name  string
	PosX  int
	PosY  int
	Edges []*Edge
}

func (n *Node) getChildNodes() (children []*Node) {
	for _, e := range n.Edges {
		if e.Start == n {
			children = append(children, e.End)
		}
	}
	return
}

type Edge struct {
	Start *Node
	End   *Node
	Name  string
}

type nodes struct {
	keys []string
}

func (n *nodes) exists(key string) bool {
	ret, _ := linq.From(n.keys).AnyWith(func(s linq.T) (bool, error) {
		return s.(string) == key, nil
	})
	return ret
}
