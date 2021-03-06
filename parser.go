// blockdiag parser

package blockdiag

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

var g = &grammar{
	rules: []*rule{
		{
			name: "diag",
			pos:  position{line: 7, col: 1, offset: 44},
			expr: &actionExpr{
				pos: position{line: 7, col: 8, offset: 51},
				run: (*parser).callondiag1,
				expr: &seqExpr{
					pos: position{line: 7, col: 8, offset: 51},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 7, col: 8, offset: 51},
							name: "_",
						},
						&actionExpr{
							pos: position{line: 7, col: 11, offset: 54},
							run: (*parser).callondiag4,
							expr: &labeledExpr{
								pos:   position{line: 7, col: 11, offset: 54},
								label: "diagtype",
								expr: &choiceExpr{
									pos: position{line: 7, col: 21, offset: 64},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 7, col: 21, offset: 64},
											val:        "blockdiag",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 7, col: 35, offset: 78},
											val:        "diagram",
											ignoreCase: false,
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 10, col: 5, offset: 142},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 10, col: 7, offset: 144},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 10, col: 11, offset: 148},
							name: "_",
						},
						&zeroOrMoreExpr{
							pos: position{line: 10, col: 13, offset: 150},
							expr: &ruleRefExpr{
								pos:  position{line: 10, col: 13, offset: 150},
								name: "diagElements",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 10, col: 27, offset: 164},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 10, col: 29, offset: 166},
							val:        "}",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 10, col: 33, offset: 170},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 10, col: 35, offset: 172},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "diagElements",
			pos:  position{line: 16, col: 1, offset: 248},
			expr: &choiceExpr{
				pos: position{line: 16, col: 16, offset: 263},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 16, col: 16, offset: 263},
						name: "diagAttr",
					},
					&ruleRefExpr{
						pos:  position{line: 16, col: 27, offset: 274},
						name: "chain",
					},
					&ruleRefExpr{
						pos:  position{line: 16, col: 35, offset: 282},
						name: "attributedNode",
					},
				},
			},
		},
		{
			name: "diagAttr",
			pos:  position{line: 18, col: 1, offset: 298},
			expr: &actionExpr{
				pos: position{line: 18, col: 12, offset: 309},
				run: (*parser).callondiagAttr1,
				expr: &seqExpr{
					pos: position{line: 18, col: 12, offset: 309},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 18, col: 12, offset: 309},
							label: "attrName",
							expr: &ruleRefExpr{
								pos:  position{line: 18, col: 21, offset: 318},
								name: "ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 18, col: 27, offset: 324},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 18, col: 29, offset: 326},
							val:        "=",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 18, col: 33, offset: 330},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 18, col: 35, offset: 332},
							label: "attrValue",
							expr: &ruleRefExpr{
								pos:  position{line: 18, col: 45, offset: 342},
								name: "attrValue",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 18, col: 55, offset: 352},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 18, col: 57, offset: 354},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 18, col: 61, offset: 358},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "attrValue",
			pos:  position{line: 23, col: 1, offset: 446},
			expr: &actionExpr{
				pos: position{line: 23, col: 13, offset: 458},
				run: (*parser).callonattrValue1,
				expr: &oneOrMoreExpr{
					pos: position{line: 23, col: 14, offset: 459},
					expr: &charClassMatcher{
						pos:        position{line: 23, col: 14, offset: 459},
						val:        "[a-zA-Z0-9]",
						ranges:     []rune{'a', 'z', 'A', 'Z', '0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "chain",
			pos:  position{line: 27, col: 1, offset: 506},
			expr: &actionExpr{
				pos: position{line: 27, col: 9, offset: 514},
				run: (*parser).callonchain1,
				expr: &seqExpr{
					pos: position{line: 27, col: 9, offset: 514},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 27, col: 9, offset: 514},
							label: "node",
							expr: &ruleRefExpr{
								pos:  position{line: 27, col: 14, offset: 519},
								name: "node",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 27, col: 19, offset: 524},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 27, col: 21, offset: 526},
							label: "nodes",
							expr: &zeroOrMoreExpr{
								pos: position{line: 27, col: 27, offset: 532},
								expr: &actionExpr{
									pos: position{line: 27, col: 28, offset: 533},
									run: (*parser).callonchain8,
									expr: &seqExpr{
										pos: position{line: 27, col: 28, offset: 533},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 27, col: 28, offset: 533},
												name: "edge",
											},
											&ruleRefExpr{
												pos:  position{line: 27, col: 33, offset: 538},
												name: "_",
											},
											&labeledExpr{
												pos:   position{line: 27, col: 35, offset: 540},
												label: "n",
												expr: &ruleRefExpr{
													pos:  position{line: 27, col: 37, offset: 542},
													name: "node",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 27, col: 42, offset: 547},
												name: "_",
											},
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 29, col: 6, offset: 571},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 29, col: 10, offset: 575},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "node",
			pos:  position{line: 45, col: 1, offset: 889},
			expr: &actionExpr{
				pos: position{line: 45, col: 8, offset: 896},
				run: (*parser).callonnode1,
				expr: &labeledExpr{
					pos:   position{line: 45, col: 8, offset: 896},
					label: "node",
					expr: &ruleRefExpr{
						pos:  position{line: 45, col: 13, offset: 901},
						name: "ident",
					},
				},
			},
		},
		{
			name: "attributedNode",
			pos:  position{line: 57, col: 1, offset: 1095},
			expr: &actionExpr{
				pos: position{line: 57, col: 18, offset: 1112},
				run: (*parser).callonattributedNode1,
				expr: &seqExpr{
					pos: position{line: 57, col: 18, offset: 1112},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 57, col: 18, offset: 1112},
							label: "node",
							expr: &ruleRefExpr{
								pos:  position{line: 57, col: 23, offset: 1117},
								name: "node",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 57, col: 28, offset: 1122},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 57, col: 30, offset: 1124},
							val:        "[",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 57, col: 34, offset: 1128},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 57, col: 36, offset: 1130},
							label: "iAttrs",
							expr: &zeroOrMoreExpr{
								pos: position{line: 57, col: 43, offset: 1137},
								expr: &ruleRefExpr{
									pos:  position{line: 57, col: 43, offset: 1137},
									name: "nodeAttr",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 57, col: 53, offset: 1147},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 57, col: 55, offset: 1149},
							val:        "]",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 57, col: 59, offset: 1153},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 57, col: 61, offset: 1155},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 57, col: 65, offset: 1159},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "nodeAttr",
			pos:  position{line: 69, col: 1, offset: 1403},
			expr: &actionExpr{
				pos: position{line: 69, col: 12, offset: 1414},
				run: (*parser).callonnodeAttr1,
				expr: &seqExpr{
					pos: position{line: 69, col: 12, offset: 1414},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 69, col: 12, offset: 1414},
							label: "attrName",
							expr: &ruleRefExpr{
								pos:  position{line: 69, col: 21, offset: 1423},
								name: "ident",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 69, col: 27, offset: 1429},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 69, col: 29, offset: 1431},
							val:        "=",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 69, col: 33, offset: 1435},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 69, col: 35, offset: 1437},
							label: "attrValue",
							expr: &ruleRefExpr{
								pos:  position{line: 69, col: 45, offset: 1447},
								name: "attrValue",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 69, col: 55, offset: 1457},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 69, col: 57, offset: 1459},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 69, col: 61, offset: 1463},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "edge",
			pos:  position{line: 74, col: 1, offset: 1558},
			expr: &labeledExpr{
				pos:   position{line: 74, col: 8, offset: 1565},
				label: "edge",
				expr: &litMatcher{
					pos:        position{line: 74, col: 14, offset: 1571},
					val:        "->",
					ignoreCase: false,
				},
			},
		},
		{
			name: "ident",
			pos:  position{line: 76, col: 1, offset: 1579},
			expr: &actionExpr{
				pos: position{line: 76, col: 9, offset: 1587},
				run: (*parser).callonident1,
				expr: &oneOrMoreExpr{
					pos: position{line: 76, col: 9, offset: 1587},
					expr: &charClassMatcher{
						pos:        position{line: 76, col: 9, offset: 1587},
						val:        "[a-zA-Z0-9_]",
						chars:      []rune{'_'},
						ranges:     []rune{'a', 'z', 'A', 'Z', '0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 80, col: 1, offset: 1634},
			expr: &zeroOrMoreExpr{
				pos: position{line: 80, col: 5, offset: 1638},
				expr: &choiceExpr{
					pos: position{line: 80, col: 6, offset: 1639},
					alternatives: []interface{}{
						&charClassMatcher{
							pos:        position{line: 80, col: 6, offset: 1639},
							val:        "[ \\t\\n\\r]",
							chars:      []rune{' ', '\t', '\n', '\r'},
							ignoreCase: false,
							inverted:   false,
						},
						&ruleRefExpr{
							pos:  position{line: 80, col: 18, offset: 1651},
							name: "Comment",
						},
					},
				},
			},
		},
		{
			name: "Comment",
			pos:  position{line: 82, col: 1, offset: 1662},
			expr: &seqExpr{
				pos: position{line: 82, col: 11, offset: 1672},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 82, col: 11, offset: 1672},
						name: "CommentPrefix",
					},
					&zeroOrMoreExpr{
						pos: position{line: 82, col: 25, offset: 1686},
						expr: &charClassMatcher{
							pos:        position{line: 82, col: 25, offset: 1686},
							val:        "[^\\n\\r]",
							chars:      []rune{'\n', '\r'},
							ignoreCase: false,
							inverted:   true,
						},
					},
					&zeroOrOneExpr{
						pos: position{line: 82, col: 34, offset: 1695},
						expr: &charClassMatcher{
							pos:        position{line: 82, col: 34, offset: 1695},
							val:        "[\\r]",
							chars:      []rune{'\r'},
							ignoreCase: false,
							inverted:   false,
						},
					},
					&charClassMatcher{
						pos:        position{line: 82, col: 39, offset: 1700},
						val:        "[\\n]",
						chars:      []rune{'\n'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "CommentPrefix",
			pos:  position{line: 84, col: 1, offset: 1706},
			expr: &choiceExpr{
				pos: position{line: 84, col: 17, offset: 1722},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 84, col: 17, offset: 1722},
						val:        "#",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 84, col: 23, offset: 1728},
						val:        "//",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 86, col: 1, offset: 1734},
			expr: &notExpr{
				pos: position{line: 86, col: 7, offset: 1740},
				expr: &anyMatcher{
					line: 86, col: 8, offset: 1741,
				},
			},
		},
	},
}

func (c *current) ondiag4(diagtype interface{}) (interface{}, error) {
	_parserDiag = NewDiag()
	return diagtype, nil
}

func (p *parser) callondiag4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.ondiag4(stack["diagtype"])
}

