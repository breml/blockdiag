package blockdiag

import (
	"fmt"
)

type grid [][]*Node

func NewGrid() grid {
	const minSize = 1

	return NewSizedGrid(minSize, minSize)
}

func NewSizedGrid(x, y int) grid {
	var g grid

	g = make([][]*Node, y)
	for i := 0; i < y; i++ {
		g[i] = make([]*Node, x)
	}

	return g
}

func (g grid) Set(x, y int, n *Node, diag *Diag) error {
	if x < 0 || y < 0 {
		return fmt.Errorf("out of bound x or y, %d, %d", x, y)
	}

	if x >= len(g[0]) {
		for i := len(g[0]); i <= x; i++ {
			g = *(g.appendCol())
			diag.Grid = g
		}
	}

	if y >= len(g) {
		for i := len(g); i <= y; i++ {
			g = *(g.appendRow())
			diag.Grid = g
		}
	}

	g[y][x] = n
	if n != nil {
		n.PosX = x
		n.PosY = y
	}

	return nil
}

func (g grid) String() string {
	ret := ""

	for y, _ := range g {
		for _, n := range g[y] {
			if n != nil {
				ret += "[" + string(n.Name[0]) + "] "
			} else {
				ret += "    "
			}
		}
		ret += "\n"
	}

	return ret
}

func (g *grid) appendRow() *grid {
	gVal := *g
	gVal = append(gVal, make([]*Node, len(gVal[0])))
	return &gVal
}

func (g *grid) appendCol() *grid {
	gVal := *g
	for i := 0; i < len(gVal); i++ {
		gVal[i] = append(gVal[i], nil)
	}
	return &gVal
}
