# Pesquisa 8 — Papel 0: Arquitectura do Micro-GPT (Transformer Generativo)

## 1. Motivação
O RAG do CROM (Pesquisas 6-7) provou que a busca vetorial via TF-IDF + PCA-LSH é ultra-rápida ($<2ms$) e precisa para retrieval. No entanto, o motor é fundamentalmente **não-generativo**: ele repete respostas pré-existentes. Um verdadeiro modelo de linguagem precisa de **compor** texto original, token a token.

Esta pesquisa documenta a implementação de um Decoder-Only Transformer (estilo GPT) treinado do zero em PyTorch.

## 2. Anatomia de um Transformer Generativo

### 2.1 Tokenização (BPE — Byte Pair Encoding)
Ao contrário do TF-IDF (que trata palavras inteiras), os LLMs modernos usam **subword tokenization**:
- "desconhecido" → ["des", "conhec", "ido"]
- Vocabulário típico: 32K-50K subwords
- Vantagem: Consegue representar QUALQUER palavra, mesmo nunca vista antes

### 2.2 Embedding Layer ($E$)
Cada token ID é convertido num vetor denso aprendido:
$$\vec{e}_i = E[\text{token\_id}_i] \in \mathbb{R}^{d_{model}}$$
Não é TF-IDF. Não é frequência. É um vetor **aprendido por backpropagation** para representar o significado do token no contexto do treino.

### 2.3 Positional Encoding
Como o Transformer não tem noção de ordem sequencial, adicionamos informação posicional:
$$\vec{h}_i^{(0)} = \vec{e}_i + \vec{p}_i$$

### 2.4 Transformer Block (×N camadas)
Cada bloco contém:

#### Self-Attention (Causal)
$$\text{Attention}(Q, K, V) = \text{softmax}\left(\frac{QK^T}{\sqrt{d_k}} + M_{\text{causal}}\right) V$$
- $Q = XW_Q$, $K = XW_K$, $V = XW_V$ (projeções lineares)
- $M_{\text{causal}}$: Máscara triangular que impede ver tokens futuros
- O modelo aprende **quais tokens prestar atenção** para prever o próximo

#### Feed-Forward Network (MLP)
$$\text{FFN}(x) = \text{GELU}(xW_1 + b_1)W_2 + b_2$$
Expansão 4x: se $d_{model}=256$, a camada interna tem $1024$ neurónios.

#### LayerNorm + Residual
$$\vec{h}^{(l+1)} = \text{LayerNorm}(\vec{h}^{(l)} + \text{Attention}(\vec{h}^{(l)}))$$

### 2.5 LM Head (Cabeça de Linguagem)
A última camada projeta o vetor de volta para o espaço do vocabulário:
$$P(\text{next\_token}) = \text{softmax}(hW_{vocab})$$

O modelo devolve a **probabilidade de cada uma das 32K palavras** ser a próxima. Amostramos dessa distribuição (com temperatura) para gerar texto.

## 3. O Treino: Next-Token Prediction

```
Corpus: "O universo é vasto e misterioso"
Input:   [O]  [universo] [é]  [vasto] [e]
Target:  [universo] [é] [vasto] [e] [misterioso]

Loss = CrossEntropy(predicted, target)
Optimizer = AdamW
```

O modelo é treinado para maximizar $P(\text{token}_{t+1} | \text{token}_1, ..., \text{token}_t)$.

## 4. Configuração do Micro-GPT CROM

| Parâmetro | Valor |
|---|---|
| `vocab_size` | 32.000 (BPE) |
| `n_embd` (d_model) | 256 |
| `n_head` | 4 |
| `n_layer` | 6 |
| `block_size` (contexto) | 256 tokens |
| `dropout` | 0.1 |
| **Total de Parâmetros** | **~10M** |

**Custo estimado de treino:** 4-8 horas em CPU moderno com 1-2GB de texto.

## 5. Comparação com o RAG (Pesquisa 6-7)

| Aspecto | RAG (K-NN) | Micro-GPT |
|---|---|---|
| Geração | Repete respostas | Compõe texto novo |
| Velocidade | ~0.5ms por resposta | ~50-100ms por token |
| Precisão factual | 100% (copia do dataset) | Pode alucinar |
| Escalabilidade | Linear com dataset | Logarítmica com treino |
| Conhecimento novo | Requer actualização do dataset | Internalizado nos pesos |

## 6. Próximos Passos
1. Implementar o tokenizer BPE
2. Construir a classe `MicroGPT` em PyTorch
3. Treinar com corpus PT
4. Exportar pesos para formato `.crom` para inferência em Go
