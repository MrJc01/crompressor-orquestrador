# Pesquisa 4: A Deduplicação de Similaridade Semântica (LSH e Hamming)

## Resumo
A arquitetura base do **Crompressor** fundamenta-se na deduplicação de borda (Edge Deduplication) operando em tempo real (O(1)) através de tabelas de Hash e operações em memória. A limitação desse modelo é a sua inflexibilidade intrínseca: o hash criptográfico SHA-256 é avesso a similaridades. Um único bit invertido em um pacote de 4KB altera o Hash completamente.

A "Deduplicação de Similaridade Semântica" propõe abandonar hashes rígidos em favor do **Locality-Sensitive Hashing (LSH)**. Ao injetarmos um vetor de características (Embedding) em um algoritmo SimHash, dados matematicamente "próximos" (ex: duas fotos do mesmo objeto sob ângulos diferentes) geram assinaturas de 64 bits altamente similares. 

Isso nos permite reciclar a hiper-otimização em C++/Go da CPU, não para encontrar "bytes exatos", mas para criar um **Banco de Dados de Padrões P2P**.

## A Mágica de Hardware: POPCNT e Distância de Hamming
Se a assinatura (Hash Semântico) da Imagem A é `10101010...` e da Imagem B é `10101011...`, a distância de Hamming entre eles é apenas 1.
Na CPU moderna (arquitetura x86-64 ou ARM), calcular isso leva quase 0 ciclo de clock utilizando a instrução de hardware `POPCNT` (Population Count), mapeada no Go como `bits.OnesCount64`.

A fórmula da similaridade torna-se:
`HammingDistance(HashA, HashB) = POPCNT(HashA XOR HashB)`

## Implicações Lossy
Aplica-se esta técnica quando a semântica (O QUÊ a foto contém) é mais importante que a integridade pixel a pixel (COMO a foto foi gerada). Esta é a fundação para o **CromGPT Multimodal** na Borda, permitindo busca vetorial e classificação de padrões em nanosegundos, sem a necessidade de um Banco de Dados Vetorial pesado (ex: Milvus, Pinecone). O próprio Dicionário CROM resolve a questão.

## 🔮 Visão de Futuro: Sistema Estilo LLM para Texto
*Anotação baseada em feedback de pesquisa inicial:*
Após os testes de LSH e validação da Distância de Hamming para deduplicação semântica, o próximo passo evolutivo lógico será iniciarmos uma Pesquisa dedicada a como essa exata arquitetura O(1) pode ser adaptada para **Texto**. A ideia é construir um "Sistema estilo LLM" puramente baseado na busca e combinação de hashes semânticos textuais. Em vez de prever o próximo *token* via matrizes colossais de atenção, buscaríamos padrões semânticos em um banco P2P global deduplicado de "conceitos textuais", gerando respostas por *matching* de contexto instantâneo na Borda.
