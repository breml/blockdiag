package blockdiag

var _parserDiag Diag

type attribute struct {
	key   string
	value string
}

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

func updateEdges() {
	for k := range _parserDiag.Edges {
		e := _parserDiag.Edges[k]
		e.Start.Edges = append(e.Start.Edges, e)
		e.End.Edges = append(e.End.Edges, e)
	}
}
