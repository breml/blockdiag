# goblockdiag

Source original:
https://bitbucket.org/blockdiag/blockdiag/src/07f96892bfda?at=default

Documentation original:
http://blockdiag.com/en/index.html

# Todo

- [ ] Map Nodes on raster
- [ ] Find start Node
- [ ] Make last ; optional
- [ ] Add attributes
  - [ ] Diagramm
  - [ ] Nodes
  - [ ] Edges
- [ ] Node Groups
- [ ] Strings in Quotes
- [ ] Diagram name
- [ ] Allow "diagram" and "blockdiag" as type
- [ ] Check with gometalinter
- [ ] Check with go cover
- [ ] Use `type grid *[][]*Node`?
- [ ] Tests for text output
- [ ] Support long block names in text output

## Text Paint

- [ ] Circular
- [ ] blockdiag {
			A -> B -> C;
			B -> D -> E -> H;
			A -> F -> E;
			F -> G -> H;
			H -> A;
			X -> Y;
		}
- [ ] blockdiag {
A -> B -> Z;
A -> C -> Z;
A -> D -> Z;
A -> E -> Z;
}
- [ ] // Todo: Go up, until on the right height, right below End