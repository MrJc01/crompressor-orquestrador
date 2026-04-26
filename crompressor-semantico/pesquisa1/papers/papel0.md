# Papel 0: Fundação da Deduplicação Semântica (Conclusão da Pesquisa 1)

## Resumo Executivo
A Pesquisa 1 do ecossistema CROM validou matematicamente a premissa de que a deduplicação de memória na Borda pode transcender a entropia de bytes (compressão tradicional) e operar no campo da similaridade de significados. Substituímos hashes criptográficos rígidos (SHA-256) por Locality-Sensitive Hashing (LSH), permitindo buscar contextos similares no cache.

## Resultados Atingidos
1. **Velocidade Brutal O(1):** Ao transformar tensores de *embeddings* em *hashes* binários de 64 bits, movemos a complexidade de busca da GPU (Dot-Product) para a CPU usando a instrução nativa `POPCNT` (`bits.OnesCount64` em Go). Provamos via benchmark latências inferiores a 1 nanossegundo por comparação.
2. **Deduplicação Dinâmica:** Frases com o mesmo peso semântico ("carro veloz" e "automovel veloz") geraram distâncias de Hamming pequenas, permitindo que um motor rudimentar identificasse a redundância antes de gastar ciclos de inferência de uma LLM.

## Limitações e a Ponte para a Pesquisa 2
Apesar do sucesso isolado, a arquitetura provou-se imatura para ambientes caóticos:
* **Busca:** Iterar sobre milhões de hashes, mesmo a 1ns cada, quebra a latência rigorosa do Crompressor. Precisamos de indexação profunda (*Multi-Index Hashing*).
* **Visão (Variância Espacial):** Dividir uma imagem e calcular o Hamming linear não suporta deslocamento de pixels (um objeto ligeiramente à esquerda muda o hash de múltiplos blocos simultaneamente).
* **Ausência de Contexto:** A deduplicação analisava apenas o *prompt* estático, sendo cega à sequência de uma conversa longa.

O **Papel 0** sela o fim da fase exploratória e introduz a **Pesquisa 2**, onde implementaremos *Overlapping Patches* focados em fóvea humana, *Multi-Index Bucketing* com 4 tabelas de 16 bits, e um *Attention Mechanism* puramente binário (intercalação 32/32) para contexto sequencial.
