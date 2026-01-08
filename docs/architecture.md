# Arquitetura

Visão de componentes e responsabilidades:

- Ingestão (Queue Consumer): componente que consome mensagens AMQP e aciona o `ConsumerReportService`.
- Serviço de Domínio (`service`): contém regras de validação, persistência mínima e orquestra chamadas externas.
- Repositório (`repository`): abstração de persistência usando GORM e tabela `reports`.
- Integração externa: cliente para a API Perspective (ver `internal/api/perspective`).
- Integração assíncrona de saída: produtor que publica o `ReportAnalysisResultMessage`.


Fluxo (fluxograma ASCII):

```
┌─────────────────────────────────────────────────────┐
│                   HTTP API (Gin)                    │
│              /api/v1/{resource}                     │
└────────────────────┬────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────┐
│              Controllers/Handlers                   │
│  (Request parsing, validation, response formatting) │
└────────────────────┬────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────┐
│                  Services                           │
│    (Business logic, domain rules, orchestration)    │
└────────────────────┬────────────────────────────────┘
                     │
         ┌───────────┴───────────┬──────────────┐
         │                       │              │
┌────────▼────────┐  ┌───────────▼────┐   ┌────▼─────────┐
│  Repositories   │  │  Queue Service │   │  Redis Cache │
│  (Data Access)  │  │  (Messaging)   │   │  (Counters)  │
└────────┬────────┘  └───────────┬────┘   └────┬─────────┘
         │                       │              │
         └───────────┬───────────┴──────────────┘
                     │
         ┌───────────▼──────────────┐
         │     PostgreSQL (GORM)    │
         │     + RabbitMQ + Redis   │
         └──────────────────────────┘
```

Explicação curta:

- `Controllers/Handlers` convertem requisições HTTP e mensagens de fila em chamadas de serviço.
- `Services` orquestram validações, persistência via `Repositories` e chamadas externas (Perspective API).
- `Queue Service` representa abstração sobre RabbitMQ (consumo e publicação de mensagens).
- `Repositories` garantem acesso ao Postgres e idempotência para evitar duplicatas (`InsertIfNotExists`).
- `Redis Cache` é opcional e pode ser usado para contadores, locks ou cache leve.

Decisões importantes:

- Persistência idempotente: `InsertIfNotExists` evita duplicatas se a mesma mensagem for reenviada.
- Comunicação assíncrona entre serviços via RabbitMQ para desacoplamento.
- Uso de GORM para abstração do Postgres (veja `internal/config/database.go`).

Escalabilidade e observabilidade:

- O consumidor pode ser escalado horizontalmente (várias instâncias), desde que idempotência seja garantida pelo repositório.
- Logging está configurado no `config` com níveis diferentes para `production` e `local`.

---

## Arquitetura (trechos extraídos do README da aplicação principal)

### Padrões de Design Utilizados

- **Repository Pattern**: abstração do acesso a dados (`internal/repository`).
- **Service Layer**: lógica de negócio centralizada (`internal/service`).
- **Dependency Injection**: instanciação e injeção no `cmd/report-processing-service/main.go`.
- **DTO Pattern**: `internal/dto` para validação/transferência de dados.
- **Producer-Consumer Pattern**: processamento assíncrono com RabbitMQ (`internal/queue`).

### Arquitetura em camadas (específica deste projeto)


O fluxograma abaixo é adaptado do exemplo recebido e mostra os principais componentes desta aplicação:

```
┌─────────────────────────────────────────────────────┐
│                   HTTP API (Gin)                    │
│              /api/v1/{resource}                     │
└────────────────────┬────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────┐
│              Controllers / Handlers                 │
│  (Request parsing, validation, response formatting) │
└────────────────────┬────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────┐
│                  Services                           │
│    (Business logic, domain rules, orchestration)    │
│    — Ex.: `ConsumerReportService` (`internal/service`)
└────────────────────┬────────────────────────────────┘
                     │
         ┌───────────┴───────────┬──────────────┐
         │                       │              │
┌────────▼────────┐  ┌───────────▼────┐   ┌────▼─────────┐
│  Repositories   │  │  Queue Service │   │  Redis Cache │
│  (Data Access)  │  │  (Messaging)   │   │  (Counters)  │
│  — `internal/repository` │  │  — `internal/queue` │   │  — `internal/redis` │
└────────┬────────┘  └───────────┬────┘   └────┬─────────┘
         │                       │              │
         └───────────┬───────────┴──────────────┘
                     │
         ┌───────────▼──────────────┐
         │     PostgreSQL (GORM)    │
         │     + RabbitMQ + Redis   │
         └──────────────────────────┘
```

### Sistema de Fila (padrão usado neste projeto)

O projeto segue o padrão Topic Exchange para roteamento flexível:

- Exchange: `topic_report`
- Routing Keys: `post.report.created` (producer), `post.report.response` (consumer)
- Queue exemplo: `q.report.response`

Fluxo de fila (resumido):

1. `ReportController` publica mensagem de criação via `producer`.
2. Mensagem vai para `RabbitMQ` (exchange `topic_report`).
3. `ReportConsumer` (goroutine) consome da fila (`q.report.response`) e delega ao `Handler`.
4. `ConsumerReportService` faz análise (Perspective API), persiste via `Repositories` e publica resultado.

