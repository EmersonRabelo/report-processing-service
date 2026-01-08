# Project Guide — Report Processing Service

Documento unificado e padronizado que centraliza informação técnica sobre o projeto.

Sumário
- Visão Geral
- Quickstart
- Arquitetura
- Estrutura de Pastas
- Contratos, DTOs e Entidades
- Banco de Dados (migrations e principais tabelas)
- Arquitetura de Filas (RabbitMQ)
- Desenvolvimento e execução
- Deploy e recomendações de produção
- Referências

---

## Visão Geral

Este serviço consome mensagens de denúncia (`report`) publicadas pela aplicação principal (First API Go — https://github.com/EmersonRabelo/first-api-go), analisa o conteúdo via Google Perspective API, persiste resultados em PostgreSQL e publica mensagens de resultado para integração assíncrona.

Fluxo resumido:

1. Aplicação principal publica `CreateReportMessage` em RabbitMQ.
2. Este serviço consome a mensagem, persiste o `report` (idempotente).
3. Texto é analisado pela Perspective API.
4. Registro é atualizado com scores; `ReportAnalysisResultMessage` é publicado.

---

## Quickstart

Pré-requisitos:
- Go 1.20+
- PostgreSQL
- RabbitMQ

Executar localmente (exemplo mínimo):

```bash
cp .env.local.example .env.local # ajustar variáveis
go run ./cmd/report-processing-service
```

Ver seção [setup.md](setup.md) para instruções completas.

---

## Arquitetura

Veja [architecture.md](architecture.md) para diagrama e explicações. Principais pontos:
- Camadas: Controllers/Handlers → Services → Repositories/Queue/Cache → Postgres
- Comunicação assíncrona via RabbitMQ (Topic Exchange)
- Idempotência no repositório via `InsertIfNotExists` (GORM `OnConflict DoNothing`)

---

## Estrutura de Pastas

Veja também: [STRUCTURE.md](STRUCTURE.md)

Esquema resumido:

```
cmd/                      # Entrypoint
db/migrations/            # SQL migrations
internal/
  api/perspective/        # cliente da Perspective API
  config/                 # carregamento de env, init DB/broker
  database/               # helpers de migrations
  dto/                    # contratos e DTOs
  entity/                 # modelos GORM
  handler/                # handlers de fila
  queue/                  # producer/consumer RabbitMQ
  repository/             # acesso a dados (GORM)
  service/                # lógica de negócio (ConsumerReportService)
router/                   # rotas HTTP (se houver)
```

---

## Contratos, DTOs e Entidades

Localização:
- `internal/dto/report/contracts` — mensagens AMQP (`CreateReportMessage`, `ReportAnalysisResultMessage`).
- `internal/dto/report` — estruturas da resposta da Perspective API.
- `internal/entity` — modelos GORM (`Report`).

Resumo das mensagens:
- `CreateReportMessage` (producer — First API Go): `id`, `post_id`, `reporter_id`, `body`, `created_at`.
- `ReportAnalysisResultMessage` (producer deste serviço): `report_id`, scores (`toxicity`, `insult`, `threat`, etc.), `language`, `analyzed_at`.

Entidade principal:
- `Report` (`internal/entity/report.go`): campos de identificação, status, timestamps, e colunas para armazenar scores da Perspective API.

---

## Banco de Dados

Tecnologia: PostgreSQL. ORM: GORM.

Versionamento: Migrations em `db/migrations/` (ex.: `000001_create_reports_table_and_indexes.up.sql`). O projeto também fornece helpers em `internal/database/migration.go` para aplicar migrations com `golang-migrate`.

Principais pontos:
- Tabela `reports` armazena registros de denúncias e scores de análise.
- Índices: `post_id`, `reporter_id`, `status`, índice único em `report_id` (com filtro para `deleted_at IS NULL`).

Exemplo de schema (trecho):

```sql
CREATE TABLE IF NOT EXISTS reports (
  report_id uuid PRIMARY KEY NOT NULL,
  post_id uuid NOT NULL,
  reporter_id uuid NOT NULL,
  status VARCHAR(30) NOT NULL DEFAULT 'pending',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  perspective_toxicity DOUBLE PRECISION NULL,
  perspective_insult DOUBLE PRECISION NULL,
  perspective_profanity DOUBLE PRECISION NULL,
  perspective_threat DOUBLE PRECISION NULL,
  perspective_identity_hate DOUBLE PRECISION NULL,
  perspective_language VARCHAR(50) NULL,
  perspective_response_at TIMESTAMP NULL
);
```

Como aplicar migrations manualmente, veja seção de Quickstart e `db/migrations`.

---

## Arquitetura de Filas (RabbitMQ)

Broker: RabbitMQ. Início de conexões em `internal/config/broker.go`.

Padrão usado: Topic Exchange — permite roteamento por `routing key`.

Configuração típica do projeto:
- Exchange: `topic_report`
- Producer (First API Go): `post.report.created`
- Consumer (este serviço): binding para `post.report.response` e fila `q.report.response`.

Fluxo detalhado:
1. First API Go publica `CreateReportMessage` no exchange.
2. RabbitMQ roteia e persiste a mensagem na fila.
3. ReportConsumer consome mensagens e chama o handler (`internal/handler/report_handler.go`).
4. Handler delega à camada de serviço (`internal/service/consumer_report_service.go`) que processa, persiste e publica resultados.

Boas práticas presentes no código:
- Ack/Nack manuseados no consumidor — erros permanentes são acked para evitar retries infinitos.
- Mensagens publicadas com `DeliveryMode: Persistent`.

---

## Desenvolvimento e execução

Run local:

```bash
go run ./cmd/report-processing-service
```

Testes: não há testes automatizados incluídos; recomenda-se adicionar `*_test.go` por pacote e rodar `go test ./...`.

Referências de implementação:
- Service principal: `internal/service/consumer_report_service.go`
- Producer: `internal/queue/producer/producer.go`
- Consumer: `internal/queue/consumer/consumer.go`
- Handler: `internal/handler/report_handler.go`

---

## Deploy e recomendações de produção

- Use secrets managers para `DB_PASSWORD`, `AMQP_PASSWORD`, `GOOGLE_PERSPECTIVE_API_TOKEN`.
- Configure liveness/readiness se for deploy em Kubernetes.
- Garanta políticas de retry/backoff em consumidores externos se necessário.

Exemplo de Dockerfile sugerido está em `docs/deployment.md`.

---

## Referências
- Projeto que produz os `CreateReportMessage`: https://github.com/EmersonRabelo/first-api-go
- Arquivos chave neste repositório: `internal/config/*`, `internal/service/*`, `internal/queue/*`, `db/migrations/*`.
