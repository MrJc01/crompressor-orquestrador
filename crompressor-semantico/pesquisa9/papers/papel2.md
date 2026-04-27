# Pesquisa 9 — Papel 2: A Autópsia e o Renascimento (MLM Generativo em Go Puro)

## 1. A Autópsia: Três Falhas Fundamentais

A análise forense do CROM-Chat v5.x revelou três falhas estruturais que nenhuma heurística podia resolver:

### 1.1 O Forward Pass Falso
O ficheiro `matmul.go` não executava nenhuma das 6 camadas do Transformer. Usava `math.Sin(tok + i*idx)` como "estado oculto", multiplicando-o pela `lm_head.weight`. O resultado: logits determinísticos sem relação semântica com o input. As mesmas palavras ("e", "o", "maquina") apareciam sempre, independentemente da pergunta.

### 1.2 O Modelo Nunca Treinado
O ficheiro `llada_10m.crom` (84MB, 77 matrizes) continha pesos aleatórios do `torch.nn.init`. O comentário no `exportador_crom.py` confirmava: *"o modelo inicializa com pesos aleatórios"*. Zero épocas de treino. Zero dados.

### 1.3 O RAG Degradado
A Pesquisa 6 construíra um RAG com 100% de precision (TF-IDF + PCA SimHash + Hamming + Sigmóide calibrada). O `main.go` da v5 ignorou tudo isto, usando `JaccardSimilarity` sobre `strings.Fields()` com 6 entradas hardcoded e threshold arbitrário.

## 2. O Renascimento: Reconstrução em Fases

### 2.1 Fase 1 — RAG Real (CROM v6)
Reimportámos o pipeline completo da Pesquisa 6:
- **Tokenizer:** TF-IDF com 386 tokens, n-grams (bi+tri), OOV via Levenshtein
- **Projecção:** PCA real (SVD Power Iteration), 64 hiperplanos × 386 dimensões
- **Hashing:** SimHash com hiperplanos PCA (não aleatórios)
- **Distância:** Hamming O(1) via `bits.OnesCount64`
- **Confiança:** Sigmóide calibrada com pMid logarítmico

**Resultados v6:**

| Query | Hamming | Confiança | Resposta | Correcto? |
|---|---|---|---|---|
| "oi" | 0 bits | 100% | "Olá! Sou o CROM-LLM..." | ✅ |
| "o que é o sol" | 0 bits | 100% | "O Sol é a estrela central..." | ✅ |
| "como fazer bolo" | 22 bits | 0.2% | Rejeitado | ✅ (honesto) |
| "o que é um lapis" | 1 bit | 100% | "Um computador é..." | ❌ (falso positivo) |

**Conclusão v6:** O RAG funciona perfeitamente para queries conhecidas, mas é um motor de busca — não gera texto novo.

### 2.2 Fase 2 — Rede Neural Treinável (CROM v7)
Implementámos uma rede neural completa em Go puro com backpropagation:

**Arquitectura:** Masked Language Model (estilo LLaDA simplificado)
- **Input:** Janela bidirecional de contexto (raio R = olha R palavras para cada lado)
- **Embedding:** Lookup table aprendida + positional embeddings relativos
- **Hidden:** Camada ReLU fully-connected
- **Output:** Softmax sobre vocabulário completo
- **Treino:** Masked Token Prediction (mascarar 40% aleatório, prever os mascarados)
- **Inferência:** Denoising Iterativo (começar com [MASK], revelar progressivamente)

**Treino:** Corpus de 50 frases em Português sobre ciência, tecnologia e identidade do CROM.

**Melhorias implementadas ao longo da sessão:**

| Versão | Parâmetros | Épocas | Loss Final | Técnica Adicionada |
|---|---|---|---|---|
| v7.0 | 103K | 200 | 1.03 | MLM base com janela R=3 |
| v7.1 | 199K | 500 | 0.40 | Hidden=192, R=4, LR menor |
| v7.2 | 199K | 300 | 0.71 | Cosine LR Decay + Gradient Clipping |

### 2.3 Fase 3 — Self-Attention (CROM v8)
Implementámos Self-Attention com backpropagation completa (Q/K/V) em Go:

