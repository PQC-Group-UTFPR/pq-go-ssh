// launch_client.go 
package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	runtimedebug "runtime/debug"
	"time"

	"golang.org/x/crypto/ssh"
)

const (
	demoUser     = "demo"
	demoPassword = "demo123"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:2222", "endereço do servidor SSH")
	cmd := flag.String("cmd", "echo 'olá SSH'; uname -a", "Mensagem / Comando remoto a executar no servidor SSH")
	flag.Parse()

	debug := os.Getenv("SSH_DEBUG")

	config := &ssh.ClientConfig{
		User: demoUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(demoPassword),
		},
		Config: ssh.Config{
			KeyExchanges: []string{
				ssh.KeyExchangeMLKEM768X25519,
				ssh.KeyExchangeCurve25519,
			},
			Ciphers: []string{
				ssh.CipherChaCha20Poly1305,
				ssh.CipherAES256GCM,
			},
		},
		HostKeyCallback: logAndAcceptHostKey,
		Timeout:         5 * time.Second,
	}

	clientssh, err := ssh.Dial("tcp", *addr, config)
	if err != nil {
		log.Fatalf("[client] falha ao conectar: %v", err)
	}
	defer clientssh.Close()
	
	if debug  == "1" {
		printSupportedAlgorithms("client")
		logSSHModule("client")
		logCryptoArtifacts("client", clientssh.Conn)
	}
	//inicia a sessão ssh
	sessionssh, err := clientssh.NewSession()
	if err != nil {
		log.Fatalf("[client] falha ao abrir sessão: %v", err)
	}
	defer sessionssh.Close()

	var stdout, stderr bytes.Buffer
	sessionssh.Stdout = &stdout
	sessionssh.Stderr = &stderr

	//Exibe a mensagem e executa o comando no servidor.
	if err := sessionssh.Run(*cmd); err != nil {
		log.Printf("[client] comando terminou com erro: %v", err)
	}
	fmt.Print(stdout.String())
	if stderr.Len() > 0 {
		fmt.Print(stderr.String())
	}
}

//Callback da configuração da sessão SSH (aceita a chave; idealmente poderia comparar com o "Known hosts" do SSH)
func logAndAcceptHostKey(hostname string, remote net.Addr, key ssh.PublicKey) error {
	log.Printf("[client] host key  de %s ACEITA POR PADRÃO. Algoritmo=%s fingerprint=%s",
		remote, key.Type(), ssh.FingerprintSHA256(key))
	return nil
}

//Exibir algoritmos Suportados (cliente)
func printSupportedAlgorithms(role string) {
	algos := ssh.SupportedAlgorithms()
	log.Printf("[%s][debug] KeyExchanges suportados: %v", role, algos.KeyExchanges)
	log.Printf("[%s][debug] HostKeys suportados: %v", role, algos.HostKeys)
//	log.Printf("[%s][debug] Ciphers suportados: %v", role, algos.Ciphers)
//	log.Printf("[%s][debug] MACs suportados: %v", role, algos.MACs)
}

//Função para debug da parte de criptografia negociada. Não inclui a autenticação da parte cliente
func logCryptoArtifacts(role string, conn ssh.Conn) {
	log.Printf("[%s] session ID (H): %s", role, hex.EncodeToString(conn.SessionID()))

	ac, ok := conn.(ssh.AlgorithmsConnMetadata)
	if !ok {
		log.Printf("[%s] esta versão de x/crypto/ssh não expõe NegotiatedAlgorithms", role)
		return
	}
	na := ac.Algorithms()
	log.Printf("[%s] KEX negociado: %s", role, na.KeyExchange)
	log.Printf("[%s] Chave de autenticação do servidor: %s", role, na.HostKey)
	//log.Printf("[%s] cifra/MAC (leitura):  %s / %q", role, na.Read.Cipher, na.Read.MAC)
	//log.Printf("[%s] cifra/MAC (escrita):  %s / %q", role, na.Write.Cipher, na.Write.MAC)
}

// logSSHModule mostra qual cópia do pacote golang.org/x/crypto/ssh foi compilada.
// Com o "replace" do go.mod, deve apontar para a cópia local em ../x-crypto.
func logSSHModule(role string) {
	info, ok := runtimedebug.ReadBuildInfo()
	if !ok {
		return
	}
	for _, dep := range info.Deps {
		if dep.Path != "golang.org/x/crypto" {
			continue
		}
		if dep.Replace != nil {
			log.Printf("[%s][debug] golang.org/x/crypto %s substituído pela cópia local: %s", role, dep.Version, dep.Replace.Path)
		} else {
			log.Printf("[%s][debug] golang.org/x/crypto %s (upstream, sem cópia local)", role, dep.Version)
		}
	}
}
