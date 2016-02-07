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
	upRight        = '\u2514' // └ http://unicode-table.com/en/2514/
	upLeft         = '\u2518' // ┘ http://unicode-table.com/en/2518/
	downRight      = '\u250C' // ┌ http://unicode-table.com/en/250C/
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
				// fmt.Println(y, x)
				outGrid[y*rowFactor+1][x*colFactor] = '['
				outGrid[y*rowFactor+1][x*colFactor+1] = rune(n.Name[0])
				outGrid[y*rowFactor+1][x*colFactor+2] = ']'
				// ret += "[" + string(n.Name[0]) + "]"
				// if x < len(diag.Grid[y])-1 && diag.Grid[y][x+1] != nil {
				// 	if len(n.getChildNodes()) > 1 {
				// 		ret += "|->"
				// 	} else {
				// 		ret += "-->"
				// 	}
				// }
				// } else {
				// 	ret += "      "
			}
		}
		//ret += "\n"
	}

	// Place Edges
	// Unicode Arrows
	// https://en.wikipedia.org/wiki/Arrow_(symbol)#Arrows_in_Unicode
	// https://en.wikipedia.org/wiki/Supplemental_Arrows-B
	// https://en.wikipedia.org/wiki/Tee_(symbol)
	// https://en.wikipedia.org/wiki/Up_tack
	// http://unicode-table.com/de/007C/ Senkrechter Strich
	// http://unicode-table.com/de/23AF/ Waagrechter Strich
	// http://unicode-table.com/de/blocks/miscellaneous-technical/
	// http://www.asciitable.com/
	// https://de.wikipedia.org/wiki/Unicodeblock_Rahmenzeichnung
	for _, e := range diag.getEdges() {
		fmt.Println(e.Start.Name, e.Start.PosX, e.Start.PosY, "|", e.End.Name, e.End.PosX, e.End.PosY)
		if e.Start.PosY == e.End.PosY && e.Start.PosX+1 == e.End.PosX {
			outGrid[e.Start.PosY*rowFactor+1][e.Start.PosX*colFactor+3] = horizontal
			if outGrid[e.Start.PosY*rowFactor+1][e.Start.PosX*colFactor+4] == empty {
				outGrid[e.Start.PosY*rowFactor+1][e.Start.PosX*colFactor+4] = horizontal
			} else {
				if outGrid[e.Start.PosY*rowFactor+1][e.Start.PosX*colFactor+4] == upLeft {
					outGrid[e.Start.PosY*rowFactor+1][e.Start.PosX*colFactor+4] = horizontalUp
				}
			}
			outGrid[e.Start.PosY*rowFactor+1][e.Start.PosX*colFactor+5] = horizontal
			outGrid[e.Start.PosY*rowFactor+1][e.Start.PosX*colFactor+6] = arrowRight
		}
		if e.Start.PosY < e.End.PosY && e.Start.PosX+1 == e.End.PosX {
			// if outGrid[(e.Start.PosY)*rowFactor+1][e.Start.PosX*colFactor+4] == horizontal {
			outGrid[(e.Start.PosY)*rowFactor+1][e.Start.PosX*colFactor+4] = horizontalDown
			// }
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
				outGrid[e.Start.PosY*rowFactor][e.Start.PosX*colFactor+4] = vertical
				outGrid[e.End.PosY*rowFactor+1][e.End.PosX*colFactor-3] = horizontalDown

				// // if outGrid[(e.Start.PosY)*rowFactor+1][e.Start.PosX*colFactor+4] == horizontal {
				// outGrid[(e.Start.PosY)*rowFactor+1][e.Start.PosX*colFactor+4] = horizontalDown
				// // }
				// for i := 1; i < (e.End.PosY-e.Start.PosY)*rowFactor+1; i++ {
				// 	outGrid[e.Start.PosY+i+1][e.Start.PosX*colFactor+4] = vertical
				// }
				// outGrid[e.End.PosY*rowFactor+1][e.Start.PosX*colFactor+4] = upRight
				// outGrid[e.End.PosY*rowFactor+1][e.Start.PosX*colFactor+5] = horizontal
				// outGrid[e.End.PosY*rowFactor+1][e.Start.PosX*colFactor+6] = arrowRight
			} else {
				// Go up until below End, go right until before End, go up and right into End
				outGrid[e.Start.PosY*rowFactor+1][e.Start.PosX*colFactor+3] = horizontal
				if outGrid[e.Start.PosY*rowFactor+1][e.Start.PosX*colFactor+4] == empty {
					outGrid[e.Start.PosY*rowFactor+1][e.Start.PosX*colFactor+4] = upLeft
				} else {
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
		// nodes = append(nodes, fmt.Sprintf("%s (%d, %d)", node.Name, node.PosX, node.PosY))
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

func (n *Node) getChildNodes() (children Nodes) {
	for _, e := range n.Edges {
		if e.Start == n && e.End != n {
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
	Start *Node
	End   *Node
	Name  string
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

func (diag *Diag) placeInGrid(n *Node, x int, y int, placedNodes map[*Node]bool) int {
	addedNodes := 0
	for _, n := range n.getChildNodes() {
		_, ok := placedNodes[n]
		if ok {
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
