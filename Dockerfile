FROM golang:1.22.2

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY .env /cmd/rate-limiter/.env

RUN go build -o /app/rate-limiter ./cmd/rate-limiter/main.go

EXPOSE 8080

CMD ["/app/rate-limiter"]

