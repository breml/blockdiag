package blockdiag

import (
	"strings"
	"testing"
)

func TestShouldParser(t *testing.T) {
	for _, test := range []struct {
		input string
		nodes []string
		edges []string
	}{
		{
			`
blockdiag {
	A;
}
`,
			[]string{"A"},
			[]string{},
		},
		{
			`
blockdiag {
	A -> B;
}
`,
			[]string{"A", "B"},
			[]string{"A|B"},
		},
	} {
		got, err := ParseReader("shouldparse.diag", strings.NewReader(test.input))
		if err != nil {
			t.Fatalf("Parse error: %t with input %s", err, test.input)
		}
		diag = got.(Diag)
		if diag.NodesString() != strings.Join(test.nodes, ", ") {
			t.Fatalf("Nodes error: %s, expected '%s', got: '%s'", test.input, strings.Join(test.nodes, ", "), diag.NodesString())
		}
		if diag.EdgesString() != strings.Join(test.edges, ", ") {
			t.Fatalf("Edges error: %s, expected '%s', got: '%s'", test.input, strings.Join(test.edges, ", "), diag.EdgesString())
		}
	}
}
