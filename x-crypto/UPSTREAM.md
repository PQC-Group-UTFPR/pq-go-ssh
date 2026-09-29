# Cópia local de golang.org/x/crypto

Este diretório contém uma cópia **parcial e inicialmente não modificada** do módulo
[`golang.org/x/crypto`](https://go.googlesource.com/crypto) (espelho no GitHub:
<https://github.com/golang/crypto>).

- Versão upstream: **v0.57.0**
- Origem: module zip oficial (`go mod download golang.org/x/crypto@v0.57.0`)
- Licença: BSD-3-Clause (ver `LICENSE` e `PATENTS` neste diretório)

O pacote SSH está em `ssh/`. Foram mantidos apenas os pacotes de que ele depende
(verificado com `go list -deps -test ./ssh/...`):

| Pacote | Uso |
|---|---|
| `ssh/` (e subpacotes `agent`, `knownhosts`, `terminal`, `test`, `internal/bcrypt_pbkdf`) | o protocolo SSH |
| `chacha20/`, `internal/poly1305/` | cifra `chacha20-poly1305@openssh.com` |
| `curve25519/` | troca de chaves `curve25519-sha256` (e parte do híbrido `mlkem768x25519-sha256`) |
| `blowfish/` | `bcrypt_pbkdf`, para chaves privadas OpenSSH com senha |
| `cryptobyte/` | codificação de mensagens em `ssh/control.go` (proxy de controle) |
| `internal/alias/` | usado por `chacha20` |
| `sha3/`, `internal/testenv/` | somente para os testes de `ssh/test` |

Os demais pacotes do módulo (`acme`, `openpgp`, `nacl`, `argon2`, ...) foram removidos.
Os demais algoritmos (AES-GCM, ML-KEM, Ed25519, ECDSA, RSA, SHA-2) vêm da biblioteca
padrão do Go (`crypto/...`). Se algum pacote removido for necessário, copie-o da
mesma versão upstream. Os pacotes `internal/...` precisam ficar aqui porque o Go só
permite importá-los de dentro do próprio módulo.

## Como é usado

`app-teste/go.mod` contém:

```
replace golang.org/x/crypto => ../x-crypto
```

Assim, os imports continuam sendo `golang.org/x/crypto/ssh`, mas o código compilado
é o deste diretório. Qualquer modificação feita aqui (ex.: novos algoritmos de
assinatura/autenticação pós-quânticos) é usada automaticamente pelo demo.

## Ver o que foi modificado em relação ao upstream

O commit que introduziu este diretório contém o código upstream completo e sem alterações;
o commit seguinte apenas removeu os pacotes não usados pelo SSH.
Para ver as modificações do grupo:

```bash
git log --oneline -- x-crypto/
git diff <commit-da-importação> -- x-crypto/
```

## Atualizar para uma nova versão upstream

```bash
cd app-teste
go mod download -json golang.org/x/crypto@vX.Y.Z   # mostra o diretório "Dir"
# copiar de "Dir" apenas os diretórios listados acima (e go.mod/go.sum/LICENSE/PATENTS)
# sobre ../x-crypto, preservando nossas mudanças (ex.: aplicando-as novamente com
# git diff/patch), rodar `go mod tidy` aqui e em ../app-teste, e atualizar a versão
# neste arquivo. Confira se novas dependências surgiram com:
#   go list -deps -test ./ssh/... | grep golang.org/x/crypto
```
