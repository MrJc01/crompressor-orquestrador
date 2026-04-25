---
title: "30 Dias de Jornada Insana: Do Zero ao Limite Matemático em Go"
author: "MrJ"
---

Esta é a Parte 2 da série sobre a construção do motor Crompressor. Na [Parte 1](URL_DA_PARTE_1), nós assumimos nosso erro contra o Zstandard e mostramos como alcançamos **99.4% de economia de rede** na sincronização de dados usando Edge Deduplication P2P. 

Hoje, quero dar um passo atrás. Quero compartilhar com a comunidade a jornada visceral e humana de sair da "bolha do CRUD" e passar um mês escovando bits, estudando entropia matemática e lutando contra a memória em Golang.

---

## O Chamado à Aventura (A Síndrome do CRUD)

Por muito tempo, o ecossistema brasileiro pareceu estagnado na criação de APIs REST e sistemas de gestão de clínicas médicas. O famoso "voce_so_sabe_crud" bateu na minha porta. Eu me sentia travado intelectualmente. Sentia falta do deslumbramento de ver as engrenagens brutas da computação rodando. Da sensação de quando, aos 12 anos, eu fuçava pastas no sistema do Windows para burlar o limite de tempo da Lan-House.

Para o filósofo Clóvis de Barros Filho, *brio* é a energia e a força interior para não aceitar a derrota, buscando a excelência e o progresso intelectual com esforço, mesmo quando é preciso reler um texto difícil dezenas de vezes. O Crompressor foi a minha busca pelo brio. Eu queria criar uma solução que desafiasse a física do armazenamento e da rede.

## 30 Dias Codando no Limite

Construir o motor Crompressor (uma ferramenta que lê, quantiza e converte padrões binários em dicionários universais via Hashing Locality-Sensitive e Árvores-B) foi brutal. 

A arquitetura original, dividida em 12 repositórios interconectados, sofreu diversas metamorfoses. Comecei aprofundando em C++, mas a escalabilidade e o ecossistema seguro do **Golang** me chamaram. Fazer concorrência massiva em Go usando *goroutines* para ler e mastigar arquivos de centenas de megabytes num motor assíncrono exigiu repensar tudo que eu achava que sabia sobre gestão de memória e coleta de lixo (GC).

**Os maiores tombos técnicos da jornada:**
1. **Entropia de Shannon Inflexível:** No começo, eu estava filtrando chunks que achava "inúteis" baseado puramente em limiares de entropia teóricos. O resultado? O motor destruía dados úteis e não conseguia sincronizar repositórios de texto. A calibragem fina (reverter limites estritos de 7.8 para permitir dicionários amplos) foi como ajustar o carburador de um carro de Fórmula 1 com ele em movimento.
2. **Race Conditions no mDNS:** Quando lancei a camada P2P usando *go-libp2p*, o sistema local entrava em colapso. O mDNS descobria peers no background e abria conexões autênticas de Handshake Soberano no exato milissegundo em que os testes de integração tentavam fazer o mesmo manualmente. Resolver essas assincronias de rede ensinou-me sobre mutexes de uma forma que nenhum tutorial jamais fez.
3. **A Matemática Quebrada da Similaridade:** Tentei usar *Hamming Distance* para medir similaridade entre pequenos blocos de texto UTF-8. A matemática dizia que um simples espaço em branco (" ") deslocaria os bits e anularia completamente a semelhança. Tive que engolir o ego e entender que para o Crompressor brilhar, o matching precisava ser de 100% (Deduplicação de Borda), onde a mágica acontece pelo descarte O(1) de arquivos inteiros, trocando-os por ponteiros de 24 bytes no Cérebro.

## Além da Sopa de Portas (O Efeito Hacker)

Lembra quando você era adolescente e via o seriado Mr. Robot achando que trabalhar com computadores era dominar a máquina, manipular a realidade do sistema operacional, e ter o controle nas pontas dos dedos?

O que a bolha corporativa fez foi tentar transformar todo programador num montador de blocos LEGO usando frameworks JS e APIs pré-prontas. Construir o Crompressor me ensinou que o hardware ainda obedece fielmente às leis matemáticas e físicas da teoria da informação. E se você entender a teoria de Lempel-Ziv, entropia termodinâmica da informação, e manipulação de Ponteiros e Memória no mais baixo nível, você retoma aquele sentimento de "poder". 

O Crompressor hoje não é só um código de P2P; é um testemunho de que a engenharia de software profunda está ao alcance de todos.

## Próximos Passos
Se esse relato lhe inspirou a sair da zona de conforto do seu framework favorito e ir investigar o código-fonte de sistemas de baixo nível, sinta-se à vontade para visitar o repositório do Crompressor.

🔗 **Explore o Código-Fonte e Estude a Matemática do Motor**: [github.com/MrJc01/crompressor](https://github.com/MrJc01/crompressor) (O motor principal está em Go e foi estruturado agora como um Orquestrador central aberto a toda a comunidade).

### O que vem a seguir
No próximo artigo desta série (Parte 3), vou mergulhar na junção do Crompressor com a Inteligência Artificial. Vou explicar didaticamente o conceito dos **Codebooks Universais** (os nossos Cérebros) e como usar essas tabelas para quantizar Modelos de Linguagem Gigantescos (LLMs) localmente. Até lá!
