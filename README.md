# Report Processing Service

Breve resumo

Este repositório contém o serviço responsável por consumir mensagens de denúncias (`report`) publicadas pela aplicação principal (First API Go), analisar o conteúdo usando a Google Perspective API, persistir resultados em PostgreSQL e publicar mensagens de análise.

Documentação completa e organizada em `docs/`.

Índice (documentação)

- **Guia unificado:** [PROJECT_GUIDE.md](docs/PROJECT_GUIDE.md) — visão geral, arquitetura, dados, filas, quickstart e recomendações de produção.
- **Mapa do repositório:** [STRUCTURE.md](docs/STRUCTURE.md) — esquema de pastas e arquivos-chave.
- **Arquitetura detalhada:** [architecture.md](docs/architecture.md)
- **Visão geral do fluxo:** [overview.md](docs/overview.md)
- **Instalação e execução (setup):** [setup.md](docs/setup.md)
- **Configuração / Variáveis de ambiente:** [configuration.md](docs/configuration.md)
- **Banco de dados e migrations:** [database.md](docs/database.md)
- **Contratos / Mensagens / DTOs:** [api.md](docs/api.md)
- **Desenvolvimento e testes:** [development.md](docs/development.md)
- **Deploy / Containers:** [deployment.md](docs/deployment.md)
- **Contribuição:** [CONTRIBUTING.md](docs/CONTRIBUTING.md)
- **Glossário:** [GLOSSARY.md](docs/GLOSSARY.md)

Relação com a aplicação principal

Esta aplicação processa as mensagens enviadas pela aplicação principal **First API Go** — repositório: https://github.com/EmersonRabelo/first-api-go

Quickstart (local)

```bash
# copiar exemplo de env e ajustar
cp .env.local.example .env.local

# executar (assumindo Postgres e RabbitMQ disponíveis)
go run ./cmd/report-processing-service
```

Executando migrations manualmente (exemplo):

```bash
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f db/migrations/000001_\ create_reports_table_and_indexes.up.sql
```

Contato e referências

- Repositório que produz `CreateReportMessage`: https://github.com/EmersonRabelo/first-api-go
- Para mais detalhes siga o guia: [docs/PROJECT_GUIDE.md](docs/PROJECT_GUIDE.md)

Para detalhes e guias passo a passo, abra os arquivos em `docs/`.