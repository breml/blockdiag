package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/breml/blockdiag"
)

func main() {
	simple := ""
	simple =
		`blockdiag {
   A -> B -> C;
   A -> D;
}`

	simple =
		`blockdiag {
   A -> B -> C;
   A -> D;
}`

	simple =
		`blockdiag {
			node_width = 128;
			A -> B -> C -> D;
			B -> F -> G -> X -> Y;
			C -> E;
			H -> I -> J -> H;
		}`

	simple =
		`blockdiag {
			A -> B -> C;
			B -> D -> E -> H;
			A -> F -> E;
			F -> G;
			X -> Y;
		}`

	got, err := blockdiag.ParseReader("simple.diag", strings.NewReader(simple))
	if err != nil {
		log.Fatal("Parse error:", err)
	}
	diag := got.(blockdiag.Diag)

	diag.PlaceInGrid()
	fmt.Printf("%s\n", diag.String())

	// fmt.Println("=", diag)

	/*
		for _, e := range diag.Edges {
			fmt.Println(e.Name)
		}
	*/

	//fmt.Println("Circular: ", diag.FindCircular())

	//fmt.Printf("Diag: %#v\n", &diag)
}
