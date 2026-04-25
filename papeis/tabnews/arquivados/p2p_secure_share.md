Title: Pitch: P2P Secure Share File Web:  Este site não tem backend  - Sem Rastros 👻🔒 100% Privado - Privacidade Inquebrável 🔒🛡️ Zero-Knowledge - OpenSource · MrJ · TabNews

Description: Em um mundo onde a privacidade de dados é constantemente mediada por servidores de terceiros, o P2P Secure Share (crom-p2pfile) surge como uma solução definitiva. O projeto é um Web App p...

Source: https://www.tabnews.com.br/MrJ/p2p-secure-share-file-web-este-site-nao-tem-backend-sem-rastros-100-por-cento-privado-privacidade-inquebravel-zero-knowledge-opensource

---

[TabNewsRelevantes](https://www.tabnews.com.br/)
[Recentes](https://www.tabnews.com.br/recentes/pagina/1)
[MrJ](https://www.tabnews.com.br/MrJ)
[1 mês atrás1 mês atrás](https://www.tabnews.com.br/MrJ/p2p-secure-share-file-web-este-site-nao-tem-backend-sem-rastros-100-por-cento-privado-privacidade-inquebravel-zero-knowledge-opensource)

Em um mundo onde a privacidade de dados é constantemente mediada por servidores de terceiros, o P2P Secure Share (crom-p2pfile) surge como uma solução definitiva. O projeto é um Web App puramente estático que permite a transferência de arquivos diretamente entre dispositivos, garantindo que nenhum bit de informação seja armazenado em infraestruturas externas.
- Acesse agora: [https://crom.run/p2pfile](https://crom.run/p2pfile)
- Código-fonte: [https://github.com/MrJc01/crom-p2pfile](https://github.com/MrJc01/crom-p2pfile)
[https://crom.run/p2pfile](https://crom.run/p2pfile)
[https://github.com/MrJc01/crom-p2pfile](https://github.com/MrJc01/crom-p2pfile)

## 🏗️ A Engenharia do "Zero-Backend"
A espinha dorsal do projeto é o WebRTC (Web Real-Time Communication), uma tecnologia que permite a comunicação direta entre navegadores. Para orquestrar o encontro inicial entre os dispositivos (signaling), o sistema utiliza a biblioteca PeerJS, que atua exclusivamente como um broker de endereços, sem nunca tocar nos arquivos transferidos.
O diferencial arquitetural reside na separação estrita entre a interface de usuário (UI) e a thread do socket (WebRTC DataChannel), garantindo que a aplicação permaneça responsiva mesmo durante transferências de grandes volumes de dados.

## 🛡️ Camadas de Segurança e Confiança Zero
O projeto implementa uma abordagem de Zero-Knowledge através de várias camadas:
- Criptografia E2EE Nativa: Todo arquivo é criptografado utilizando a Web Crypto API com o algoritmo AES-GCM de 256 bits antes mesmo de sair do dispositivo de origem.
- Handshake Optical-OOB (Out-of-Band): A chave de criptografia é gerada localmente e embutida em um fragmento de URL (#) contido em um QR Code. Como fragmentos de URL não são enviados ao servidor pelo navegador, a chave permanece isolada fisicamente entre as duas telas durante o escaneamento.
- Autenticação por PIN: O sistema exige a validação de um código de segurança em ambos os dispositivos para destravar o túnel criptográfico, prevenindo ataques de interceptação (MITM).

```
#
```

## 🚀 Superando Limites: Chunking e Backpressure
Transferir arquivos grandes diretamente no navegador apresenta desafios de memória (Heap de RAM). O P2P Secure Share soluciona isso através de uma estratégia de Chunking:
- Fragmentação em Fatias: O sistema quebra binários em pedaços de 64KB usando FileReader.
- Controle de Fluxo (Backpressure): Cada "chunk" é enviado com um IV (Vetor de Inicialização) único e o receptor confirma o recebimento para evitar o estouro do buffer do WebRTC.
- Aglutinação Segura: No destino, os pedaços são decifrados individualmente e remontados em um Blob final para download.

```
FileReader
```


```
Blob
```

## 🎨 Stack Tecnológica e Filosofia Local-First
A interface, desenvolvida com React 18 e Tailwind CSS, oferece uma experiência fluida e focada no mobile.
- Vite: Para builds estáticos de alta performance.
- jsQR: Para detecção de QR Code em tempo real via câmera.
- Soberania Digital: O projeto utiliza apenas APIs nativas do navegador, eliminando dependências de nuvem para o processamento de dados.

1. Abra o site no seu PC.
2. Escaneie o QR Code com o seu celular.
3. Confirme o PIN de segurança.
4. Arraste e solte arquivos para transferir instantaneamente.
Fonte: [https://crom.run/blog/p2p-secure-share-file-web-este-site-nao-tem-backend-sem-rastros-100-privado-privacidade-inquebravel-zero-knowledge-opensource](https://crom.run/blog/p2p-secure-share-file-web-este-site-nao-tem-backend-sem-rastros-100-privado-privacidade-inquebravel-zero-knowledge-opensource)
[RodrigoMedeiros](https://www.tabnews.com.br/RodrigoMedeiros)
[1 mês atrás1 mês atrás](https://www.tabnews.com.br/RodrigoMedeiros/f4cbb775-ad12-4f0f-a89a-3b02e14084db)
Muito bom! realmente é bem privado. Será que é possível colocar uma opção pra remover a criptografia, para simplesmente transferir arquivos mais rápido?
[MrJ](https://www.tabnews.com.br/MrJ)
- AutorAutor
- 
[1 mês atrás1 mês atrás](https://www.tabnews.com.br/MrJ/53171c06-f05b-4a68-91a3-84a0e06a77e2)
Acredito que sim, mas ficaria publico para qualquer acesso as informações:
[https://github.com/MrJc01/crom-p2pfile](https://github.com/MrJc01/crom-p2pfile)
[Contato](https://www.tabnews.com.br/contato)
[FAQ](https://www.tabnews.com.br/faq)
[GitHub](https://github.com/filipedeschamps/tabnews.com.br)
[Museu](https://www.tabnews.com.br/museu)
[RSS](https://www.tabnews.com.br/recentes/rss)
[Sobre](https://www.tabnews.com.br/filipedeschamps/tentando-construir-um-pedaco-de-internet-mais-massa)
[Status](https://www.tabnews.com.br/status)
[Termos de Uso](https://www.tabnews.com.br/termos-de-uso)
[curso.dev](https://curso.dev)

