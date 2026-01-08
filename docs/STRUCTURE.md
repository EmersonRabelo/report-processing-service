# Estrutura e Referências do Projeto

Este documento reúne a estrutura do repositório, descreve contratos/DTOs/entidades, mostra o modelo de dados (migrations) e explica a arquitetura de filas usada por esta aplicação.

Importante: a aplicação principal que publica mensagens de `report` para este serviço é **First API Go** — repositório: https://github.com/EmersonRabelo/first-api-go

---

## Mapa rápido de pastas

```
report-processing-service/
├── cmd/
│   └── report-processing-service/
│       └── main.go
├── db/
│   └── migrations/
├── internal/
│   ├── api/
│   │   └── perspective/
│   ├── config/
│   ├── database/
│   ├── dto/
   │   └── report/
│   ├── entity/
│   ├── handler/
│   ├── queue/
│   ├── repository/
│   └── service/
├── router/
├── go.mod
└── README.md
```

---

## Contratos, DTOs e Entidades

Onde procurar:

- Contratos / mensagens (DTOs): [internal/dto/report/contracts](internal/dto/report/contracts)
  - `create_report_message.go` — formato da mensagem enviada pela aplicação principal (producer).
  - `report_analysis_result_message.go` — formato da mensagem publicada por este serviço após análise.
- DTOs específicos: `internal/dto/report/perspective_update.go` (se aplicável para atualização de campos de análise).
- Entidade de domínio: `internal/entity/report.go` — modelo GORM que representa a tabela `reports`.

O que cada camada representa:

- Contratos/DTOs: estruturas usadas para validação e serialização das mensagens trocadas entre serviços (AMQP) e para entrada/saída HTTP.
- Entidades: estruturas que mapeiam para tabelas do banco (GORM) e contém comportamento de modelagem do domínio.

Exemplo (campos importantes na mensagem `CreateReportMessage`):

- `Id` (UUID) — id do report
- `PostId` (UUID) — id do post denunciado
- `ReporterId` (UUID) — id do usuário que reportou
- `Body` (string) — texto analisado
- `CreatedAt` (timestamp)

Exemplo (mensagem de resultado `ReportAnalysisResultMessage`):

- `ReportId` (UUID)
- Scores: `Toxicity`, `SevereToxicity`, `IdentityAttack`, `Insult`, `Profanity`, `Threat` (float)
- `Language` (string)
- `AnalyzedAt` (timestamp)

---

## Banco de Dados

Tecnologia usada:

- PostgreSQL como banco relacional.
- GORM como ORM (configurado em `internal/config/database.go`).

Versionamento de schema (migrations):

- Os scripts de migrations ficam em `db/migrations/` e são aplicados manualmente ou via ferramentas de migrations (ex.: golang-migrate).
- Migrations presentes (exemplos):
  - `000001_ create_reports_table_and_indexes.up.sql`
  - `000002_alter_reports_table_perspective_identity_hate_column_name.up.sql`
  - `000003_alter_reports_table_add_new_column_perspective_severe_toxicity.up.sql`

Como aplicar manualmente (exemplo):

```bash
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f db/migrations/000001_\ create_reports_table_and_indexes.up.sql
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f db/migrations/000002_alter_reports_table_perspective_identity_hate_column_name.up.sql
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f db/migrations/000003_alter_reports_table_add_new_column_perspective_severe_toxicity.up.sql
```

Principais tabelas (apresentação):

### `reports`

Definição principal (extraída de `db/migrations/000001_ create_reports_table_and_indexes.up.sql`):

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

Índices e comentários importantes também são criados na migration (ex.: índices em `post_id`, `reporter_id`, `status`, e índice único em `report_id` com filtro `WHERE deleted_at IS NULL`).

Observações:

- O PostgreSQL é a fonte da verdade. Redis (quando usado) é cache/contadores apenas.
- GORM é usado para acesso/abstração no código (`internal/config/database.go` e repositórios em `internal/repository`).

---

## Arquitetura de Filas (RabbitMQ)

Resumo e responsabilidades:

- Broker: RabbitMQ (conexão e canal inicializados em `internal/config/broker.go`).
- Padrão: Topic Exchange para roteamento flexível (`topic_report` é o exemplo adotado).
- Producers: a aplicação principal (First API Go) publica eventos do tipo `post.report.created` quando um usuário denuncia um post.
- Consumers: este serviço (report-processing-service) consome mensagens de filas configuradas (ex.: `q.report.response`) e processa análise com a Perspective API.

Exemplo de configuração/roteamento usado neste projeto:

- Exchange: `topic_report`
- Routing key (producer): `post.report.created`
- Routing key (consumer): `post.report.response` (ou binding correspondente)
- Queue exemplo: `q.report.response`

Fluxo simplificado:

1. Usuário reporta um post na aplicação principal (First API Go) → a aplicação publica `CreateReportMessage` em RabbitMQ (exchange `topic_report`, routing key `post.report.created`).
2. RabbitMQ roteia a mensagem para a fila apropriada (`q.report.response`).
3. `ReportConsumer` (nesta aplicação) consome a mensagem e chama o handler (`internal/handler/report_handler.go`).
4. O handler delega ao `ConsumerReportService` (`internal/service/consumer_report_service.go`) que:
   - persiste o report (idempotente)
   - chama a Perspective API para análise
   - atualiza o registro no Postgres com os scores
   - publica `ReportAnalysisResultMessage` se necessário (para outros consumidores)

Arquivo de configuração do broker: `internal/config/broker.go` (ex.: construção do URL AMQP usando `AMQP_HOST`, `AMQP_PORT`, `AMQP_USER`, `AMQP_PASSWORD`).

---

## Referências e arquivos importantes

- Código de consumo e serviço: `internal/service/consumer_report_service.go`
- Handlers de fila: `internal/handler/report_handler.go`
- Repositório: `internal/repository/report_repository.go`
- Migrations: `db/migrations/`
- Configurações: `internal/config/config.go`, `internal/config/database.go`, `internal/config/broker.go`
- Repositório que produz `CreateReportMessage`: https://github.com/EmersonRabelo/first-api-go

---

Se quiser, eu adiciono também uma versão navegável com links diretos para cada arquivo (`docs/STRUCTURE.md` já tem caminhos, mas posso gerar `docs/FILES.md` com links prontos para abrir em VS Code / GitHub). 
