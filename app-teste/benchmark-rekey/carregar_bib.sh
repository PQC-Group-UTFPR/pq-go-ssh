#!/bin/bash
rodar_teste() {
    local gb=${1:-1}
    echo "[DOCKER] Executando benchmark de ${gb} GB..."
    docker run --rm -e TRANSFER_GB="$gb" pq-ssh-bench | tee "resultado_${gb}gb.txt"
    echo "[OK] Resultado salvo em resultado_${gb}gb.txt"
}

gb1()   { rodar_teste 1; }
gb5()   { rodar_teste 5; }
gb10()  { rodar_teste 10; }
gb128() { rodar_teste 128; }
