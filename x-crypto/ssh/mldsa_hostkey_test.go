package ssh

import (
	"crypto/mldsa"
	"crypto/rand"
	"testing"
)

// novaHostKeyMLDSA44 gera uma chave ML-DSA-44 e devolve o signer SSH dela.
func novaHostKeyMLDSA44(t *testing.T) Signer {
	t.Helper()
	priv, err := mldsa.GenerateKey(mldsa.MLDSA44())
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	signer, err := NewSignerFromKey(priv)
	if err != nil {
		t.Fatalf("NewSignerFromKey: %v", err)
	}
	return signer
}

// TestMLDSA44AssinaturaIdaEVolta assina, serializa a chave pública, lê de volta
// e verifica. Também verifica se possíveis modificações indevidas são
// identificadas e recusadas.
func TestMLDSA44AssinaturaIdaEVolta(t *testing.T) {
	signer := novaHostKeyMLDSA44(t)
	dados := []byte("hash do handshake")

	sig, err := signer.Sign(rand.Reader, dados)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if sig.Format != KeyAlgoMLDSA44 {
		t.Fatalf("formato da assinatura = %q, esperado %q", sig.Format, KeyAlgoMLDSA44)
	}
	if len(sig.Blob) != mldsa.MLDSA44SignatureSize {
		t.Fatalf("tamanho da assinatura = %d, esperado %d", len(sig.Blob), mldsa.MLDSA44SignatureSize)
	}

	// A chave pública faz o caminho da rede: Marshal no servidor, Parse no cliente.
	pub, err := ParsePublicKey(signer.PublicKey().Marshal())
	if err != nil {
		t.Fatalf("ParsePublicKey: %v", err)
	}
	if pub.Type() != KeyAlgoMLDSA44 {
		t.Fatalf("tipo da chave = %q, esperado %q", pub.Type(), KeyAlgoMLDSA44)
	}

	if err := pub.Verify(dados, sig); err != nil {
		t.Fatalf("assinatura válida foi rejeitada: %v", err)
	}

	t.Run("mensagem alterada", func(t *testing.T) {
		if err := pub.Verify([]byte("outro hash"), sig); err == nil {
			t.Fatal("assinatura aceita para uma mensagem diferente")
		}
	})

	t.Run("assinatura alterada", func(t *testing.T) {
		adulterada := *sig // copia a assinatura
		adulterada.Blob = append([]byte(nil), sig.Blob...)
		adulterada.Blob[0] ^= 0xff // inverte os bits do primeiro byte
		if err := pub.Verify(dados, &adulterada); err == nil {
			t.Fatal("assinatura adulterada foi aceita")
		}
	})
}

// handshakeHostKey faz um handshake SSH completo em memória. O servidor usa hostKey;
// o cliente aceita os algoritmos de clientAlgos (nil = lista padrão) e só confia
// na chave esperada. Retorna nil se o handshake deu certo.
func handshakeHostKey(t *testing.T, hostKey Signer, esperada PublicKey, clientAlgos []string) error {
	t.Helper()

	c1, c2, err := netPipe()
	if err != nil {
		t.Fatalf("netPipe: %v", err)
	}
	defer c1.Close()
	defer c2.Close()

	serverConfig := &ServerConfig{NoClientAuth: true}
	serverConfig.AddHostKey(hostKey)
	go NewServerConn(c1, serverConfig)

	clientConfig := &ClientConfig{
		User:              "teste",
		HostKeyAlgorithms: clientAlgos,
		HostKeyCallback:   FixedHostKey(esperada),
	}
	_, _, _, err = NewClientConn(c2, "", clientConfig)
	return err
}

// TestHostKeyMLDSA44 testa a autenticação do servidor com ML-DSA-44 num handshake real.
func TestHostKeyMLDSA44(t *testing.T) {
	hostKey := novaHostKeyMLDSA44(t)
	soMLDSA := []string{KeyAlgoMLDSA44}

	t.Run("cliente aceita ML-DSA-44", func(t *testing.T) {
		if err := handshakeHostKey(t, hostKey, hostKey.PublicKey(), soMLDSA); err != nil {
			t.Fatalf("handshake falhou: %v", err)
		}
	})

	t.Run("cliente com lista padrão recusa", func(t *testing.T) {
		if err := handshakeHostKey(t, hostKey, hostKey.PublicKey(), nil); err == nil {
			t.Fatal("cliente com a lista padrão negociou ML-DSA-44")
		}
	})

	t.Run("chave de host diferente é recusada", func(t *testing.T) {
		outra := novaHostKeyMLDSA44(t)
		if err := handshakeHostKey(t, hostKey, outra.PublicKey(), soMLDSA); err == nil {
			t.Fatal("cliente aceitou uma chave de host diferente da esperada")
		}
	})
}
