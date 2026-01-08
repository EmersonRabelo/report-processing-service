# Deploy

Opções de deploy típicas:

- Containerização (Docker) — crie uma imagem com `go build` e rode em Kubernetes ou serviço de container.
- Máquina virtual ou VM — build e execução direta do binário.

Recomendações para produção:

- Use variáveis de ambiente seguras (secrets) para `DB_PASSWORD`, `AMQP_PASSWORD` e `GOOGLE_PERSPECTIVE_API_TOKEN`.
- Não exponha a API de admin/metrics sem autenticação.
- Configure readiness/liveness probes se rodando em Kubernetes.

Exemplo mínimo de Dockerfile (sugerido):

```dockerfile
FROM golang:1.20-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /bin/report-processing-service ./cmd/report-processing-service

FROM alpine:3.18
COPY --from=build /bin/report-processing-service /bin/report-processing-service
EXPOSE 8080
ENTRYPOINT ["/bin/report-processing-service"]
```

Vars de ambiente precisam ser fornecidas no runtime (via `env` ou secrets manager).
