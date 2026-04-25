# 3. Comparação Injusta: A Derrota na Compressão Clássica

Para manter o rigor e a transparência científica do ecossistema, resgatamos os primeiros benchmarks (V1 a V4) executados no laboratório do CROM, onde tentamos forçar a engine a agir como um compactador de arquivo único comum (como GZIP e Zstandard). 

**Os resultados foram catastróficos para o CROM** – e isso foi essencial para compreendermos sua real natureza.

## A Metodologia do Fracasso

Nos Benchmarks Iniciais, pegamos arquivos totalmente não-correlacionados, treinamos um cérebro neles próprios (Overfitting) e tentamos comprimir. Como a Similaridade de Bits (Hamming) exigia que um chunk inteiro de 128 bytes (1024 bits) fosse idêntico para valer a pena, a maioria dos chunks deu um hit de similaridade baixo.
O que o motor fez? Ele transformou quase todo o arquivo em **Literais** (pedaços crús de dados) e ainda adicionou 24 bytes de cabeçalho da *Chunk Table* para cada um.

## Resultados do Benchmark V4 (Arquivos Isolados vs ZSTD)

| Arquivo Testado | Ferramenta | Tamanho Final | % Original | Tempo |
|:---|:---|:---|:---|:---|
| **Logs API** (260 MB) | ZSTD | **20.5 MB** | 8% | 2.1s |
| | GZIP | 35.8 MB | 13% | 5.3s |
| | CROM (128 bytes) | *325.0 MB* | **125%** | 4.8s |
| **VM Dump Binário** (517MB)| ZSTD | **180.2 MB** | 34% | 3.5s |
| | GZIP | 220.1 MB | 42% | 14.1s |
| | CROM (128 bytes) | *630.0 MB* | **121%** | 8.2s |

### A Análise Honesta do Erro

Observe os números: Em todos os cenários de arquivo isolado, **o Crompressor AUMENTOU o tamanho do arquivo em mais de 20%**. 
Isso confirmou que o CROM é estruturalmente incompetente para comprimir textos contínuos sem um dicionário dinâmico integrado intra-arquivo (como os algoritmos de Lempel-Ziv usados no ZSTD). O CROM *precisa* da tabela estática de Cérebro. 

Isso nos levou à revelação crucial:
> ❌ **O Crompressor NÃO É um "novo ZIP"**. Ele jamais deve ser usado para compactar um único arquivo isoladamente em sua máquina.

Mas se ele falha tão feio aí, onde está o valor matemático que criamos? A resposta está na próxima página, no documento 04, onde migramos de *Compressão Isolada* para *Deduplicação de Borda P2P*.
