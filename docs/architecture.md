# Arquitetura

Visão de componentes e responsabilidades:

- Ingestão (Queue Consumer): componente que consome mensagens AMQP e aciona o `ConsumerReportService`.
- Serviço de Domínio (`service`): contém regras de validação, persistência mínima e orquestra chamadas externas.
- Repositório (`repository`): abstração de persistência usando GORM e tabela `reports`.
- Integração externa: cliente para a API Perspective (ver `internal/api/perspective`).
- Integração assíncrona de saída: produtor que publica o `ReportAnalysisResultMessage`.

Diagrama lógico (texto):

Producer -> RabbitMQ -> Consumer (this service) -> DB (Postgres)
Consumer -> Perspective API (HTTP) -> Consumer -> DB update -> Producer -> RabbitMQ

Decisões importantes:

- Persistência idempotente: `InsertIfNotExists` evita duplicatas se a mesma mensagem for reenviada.
- Comunicação assíncrona entre serviços via RabbitMQ para desacoplamento.
- Uso de GORM para abstração do Postgres.

Escalabilidade e observabilidade:

- O consumidor pode ser escalado horizontalmente (várias instâncias), desde que idempotência seja garantida pelo repositório.
- Logging está configurado no `config` com níveis diferentes para `production` e `local`.
