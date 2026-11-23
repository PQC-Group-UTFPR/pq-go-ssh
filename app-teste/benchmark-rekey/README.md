# Benchmark de Rekeying Pós-Quântico (`mlkem768x25519-sha256`) — Issue #5

Ambiente isolado em **Docker** para simular transferências contínuas de grandes volumes via SSH pós-quântico (`golang.org/x/crypto/ssh`) e medir o impacto da renovação de chaves (*rekeying*) antes da implementação do mecanismo de **PSK e Session Tickets**.

## Arquivos

* `bench_ssh.go`: Cliente e servidor SSH em *loopback* (`127.0.0.1:2222`) realizando transferência em memória (`zeroReader` para `io.Discard`) com política de rotação configurável (`RekeyThreshold`).
* `Dockerfile`: Empacota o ambiente Go com o `x-crypto` local do repositório e recebe o volume da transferência via variável `TRANSFER_GB`.
* `rodar_teste.ps1` e `carregar_bib.ps1`: Atalhos automatizados para PowerShell (Windows).
* `carregar_bib.sh`: Atalhos automatizados equivalentes para Bash (Linux/macOS).
* `resultado_*.txt`: Logs das execuções coletadas na linha de base (Fase 1).

## Como Executar

### 1. Build da imagem Docker
A partir da **raiz do repositório (`pq-go-ssh`)**, execute:

    docker build -t pq-ssh-bench -f app-teste/benchmark-rekey/Dockerfile .

### 2. Execução direta via Docker
Passe a quantidade desejada de gigabytes pela variável `TRANSFER_GB`:

    docker run --rm -e TRANSFER_GB=1 pq-ssh-bench
    docker run --rm -e TRANSFER_GB=5 pq-ssh-bench
    docker run --rm -e TRANSFER_GB=10 pq-ssh-bench

### 3. Atalhos automatizados (Windows - PowerShell)
Entre na pasta `app-teste/benchmark-rekey`, libere a execução na sessão atual e carregue os atalhos:

    Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass
    . .\carregar_bib.ps1

### 4. Atalhos automatizados (Linux / macOS - Bash)
Entre na pasta `app-teste/benchmark-rekey` e carregue os atalhos:

    source ./carregar_bib.sh

Comandos disponíveis em ambos os sistemas:
* `gb1` — transfere **1 GB** e salva em `resultado_1gb.txt`
* `gb5` — transfere **5 GB** e salva em `resultado_5gb.txt`
* `gb10` — transfere **10 GB** e salva em `resultado_10gb.txt`
* `gb128` — transfere **128 GB** e salva em `resultado_128gb.txt`
