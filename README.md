# pq-go-ssh
Experimenting with PQC in SSH

# Demo - SSH client-server em Go  

Veja em `app-teste/` uma demonstração de uma conexão SSH em Go com opção de depuração do protocolo.


## Disclaimer 
- Em `app-teste/` A função de criação do servidor e cliente SSH, a informação e exemplo de depuração SSH, bem como este README foram gerados inicialmente a partir de prompt ao modelo Sonnet 5 / Claude. A estrutura foi modificada posteriormente (//Descrever), código também modificado, e inserido em um repositório dockerizado almejando reprodutibilidade.
- Não está usando autenticação por chave pública. A senha (`demo`/`demo123`) e a host key (chave do servidor) Ed25519 é efêmera, gerada a cada execução. 

## Instalação

Requer linguagem Go instalada (TODO: dockerizar)

## Execução do exemplo app-teste

```bash
cd pq-go-ssh/app-teste/
go mod tidy 

# terminal 1
SSH_DEBUG=1 go run launch_server.go

# terminal 2
SSH_DEBUG=1 go run launch_client.go 
```


