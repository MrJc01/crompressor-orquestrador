package main

import (
	"fmt"
	"github.com/MrJc01/crompressor-semantico/pesquisa9/pkg/neural"
)

func main() {
	corpus := `
		A inteligência artificial é a área da ciência da computação.
		A computação é fundamental para a computação moderna.
		O computador processa a informação artificialmente.
		`

	fmt.Println("=== TESTE BPE ===")
	tokenizer := neural.NewBPETokenizer(50) // Limite pequeno para forçar subwords
	tokenizer.Train(corpus)

	testWord := "computador"
	encoded := tokenizer.Encode(testWord)
	fmt.Printf("\nEncode('%s') -> IDs: %v\n", testWord, encoded)
	fmt.Printf("DecodeTokens -> %v\n", tokenizer.DecodeTokens(encoded))
	fmt.Printf("Decode -> '%s'\n", tokenizer.Decode(encoded))

	testWord2 := "artificialmente"
	encoded2 := tokenizer.Encode(testWord2)
	fmt.Printf("\nEncode('%s') -> IDs: %v\n", testWord2, encoded2)
	fmt.Printf("DecodeTokens -> %v\n", tokenizer.DecodeTokens(encoded2))
	fmt.Printf("Decode -> '%s'\n", tokenizer.Decode(encoded2))

	testWord3 := "descomplicado" // OOV
	encoded3 := tokenizer.Encode(testWord3)
	fmt.Printf("\nEncode('%s') (OOV) -> IDs: %v\n", testWord3, encoded3)
	fmt.Printf("DecodeTokens -> %v\n", tokenizer.DecodeTokens(encoded3))
	fmt.Printf("Decode -> '%s'\n", tokenizer.Decode(encoded3))
}
