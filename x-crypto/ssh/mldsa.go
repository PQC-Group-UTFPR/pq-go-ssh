package ssh

import (
	"crypto"
	"crypto/mldsa"
	"fmt"
)

// Nomes SSH dos algoritmos ML-DSA puros (sem híbrido), para chave de host e de usuário.
// Não existe nome padronizado para ML-DSA puro no SSH (provisório)
const (
	KeyAlgoMLDSA44 = "ssh-mldsa44@pq-go-ssh" // nivel 2
	KeyAlgoMLDSA65 = "ssh-mldsa65@pq-go-ssh" // nivel 3
	KeyAlgoMLDSA87 = "ssh-mldsa87@pq-go-ssh" // nivel 5
)

var mldsaParams = map[string]mldsa.Parameters{
	KeyAlgoMLDSA44: mldsa.MLDSA44(),
	KeyAlgoMLDSA65: mldsa.MLDSA65(),
	KeyAlgoMLDSA87: mldsa.MLDSA87(),
}

func mldsaAlgoName(p mldsa.Parameters) string {
	for nome, ficha := range mldsaParams {
		if ficha == p {
			return nome
		}
	}
	return ""
}

// mldsaPublicKey embrulha uma chave pública ML-DSA da biblioteca padrão
type mldsaPublicKey struct {
	pk *mldsa.PublicKey
}

func (k *mldsaPublicKey) Type() string {
	return mldsaAlgoName(k.pk.Parameters())
}

func (k *mldsaPublicKey) Marshal() []byte {
	w := struct {
		Name     string
		KeyBytes []byte
	}{
		k.Type(),
		k.pk.Bytes(),
	}
	return Marshal(&w)
}

func parseMLDSA(in []byte, algo string) (out PublicKey, rest []byte, err error) {
	params, ok := mldsaParams[algo]
	if !ok {
		return nil, nil, fmt.Errorf("ssh: algoritmo ML-DSA desconhecido: %s", algo)
	}
	var w struct {
		KeyBytes []byte
		Rest     []byte `ssh:"rest"`
	}
	if err := Unmarshal(in, &w); err != nil {
		return nil, nil, err
	}
	if l := len(w.KeyBytes); l != params.PublicKeySize() {
		return nil, nil, fmt.Errorf("ssh: tamanho inválido %d para chave pública %s", l, params)
	}
	pk, err := mldsa.NewPublicKey(params, w.KeyBytes)
	if err != nil {
		return nil, nil, err
	}
	return &mldsaPublicKey{pk}, w.Rest, nil
}

func (k *mldsaPublicKey) Verify(data []byte, sig *Signature) error {
	if sig.Format != k.Type() {
		return fmt.Errorf("ssh: assinatura do tipo %s para chave do tipo %s", sig.Format, k.Type())
	}
	params := k.pk.Parameters()
	if l := len(sig.Blob); l != params.SignatureSize() {
		return fmt.Errorf("ssh: tamanho inválido %d para assinatura %s", l, params)
	}
	if err := mldsa.Verify(k.pk, data, sig.Blob, nil); err != nil {
		return fmt.Errorf("ssh: assinatura %s não confere", params)
	}
	return nil
}

func (k *mldsaPublicKey) CryptoPublicKey() crypto.PublicKey {
	return k.pk
}
