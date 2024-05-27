# Makefile

# Nome do serviço Redis no Docker Compose
REDIS_SERVICE = redis

# Comando para executar o Redis CLI no contêiner Redis
REDIS_CLI = docker-compose exec $(REDIS_SERVICE) redis-cli

.PHONY: flushdb

flushdb:
	$(REDIS_CLI) FLUSHDB


.PHONY: listkeys

listkeys:
	$(REDIS_CLI) KEYS '*'

1:
	go run cmd/rate-limiter/main.go

