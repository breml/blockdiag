package blockdiag

var diag Diag

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

func updateEdges() {
	for k := range diag.Edges {
		e := diag.Edges[k]
		e.Start.Edges = append(e.Start.Edges, e)
		e.End.Edges = append(e.End.Edges, e)
	}
}
