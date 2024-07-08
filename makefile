# Makefile

REDIS_SERVICE = redis

REDIS_CLI = docker-compose exec $(REDIS_SERVICE) redis-cli

.PHONY: flushdb

flushdb:
	$(REDIS_CLI) FLUSHDB


.PHONY: listkeys

listkeys:
	$(REDIS_CLI) KEYS '*'

1:
	go run cmd/rate-limiter/main.go


cleanDocker:
	docker-compose down --volumes --rmi all
	docker system prune -a


build:
	docker-compose up --build


up:
	docker-compose up


bash:
	docker exec -it rate-limiter-app-1 sh

test:
	go test -v ./...

