FROM golang:1.24.5-alpine3.22 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o invista-ai-cli ./main.go

FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/invista-ai-cli .
COPY --from=builder /app/data ./data
COPY --from=builder /app/public ./public
COPY --from=builder /app/README.md ./README.md

# Garantir que o executável tem permissões de execução
RUN chmod +x /app/invista-ai-cli

ENTRYPOINT ["/app/invista-ai-cli"] 

