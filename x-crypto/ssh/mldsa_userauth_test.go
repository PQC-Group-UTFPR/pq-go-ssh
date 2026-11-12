package ssh

import (
	"bytes"
	"crypto/ed25519"
	"crypto/mldsa"
	"crypto/rand"
	"fmt"
	"testing"
)

// loginMLDSA44 sobe um servidor e um cliente ssh ligados por netPipe e tenta logar
// como "joao" com clientKey. O servidor so autoriza authorizedKey e aceita os
// algoritmos de serverAlgos (nil = defaultPubKeyAuthAlgos). Retorna nil se o login
// deu certo.
func loginMLDSA44(t *testing.T, clientKey, authorizedKey *mldsa.PrivateKey, serverAlgos []string) error {
	t.Helper()

	c1, c2, err := netPipe()
	if err != nil {
		t.Fatalf("netPipe: %v", err)
	}
	defer c1.Close()
	defer c2.Close()

	autorizada, err := NewPublicKey(authorizedKey.PublicKey())
	if err != nil {
		t.Fatalf("NewPublicKey: %v", err)
	}

	serverConfig := &ServerConfig{
		PublicKeyAuthAlgorithms: serverAlgos,
		PublicKeyCallback: func(conn ConnMetadata, key PublicKey) (*Permissions, error) {
			if conn.User() == "joao" && bytes.Equal(key.Marshal(), autorizada.Marshal()) {
				return nil, nil
			}
			return nil, fmt.Errorf("chave de %q nao autorizada", conn.User())
		},
	}
	serverConfig.AddHostKey(testSigners["ed25519"])
	go NewServerConn(c1, serverConfig)

	signer, err := NewSignerFromSigner(clientKey)
	if err != nil {
		t.Fatalf("NewSignerFromSigner: %v", err)
	}
	clientConfig := &ClientConfig{
		User:            "joao",
		Auth:            []AuthMethod{PublicKeys(signer)},
		HostKeyCallback: FixedHostKey(testSigners["ed25519"].PublicKey()),
	}
	_, _, _, err = NewClientConn(c2, "", clientConfig)
	return err
}

// TestLoginUsuarioMLDSA44 testa o login de usuario com chave ML-DSA-44: a chave
// autorizada entra, outra chave ML-DSA-44 e recusada, e um servidor com a lista
// default recusa porque o ML-DSA-44 nao esta nela.
func TestLoginUsuarioMLDSA44(t *testing.T) {
	joao, _ := mldsa.GenerateKey(mldsa.MLDSA44())
	intrusa, _ := mldsa.GenerateKey(mldsa.MLDSA44())
	soMLDSA := []string{KeyAlgoMLDSA44}

	t.Run("chave autorizada", func(t *testing.T) {
		if err := loginMLDSA44(t, joao, joao, soMLDSA); err != nil {
			t.Fatalf("login falhou: %v", err)
		}
	})

	t.Run("chave nao autorizada", func(t *testing.T) {
		if err := loginMLDSA44(t, intrusa, joao, soMLDSA); err == nil {
			t.Fatal("login com chave nao autorizada foi aceito")
		}
	})

	t.Run("lista default", func(t *testing.T) {
		if err := loginMLDSA44(t, joao, joao, nil); err == nil {
			t.Fatal("servidor com a lista default aceitou ML-DSA-44")
		}
	})
}

// TestLoginUsuarioEd25519Recusado testa que um servidor que so aceita ML-DSA-44
// recusa uma chave Ed25519. O callback aceita qualquer chave, entao a recusa so
// pode vir da lista de algoritmos.
func TestLoginUsuarioEd25519Recusado(t *testing.T) {
	c1, c2, err := netPipe()
	if err != nil {
		t.Fatalf("netPipe: %v", err)
	}
	defer c1.Close()
	defer c2.Close()

	serverConfig := &ServerConfig{
		PublicKeyAuthAlgorithms: []string{KeyAlgoMLDSA44},
		PublicKeyCallback: func(conn ConnMetadata, key PublicKey) (*Permissions, error) {
			return nil, nil
		},
	}
	serverConfig.AddHostKey(testSigners["ed25519"])
	go NewServerConn(c1, serverConfig)

	_, chavePrivada, _ := ed25519.GenerateKey(rand.Reader)
	signer, err := NewSignerFromKey(chavePrivada)
	if err != nil {
		t.Fatalf("signer: %v", err)
	}

	clientConfig := &ClientConfig{
		User:            "joao",
		Auth:            []AuthMethod{PublicKeys(signer)},
		HostKeyCallback: FixedHostKey(testSigners["ed25519"].PublicKey()),
	}
	_, _, _, err = NewClientConn(c2, "", clientConfig)
	if err == nil {
		t.Fatal("servidor so-ML-DSA aceitou chave Ed25519")
	}
}
