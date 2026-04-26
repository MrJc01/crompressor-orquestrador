# Pesquisa 5: Atenção Esparsa via LSH Treinado e Multi-Cérebros

## Introdução: O Fim da Caixa Preta Monolítica
A arquitetura fundamental de Large Language Models (LLMs) tradicionais apoia-se em pesadas matrizes de atenção densa ($O(N^2)$). Na Pesquisa 4, o CROM provou que inferências de vizinhança na ordem de $O(1)$ são possíveis em hardware restrito, mas falhou ao tentar comprimir a complexidade de intenções humanas (fato vs. tom) em um único bloco de hashes aleatórios. A "Miopia Semântica" expôs que a verdadeira compressão de inteligência não pode ser estocástica.

A **Pesquisa 5** introduz a **Atenção Esparsa via LSH Treinado (PCA-LSH)**. Ao invés de lançar hiperplanos às cegas, o CROM passará a "enxergar" o espaço vetorial.

## Multicamadas de Atenção Binária
Inspirados no conceito de *Attention Heads* dos transformers clássicos, substituímos a força bruta flutuante por assinaturas paralelas rápidas: os **Multi-Cérebros**.

Cada dado processado não será mais um único inteiro de 64 bits, mas sim um complexo esparso:
1.  **Head A (Entidade/Substantivos):** Concentra a predição topológica nas fundações do objeto (quem/o quê). Distâncias altas aqui ativam uma rejeição precoce, poupando CPU.
2.  **Head B (Relação/Contexto):** Avalia os vetores de ação (verbos) e o "tom" pragmático da sentença.
3.  **Head C (Visual/Padrão):** Engajada em inputs multimodais para texturas e bordas.

## O "Consenso" e a Memória Ponderada
O "Match" neural passa a ser um sistema de votação. Se o texto fala sobre um *Banco* (Entidade alinhada), mas o contexto trata de um *Banco de Praça* ao invés de um *Banco Central* (Contexto distinto), o dissenso entre as cabeças aborta a identificação e invoca um novo bloco lógico.

Apoiando essa atenção dinâmica, a arquitetura introduz a **Memória Evolutiva Ponderada**. Diferente de um decaimento onde memórias esquecem dados de forma plana, a nova topologia protege os "bits de alta relevância estatística", permitindo que o motor simule uma retenção orgânica de curto prazo, limpando apenas a periferia irrelevante.

A era da "adivinhação estocástica" acabou. Bem-vindo à Atenção Binária Consciente.
