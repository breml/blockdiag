package blockdiag

import (
	"strings"
	"testing"
)

func TestShouldParser(t *testing.T) {
	for _, test := range []struct {
		description string
		input       string
		nodes       []string
		edges       []string
	}{
		{
			"Empty diagram",
			`
blockdiag {}
`,
			[]string{},
			[]string{},
		},
		{
			"Single Node",
			`
blockdiag {
	A;
}
`,
			[]string{"A"},
			[]string{},
		},
		{
			// TODO Add test case for node chain without tailing ;
			"Node chain",
			`
blockdiag {
	A -> B;
}
`,
			[]string{"A", "B"},
			[]string{"A|B"},
		},
		{
			"Multiple chains, using same nodes",
			`
blockdiag {
	A -> B -> C;
	A -> D;
}
`,
			[]string{"A", "B", "C", "D"},
			[]string{"A|B", "A|D", "B|C"},
		},
		{
			"Self reference",
			`
blockdiag {
	A -> A;
}
`,
			[]string{"A"},
			[]string{"A|A"},
		},
	} {
		got, err := ParseReader("shouldparse.diag", strings.NewReader(test.input))
		if err != nil {
			t.Fatalf("%s: parse error: %t with input %s", test.description, err, test.input)
		}
		gotDiag := got.(Diag)
		if gotDiag.NodesString() != strings.Join(test.nodes, ", ") {
			t.Fatalf("%s: nodes error: %s, expected '%s', got: '%s'", test.description, test.input, strings.Join(test.nodes, ", "), gotDiag.NodesString())
		}
		if gotDiag.EdgesString() != strings.Join(test.edges, ", ") {
			t.Fatalf("%s edges error: %s, expected '%s', got: '%s'", test.description, test.input, strings.Join(test.edges, ", "), gotDiag.EdgesString())
		}
	}
}
