package main

import (
	"fmt"
	"github.com/MrJc01/crompressor-semantico/pesquisa6/pkg/nlp"
)

func main() {
	tokens := nlp.ExtrairTokens("meu nome é jorge")
	fmt.Printf("%q\n", tokens)
}
