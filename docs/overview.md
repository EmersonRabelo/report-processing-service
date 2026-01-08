# Visão Geral

Este serviço processa denúncias (reports) recebidas via fila de mensagens, analisa o conteúdo usando a API Perspective do Google e persiste os resultados em banco de dados. Em seguida publica o resultado da análise para que outros serviços possam reagir.

Fluxo simplificado:

1. Um produtor publica uma mensagem de criação de `Report` na fila (AMQP / RabbitMQ).
2. O consumidor deste serviço consome a mensagem e valida o payload.
3. O serviço persiste o registro na tabela `reports` (caso ainda não exista).
4. O texto é enviado à API Perspective para análise de atributos de toxicidade.
5. O serviço atualiza o registro com os scores retornados e marca o status como `done`.
6. É publicada uma mensagem de resultado de análise com os scores para outras filas/consumidores.

Principais responsabilidades do serviço:

- Garantir persistência consistente de cada `Report` (evitar duplicatas).
- Enriquecer o `Report` com scores de análise de conteúdo.
- Publicar eventos de resultado para integração assíncrona.

Arquivos relevantes:

- `internal/service/consumer_report_service.go` — lógica principal do pipeline.
- `internal/dto/report/contracts` — modelos das mensagens trocadas.
- `internal/config` — inicialização de DB e broker.
