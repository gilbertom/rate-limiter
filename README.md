# Rate Limiter

<p align="center">
  <img src="https://blog.golang.org/gopher/gopher.png" alt="">
</p>


Rate Limiter é uma aplicação desenvolvida em Go, cujo objetivo é controlar a quantidade de requisições realizadas para um serviço web.

A quantidade de requisições (por segundo) pode ser limitada tanto por IP quanto por Token de Acesso. O Token de acesso é identificado através do Header `API_KEY`.

## Variáveis de Configuração

Lista de variáveis para configuração do Rate Limiter encontrado no arquivo .env:

- **SERVER_PORT** (default: 8080)
  - Porta da aplicação web que receberá as requisições do tipo GET.
  
- **REDIS_ADDR** (default: localhost)
  - Endereço do Redis. Neste projeto, o Redis é utilizado num container Docker.
  
- **REDIS_PORT** (default: 6379)
  - Porta de conexão com o Redis.
  
- **MAX_REQUESTS_ALLOWED_BY_IP** (inteiro positivo)
  - Define a quantidade máxima de requisições por segundo que um determinado IP pode enviar.
  
- **MAX_REQUESTS_ALLOWED_BY_TOKEN** (inteiro positivo)
  - Define a quantidade máxima de requisições por segundo que um determinado Token pode enviar.
  
- **ON_REQUEST_EXCEEDED_BLOCK_BY** (IP|TOKEN)
  - Se definido como IP, o tempo para liberar o bloqueio será o tempo definido na variável `TIME_TO_RELEASE_REQUESTS_IP`. Se definido como TOKEN, o tempo de liberação do bloqueio será o definido na variável `TIME_TO_RELEASE_REQUESTS_TOKEN`.
  
- **TIME_TO_RELEASE_REQUESTS_IP** (inteiro positivo)
  - Quantidade de segundos que um determinado IP ficará bloqueado para novas requisições.
  
- **TIME_TO_RELEASE_REQUESTS_TOKEN** (inteiro positivo)
  - Quantidade de segundos que um determinado Token ficará bloqueado para novas requisições.
  
**Nota:** Caso as requisições sejam enviadas informando um Token, a quantidade de segundos definida na variável `MAX_REQUESTS_ALLOWED_BY_TOKEN` sobrepõe a quantidade de segundos definida na variável `MAX_REQUESTS_ALLOWED_BY_IP`.

## Índice

- [Instalação](#instalação)
- [Como Usar](#como-usar)
- [Execução dos Testes Unitários](#execução-dos-testes-unitários)
- [Contato](#contato)
- [Agradecimentos](#agradecimentos)

## Instalação

Para instalar o Rate Limiter, siga os passos abaixo:

```sh
git clone https://github.com/gilbertom/rate-limiter
docker-compose up --build
```

Nota: Ao executar o comando `docker-compose up --build`, os testes unitários são executados automaticamente e a aplicação é iniciada, ficando disponível na porta `8080`.


## Como Usar

Exemplo de uso:
```sh
curl -X GET 'http://localhost:8080' --header 'API_KEY: TOKEN 123'
```
Response

  Em caso de sucesso:
  ```json
  {
    "code": 200,
    "message": "Success"
  }
  ```

  Caso exceda o limite máximo de requisições:

  ```json
  {
    "code": 429,
    "message": "You have reached the maximum number of requests or actions allowed within a certain time frame"
  }
  ```

## Execução dos Testes Unitários
Para executar os testes unitários manualmente, acesse o container do serviço `app` e abra um shell interativo dentro do container:

```sh
docker-compose exec app sh
go test -v -count=1 ./...
```

## Contato
Para entrar em contato com o desenvolvedor deste projeto:
[gilbertomakiyama@gmail.com](mailto:gilbertomakiyama@gmail.com)

## Agradecimentos
Gostaria de expressar minha sincera gratidão a todo o time do curso de Pós-Graduação em Go Avançado da FullCycle pelo empenho, dedicação e excelência no ensino. Suas contribuições foram fundamentais para o meu desenvolvimento e sucesso. Muito obrigado!