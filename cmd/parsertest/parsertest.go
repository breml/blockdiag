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
	A -> C;
	# B -> D;
}`

	got, err := blockdiag.ParseReader("simple.diag", strings.NewReader(simple))
	if err != nil {
		log.Fatal("Parse error:", err)
	}
	diag := got.(blockdiag.Diag)

	diag.PlaceInGrid()
	fmt.Printf("%s\n", diag.String())
}
