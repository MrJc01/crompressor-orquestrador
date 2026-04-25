# Codebooks para IA: quando compressão de dados encontra redes neurais — Product Quantization explicado

**TL;DR**: O mesmo dicionário (codebook) que comprime arquivos pode comprimir pesos de redes neurais. Isso não é invenção minha — Google e Meta já publicaram papers. Eu tentei implementar do zero em Go. Eis o que funcionou e o que não funcionou.

---

## O insight: compressão e IA são o mesmo problema

Compressão: representar dados com menos bits.
Rede neural: representar conhecimento com menos parâmetros.

As duas operam sobre a mesma matemática: **encontrar representações eficientes** em espaços de alta dimensão.

### O peso de um modelo

Um modelo como o Llama-3 8B tem **8 bilhões de parâmetros**. Cada um é um float32 (4 bytes). Total: **32 GB** de pesos.

```
Llama-3 8B:
  - 32 camadas Transformer
  - Cada camada: Attention (Q, K, V, O) + FFN (up, gate, down)
  - Total: ~8 bilhões de floats × 4 bytes = 32 GB
```

Para rodar em CPU com 8GB RAM, precisa comprimir. As técnicas padrão:

| Técnica | Ratio | Qualidade | Método |
|---------|-------|-----------|--------|
| FP16 | 2× | ~100% | Trocar float32 por float16 |
| INT8 | 4× | ~97% | Quantização linear |
| INT4 (GPTQ) | 8× | ~90% | Quantização calibrada |
| **VQ/PQ** | **10-16×** | **~85-95%** | Codebook de centroides |

## O que é Product Quantization (PQ)

Em vez de quantizar cada peso individualmente (INT4/INT8), PQ agrupa pesos em vetores e substitui cada grupo pelo centroide mais próximo em um codebook:

```
Pesos originais (vetor de 256 floats):
[0.123, -0.456, 0.789, ..., 0.321]

Codebook com 256 centroides:
  centroide[42] = [0.120, -0.460, 0.792, ..., 0.318]  ← mais próximo!

Armazenamento:
  ID: 42 (1 byte em vez de 1024 bytes)
  Compressão: 1024× !!
```

### Estado da arte (2024-2026)

- **AQLM** (Google, 2024): Additive Quantization para LLMs — múltiplos codebooks sobrepostos
- **QuIP#** (Cornell, 2024): Quantization with Incoherence Processing — usa rotação dos pesos
- **GPTQ** (IST/ETH, 2023): Post-Training Quantization calibrada por camada

## O que eu fiz: CromLinear

Implementei uma camada `CromLinear` que substitui `nn.Linear` do PyTorch:

```python
class CromLinear(nn.Module):
    def __init__(self, in_features, out_features, n_codes=256):
        self.codebook = nn.Parameter(torch.randn(n_codes, in_features))
        self.indices = nn.Parameter(torch.zeros(out_features, dtype=torch.long))
        self.bias = nn.Parameter(torch.zeros(out_features))

    def forward(self, x):
        # Lookup: cada linha de "pesos" é um centroide do codebook
        weights = self.codebook[self.indices]  # shape: [out, in]
        return x @ weights.T + self.bias
```

### CromGPT 125M — o que funcionou

- Modelo Transformer com 12 camadas, 768 hidden, 12 heads
- Todas as `nn.Linear` substituídas por `CromLinear`
- Treinado na Wikipedia PT (Vast.ai GPU)
- **Resultado**: gera texto coerente em português

```
Prompt: "O Brasil é um país"
CromGPT: "O Brasil é um país de contrastes, onde a diversidade cultural
se manifesta em cada região, desde o sertão nordestino até as metrópoles
do sudeste..."
```

### Hybrid PTQ — o que NÃO funcionou (ainda)

Tentei aplicar PQ em modelos SOTA pré-treinados (Llama-3, Phi-3) sem re-treinar. O resultado:

```
Antes (FP16): "O céu é azul por causa do espalhamento Rayleigh..."
Depois (PQ):  "O céu azul azul do azul Rayleigh Rayleigh..."
```

**Diagnóstico**: comprimir as camadas de Attention destrói a coerência. A solução (em progresso):

```
Camada          Precisão    Motivo
──────────────────────────────────────
Embedding       FP16        Vocabulário original
Attention QKV   FP16        Crítico para coerência
FFN (up/down)   PQ/VQ       Compressível (alta redundância)
LM Head         FP16        Output de tokens
```

Preservar Attention em FP16 + comprimir só FFN = **Hybrid PTQ**. AQLM e QuIP# fazem algo similar.

---

## A conexão com compressão de dados

O codebook do Crompressor (motor de compressão) e o codebook do CromLinear (motor neural) são **a mesma estrutura de dados**: um array de centroides treinados por K-Means.

```
Codebook de compressão:  128 bytes × 16384 centroides = 2MB
Codebook de IA:          256 floats × 256 centroides × 32 camadas = ~8MB
```

A visão de longo prazo: um único codebook universal que serve tanto para comprimir dados quanto para comprimir modelos neurais. **Isso ainda não funciona.** É pesquisa em andamento.

---

## Referências reais (não sou eu inventando)

1. **AQLM**: Egiazarian et al., "AQLM: Extreme Compression of LLMs via Additive Quantization", ICML 2024
2. **QuIP#**: Tseng et al., "QuIP#: Even Better LLM Quantization with Hadamard Incoherence", NeurIPS 2024
3. **GPTQ**: Frantar et al., "GPTQ: Accurate Post-Training Quantization for GPT", ICLR 2023
4. **Product Quantization**: Jegou et al., "Product Quantization for Nearest Neighbor Search", IEEE TPAMI 2011

---

## Como contribuir

Se você tem experiência com:
- **Compressão de modelos**: ajuda a validar o Hybrid PTQ
- **PyTorch internals**: otimizar o CromLinear para ser competitivo com nn.Linear
- **Matemática formal**: revisar os 7 papéis em crompressor-matematica

Repositório: https://github.com/MrJc01/crompressor-neuronio

---

*MrJ — [crom.run](https://crom.run) — Abril 2026*
