package blockdiag

import (
	"testing"
)

func TestNodeNil(t *testing.T) {
	var node *Node

	node.getChildNodes(false)
	node.getParentNodes(false)
}
