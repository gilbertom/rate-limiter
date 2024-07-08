```markdown
# Rate Limiter

Rate Limiter é uma aplicação desenvolvida em Go cujo objetivo é controlar a quantidade de requisições realizadas para um serviço web.

A quantidade de requisições (por segundo) pode ser limitada tanto por IP quanto Token de Acesso. O Token de acesso é identificado através do Header API_KEY.

Lista de variáveis para configuração do Rate Limiter:

SERVER_PORT (default 8080)
Porta da aplicação web que receberá as requisições do tipo GET.

REDIS_ADDR (default localhost)
Endereço do Redis. Neste projeto o Redis é utilizado num container Docker

REDIS_PORT (default 6379)
Porta de conexão com o Redis

MAX_REQUESTS_ALLOWED_BY_IP
Define a quantidade máxima de requisições por segundo que um determinado IP pode enviar.

MAX_REQUESTS_ALLOWED_BY_TOKEN
Define a quantidade máxima de requisições por segundo que um determinado Token pode enviar.

LIMIT_REQUESTS_BY_IP

ON_REQUEST_EXCEEDED_BLOCK_BY=IP
TIME_TO_RELEASE_REQUESTS_IP=10
TIME_TO_RELEASE_REQUESTS_TOKEN=20


![Build Status](https://img.shields.io/travis/username/GoApp.svg)
![Coverage Status](https://img.shields.io/coveralls/github/username/GoApp.svg)
![Go Version](https://img.shields.io/badge/go-%3E%3D1.15-blue.svg)

## Índice

- [Instalação](#instalação)
- [Uso](#uso)
- [Contribuição](#contribuição)
- [Licença](#licença)

## Instalação

Para instalar o GoApp, siga os passos abaixo:

```sh
git clone https://github.com/username/GoApp.git
cd GoApp
go build
go test ./...