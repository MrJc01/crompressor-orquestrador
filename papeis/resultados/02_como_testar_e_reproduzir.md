# 2. Como Testar e Reproduzir os Resultados

Quer validar por conta própria o potencial de deduplicação de borda do Crompressor? Este guia mostrará exatamente como executar o ciclo completo de treinamento e empacotamento, provando os testes que demonstramos nos relatórios de 99% de redução.

## Pré-requisitos
- Ter o Go (1.20+) instalado na sua máquina.
- Clonar o repositório principal:
  ```bash
  git clone https://github.com/MrJc01/crompressor.git
  cd crompressor
  go build -o crompressor_bin ./cmd/crompressor
  ```

## Passo 1: Preparando seus Dados
Você precisa de dados que representem atualizações incrementais (ex: duas versões de um banco de dados, repositório ou de uma ISO de sistema).
Crie um arquivo base e depois crie uma "nova versão" dele.
```bash
# Baixe ou crie um arquivo gigante (ex: logs de uma semana inteira)
cat /var/log/syslog* > base_logs.log

# Crie a "nova versão" simulando que o dia passou (dados altamente similares)
cat base_logs.log /var/log/auth.log > novos_logs.log
```

## Passo 2: Treinando o Cérebro (Codebook)
O Cérebro deve ser treinado com os **dados base**, e o ideal é usar o tamanho de bloco agressivo (4KB) para maximizar o ganho de rede.

```bash
./crompressor_bin train -i base_logs.log -o cerebro_logs.cromdb --chunk-size 4096 --size 16384
```
*O que ocorreu?* O motor analisou seu arquivo de base, extraiu até 16.384 padrões fundamentais de 4KB e gerou o dicionário descentralizado `cerebro_logs.cromdb`.

## Passo 3: Empacotamento / Sincronização
Agora, imagine que os `novos_logs.log` nasceram hoje e você quer mandá-los pela rede para o outro servidor que já possui o `cerebro_logs.cromdb`. 

Vamos compactar a nova versão referenciando o Cérebro:
```bash
./crompressor_bin pack -i novos_logs.log -o payload.crom -c cerebro_logs.cromdb --chunk-size 4096 --mode edge
```

## Passo 4: Validando a Mágica
Analise o tamanho dos arquivos e choque-se com o resultado matemático:
```bash
ls -lh novos_logs.log payload.crom
```
O arquivo `payload.crom` será, matematicamente, na melhor das hipóteses, mais de **99% menor** que o arquivo original `novos_logs.log`. 

Você acaba de provar como funciona uma **CDN Descentralizada com Deduplicação Extrema**.
