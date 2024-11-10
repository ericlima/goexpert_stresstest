# Etapa de build
FROM golang:1.20 AS builder

WORKDIR /app

COPY . .

# Construção do binário com GOOS e GOARCH para compatibilidade
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

# Etapa final para o container leve
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/app .

# Garantir que o binário seja executável
RUN chmod +x ./app

ENTRYPOINT ["./app"]
