# goblockdiag

Source original:
https://bitbucket.org/blockdiag/blockdiag/src/07f96892bfda?at=default

Documentation original:
http://blockdiag.com/en/index.html

# Todo

- [X] Map Nodes on raster
- [X] Find start Node
- [ ] Make last ; optional
- [ ] Add attributes
  - [X] Diagramm
  - [ ] Nodes
  - [ ] Edges
- [ ] Node Groups
- [ ] Strings in Quotes
- [ ] Diagram name
- [X] Allow "diagram" and "blockdiag" as type
- [ ] Check with gometalinter
- [ ] Check with go cover
- [ ] Use `type grid *[][]*Node`?
- [ ] Tests for text output
- [ ] Support long block names in text output
- [ ] Code cleanup
- [ ] Refactoring, if placeing part of edge, save what directions are already covered and add only new endpoints, decide the char in a later step
- [ ] Self reference A -> A
- [ ] Circular A -> B -> A
- [ ] Refactoring API, only make neccessary functions, methods, types, etc. public
- [ ] Split blockdiag into multiple files
- [ ] Split blockdiag_test into multiple files
- [ ] Add tests for getChildNodes and getParentNodes
- [ ] Test implementation of PlaceInGrid with Breadth First
- [X] Self reference

```
A -> A;

 ▼──┐
[A]─┘
```

- [ ] Circular

```	
	blockdiag {
			A -> B -> C;
			B -> D -> E -> H;
			A -> F -> E;
			F -> G -> H;
			H -> A;
		}
```

- [X] Go up, until on the right height, right below End

- [ ] Fix paint error:

```
`blockdiag {
	A -> B -> D;
	A -> X -> B;
	X -> D;
	X -> R;
	D -> G;
	D -> H;
	X -> H;
}`
```