func (c *current) ondiag1(diagtype interface{}) (interface{}, error) {
	updateEdges()
	_parserDiag.Name = "test"
	return _parserDiag, nil
}

func (p *parser) callondiag1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.ondiag1(stack["diagtype"])
}

func (c *current) ondiagAttr1(attrName, attrValue interface{}) (interface{}, error) {
	_parserDiag.Attributes[attrName.(string)] = attrValue.(string)
	return nil, nil
}

func (p *parser) callondiagAttr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.ondiagAttr1(stack["attrName"], stack["attrValue"])
}

func (c *current) onattrValue1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonattrValue1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onattrValue1()
}

func (c *current) onchain8(n interface{}) (interface{}, error) {
	return n, nil
}

func (p *parser) callonchain8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onchain8(stack["n"])
}

func (c *current) onchain1(node, nodes interface{}) (interface{}, error) {
	nodeA := node.(*Node)

	for _, n := range toIfaceSlice(nodes) {
		nodeB := n.(*Node)
		edge := nodeA.Name + "|" + nodeB.Name

		if e, present := _parserDiag.Edges[edge]; !present {
			e = &Edge{Name: edge, Start: nodeA, End: nodeB}
			_parserDiag.Edges[edge] = e
		}
		nodeA = nodeB
	}
	return nil, nil
}

