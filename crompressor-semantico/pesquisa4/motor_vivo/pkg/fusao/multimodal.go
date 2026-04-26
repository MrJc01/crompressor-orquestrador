package fusao

// GerarHashHibrido combina o Hash do Texto e o Hash da Visão num vetor só.
// Ele extrai os 32 bits mais significativos da Imagem e une aos 32 bits 
// menos significativos do Texto. O resultado é um hash consolidado onde 
// "Carro" + "Velocidade" terão uma colisão de Bucket distinta de 
// "Carro" + "Cor".
func GerarHashHibrido(hashImagem, hashTexto uint64) uint64 {
	// A Imagem compõe o Contexto Maior (High Bits)
	contextoVisual := hashImagem & 0xFFFFFFFF00000000
	// A Pergunta compõe a Intenção Fina (Low Bits)
	intencaoTexto := hashTexto & 0x00000000FFFFFFFF
	
	return contextoVisual | intencaoTexto
}