### Redis (cache)

- Uso típico: contadores em tempo real (likes/replies), locks ou cache leve (`internal/redis`).
- Chaves sugeridas: `post:{postId}:likes`, `post:{postId}:replies`.

### Responsabilidades por camada (mapeamento para arquivos)

- `cmd/report-processing-service/main.go` — inicialização, carregamento de configurações, criação de conexões (DB, RabbitMQ), injeção de dependências.
- `internal/config` — `config.go`, `database.go`, `broker.go` (carregamento de variáveis de ambiente e inicialização de infra).
- `controller/` / `handler/` — handlers HTTP e handlers de fila (ex.: `internal/handler/report_handler.go`).
- `internal/service` — regras de negócio e orquestração (ex.: `internal/service/consumer_report_service.go`).
- `internal/repository` — acesso a dados e garantias de idempotência (ex.: `InsertIfNotExists`).
- `internal/queue` — produtores e consumidores RabbitMQ (`producer.go`, `consumer.go`).
- `internal/api/perspective` — cliente HTTP para a Perspective API.

### Regras de decisão e fluxo de moderação (resumo)

- A aplicação aplica regras com base nos scores da Perspective API (ex.: `THREAT`, `IDENTITY_ATTACK`, `SEVERE_TOXICITY`, `TOXICITY`, `INSULT`) para decidir flags do `Post` (`visible`, `limited`, `hidden_pending_review`, `removed`).

---

Referências rápidas:

- Arquivos de configuração: `internal/config/config.go`, `internal/config/database.go`, `internal/config/broker.go`.
- Lógica de consumo: `internal/service/consumer_report_service.go`, `internal/queue/consumer/consumer.go`, `internal/handler/report_handler.go`.

---

## Estrutura de pastas (esquema do projeto)

Para facilitar navegação rápida no repositório, abaixo há um esquema simplificado da estrutura de pastas atual:

```
report-processing-service/
├── cmd/
│   └── report-processing-service/
│       └── main.go
├── db/
│   └── migrations/
│       ├── 000001_create_reports_table_and_indexes.up.sql
│       ├── 000002_alter_reports_table_perspective_identity_hate_column_name.up.sql
│       └── 000003_alter_reports_table_add_new_column_perspective_severe_toxicity.up.sql
├── internal/
│   ├── api/
│   │   └── perspective/
│   │       └── client.go
│   ├── config/
│   │   ├── broker.go
│   │   ├── config.go
│   │   └── database.go
│   ├── database/
│   │   └── migration.go
│   ├── dto/
│   │   └── report/
│   │       ├── perspective_update.go
│   │       └── contracts/
│   │           ├── create_report_message.go
│   │           └── report_analysis_result_message.go
│   ├── entity/
│   │   └── report.go
│   ├── handler/
│   │   └── report_handler.go
│   ├── queue/
│   │   ├── consumer/
│   │   │   └── consumer.go
│   │   └── producer/
│   │       └── producer.go
│   ├── repository/
│   │   └── report_repository.go
│   └── service/
│       └── consumer_report_service.go
├── router/
│   └── router.go
├── go.mod
└── README.md
```

Este esquema resume onde procurar cada responsabilidade discutida na seção anterior.

---

## Descrição dos pacotes

Abaixo segue uma breve descrição do propósito de cada pasta/`package` principal do projeto (onde aplicável referenciei o arquivo mais relevante):

- `cmd/report-processing-service/` : ponto de entrada da aplicação — carrega configuração, inicializa conexões (DB, RabbitMQ) e inicia servidor/consumidores (`cmd/report-processing-service/main.go`).
- `db/migrations/` : scripts SQL de versionamento do schema; aplicados manualmente ou via ferramenta de migrations.
- `internal/api/perspective` : cliente HTTP para a Google Perspective API (`client.go`).
- `internal/config` : carga de configuração e inicialização de infra (DB, broker) — veja `config.go`, `database.go`, `broker.go`.
- `internal/database` : código auxiliar para executar migrations/rotinas relacionadas ao banco (`migration.go`).
- `internal/dto` : Data Transfer Objects e contratos das mensagens (validação e formato das payloads).
- `internal/entity` : modelos de domínio (estruturas GORM que representam tabelas), ex.: `report.go`.
- `internal/handler` : handlers de fila/integração que processam mensagens recebidas (`report_handler.go`).
- `internal/queue` : abstração de produtor/consumidor RabbitMQ (`producer/producer.go`, `consumer/consumer.go`).
- `internal/repository` : implementação do acesso a dados (GORM) e garantias de idempotência (`report_repository.go`).
- `internal/service` : camada de domínio/serviços que contém regras de negócio e orquestração (ex.: `consumer_report_service.go`).
- `router` : configuração de rotas HTTP (Gin) e grupos de endpoints (`router.go`).
- `go.mod` : declaração do módulo Go e dependências.
- `README.md` (raiz) : índice e ponte para a documentação (atualizado para apontar para `docs/`).

Se desejar, posso também gerar uma tabela com links diretos para os arquivos mencionados (um `docs/STRUCTURE.md`) para navegação rápida.

