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
		{
			"Comment",
			`
# Comment
blockdiag # Comment
{
# Comment
	A; # Comment
# Comment
} # Comment
`,
			[]string{"A"},
			[]string{},
		},
		{
			"Multi Char Node Names",
			`
blockdiag
{
	MultiCharNodeName1;
}
`,
			[]string{"MultiCharNodeName1"},
			[]string{},
		},
	} {
		got, err := ParseReader("shouldparse.diag", strings.NewReader(test.input))
		if err != nil {
			t.Fatalf("%s: parse error: %t with input %s", test.description, err, test.input)
		}
		gotDiag, ok := got.(Diag)
		if !ok {
			t.Fatalf("%s: assertion error: %s should parse to diag", test.description, test.input)
		}
		if gotDiag.NodesString() != strings.Join(test.nodes, ", ") {
			t.Fatalf("%s: nodes error: %s, expected '%s', got: '%s'", test.description, test.input, strings.Join(test.nodes, ", "), gotDiag.NodesString())
		}
		if gotDiag.EdgesString() != strings.Join(test.edges, ", ") {
			t.Fatalf("%s edges error: %s, expected '%s', got: '%s'", test.description, test.input, strings.Join(test.edges, ", "), gotDiag.EdgesString())
		}
	}
}

func TestShouldNotParse(t *testing.T) {
	for _, test := range []struct {
		description string
		input       string
	}{
		{
			"No block",
			`
blockdiag
`,
		},
	} {
		_, err := ParseReader("shouldnotparse.diag", strings.NewReader(test.input))
		if err == nil {
			t.Fatalf("%s: should not parse, but didn't give an error with input %s", test.description, test.input)
		}
	}
}

func TestCircular(t *testing.T) {
	for _, test := range []struct {
		input    string
		circular bool
	}{
		{
			`
blockdiag{
	A;
}
`,
			false,
		},
		{
			`
blockdiag{
	A -> B -> C;
}
`,
			false,
		},
		{
			`
blockdiag{
	A -> A;
}
`,
			false,
		},
		{
			`
blockdiag{
	A -> B -> C -> A;
}
`,
			true,
		},
	} {
		got, err := ParseReader("shouldnotparse.diag", strings.NewReader(test.input))
		if err != nil {
			t.Fatalf("should not parse, but didn't give an error with input %s", test.input)
		}
		gotDiag, ok := got.(Diag)
		if !ok {
			t.Fatalf("assertion error: %s should parse to diag", test.input)
		}
		if gotDiag.FindCircular() != test.circular {
			t.Fatalf("expect %s to be circular == %t", test.input, test.circular)
		}
	}
}
