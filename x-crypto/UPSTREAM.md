# Cópia local de golang.org/x/crypto

Este diretório contém uma cópia **completa e inicialmente não modificada** do módulo
[`golang.org/x/crypto`](https://go.googlesource.com/crypto) (espelho no GitHub:
<https://github.com/golang/crypto>).

- Versão upstream: **v0.57.0**
- Origem: module zip oficial (`go mod download golang.org/x/crypto@v0.57.0`)
- Licença: BSD-3-Clause (ver `LICENSE` e `PATENTS` neste diretório)

O pacote SSH está em `ssh/`. Os demais pacotes foram mantidos porque `ssh/` depende
de pacotes internos do módulo (ex.: `internal/poly1305`), que o Go só permite
importar de dentro do próprio módulo.

## Como é usado

`app-teste/go.mod` contém:

```
replace golang.org/x/crypto => ../x-crypto
```

Assim, os imports continuam sendo `golang.org/x/crypto/ssh`, mas o código compilado
é o deste diretório. Qualquer modificação feita aqui (ex.: novos algoritmos de
assinatura/autenticação pós-quânticos) é usada automaticamente pelo demo.

## Ver o que foi modificado em relação ao upstream

O commit que introduziu este diretório contém o código upstream sem alterações.
Para ver as modificações do grupo:

```bash
git log --oneline -- x-crypto/
git diff <commit-da-importação> -- x-crypto/
```

## Atualizar para uma nova versão upstream

```bash
cd app-teste
go mod download -json golang.org/x/crypto@vX.Y.Z   # mostra o diretório "Dir"
# copiar o conteúdo de "Dir" sobre ../x-crypto (preservando nossas mudanças,
# ex.: aplicando-as novamente com git diff/patch) e atualizar a versão neste arquivo
```
