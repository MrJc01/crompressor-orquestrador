# Pesquisa 6 — Papel 0: Do Mock à Busca Vetorial Real (RAG de Borda)

## Resumo
A Pesquisa 6 transformou o CROM-LLM de um protótipo com respostas hardcoded num Motor de Busca Vetorial K-NN autêntico, operando sobre TF-IDF real e projeções PCA calibradas via SVD. Este papel documenta as descobertas, os fracassos e a revolução matemática que emergiu do processo.

## 1. O Estado Inicial: A Fraude dos Mocks
O motor "conversacional" da Pesquisa 5 era uma mentira técnica:
- Respostas pré-programadas via `if strings.Contains("oi")`
- Vetores gerados por Feature Hashing aleatório (FNV-1a)
- Um ficheiro `knowledge_base.json` com 9 pares pergunta/resposta escritos à mão
- **Conclusão:** O sistema não era uma IA. Era um `switch/case` glorificado.

## 2. Fase 1: TF-IDF Nativo em Go
### 2.1 Decisão Arquitetural
Substituímos o Feature Hashing por TF-IDF (Term Frequency × Inverse Document Frequency) real:
- **Python** (`ingestao_chat.py`): Treina o vocabulário, calcula os pesos IDF de cada token, exporta `vocabulario.json`
- **Go** (`tokenizer_tfidf.go`): Carrega o vocabulário e gera vetores esparsos exatos em tempo real

### 2.2 Bug Crítico: Desalinhamento de Tokenização
O Python usava `re.findall(r'\b[a-z0-9]+\b')` (remove acentos), enquanto o Go usava `unicode.IsLetter()` (mantém acentos). A palavra "é" existia no Go mas não no Python, corrompendo os vetores.

**Solução:** Importar `golang.org/x/text` e replicar exactamente a regex do Python no Go.

### 2.3 Descoberta: Colisão Semântica Aguda
Com apenas 73 tokens e 10 interações, o espaço latente era tão comprimido que:
- "Gravidade" → desconhecido → vetor de zeros → colide com "Dólar" (20 bits)
- A Sigmóide "preguiçosa" aceitava 20 bits como 70% de confiança
- **Resultado:** O motor "alucinava" por falta de dados, não por bug de código

## 3. Fase 2: Escala SQuAD e Guilhotina na Sigmóide
### 3.1 Ingestão do SQuAD v1.1
O script Python foi refatorado para baixar 5.000 amostras de Q&A do Stanford Question Answering Dataset. O vocabulário escalou de 73 para **5.920 tokens**.

### 3.2 Motor de Rejeição
A Sigmóide foi endurecida:
- $f(x) = \frac{1}{1 + e^{k(x - p_{mid})}}$ com $k=1.0$ e $p_{mid}=8.0$
- Qualquer distância > 12 bits → Confiança < 5% → **Rejeição absoluta**

### 3.3 O Gauntlet OOD (Out-of-Distribution)
50 desafios cegos (Star Trek, Elden Ring, receitas) foram injectados no motor.
**Resultado:** 100% de Rejeições Corretas. Zero Falsos Positivos.

## 4. Fase 3: A Revolução do PCA Real (SVD)
### 4.1 O Problema dos Hiperplanos Aleatórios
A função `gerar_matriz_fallback()` criava hiperplanos de projeção com `random.uniform(-1, 1)`. Estes eixos não respeitavam a distribuição real dos dados.

**Consequência Matemática:** A distância de Hamming esperada entre dois vetores com ângulo $\theta$ é:
$$d_{esperada} = 64 \times \frac{\theta}{180°}$$

Para "oi" (1 token) vs "oi oii oie" (3 tokens), o ângulo é $\arccos(1/\sqrt{3}) = 54.7°$, dando $d \approx 19.4$ bits. Mesmo para a mesma *intent*, o motor rejeitava a 17 bits!

### 4.2 A Solução: SVD Puro em Python (Power Iteration)
Implementámos PCA real sem NumPy:
1. **Centróide:** $\mu = \frac{1}{N}\sum_{i=1}^{N} \vec{v}_i$ (média de todas as amostras)
2. **Centralização:** $\vec{v}'_i = \vec{v}_i - \mu$ (subtrair a média)
3. **Power Iteration × 20:** Para cada um dos 64 componentes, extraímos o autovetor de máxima variância
4. **Deflação:** Removemos a variância capturada antes de extrair o próximo componente

### 4.3 Resultados: A Morte da Distância Espúria

| Query | Antes (Random) | Depois (PCA Real) |
|---|---|---|
| "oi" vs DB "oi" | 17 bits ❌ | **0 bits** ✅ |
| "universo" vs DB "universo" | 15 bits ❌ | **0 bits** ✅ |
| "oi" vs "universo" (divergência) | ~20 bits | **27 bits** ✅ |

O PCA Real colapsou conceitos idênticos para 0 bits e afastou conceitos diferentes para 27+ bits. A "zona cinzenta" foi eliminada.

## 5. O Muro: Retrieval vs. Geração
### 5.1 O Que o Motor Não Faz
O CROM-LLM Pesquisa 6 é um **Motor de Busca Vetorial (RAG)**, não um modelo generativo:
- Não compõe frases novas
- Não entende contexto entre turnos de conversa
- Não pode responder sobre "galáxias" se a palavra não está no dataset
- Está limitado a **repetir** as respostas pré-existentes mais próximas

### 5.2 O Que Seria Necessário Para "Pensar"
Para gerar texto original, o motor precisaria de:
- Uma camada `nn.Embedding` aprendida (não TF-IDF estático)
- Blocos de Self-Attention (Transformer)
- Treino por Backpropagation sobre milhões de tokens
- Uma "LM Head" que preveja o próximo token numa distribuição de 50K+ palavras

## 6. Métricas Finais da Pesquisa 6

| Métrica | Valor |
|---|---|
| Vocabulário | 386 tokens (offline) / 5.920 (SQuAD) |
| Dataset | 124 intents (offline) / 5.006 (SQuAD) |
| Latência de Inferência | 200ns - 2ms |
| PCA | SVD Real (Power Iteration, 64 componentes) |
| Precisão OOD Gauntlet | 100% (50/50 rejeições corretas) |
| Falsos Positivos | 0 |
| Sigmóide | k=1.0, pMid=8.0, corte=12 bits |

## 7. Próximos Passos (Pesquisas 7, 8, 9)
1. **Pesquisa 7 (RAG Massivo):** Escalar para 500K+ entradas com Wikipedia PT + SQuAD completo
2. **Pesquisa 8 (Micro-GPT):** Treinar um Transformer generativo de 10-50M parâmetros em PyTorch
3. **Pesquisa 9 (Micro-Diffusion):** Explorar a arquitectura LLaDA (masking + denoising)
