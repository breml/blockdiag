package blockdiag

type Diag struct {
	Name  string
	Nodes map[string]*Node
	Edges map[string]*Edge
}

type Node struct {
	Name  string
	PosX  int
	PosY  int
	Edges []*Edge
}

type Edge struct {
	Start *Node
	End   *Node
	Name  string
}
