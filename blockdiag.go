package blockdiag

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ahmetalpbalkan/go-linq"
)

type Diag struct {
	Name       string
	Nodes      map[string]*Node
	Edges      map[string]*Edge
	Attributes map[string]string
	Circular   []*nodes
	Grid       grid
}

func NewDiag() Diag {
	diag := Diag{}

	diag.Nodes = make(map[string]*Node)
	diag.Edges = make(map[string]*Edge)
	diag.Attributes = make(map[string]string)
	diag.Grid = NewGrid()

	return diag
}

const (
	empty          = ' '
	arrowRight     = '>'
	horizontal     = '\u2500' // ─ http://unicode-table.com/en/2500/
	vertical       = '\u2502' // │ http://unicode-table.com/en/2502/
	horizontalUp   = '\u2534' // ┴ http://unicode-table.com/en/2534/
	horizontalDown = '\u252C' // ┬ http://unicode-table.com/en/252C/
	verticalRight  = '\u251C' // ├ http://unicode-table.com/en/251C/
	verticalLeft   = '\u2524' // ┤ http://unicode-table.com/en/2524/
	upRight        = '\u2514' // └ http://unicode-table.com/en/2514/
	upLeft         = '\u2518' // ┘ http://unicode-table.com/en/2518/
	downRight      = '\u250C' // ┌ http://unicode-table.com/en/250C/
	fourWay        = '\u253C' // ┼ http://unicode-table.com/en/253C/
)

