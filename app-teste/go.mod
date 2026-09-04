module ssh-pqc-demo

// go1.24+ é necessário porque é a partir dessa versão do toolchain que o
// x/crypto/ssh compila o suporte ao key exchange híbrido pós-quântico
// mlkem768x25519-sha256 (arquivo com "//go:build go1.24").
// Recomenda-se usar o Go mais atual possível (1.27.x no momento em que
// este projeto foi escrito) para ter também crypto/mldsa, crypto/tls com
// ML-DSA e demais novidades do stdlib, mesmo que o pacote ssh em si ainda
// não use ML-DSA (veja o README.md).
go 1.26.0

require golang.org/x/crypto v0.56.0

require golang.org/x/sys v0.47.0 // indirect
