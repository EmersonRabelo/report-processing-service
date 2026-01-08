# API e Contratos de Mensagem

O serviço consome e produz mensagens AMQP usando contratos definidos em `internal/dto/report/contracts`.

Principais contratos:

- `CreateReportMessage` — mensagem de criação de um report (consumida pelo serviço).
- `ReportAnalysisResultMessage` — mensagem publicada quando a análise é concluída.

Resumo dos campos importantes (exemplo):

- `CreateReportMessage`:
  - `Id` (UUID) — identificador do report.
  - `PostId` (UUID) — identificador do post reportado.
  - `ReporterId` (UUID) — usuário que fez a denúncia.
  - `Body` (string) — conteúdo textual analisado.
  - `CreatedAt` (timestamp)

- `ReportAnalysisResultMessage`:
  - `ReportId` (UUID)
  - `Toxicity`, `SevereToxicity`, `IdentityAttack`, `Insult`, `Profanity`, `Threat` — valores float (pontuações da API Perspective).
  - `Language` — linguagem detectada (opcional).
  - `AnalyzedAt` — timestamp da análise.

Onde olhar o contrato exato:

- `internal/dto/report/contracts/create_report_message.go`
- `internal/dto/report/contracts/report_analysis_result_message.go`

Boas práticas para integração:

- Garanta que `Id`, `PostId` e `ReporterId` sejam UUIDs válidos.
- `Body` não deve ser vazio; o consumidor rejeita mensagens inválidas.
