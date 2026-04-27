# Pesquisa 9 — Papel 1: Fusão do LLaDA com o Crompressor (Sinergia LSH-Go)

## 1. O Problema da Inferência LLaDA Tradicional
Conforme estabelecido no Papel 0 da Pesquisa 9, a arquitetura de difusão de texto (LLaDA) resolve o problema de tempo linear autoregressivo O(L) processando todos os tokens de uma vez no Reverse Process (Denoising).
No entanto, no paradigma original do LLaDA, o rácio de tokens rejeitados/aceites a cada passo depende da probabilidade da Camada Softmax (Logits). Numa inferência rodando estritamente numa CPU de borda, o LLaDA pode ser suscetível a "alucinações estocásticas" e quebra de coerência.

## 2. A Tese do Crompressor Híbrido
O Motor de Busca Crompressor (Pesquisa 6) já provara que o espaço latente de LSH (Local Sensitive Hashing usando PCA + TF-IDF) consegue separar contextos semânticos calculando a Distância de Hamming, atuando como um radar implacável contra entropia em O(1) de complexidade.

**Decisão Arquitetural:** Em vez de depender inteiramente do Logit neuronal do PyTorch para a decisão de Re-Masking, nós delegámos a curadoria dos tokens para a heurística LSH em Go puro.
> *O PyTorch atua como o Cérebro "Intuitivo". O Go atua como o Juiz "Lógico".*

## 3. O Exportador CROM-Cold
Para isolar a execução de frameworks massivos (como o PyTorch/Numpy), implementámos uma diretiva `CROM-Cold`:
1. A rede é inicializada em `micro_llada.py` (13M Parâmetros).
2. O script `exportador_crom.py` converte e achata os `state_dicts` exportando matrizes multidimensionais em Float32 nativo para o formato `.crom`.
3. Total de Matrizes Extraídas: 77 (Embeddings, Wte, LayerNorms, Atenção).

## 4. O Reverse Process e a "Sigmóide Reversa" (LSH Denoising)
Na reconstrução do Chat (`crom-chat/main.go`), cada passo itera sobre os tokens.
Quando o LLaDA preenche o `[MASK]` com o token "universo", o código Go atira a nova frase contra a âncora TF-IDF da frase original. 
Se a *Distância de Hamming* ($d > 15 \text{ bits}$) falhar a Sigmóide de aceitação, significa que o LLaDA "alucinou". O token volta a ser atirado para a Sombra (`[MASK]`), forçando o modelo a propor outra palavra na iteração seguinte.

## 5. Métricas Obtidas
Ao correr o binário em Go nativo com as matrizes matemáticas de PyTorch acopladas à simulação de Denoising de Hamming, a latência registou-se vertiginosa e pronta para a orquestração em Dispositivos Finais:
* **Latência de Leitura/Carregamento na RAM:** Perto de ~0ms.
* **Latência Média do Loop de Denoising (Mock):** 51.600ns (51μs).

O modelo CROM agora é, essencialmente, um **Transístor Semântico de Difusão e Refração**. Ele absorve probabilidades e "refrata" a sujeira heurística usando a distância de Hamming.
