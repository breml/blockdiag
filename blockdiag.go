package blockdiag

import (
	"fmt"
	"sort"
	"strings"

	linq "github.com/ahmetalpbalkan/go-linq"
)

type Diag struct {
	Name       string
	Nodes      map[string]*Node
	Edges      map[string]*Edge
	Attributes map[string]string
	Circular   []*nodes
	Grid       [][]*Node
}

func (diag *Diag) GoString() string {
	var edges []string
	var ret string

	for _, edge := range diag.Edges {
		edges = append(edges, edge.Name)
	}
	sort.Strings(edges)

	ret += fmt.Sprintln("Name:", diag.Name)
	ret += fmt.Sprintln("Nodes:", diag.NodesString())
	ret += fmt.Sprintln("Edges:", diag.EdgesString())
	ret += fmt.Sprintln("Circulars:", diag.CircularString())
	ret += fmt.Sprintln("Attributes:", diag.AttributesString())
	return ret
}

func (diag *Diag) NodesString() string {
	var nodes []string

	for _, node := range diag.Nodes {
		nodes = append(nodes, node.Name)
	}
	sort.Strings(nodes)

	return strings.Join(nodes, ", ")
}

func (diag *Diag) EdgesString() string {
	var edges []string

	for _, edge := range diag.Edges {
		edges = append(edges, edge.Name)
	}
	sort.Strings(edges)

	return strings.Join(edges, ", ")
}

func (diag *Diag) CircularString() string {
	var circulars []string

	for _, circular := range diag.Circular {
		var circularPath []string
		for _, node := range circular.keys {
			circularPath = append(circularPath, node)
		}
		circulars = append(circulars, strings.Join(circularPath, " -> "))
	}
	sort.Strings(circulars)

	return strings.Join(circulars, "\n")
}

func (diag *Diag) AttributesString() string {
	var attributes []string

	for key, value := range diag.Attributes {
		attributes = append(attributes, key+"="+value)
	}
	sort.Strings(attributes)

	return strings.Join(attributes, "\n")
}

func (diag *Diag) FindCircular() bool {
	for _, n := range diag.Nodes {
		visitedNodes := &nodes{}

		if !visitedNodes.exists(n.Name) {
			visitedNodes.keys = append(visitedNodes.keys, n.Name)
		}
		for _, c := range n.getChildNodes() {
			diag.subFindCircular(c, visitedNodes)
		}
	}

	if len(diag.Circular) > 0 {
		return true
	}
	return false
}

func (diag *Diag) subFindCircular(n *Node, visitedNodes *nodes) {
	if visitedNodes.exists(n.Name) {
		circularNodes := &nodes{}
		for _, p := range visitedNodes.keys {
			circularNodes.keys = append(circularNodes.keys, p)
		}
		circularNodes.keys = append(circularNodes.keys, n.Name)

		diag.Circular = append(diag.Circular, circularNodes)
		return
	}
	visitedNodes.keys = append(visitedNodes.keys, n.Name)

	for _, c := range n.getChildNodes() {
		diag.subFindCircular(c, visitedNodes)
	}
	visitedNodes.keys = visitedNodes.keys[:len(visitedNodes.keys)-1]
}

func (diag *Diag) getStartNodes() Nodes {
	var startNodes Nodes

	for _, n := range diag.Nodes {
		startNode := true
		for _, e := range n.Edges {
			if e.End == n {
				startNode = false
				break
			}
		}
		if startNode {
			startNodes = append(startNodes, n)
		}
	}

	return startNodes
}

type Node struct {
	Name  string
	PosX  int
	PosY  int
	Edges []*Edge
}

func (n *Node) getChildNodes() (children []*Node) {
	for _, e := range n.Edges {
		if e.Start == n && e.End != n {
			children = append(children, e.End)
		}
	}
	return
}

type Nodes []*Node

func (nodes Nodes) Len() int {
	return len(nodes)
}

func (nodes Nodes) Less(i, j int) bool {
	return nodes[i].Name < nodes[j].Name
}

func (nodes Nodes) Swap(i, j int) {
	nodes[i], nodes[j] = nodes[j], nodes[i]
}

func (nodes Nodes) String() string {
	var s, delim string
	sort.Sort(nodes)
	for _, n := range nodes {
		s += delim + n.Name
		delim = ", "
	}
	return s
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
