# Instalação e execução local

Pré-requisitos:

- Go 1.20+ instalado
- PostgreSQL (versão compatível)
- RabbitMQ (ou broker AMQP compatível)

Passos básicos:

1. Clone o repositório:

```bash
git clone <repo-url>
cd report-processing-service
```

2. Configure variáveis de ambiente (exemplo `.env.local`):

Veja `configuration.md` para a lista completa. Um exemplo mínimo:

```env
APP_ENV=local
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=1234
DB_NAME=go-api
AMQP_HOST=localhost
AMQP_PORT=5672
AMQP_USER=guest
AMQP_PASSWORD=guest
GOOGLE_PERSPECTIVE_API_TOKEN=
```

3. Executar migrations (usar psql ou a ferramenta de sua escolha):

```bash
# Exemplo utilizando psql local
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f db/migrations/000001_\ create_reports_table_and_indexes.up.sql
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f db/migrations/000002_alter_reports_table_perspective_identity_hate_column_name.up.sql
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f db/migrations/000003_alter_reports_table_add_new_column_perspective_severe_toxicity.up.sql
```

4. Executar o serviço:

```bash
go run ./cmd/report-processing-service
# ou build
go build -o bin/report-processing-service ./cmd/report-processing-service
./bin/report-processing-service
```

Observação: a aplicação usa arquivos `.env.<env>` se presentes; por exemplo `.env.local`.
