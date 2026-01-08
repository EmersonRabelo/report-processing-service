# Desenvolvimento

Como desenvolver e testar localmente.

Preparação do ambiente:

1. Garanta que as variáveis de ambiente estejam definidas (veja `configuration.md`).
2. Inicie serviços dependentes (Postgres, RabbitMQ).

Rodar o serviço em modo de desenvolvimento:

```bash
go run ./cmd/report-processing-service
```

Executar build:

```bash
go build -o bin/report-processing-service ./cmd/report-processing-service
```

Testes:

- Atualmente não há testes unitários no repositório. Para adicionar testes, crie pacotes `*_test.go` próximos às unidades a serem testadas e utilize `go test ./...`.

Depuração:

- Use `delve` (`dlv`) para debugar localmente ou logs por `fmt`/`log` conforme necessário.

Código relevante para leitura:

- `internal/service/consumer_report_service.go` — fluxo principal de criação e análise.
- `internal/repository/report_repository.go` — persistência e garantias de idempotência.
- `internal/queue/consumer` e `internal/queue/producer` — integração com RabbitMQ.
