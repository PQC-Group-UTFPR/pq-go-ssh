//launch_server.go
package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	runtimedebug "runtime/debug"
	"os/exec"

	"golang.org/x/crypto/ssh"
)

const (
	demoUser     = "demo"
	demoPassword = "demo123"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:2222", "endereço TCP para escutar")
	flag.Parse()

	debug := os.Getenv("SSH_DEBUG") 
	
	//gera um par de chaves de autenticação
	hostSigner, fingerprint, err := generateHostKey()
	if err != nil {
		log.Fatalf("[server] falha ao gerar chave do host: %v", err)
	}
	log.Printf("[server] host key: tipo=%s fingerprint=%s", hostSigner.PublicKey().Type(), fingerprint)

	config := &ssh.ServerConfig{
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
		//TODO: mover isso para uma função (similar ao do cliente)
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if c.User() == demoUser && string(pass) == demoPassword {
				return nil, nil
			}
			return nil, fmt.Errorf("credenciais inválidas para %q", c.User())
		},
	}
	config.AddHostKey(hostSigner)

	if debug == "1" {
		printSupportedAlgorithms("server")
		logSSHModule("server")
	}

	//Servidor ouvindo conexões:
	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("[server] falha ao escutar em %s: %v", *addr, err)
	}
	defer listener.Close()
	log.Printf("[server] escutando em %s (defina SSH_DEBUG=1 para ver mais artefatos)", *addr)

	for {
		nConn, err := listener.Accept()
		if err != nil {
			log.Printf("[server] erro ao aceitar conexão: %v", err)
			continue
		}
		go handleConn(nConn, config)
	}
}

//Tratamento da conexão vinda do cliente para criar uma sessão
func handleConn(nConn net.Conn, config *ssh.ServerConfig) {
	defer nConn.Close()

	sshConn, chans, reqs, err := ssh.NewServerConn(nConn, config)
	if err != nil {
		log.Printf("[server] handshake SSH falhou com %s: %v", nConn.RemoteAddr(), err)
		return
	}
	defer sshConn.Close()

	logCryptoArtifacts("server", sshConn.Conn)

	go ssh.DiscardRequests(reqs)

	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "tipo de canal não suportado")
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Printf("[server] falha ao aceitar canal: %v", err)
			continue
		}
		go handleSession(channel, requests)
	}
}

// exibe a mensagem e executa um comando do cliente no servidor.
func handleSession(channel ssh.Channel, requests <-chan *ssh.Request) {
	defer channel.Close()
	for req := range requests {
		switch req.Type {
		case "exec":
			var payload struct{ Command string }
			if err := ssh.Unmarshal(req.Payload, &payload); err != nil {
				if req.WantReply {
					req.Reply(false, nil)
				}
				continue
			}
			if req.WantReply {
				req.Reply(true, nil)
			}

			cmd := exec.Command("sh", "-c", payload.Command)
			cmd.Stdout = channel
			cmd.Stderr = channel.Stderr()
			_ = cmd.Run()

			exitStatus := struct{ Status uint32 }{0}
			if cmd.ProcessState != nil {
				exitStatus.Status = uint32(cmd.ProcessState.ExitCode())
			}
			channel.SendRequest("exit-status", false, ssh.Marshal(&exitStatus))
			return

		case "shell":
			if req.WantReply {
				req.Reply(true, nil)
			}
			fmt.Fprintln(channel, "Bem-vindo ao servidor SSH de demonstração.")
			return

		default:
			if req.WantReply {
				req.Reply(false, nil)
			}
		}
	}
}

//Gera um par de chaves ed25519.
//TODO: receber o nome do algoritmo como parâmetro da função
func generateHostKey() (ssh.Signer, string, error) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, "", err
	}
	signer, err := ssh.NewSignerFromKey(priv)
	if err != nil {
		return nil, "", err
	}
	return signer, ssh.FingerprintSHA256(signer.PublicKey()), nil
}

//Exibe algoritmos suportados
func printSupportedAlgorithms(role string) {
	algos := ssh.SupportedAlgorithms()
	log.Printf("[%s][debug] KeyExchanges suportados: %v", role, algos.KeyExchanges)
	log.Printf("[%s][debug] HostKeys suportados: %v", role, algos.HostKeys)
}

// logCryptoArtifacts registra qual o usuário que conectou e imprime os artefatos criptográficos da conexão 
func logCryptoArtifacts(role string, conn ssh.Conn) {
	log.Printf("[%s] usuário autenticado: %s", role, conn.User())
	log.Printf("[%s] session ID (H): %s", role, hex.EncodeToString(conn.SessionID()))
	
	ac, ok := conn.(ssh.AlgorithmsConnMetadata)
	if !ok {
		log.Printf("[%s] esta versão de x/crypto/ssh não expõe NegotiatedAlgorithms", role)
		return
	}
	na := ac.Algorithms()
	log.Printf("[%s] KEX negociado: %s", role, na.KeyExchange)
	log.Printf("[%s] algoritmo de host key negociado: %s", role, na.HostKey)
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
