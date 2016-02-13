package blockdiag

import (
	"testing"
)

func TestDiagNil(t *testing.T) {
	var diag *Diag

	diag.AttributesString()
	diag.CircularString()
	diag.EdgesString()
	diag.FindCircular()
	diag.GoString()
	diag.GridString()
	diag.NodesString()
	diag.PlaceInGrid()
	diag.String()
	diag.getEdges()
	diag.getNodes()
	diag.getStartNodes()
	// diag.moveDependingNodesRight(node, placedNodes, count)
}
