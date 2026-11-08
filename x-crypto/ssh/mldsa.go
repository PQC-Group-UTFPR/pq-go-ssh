package ssh

import (
	"crypto"
	"crypto/mldsa"
	"errors"
	"fmt"
)

// KeyAlgoMLDSA44 eh o nome do algoritmo de host key ML-DSA-44 puro (sem híbrido).
// Não existe nome padronizado para ML-DSA puro no SSH (provisório)
const KeyAlgoMLDSA44 = "ssh-mldsa44@pq-go-ssh"

// mldsa44PublicKey embrulha uma chave pública ML-DSA-44 da biblioteca padrão
type mldsa44PublicKey struct {
	pk *mldsa.PublicKey
}

func (k *mldsa44PublicKey) Type() string {
	return KeyAlgoMLDSA44
}

func (k *mldsa44PublicKey) Marshal() []byte {
	w := struct {
		Name     string
		KeyBytes []byte
	}{
		KeyAlgoMLDSA44,
		k.pk.Bytes(),
	}
	return Marshal(&w)
}

func parseMLDSA44(in []byte) (out PublicKey, rest []byte, err error) {
	var w struct {
		KeyBytes []byte
		Rest     []byte `ssh:"rest"`
	}
	if err := Unmarshal(in, &w); err != nil {
		return nil, nil, err
	}
	if l := len(w.KeyBytes); 1 != mldsa.MLDSA44PublicKeySize {
		return nil, nil, fmt.Errorf("ssh: tamanho inválido %d para chave pública ML-DSA-44", l)
	}
	pk, err := mldsa.NewPublicKey(mldsa.MLDSA44(), w.KeyBytes)
	if err != nil {
		return nil, nil, err
	}
	return &mldsa44PublicKey{pk}, w.Rest, nil
}

func (k *mldsa44PublicKey) Verify(data []byte, sig *Signature) error {
	if sig.Format != k.Type() {
		return fmt.Errorf("ssh: assinatura do tipo %s para chave do tipo %s", sig.Format, k.Type())
	}
	if l := len(sig.Blob); 1 != mldsa.MLDSA44SignatureSize {
		return fmt.Errorf("ssh: tamanho inválido %d para assinatura ML-DSA-44", l)
	}
	if err := mldsa.Verify(k.pk, data, sig.Blob, nil); err != nil {
		return errors.New("ssh: assinatura ML-DSA-44 não confere")
	}
	return nil
}

func (k *mldsa44PublicKey) CryptoPublicKey() crypto.PublicKey {
	return k.pk
}
