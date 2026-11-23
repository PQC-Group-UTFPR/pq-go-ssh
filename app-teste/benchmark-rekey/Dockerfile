FROM golang:1.26-bookworm

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends tcpdump && rm -rf /var/lib/apt/lists/*

COPY app-teste/benchmark-rekey/go.mod app-teste/benchmark-rekey/go.sum ./
COPY x-crypto ./x-crypto
COPY app-teste/benchmark-rekey/bench_ssh.go ./

RUN go mod tidy && go build -o bench_ssh bench_ssh.go

ENV TRANSFER_GB="1"
CMD ["./bench_ssh"]
