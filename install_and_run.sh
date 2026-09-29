#!/bin/bash

#entradas para este script: DOCKERNAME (sshclient, sshserver), e IPSERVER (e.g., 172.17.0.2:22222)
DOCKERNAME=$1
IPSERVER=$2

sudo apt install -y wget unzip git pkg-config 

echo -e "\nInstaling Go 1.26...."
wget https://go.dev/dl/go1.26.0.linux-amd64.tar.gz -O /tmp/go1.26.0.tar.gz || { echo "Failed to download Go 1.26. Exiting" ; exit 1 ;}

sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf /tmp/go1.26.0.tar.gz
export PATH="/usr/local/go/bin:${PATH}"
go version || { echo "Failed to install Go 1.26. Exiting 1." ; exit 1 ;}

cd app-teste/
go mod tidy

if [ "${DOCKERNAME}" == "sshclient" ] ; then 
	SSH_DEBUG=1 go run launch_client.go -addr $IPSERVER
else 
	echo "				Executando Servidor SSH. Ctrl+C para finalizar" 
	SSH_DEBUG=1 go run launch_server.go -addr $IPSERVER
fi