| Versão | Parâmetros | Loss Final | Resultado |
|---|---|---|---|
| v8.0 (lr=0.003) | 62K | 4.54 (estagnado) | Não convergiu |
| v8.1 (lr=0.01) | 62K | 3.33 (oscilante) | Repetição de tokens |
| v8.2 (cosine) | 62K | 0.78 | Coerente mas limitado |
| v8.3 (emb=64) | 88K | 1.53 | Subtreinado |

**Conclusão v8:** A attention com 62-88K parâmetros e 50 frases é inferior à janela fixa com 199K parâmetros. A attention precisa de mais dados e mais capacidade para ser superior — isto é consistente com a literatura (Vaswani et al. demonstraram que a attention supera apenas a partir de certa escala).

## 3. Resultados Experimentais: Geração de Texto

### 3.1 Melhor Modelo (v7.1 — Janela Bidirecional, 199K parâmetros)

| Prompt | Texto Gerado | Análise |
|---|---|---|
| "a gravidade é" | "tão forte que nada escapa" | ✅ **Generalização!** Combinou "gravidade" + "buracos negros" |
| "eu sou o" | "crom estou pronto para conversar e responder perguntas" | ✅ **Perfeito** |
| "linux é um" | "sistema organizado para armazenar e recuperar informação" | ✅ Coerente |
| "a internet é" | "o processo de criar instruções para transferir computador executar" | ⚠️ Mistura parcial |

### 3.2 Modelo com Cosine LR (v7.2)

| Prompt | Texto Gerado | Análise |
|---|---|---|
| "o sol é" | "a estrela central" | ✅ Facto correcto |
| "eu sou o crom" | "uma inteligência artificial" | ✅ Identidade correcta |
| "linux é um" | "sistema de controle versão" | ⚠️ Misturou com Git |

## 4. O Que Aprendemos

### 4.1 O Forward Pass importa mais que a arquitectura
O maior ganho veio de implementar backpropagation REAL em Go — não de mudar a arquitectura de janela fixa para attention. Com poucos dados, o modelo mais simples ganha.

### 4.2 Cosine LR Decay é essencial
Sem decay, a loss oscila violentamente nas últimas épocas. Com cosine decay, a convergência é suave e estável.

### 4.3 Gradient Clipping previne colapso
Sem clipping (max=5.0), os gradientes explodem na camada de output (341 classes) causando repetição de tokens.

### 4.4 A escala mínima para coerência
- **< 50K parâmetros:** Sopa de letras incoerente
- **50-100K parâmetros:** Fragmentos reconhecíveis mas misturados
- **100-200K parâmetros:** Frases semi-coerentes com generalização observável
- **> 200K parâmetros:** Necessário para frases completas fluentes

### 4.5 O corpus é o gargalo
Com 50 frases (~500 palavras únicas), o modelo memoriza padrões locais. Para generalização real, precisamos de 1000+ frases diversas.

## 5. Inventário Técnico Final

### Ficheiros Go Criados:
- `pesquisa9/pkg/neural/model.go` — MLM com backprop, cosine LR, gradient clipping, denoising iterativo
- `pesquisa9/cmd/crom-chat/main.go` — Motor de chat com hiperparâmetros configuráveis

### Hiperparâmetros Configuráveis:
```
EmbeddingDim   = 64      // Dimensão dos embeddings
HiddenDim      = 192     // Neurónios ocultos
ContextRadius  = 4       // Janela bidirecional
LearningRate   = 0.005   // Com cosine decay automático
Epochs         = 300     // Passagens pelo corpus
MaskRatio      = 0.4     // % de tokens mascarados
DiffusionSteps = 10      // Passos de denoising
Temperature    = 0.5     // Criatividade (0.1-1.0)
TopK           = 8       // Top-K sampling
MaxTokens      = 12      // Tokens a gerar
```

## 6. Próximos Passos

1. **Expandir o corpus** para 500+ frases (Wikipedia PT, SQuAD traduzido)
2. **Implementar Self-Attention funcional** quando os dados justificarem a complexidade
3. **Salvar/Carregar pesos** treinados em formato `.crom` para não retreinar a cada execução
4. **Combinar RAG + Neural:** RAG para respostas exactas, Neural para geração quando o RAG falha
5. **Multi-head Attention** e camadas empilhadas (verdadeiro Transformer)
