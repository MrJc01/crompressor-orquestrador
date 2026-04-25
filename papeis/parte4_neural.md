# Parte 4 — Pesquisa Neural: Do CromGPT ao PTQ

## 4.1 A Tese Central

A compressão de dados e a inteligência artificial compartilham uma base matemática: **ambas buscam representações eficientes de informação**. O Crompressor Neural explora essa dualidade.

A tese: se um codebook CROM pode representar padrões de dados com eficiência, então pode substituir os pesos densos de uma rede neural — comprimindo o modelo sem perder inteligência.

## 4.2 Pesquisa 0 — Inteligência Dimensional (5D Active Inference)

**Repo**: `crompressor-neuronio/pesquisa0/`

Investigação teórica sobre percepção em múltiplas dimensões. 6 laboratórios:

1. **Percepção Temporal** — como agentes comprimem sequências temporais
2. **Observadores** — múltiplos pontos de vista sobre a mesma informação
3. **Simulação de Realidades** — worlds como codebooks de experiência
4. **Dimensões** — representação em espaços 5D
5. **IA Dimensional** — agentes que operam em espaços comprimidos
6. **Integração CROM** — como isso se conecta ao motor de compressão

**Resultado**: Framework teórico para Active Inference usando codebooks dimensionais.

## 4.3 Pesquisa 1 — CromGPT (CromLinear Transformer)

### O Conceito
Substituir `nn.Linear` (multiplicação de matrizes densa) por `CromLinear` (lookup em codebook + interpolação):

```python
# nn.Linear tradicional: y = Wx + b
# W tem shape [out_features, in_features] = milhões de parâmetros

# CromLinear: y = codebook_lookup(quantize(x)) + bias
# codebook tem K centroides de dimensão D = muito menos parâmetros
```

### Treinamento
- **Modelo**: Transformer 125M parâmetros
- **Dados**: Wikipedia Português
- **Hardware**: GPU cloud (Vast.ai)
- **Resultado**: Modelo gera texto coerente em português, com loss comparável ao baseline nn.Linear

### Formato .crom v3
Formato customizado para serializar modelos neurais comprimidos:
- Codebook de pesos por camada
- Índices de quantização
- Biases em full precision

## 4.4 Pesquisa 2 — Validação Qualitativa

10 amostras de chat interativo com o CromGPT 125M:
- Avaliação da coerência semântica
- Comparação com baseline
- Documentação em `papel2.md`

**Conclusão**: O modelo mantém ~95% da coerência do baseline com 4x menos parâmetros efetivos.

## 4.5 Pesquisa 3 — Post-Training Quantization (PTQ) de Modelos SOTA

### O Problema
Rodar Llama-3 ou Phi-3 em CPU com <8GB RAM.

### Abordagem: Hybrid PTQ
Em vez de quantizar todo o modelo uniformemente:

1. **Preservar** camadas de Attention e LM-Head em FP16 (críticas para coerência)
2. **Comprimir** apenas camadas FFN com codebooks CROM
3. **Usar** kernels AVX nativos do PyTorch para as camadas preservadas

```
Camada          Precisão    Motivo
───────────────────────────────────────
Embedding       FP16        Vocabulário original
Attention (QKV) FP16        Mecanismo de atenção (crítico)
FFN (up/down)   CROM VQ     Compressível (redundância alta)
LM Head         FP16        Output token (crítico)
```

### Pipeline C++ Custom
```
crompressor-neuronio/pesquisa3/
├── kernels/
│   ├── crom_linear_cpu.cpp   # Kernel AVX para forward pass
│   └── setup.py              # PyTorch extension build
├── ptq_compressor.py          # Pipeline de quantização
└── local_chat.py              # Chat interativo local
```

## 4.6 crompressor-ia — Edge Deployment

**Problema**: Rodar IA em dispositivos sem GPU.
**Solução**: Integrar com `llama.cpp` via wrapper Go/Bash.

Funcionalidades:
- PTY wrapper sobre `llama-cli`
- KV cache persistente entre sessões
- Modelo comprimido via codebook CROM

## 4.7 Conexão com a Matemática

O `crompressor-matematica` fundamenta tudo isso com provas formais:

| Papel | Tema | Resultado |
|-------|------|-----------|
| papel0 | Fundamentos VQ e Rate-Distortion | Prova de convergência do K-Means |
| papel1 | Limites de compressão codebook-based | Lower bound teórico |
| papel2 | CromGPT resultados empíricos | 125M modelo validado |
| papel3 | Análise de entropia em codebooks | Classificação de domínios |
| papel4 | CDC vs Fixed chunking | Análise de sincronização |
| papel5 | 15 Testes de Impossibilidade | Stress-test dos limites teóricos |
| papel6 | Bifurcação de Shannon | Resposta ao peer-review |

---

**Próximo**: [Parte 5 — Cronologia](./parte5_cronologia.md)
