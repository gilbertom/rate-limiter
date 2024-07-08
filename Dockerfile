# Use uma imagem base apropriada
FROM golang:1.22.2

# Crie um diretório de trabalho
WORKDIR /app

# Copie os arquivos go.mod e go.sum e baixe as dependências
COPY go.mod go.sum ./
RUN go mod download

# Copie o restante do código fonte
COPY . .

# Copie o arquivo .env para o diretório /cmd/rate-limiter
COPY .env /cmd/rate-limiter/.env

# Compilar a aplicação
RUN go build -o /app/rate-limiter ./cmd/rate-limiter/main.go

# Expor a porta que a aplicação irá rodar
EXPOSE 8080

CMD ["/app/rate-limiter"]

