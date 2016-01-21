package blockdiag

import (
	"sort"
	"strings"
	"testing"
)

func TestShouldParser(t *testing.T) {
	for _, test := range []struct {
		description string
		input       string
		nodes       []string
		edges       []string
		attributes  map[string]string
	}{
		{
			"Empty diagram",
			`
blockdiag {}
`,
			[]string{},
			[]string{},
			map[string]string{},
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
			map[string]string{},
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
			map[string]string{},
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
			map[string]string{},
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
			map[string]string{},
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
			map[string]string{},
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
			map[string]string{},
		},
		{
			"Digramm Attributes",
			`
blockdiag
{
	node_width = 128;
	A;
}
`,
			[]string{"A"},
			[]string{},
			map[string]string{
				"node_width": "128",
			},
		},
		{
			"Digramm type 'diagram'",
			`
diagram
{
	A;
}
`,
			[]string{"A"},
			[]string{},
			map[string]string{},
		},
	} {
		got, err := ParseReader("shouldparse.diag", strings.NewReader(test.input))
		if err != nil {
			t.Fatalf("%s: parse error: %s with input %s", test.description, err, test.input)
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

		var attributes []string
		for key, value := range test.attributes {
			attributes = append(attributes, key+"="+value)
		}
		sort.Strings(attributes)
		if gotDiag.AttributesString() != strings.Join(attributes, "\n") {
			t.Fatalf("%s attributes error: %s, expected '%s', got: '%s'", test.description, test.input, strings.Join(attributes, "\n"), gotDiag.AttributesString())
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

func TestGetStartNodes(t *testing.T) {
	for _, test := range []struct {
		input      string
		startNodes []string
	}{
		{
			`
blockdiag{
	A -> B -> C;
}
`,
			[]string{"A"},
		},
		{
			`
blockdiag {
	A -> B -> C;
	D;
	E -> F;
}
`,
			[]string{"A", "D", "E"},
		},
		{
			`
blockdiag {
	D;
	E -> F;
	A -> B -> C;
}
`,
			[]string{"A", "D", "E"},
		},
	} {
		got, err := ParseReader("placeingrid.diag", strings.NewReader(test.input))
		if err != nil {
			t.Fatalf("should not parse, but didn't give an error with input %s", test.input)
		}
		gotDiag, ok := got.(Diag)
		if !ok {
			t.Fatalf("assertion error: %s should parse to diag", test.input)
		}
		// if gotDiag.PlaceInGrid() != test.circular {
		// 	t.Fatalf("expect %s to be circular == %t", test.input, test.circular)
		// }
		startNodes := gotDiag.getStartNodes()
		if len(startNodes) != len(test.startNodes) {
			t.Fatalf("Start Nodes count wrong, expected: %s, got: %s", strings.Join(test.startNodes, ", "), startNodes)
		}
		sort.Strings(test.startNodes)
		for i, n := range startNodes {
			if n.Name != test.startNodes[i] {
				t.Fatalf("Start Nodes do not match, expected: %s, got: %s", strings.Join(test.startNodes, ", "), startNodes)
			}
		}
	}
}

func TestPlaceInGrid(t *testing.T) {
	for _, test := range []struct {
		input string
	}{
		{
			`
blockdiag{
	A -> B -> C;
}
`,
		},
		{
			`
blockdiag{
	A -> B -> C;
	B -> D;
	A -> E -> C;
}
`,
		},
	} {
		got, err := ParseReader("placeingrid.diag", strings.NewReader(test.input))
		if err != nil {
			t.Fatalf("should not parse, but didn't give an error with input %s", test.input)
		}
		gotDiag, ok := got.(Diag)
		if !ok {
			t.Fatalf("assertion error: %s should parse to diag", test.input)
		}
		// if gotDiag.PlaceInGrid() != test.circular {
		// 	t.Fatalf("expect %s to be circular == %t", test.input, test.circular)
		// }
		gotDiag.PlaceInGrid()
		t.Logf("%s\n", gotDiag.GridString())
	}
}
