---
title: "30 Dias Insanos: Do Zero ao Limite Matemático em Go"
author: "MrJ"
---

Esta é a Parte 2 da série oficial sobre a construção do Crompressor. Se você não leu a [Parte 1 (A Ilusão da Compressão e a vitória de 99.4% no P2P)](https://github.com/MrJc01/crompressor/tree/main/papeis/tabnews/crompressor/01_a_ilusao_da_compressao.md), eu recomendo fortemente que pare e leia os fundamentos do nosso ecossistema e os números antes de prosseguir.

Hoje, não falarei das engrenagens do motor prontas, mas quero convidar você, dev do TabNews, para uma volta ao começo. O relato visceral, não censurado, de sair do CRUD, brigar com vazamentos de memória na madrugada, perder o controle da IA orquestradora e falhar incansavelmente até o código funcionar.

---

## 1. O Chamado à Aventura (e a Fuga do CRUD)

Por volta dos meus 14 anos, assisti a *Mr. Robot*. Aquela ideia do hacker, o "Batman de home office", explodiu minha cabeça. Contudo, ao longo dos anos, com a entrada brutal no mercado de trabalho tradicional, muito desse encantamento sumiu. Fui sugado pelas regras das empresas: fazer APIs em Node.js, levantar bancos de dados, escrever controllers simples e viver uma rotina corporativa (o tão falado "Você só sabe CRUD").

Em março de 2026, meu "brio" interno (aquela energia para buscar excelência e não aceitar o progresso estagnado que o filósofo Clóvis de Barros sempre defende) gritou. Eu precisava fazer algo *hardcore*. Algo que espremesse a máquina e me lembrasse do porquê entrei na área de tecnologia.

Foi aí que a ideia absurda veio: **E se um algoritmo de compressão já possuísse o dicionário dos meus dados antes de comprimi-los?**

```bash
# O primeiro respiro (26 de Março de 2026)
git init crompressor
echo "start" > README.md && git add . && git commit -m "start"
```

---

## 2. A Orquestração de IAs e as "Mentiras" do Algoritmo

Muitos desenvolvedores têm medo de admitir que usam IA, mas a verdade crua é: eu usei Inteligência Artificial desde o primeiríssimo dia. Não para escrever código cego, mas para entender profundamente sistemas complexos externos e aplicar ideias parecidas no motor do CROM. 

Para esse projeto massivo, eu usei o **Antigravity** e orquestrei mais de **30 contas do Google** (sendo uma delas paga). A dinâmica era clara:
* O **Claude 4.6** foi absolutamente brilhante para explicar conceitos obtusos e rascunhar as bases do que eu precisava.
* O **Gemini 3.1**, por outro lado, precisava de muito mais foco e restrições de detalhes para entender a arquitetura (e, por ser um sistema mais novo e com pouca documentação de orquestração interna, tentou mentir e alucinar diversas vezes dizendo que tinha feito algo que não fez).

Como construí dezenas de outros projetos na vida, orquestrar IAs não é um problema quando você sabe exatamente o que está fazendo. O problema é que, no Crompressor, eu estava entrando num pântano matemático que levaria uns 10 anos de estudo acadêmico tradicional antes da era da IA. Só o que eu li, reli e explorei sobre o [Teorema do Limite de Shannon](https://pt.wikipedia.org/wiki/Teorema_de_Shannon-Hartley) daria um livro.

Mas quer saber? Eu agradeço por cada "cara batida na parede" com a IA alucinando ou o código quebrando. Porque a cada batida, o projeto revelava uma nova direção. Quando a [busca linear falhava](https://github.com/MrJc01/crompressor-orquestrador/blob/main/crompressor-neuronio/pesquisa0/papers/papel0.md), eu era forçado a explorar a [Distância de Hamming](https://pt.wikipedia.org/wiki/Dist%C3%A2ncia_de_Hamming). Quando a memória RAM estourava, eu era forçado a abandonar a compressão tradicional e abraçar os Codebooks P2P. Explodimos todas as possibilidades e isso é só a ponta do iceberg.

---

## 3. Os Falsos Positivos e a Luta Matemática (A Distância de Hamming)

Na terceira semana, eu já lidava com dezenas de `goroutines` espalhadas no meu sistema operacional processando fluxos de texto (I/O streaming pesado em Golang). Tudo parecia perfeito. E então veio a facada.

Nós decidimos tentar comprimir textos longos via UTF-8, medindo a similaridade usando a matemática de `Hamming Distance` (medindo bits conflitantes). A hipótese parecia sólida. No entanto, em arquivos de texto, a alteração de **uma única vírgula ou um espaço** numa frase deslocava todo o offset da memória. Como os blocos do Crompressor na época eram de 128 bytes super-rígidos, o deslocamento destruía a similaridade do bloco inteiro. 
Nosso *Hit Rate* (taxa de acerto do Codebook) despencou para vergonhosos **0%**. 

Foram dias analisando esses falsos positivos. Descobrimos que o caminho era pivotar o motor para uma camada pura de O(1). Em vez de aceitar pedaços parcialmente parecidos de texto (que geram inflação no arquivo), nós migramos a lógica para a agressividade da rede. Nós adotamos o "match exato" (Deduplicação). A similaridade devia ser idêntica. Só que para compensar, subimos a "lente" do Chunk de 128 bytes para os blocos gigantes de 4096 bytes. Essa decisão de arquitetura salvou o projeto.

---

## 4. O Caminho é Longo e Isso Não é Nem o Começo

Testamos exaustivamente todas as matrizes, todas as métricas de compressão cruzada com Zstandard, todos os fluxos de rede UDP/TCP do repositório. O repositório [github.com/MrJc01/crompressor](https://github.com/MrJc01/crompressor) hoje contém a versão estável, capaz de estilhaçar gigabytes na rede usando o conceito de Cérebro Distribuído.

Mas, em retrospecto e olhando tudo que conquistei em apenas 30 dias com o Crompressor, o mais fascinante é a sensação de que essa ferramenta ainda não faz nem 1% do que planejo explorar com ela.

Neste exato instante, meu laboratório está com repositórios satélites sendo lapidados (como visto na Parte 1, sobre as *Simulações do Sistema Solar* e da compressão do Algoritmo *A\* de caminhos em Mapas Gigantes*, onde reduzimos GBs em MBs usando IDs da Chunk Table). Nós estamos no limiar de quebrar paradigmas em P2P seguro, cache de estado de GPUs em nuvem e fragmentação de Vídeos.

A jornada do desenvolvimento é dolorosa, cansativa e crua. Contudo, quando o Terminal apita e acusa uma rede de 500 megas sincrônica baixando para míseros centenas de kilobytes provando sua teoria, você entende o porquê escolhemos essa profissão.

Na última parte desta série **(Parte 3)**, não vou continuar discursando. Vou convocar vocês e pedir a ajuda da comunidade para consolidarmos esses resultados abertamente na comunidade científica. Até lá.
