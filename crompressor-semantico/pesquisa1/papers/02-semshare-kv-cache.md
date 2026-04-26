# Paper Annotation: SemShareKV e Cache LSH em LLMs

## O Problema do Gargalo de Memória (Memory-Bound)
Rodar modelos gigantes (como o Llama-3 de bilhões de parâmetros) em CPUs ou máquinas de Borda esbarra num problema físico de memória conhecido como **KV Cache**. Durante a geração de texto, o modelo armazena os pares de Chave-Valor (*Key-Value*) de todos os tokens gerados anteriormente para não precisar recalculá-los a cada nova palavra. Com prompts longos e multi-turnos, esse cache engole gigabytes de RAM.

## Compartilhamento de Cache Baseado em Semântica (LSH)
Trabalhos recentes, impulsionados pela otimização extrema de inferência, têm pesquisado formas de **Deduplicar o KV Cache**.

Em arquiteturas como o *SemShareKV*, os pesquisadores notaram que usuários frequentemente enviam *prompts* que são semanticamente idênticos, embora possuam palavras levemente diferentes (ex: "Me dê a receita de bolo de fuba" vs "Como faço bolo de fuba").
Se um modelo tradicional as processar, ele cria dois KV Caches enormes do zero.

Aplica-se o **LSH (Locality-Sensitive Hashing)** no *embedding* do *prompt* inicial. Se o *prompt* A e o *prompt* B geram Hashes Semânticos que caem dentro de um limiar estreito de **Distância de Hamming**, o motor de inferência **DEDUPLICA** o cálculo e simplesmente injeta o KV Cache antigo na nova sessão.

## Conexão com o Crompressor Semântico
Essa abordagem valida 100% o que acabamos de escrever em `pkg/hamming` e `pkg/lsh`. O nosso motor O(1) de Distância de Hamming tem a capacidade de atuar como o **Gatekeeper de Memória** de um motor de inferência LLM nativo em Go (CromGPT). Antes de rodar as pesadas multiplicações de matrizes, tiramos o *Hash* do texto. Se o LSH apontar que a "pergunta" já foi feita (Deduplicação Semântica), buscamos o vetor pronto no nosso banco P2P e retornamos em milissegundos.