func (p *parser) callonchain1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onchain1(stack["node"], stack["nodes"])
}

func (c *current) onnode1(node interface{}) (interface{}, error) {
	var n *Node
	var present bool

	name := string(c.text)
	if n, present = _parserDiag.Nodes[name]; !present {
		n = &Node{Name: name}
		_parserDiag.Nodes[name] = n
	}
	return n, nil
}

func (p *parser) callonnode1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onnode1(stack["node"])
}

func (c *current) onattributedNode1(node, iAttrs interface{}) (interface{}, error) {
	nodeA := node.(*Node)
	if nodeA.Attributes == nil {
		nodeA.Attributes = make(map[string]string)
	}
	for _, iAttr := range toIfaceSlice(iAttrs) {
		attr := iAttr.(attribute)
		nodeA.Attributes[attr.key] = attr.value
	}
	return nil, nil
}

func (p *parser) callonattributedNode1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onattributedNode1(stack["node"], stack["iAttrs"])
}

func (c *current) onnodeAttr1(attrName, attrValue interface{}) (interface{}, error) {
	attr := attribute{key: attrName.(string), value: attrValue.(string)}
	return attr, nil
}

func (p *parser) callonnodeAttr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onnodeAttr1(stack["attrName"], stack["attrValue"])
}

func (c *current) onident1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonident1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onident1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errNoMatch is returned if no match could be found.
	errNoMatch = errors.New("no match found")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner  error
	pos    position
	prefix string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	recover bool
	debug   bool
	depth   int

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position)
}

func (p *parser) addErrAt(err error, pos position) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, prefix: buf.String()}
	p.errs.add(pe)
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n > 0 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// make sure this doesn't go out silently
			p.addErr(errNoMatch)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint
	var ok bool

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position)
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	// can't match EOF
	if cur == utf8.RuneError {
		return nil, false
	}
	start := p.pt
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(not.expr)
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	var vals []interface{}

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}
