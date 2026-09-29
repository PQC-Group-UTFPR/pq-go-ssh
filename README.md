# pq-go-ssh
Experimenting with PQC in SSH

# Demo - SSH client-server em Go  

Veja em `app-teste/` uma demonstração de uma conexão SSH em Go com opção de depuração do protocolo.

## Código-fonte do pacote SSH (`x-crypto/`)

O diretório `x-crypto/` contém uma cópia do módulo `golang.org/x/crypto` (v0.57.0),
reduzida ao pacote `golang.org/x/crypto/ssh` e às suas dependências. O demo usa essa cópia local por meio
de uma diretiva `replace` em `app-teste/go.mod`:

```
replace golang.org/x/crypto => ../x-crypto
```

Portanto, para pesquisar/implementar autenticação pós-quântica, modifique o código em
`x-crypto/ssh/` (ex.: `keys.go`, `certs.go`, `client_auth.go`, `server.go`) e execute
o demo normalmente. Com `SSH_DEBUG=1`, cliente e servidor mostram qual cópia do pacote
foi compilada. Veja `x-crypto/UPSTREAM.md` para detalhes da versão e de como atualizar.

Testes do pacote SSH:

```bash
cd x-crypto
go test ./ssh/...
```


## Disclaimer 
- Em `app-teste/` A função de criação do servidor e cliente SSH, a informação e exemplo de depuração SSH, bem como este README foram gerados inicialmente a partir de prompt ao modelo Sonnet 5 / Claude. A estrutura foi modificada posteriormente (//Descrever), código também modificado, e inserido em um repositório dockerizado almejando reprodutibilidade.
- Não está usando autenticação por chave pública. A senha (`demo`/`demo123`) e a host key (chave do servidor) Ed25519 é efêmera, gerada a cada execução. 


## Uso via Docker

Após instalar docker, use

`sudo docker build --no-cache -t pq-go-ssh .`

- Em um novo terminal, o servidor:
`sudo docker run   --name sshserver -e NAME=sshserver -e IPSERVER=172.17.0.2:22222 -t pq-go-ssh` 

- Em um novo terminal, o cliente:
`sudo docker run --name sshclient -e NAME=sshclient -e IPSERVER=172.17.0.2:22222 -t pq-go-ssh` 


## Instalação manual

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


