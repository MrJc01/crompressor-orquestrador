# Paper Annotation: Reformer & LSH Attention

## O Problema do Transformer Clássico
A arquitetura base de *Large Language Models* (LLMs) depende do mecanismo de **Atenção (Self-Attention)**, cuja complexidade computacional cresce quadraticamente com o tamanho do contexto ($O(N^2)$). Se quisermos que o LLM analise um livro inteiro (ex: 64K tokens), a memória de GPU necessária explode, pois cada token precisa calcular seu "peso" em relação a todos os outros tokens da sequência.

## A Solução com Locality-Sensitive Hashing (LSH)
O artigo que introduziu o **Reformer** (Kitaev et al.) propôs a substituição do *Dot-Product Attention* denso pelo **LSH Attention**.

A tese é simples: a função Softmax no mecanismo de atenção é fortemente dominada pelos tokens mais similares (que produzem os maiores resultados numéricos). Em vez de calcular tudo contra tudo, podemos calcular a atenção *apenas entre tokens similares*.

1. O modelo usa **LSH (SimHash / Random Projections)** para "hashear" os *embeddings* de todos os tokens.
2. Tokens que caem no mesmo "bucket" de hash são agrupados.
3. O LLM calcula a Atenção apenas dentro do próprio bucket.

Isso derruba a complexidade matemática de $O(N^2)$ para $O(N \log N)$, viabilizando contextos absurdamente longos.

## Conexão com o Crompressor Semântico
No CROM, aplicamos a lógica do LSH não apenas para acelerar cálculos internos de inferência, mas como um **mecanismo de armazenamento (Deduplicação)**. Se dois blocos de texto geram vetores que colidem no LSH, eles ocupam exatamente o mesmo "espaço conceptual" e um deles pode ser deduplicado da memória cache, salvando RAM na Borda.
