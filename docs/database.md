# Banco de Dados

Este projeto persiste informações de `reports` em PostgreSQL. As migrations estão no diretório `db/migrations`.

Arquivos de migrations existentes:

- `000001_ create_reports_table_and_indexes.up.sql` — criação da tabela inicial e índices.
- `000002_alter_reports_table_perspective_identity_hate_column_name.up.sql` — ajuste de coluna.
- `000003_alter_reports_table_add_new_column_perspective_severe_toxicity.up.sql` — adiciona coluna `perspective_severe_toxicity`.

Como aplicar migrations manualmente:

1. Crie o banco (se necessário):

```bash
createdb -h $DB_HOST -U $DB_USER $DB_NAME
```

2. Aplique os scripts SQL com `psql`:

```bash
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f db/migrations/000001_\ create_reports_table_and_indexes.up.sql
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f db/migrations/000002_alter_reports_table_perspective_identity_hate_column_name.up.sql
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f db/migrations/000003_alter_reports_table_add_new_column_perspective_severe_toxicity.up.sql
```

Automação:

- Se preferir uma ferramenta de migrations (por exemplo `golang-migrate`), adapte os scripts conforme necessário e aplique via `migrate`.

Observações sobre a modelagem:

- A tabela `reports` contém colunas para armazenar os scores retornados pela API Perspective para diversos atributos (toxicity, insult, threat, etc.) e metadados do processo (status, created_at, updated_at, perspective_response_at).
