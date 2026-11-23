package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"
)

type zeroReader struct{}

func (zeroReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

func main() {
	gbFlag := flag.Int64("gb", 0, "Quantidade de GB para transferir")
	flag.Parse()

	gb := *gbFlag
	if gb <= 0 {
		if envGB := os.Getenv("TRANSFER_GB"); envGB != "" {
			if parsed, err := strconv.ParseInt(envGB, 10, 64); err == nil {
				gb = parsed
			}
		}
	}
	if gb <= 0 {
		gb = 1
	}

	addr := "127.0.0.1:2222"
	ready := make(chan struct{})

	// 1. Sobe o Servidor SSH Pós-Quântico em background
	go runServer(addr, ready)
	<-ready

	// 2. Conecta o Cliente e transfere os GBs
	runClient(addr, gb)
}

func runServer(addr string, ready chan struct{}) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("[SERVER] Erro chave: %v", err)
	}
	signer, err := ssh.NewSignerFromKey(priv)
	if err != nil {
		log.Fatalf("[SERVER] Erro signer: %v", err)
	}

	config := &ssh.ServerConfig{
		Config: ssh.Config{
			KeyExchanges:   []string{ssh.KeyExchangeMLKEM768X25519},
			RekeyThreshold: 1 << 30, // Obriga rekeying completo a cada 1 GB
		},
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			return nil, nil
		},
	}

	config.AddHostKey(signer)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("[SERVER] Erro listen: %v", err)
	}
	defer listener.Close()
	close(ready)

	nConn, err := listener.Accept()
	if err != nil {
		return
	}
	defer nConn.Close()

	_, chans, reqs, err := ssh.NewServerConn(nConn, config)
	if err != nil {
		log.Fatalf("[SERVER] Erro handshake: %v", err)
	}
	go ssh.DiscardRequests(reqs)

	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unsupported")
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			continue
		}
		go func(ch ssh.Channel, in <-chan *ssh.Request) {
			defer ch.Close()
			for req := range in {
				if req.WantReply {
					req.Reply(true, nil)
				}
				if req.Type == "exec" {
					// Descarta os bytes recebidos em alta velocidade (simula /dev/null)
					_, _ = io.Copy(io.Discard, ch)
					return
				}
			}
		}(channel, requests)
	}
}

func runClient(addr string, gb int64) {
	config := &ssh.ClientConfig{
		User: "bench",
		Auth: []ssh.AuthMethod{ssh.Password("bench")},
		Config: ssh.Config{
			KeyExchanges:   []string{ssh.KeyExchangeMLKEM768X25519},
			RekeyThreshold: 1 << 30, // Obriga rekeying completo a cada 1 GB
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	startHandshake := time.Now()
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatalf("[CLIENT] Falha ao conectar: %v", err)
	}
	defer client.Close()
	handshakeDuration := time.Since(startHandshake)

	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("[CLIENT] Falha na sessão: %v", err)
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatalf("[CLIENT] Falha no StdinPipe: %v", err)
	}

	if err := session.Start("benchmark-sink"); err != nil {
		log.Fatalf("[CLIENT] Falha ao iniciar comando: %v", err)
	}

	totalBytes := gb * 1024 * 1024 * 1024
	fmt.Printf("\n==================================================\n")
	fmt.Printf("[BENCHMARK] KEX: %s\n", ssh.KeyExchangeMLKEM768X25519)
	fmt.Printf("[BENCHMARK] Handshake Inicial: %v\n", handshakeDuration)
	fmt.Printf("[BENCHMARK] Transferindo %d GB (%d bytes)...\n", gb, totalBytes)
	fmt.Printf("==================================================\n")

	startTransfer := time.Now()
	written, err := io.CopyN(stdin, zeroReader{}, totalBytes)
	if err != nil {
		log.Fatalf("[CLIENT] Erro durante transferência: %v", err)
	}
	_ = stdin.Close()
	_ = session.Wait()
	transferDuration := time.Since(startTransfer)

	throughputMBs := (float64(written) / (1024 * 1024)) / transferDuration.Seconds()

	fmt.Printf("==================================================\n")
	fmt.Printf("[RESULTADO] Volume: %d GB\n", gb)
	fmt.Printf("[RESULTADO] Tempo Total Transferência: %v\n", transferDuration)
	fmt.Printf("[RESULTADO] Vazão Média (Throughput): %.2f MB/s\n", throughputMBs)
	fmt.Printf("==================================================\n\n")
}