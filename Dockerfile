# syntax=docker/dockerfile:1
FROM golang:1.24.5-alpine3.22 AS builder

WORKDIR /app

# Copia os arquivos de dependências primeiro para cache eficiente
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copia o restante do código
COPY . .

# Compila o binário
RUN go build -o invista-ai-cli ./main.go

# Imagem final, menor, apenas com o binário
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/invista-ai-cli .
COPY --from=builder /app/data ./data
COPY --from=builder /app/public ./public
COPY --from=builder /app/README.md ./README.md

# Comando padrão ao rodar o container
ENTRYPOINT ["/app/invista-ai-cli"] 