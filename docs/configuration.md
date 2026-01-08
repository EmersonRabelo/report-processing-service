# Configuração

As configurações são carregadas por arquivo `.env.<APP_ENV>` (se existir) ou pelas variáveis do ambiente do sistema. As variáveis usadas pelo sistema são:

- `APP_ENV` — ambiente da aplicação (`local`, `development`, `test`, `production`). Padrão: `local`.
- `SERVER_PORT` — porta HTTP do servidor. Padrão: `8080`.

Banco de dados (Postgres):
- `DB_HOST` — host do Postgres (padrão: `localhost`).
- `DB_PORT` — porta do Postgres (padrão: `5432`).
- `DB_USER` — usuário do DB (padrão: `postgres`).
- `DB_PASSWORD` — senha do DB (padrão: `1234`).
- `DB_NAME` — nome do banco (padrão: `go-api`).
- `DB_SSL_MODE` — modo SSL para conexão (padrão: `disable`).

Broker AMQP (RabbitMQ):
- `AMQP_HOST` — host do RabbitMQ (padrão: `localhost`).
- `AMQP_PORT` — porta do RabbitMQ (padrão: `5672`).
- `AMQP_USER` — usuário AMQP (padrão: `guest`).
- `AMQP_PASSWORD` — senha AMQP (padrão: `guest`).

Integração com Perspective API:
- `GOOGLE_PERSPECTIVE_API_BASE_URL` — URL base da API (padrão presente no código).
- `GOOGLE_PERSPECTIVE_API_TOKEN` — token (API key) para acessar a Perspective API.

Notas:
- Valores padrão são fornecidos no código (`internal/config/config.go`) quando a variável não está definida.
- Para ambientes de produção, não deixe `GOOGLE_PERSPECTIVE_API_TOKEN` vazio.