func (diag *Diag) String() string {
	var outGrid [][]rune

	const (
		rowFactor = 2
		colFactor = 7
	)

	// Prepare Output Grid
	outGrid = make([][]rune, len(diag.Grid)*rowFactor)
	for y := 0; y < len(outGrid); y++ {
		outGrid[y] = make([]rune, len(diag.Grid[0])*colFactor)
		for x, _ := range outGrid[y] {
			outGrid[y][x] = ' '
		}
	}

	// Place Nodes
	for y, _ := range diag.Grid {
		for x, n := range diag.Grid[y] {
			if n != nil {
				outGrid[y*rowFactor+1][x*colFactor] = '['
				outGrid[y*rowFactor+1][x*colFactor+1] = rune(n.Name[0])
				outGrid[y*rowFactor+1][x*colFactor+2] = ']'
			}
		}
	}

	// Place Edges
	for _, e := range diag.getEdges() {
		if e.Start.PosY == e.End.PosY && e.Start.PosX < e.End.PosX {
			outGrid[e.Start.PosY*rowFactor+1][e.Start.PosX*colFactor+3] = horizontal
			switch outGrid[e.Start.PosY*rowFactor+1][e.Start.PosX*colFactor+4] {
			case empty:
				outGrid[e.Start.PosY*rowFactor+1][e.Start.PosX*colFactor+4] = horizontal
			case upLeft:
				outGrid[e.Start.PosY*rowFactor+1][e.Start.PosX*colFactor+4] = horizontalUp
			}
			for i := 1; i < (e.End.PosX-e.Start.PosX-1)*colFactor+2; i++ {
				outGrid[e.Start.PosY*rowFactor+1][e.Start.PosX*colFactor+4+i] = horizontal
			}
			outGrid[e.Start.PosY*rowFactor+1][e.End.PosX*colFactor-1] = arrowRight
		}
		if e.Start.PosY < e.End.PosY && e.Start.PosX+1 == e.End.PosX {
			switch outGrid[(e.Start.PosY)*rowFactor+1][e.Start.PosX*colFactor+4] {
			case horizontal:
				outGrid[(e.Start.PosY)*rowFactor+1][e.Start.PosX*colFactor+4] = horizontalDown
			case horizontalUp:
				outGrid[(e.Start.PosY)*rowFactor+1][e.Start.PosX*colFactor+4] = fourWay
			}
			for i := 1; i < (e.End.PosY-e.Start.PosY)*rowFactor+1; i++ {
				switch outGrid[e.Start.PosY+i+1][e.Start.PosX*colFactor+4] {
				case empty:
					outGrid[e.Start.PosY+i+1][e.Start.PosX*colFactor+4] = vertical
				case upRight:
					outGrid[e.Start.PosY+i+1][e.Start.PosX*colFactor+4] = verticalRight
				}
			}
			outGrid[e.End.PosY*rowFactor+1][e.Start.PosX*colFactor+4] = upRight
			outGrid[e.End.PosY*rowFactor+1][e.Start.PosX*colFactor+5] = horizontal
			outGrid[e.End.PosY*rowFactor+1][e.Start.PosX*colFactor+6] = arrowRight
		}
		if e.Start.PosY > e.End.PosY {
			if e.Start.PosX+1 == e.End.PosX {
				// Go directly up and right
				outGrid[e.Start.PosY*rowFactor+1][e.Start.PosX*colFactor+3] = horizontal
				outGrid[e.Start.PosY*rowFactor+1][e.Start.PosX*colFactor+4] = upLeft
				for i := 0; i < (e.Start.PosY-e.End.PosY)*rowFactor-1; i++ {
					switch outGrid[e.Start.PosY*rowFactor-i][e.Start.PosX*colFactor+4] {
					case empty:
						outGrid[e.Start.PosY*rowFactor-i][e.Start.PosX*colFactor+4] = vertical
					case upLeft:
						outGrid[e.Start.PosY*rowFactor-i][e.Start.PosX*colFactor+4] = verticalLeft
					}
				}
				outGrid[e.End.PosY*rowFactor+1][e.End.PosX*colFactor-3] = horizontalDown
			} else {
				// Check if the straight way is free
				straight := true
				for i := 1; i < (e.End.PosX - e.Start.PosX); i++ {
					if diag.Grid[e.Start.PosY][e.Start.PosX+i] != nil {
						straight = false
						break
					}
				}

				if straight {
					for i := 0; i < (e.End.PosX-e.Start.PosX-1)*colFactor+1; i++ {
						outGrid[e.Start.PosY*rowFactor+1][e.Start.PosX*colFactor+3+i] = horizontal
					}
					outGrid[e.Start.PosY*rowFactor+1][e.End.PosX*colFactor-3] = upLeft
					outGrid[e.End.PosY*rowFactor+2][e.End.PosX*colFactor-3] = vertical
					outGrid[e.End.PosY*rowFactor+1][e.End.PosX*colFactor-3] = horizontalDown

				} else {
					// Go up until below End, go right until before End, go up and right into End
					outGrid[e.Start.PosY*rowFactor+1][e.Start.PosX*colFactor+3] = horizontal
					switch outGrid[e.Start.PosY*rowFactor+1][e.Start.PosX*colFactor+4] {
					case empty:
						outGrid[e.Start.PosY*rowFactor+1][e.Start.PosX*colFactor+4] = upLeft
					case upLeft:
						outGrid[e.Start.PosY*rowFactor+1][e.Start.PosX*colFactor+4] = horizontalUp
					}

					// Todo: Go up, until on the right height, right below End

					outGrid[e.Start.PosY*rowFactor][e.Start.PosX*colFactor+4] = downRight

					for i := 1; i < (e.End.PosX-e.Start.PosX-1)*colFactor; i++ {
						outGrid[e.Start.PosY*rowFactor][e.Start.PosX*colFactor+4+i] = horizontal
					}

					outGrid[e.Start.PosY*rowFactor][e.End.PosX*colFactor-3] = upLeft

					// Findout correct junction
					outGrid[e.End.PosY*rowFactor+1][e.End.PosX*colFactor-3] = horizontalDown
				}

			}
		}
	}

	// Prepare Output String
	ret := ""
	for _, rs := range outGrid {
		ret += string(rs) + "\n"
	}

	return ret
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

func (diag *Diag) GridString() string {
	return diag.Grid.String()
}

func (diag *Diag) FindCircular() bool {
	diag.Circular = nil

	for _, n := range diag.getNodes() {
		visitedNodes := &nodes{}

		if !visitedNodes.exists(n.Name) {
			visitedNodes.keys = append(visitedNodes.keys, n.Name)
		}
		for _, c := range n.getChildNodes(false) {
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

		closingEdgeStart := diag.Nodes[circularNodes.keys[len(circularNodes.keys)-2]]
		closingEdgeEnd := diag.Nodes[circularNodes.keys[len(circularNodes.keys)-1]]
		for _, e := range closingEdgeStart.Edges {
			if e.End == closingEdgeEnd {
				e.closeCircle = true
				break
			}
		}

		diag.Circular = append(diag.Circular, circularNodes)
		return
	}
	visitedNodes.keys = append(visitedNodes.keys, n.Name)

	for _, c := range n.getChildNodes(false) {
		diag.subFindCircular(c, visitedNodes)
	}
	visitedNodes.keys = visitedNodes.keys[:len(visitedNodes.keys)-1]
}

func (diag *Diag) getNodes() Nodes {
	var nodes Nodes
	for _, n := range diag.Nodes {
		nodes = append(nodes, n)
	}
	sort.Sort(nodes)
	return nodes
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

	sort.Sort(startNodes)

	return startNodes
}

func (diag *Diag) getEdges() Edges {
	var edges Edges
	for _, e := range diag.Edges {
		edges = append(edges, e)
	}
	sort.Sort(edges)
	return edges
}

type Node struct {
	Name  string
	PosX  int
	PosY  int
	Edges []*Edge
}

func (n *Node) getChildNodes(includeCloseCircle bool) (children Nodes) {
	for _, e := range n.Edges {
		if e.Start == n && e.End != n && (!e.closeCircle || includeCloseCircle) {
			children = append(children, e.End)
		}
	}
	sort.Sort(children)
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
	Start       *Node
	End         *Node
	Name        string
	closeCircle bool
}

type Edges []*Edge

func (edges Edges) Len() int {
	return len(edges)
}

func (edges Edges) Less(i, j int) bool {
	return edges[i].Name < edges[j].Name
}

func (edges Edges) Swap(i, j int) {
	edges[i], edges[j] = edges[j], edges[i]
}

func (edges Edges) String() string {
	var s, delim string
	sort.Sort(edges)
	for _, e := range edges {
		s += delim + e.Name
		delim = ", "
	}
	return s
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

func (diag *Diag) PlaceInGrid() {
	var x, y int

	diag.FindCircular()

	placedNodes := make(map[*Node]bool)

	for _, n := range diag.getStartNodes() {
		_, ok := placedNodes[n]
		if ok {
			continue
		}
		placedNodes[n] = true
		err := diag.Grid.Set(x, y, n, diag)
		if err != nil {
			panic("Set failed")
		}
		y += diag.placeInGrid(n, x+1, y, placedNodes)
		y++
	}

	if len(placedNodes) != len(diag.Nodes) {
		for _, n := range diag.getNodes() {
			_, ok := placedNodes[n]
			if ok {
				continue
			}
			placedNodes[n] = true
			err := diag.Grid.Set(x, y, n, diag)
			if err != nil {
				panic("Set failed")
			}
			y += diag.placeInGrid(n, x+1, y, placedNodes)
			y++
		}
	}
}

func (diag *Diag) placeInGrid(node *Node, x int, y int, placedNodes map[*Node]bool) int {
	addedNodes := 0
	for _, n := range node.getChildNodes(false) {
		_, ok := placedNodes[n]
		if ok {
			if node.PosX >= n.PosX {
				diag.moveDependingNodesRight(n, placedNodes, node.PosX-n.PosX+1)
			}
			continue
		}
		placedNodes[n] = true
		err := diag.Grid.Set(x, y+addedNodes, n, diag)
		if err != nil {
			panic("Set failed")
		}

		addedNodes += diag.placeInGrid(n, x+1, y+addedNodes, placedNodes)
		addedNodes++
	}

	if addedNodes > 0 {
		return addedNodes - 1
	}
	return addedNodes
}

func (diag *Diag) moveDependingNodesRight(node *Node, placedNodes map[*Node]bool, count int) {
	for _, n := range node.getChildNodes(false) {
		diag.moveDependingNodesRight(n, placedNodes, count)
	}
	oldX := node.PosX
	err := diag.Grid.Set(oldX+count, node.PosY, node, diag)
	if err != nil {
		panic("Set failed")
	}
	err = diag.Grid.Set(oldX, node.PosY, nil, diag)
	if err != nil {
		panic("Set failed")
	}
}
