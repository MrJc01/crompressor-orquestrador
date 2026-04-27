# Pesquisa 9 — Papel 0: Arquitectura LLaDA (Difusão Mascarada para Texto)

## 1. O Paradigma da Difusão Aplicada a Texto

### 1.1 Contexto Histórico
A difusão revolucionou a geração de imagens (Stable Diffusion, DALL-E 3) em 2022-2024. A ideia central:
1. **Forward Process:** Adicionar ruído progressivamente a uma imagem limpa até virar ruído puro
2. **Reverse Process:** Treinar uma rede neural para **reverter** o ruído, reconstruindo a imagem

Em 2025, investigadores aplicaram este princípio ao **texto**, substituindo "ruído gaussiano" por **masking de tokens**.

### 1.2 O Problema com Texto
Texto é **discreto** (palavras), não contínuo (pixels). Não podemos "adicionar ruído gaussiano" a uma palavra. Solução:
- **Forward:** Substituir tokens por `[MASK]` (em vez de adicionar ruído)
- **Reverse:** Treinar o modelo para prever os tokens mascarados (em vez de "limpar" o ruído)

## 2. LLaDA: Large Language Diffusion with mAsking

### 2.1 Arquitectura
O LLaDA usa um **Transformer bidirecional** (sem máscara causal) como backbone:

```
Forward Process (Treino):
  "O universo é vasto" 
  → (t=0.3) "O [MASK] é vasto"
  → (t=0.6) "[MASK] [MASK] é [MASK]"
  → (t=1.0) "[MASK] [MASK] [MASK] [MASK]"

Reverse Process (Inferência):
  "[MASK] [MASK] [MASK] [MASK]"
  → Passo 1: "[MASK] universo [MASK] [MASK]"   (prevê os mais fáceis)
  → Passo 2: "O universo é [MASK]"              (refina)
  → Passo 3: "O universo é vasto"               (converge)
```

### 2.2 Diferenças Fundamentais em Relação ao GPT

| Aspecto | GPT (Autoregressive) | LLaDA (Diffusion) |
|---|---|---|
| **Direcção** | Esquerda → Direita (causal) | **Bidirecional** (vê tudo) |
| **Geração** | 1 token de cada vez | **Todos de uma vez**, refinados iterativamente |
| **Atenção** | Causal mask (triângulo) | Full attention (sem máscara) |
| **Treino** | $P(x_t | x_1, ..., x_{t-1})$ | $P(x_t | x_{\text{unmasked}})$ |
| **"Reversal Curse"** | Sim (não aprende A↔B simétrico) | **Não** (contexto bidirecional) |

### 2.3 O Treino
1. Amostrar um ratio de masking $t \sim U[0, 1]$
2. Mascarar $t$% dos tokens da sequência
3. O Transformer bidirecional prevê os tokens mascarados
4. Loss = Cross-Entropy entre predição e tokens originais

### 2.4 A Inferência (Denoising Iterativo)
1. Começar com sequência 100% mascarada
2. Em cada passo:
   a. O modelo prevê todos os tokens
   b. Aceitar os tokens com alta confiança
   c. **Re-mascarar** os de baixa confiança
3. Repetir ~20 passos até convergir

## 3. Vantagens para o CROM

### 3.1 Geração Paralela
Ao contrário do GPT que gera token a token ($O(L)$ passos para $L$ tokens), o LLaDA gera todos simultaneamente. Com 20 passos de refinamento, uma resposta de 100 tokens requer **20 forward passes** (vs 100 no GPT).

### 3.2 Sinergia com o Motor LSH Existente
O conceito de "masking" é análogo ao que já fazemos com o SimHash:
- Os bits "1" são "tokens revelados"
- Os bits "0" são "tokens mascarados"
- O motor de Hamming Distance pode ser reutilizado para guiar o denoising

### 3.3 Go-Native
A inferência de difusão não requer backward pass. É apenas forward pass repetido — compatível com matmul puro em Go.

## 4. Modelo Proposto: Micro-LLaDA CROM

| Parâmetro | Valor |
|---|---|
| `vocab_size` | 32.000 (BPE) |
| `n_embd` | 256 |
| `n_head` | 4 (bidirectional) |
| `n_layer` | 6 |
| `block_size` | 256 |
| `diffusion_steps` | 20 |
| **Total de Parâmetros** | **~10M** |

## 5. Riscos e Limitações
- **Pesquisa de Ponta:** Poucos exemplos "from scratch" disponíveis
- **Qualidade vs Velocidade:** As amostras de difusão podem ser menos coerentes que GPT autoregressive
- **Complexidade:** O scheduling do masking ratio e a estratégia de remasking requerem experimentação

## 6. Próximos Passos
1. Estudar o código-fonte do LLaDA no GitHub
2. Implementar o Forward Process (masking schedule)
3. Implementar o Reverse Process (denoising com Transformer bidirecional)
4. Comparar com o Micro-GPT da Pesquisa 8
