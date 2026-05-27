# Migrations

Este projeto usa migrations SQL versionadas no formato esperado pelo `golang-migrate`.

## Conceito

Migration e a evolucao controlada do schema do banco. Cada mudanca estrutural deve ser registrada em arquivos SQL versionados.

Neste projeto, cada migration tem dois arquivos:

- `.up.sql`: aplica a mudanca
- `.down.sql`: desfaz a mudanca

Exemplo:

```text
migrations/
  000001_create_webhook_events.up.sql
  000001_create_webhook_events.down.sql
```

O banco de dados em si deve existir antes de rodar as migrations. A migration cria tabelas, indices, constraints e outras estruturas dentro desse banco.

## Instalacao da CLI

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Executando localmente

Configure a URL do banco:

```bash
export DATABASE_URL='postgres://usuario:senha@localhost:5432/webhook_delivery?sslmode=disable'
```

No PowerShell:

```powershell
$env:DATABASE_URL='postgres://usuario:senha@localhost:5432/webhook_delivery?sslmode=disable'
```

Aplique as migrations:

```bash
migrate -path migrations -database "$DATABASE_URL" up
```

No PowerShell:

```powershell
migrate -path migrations -database $env:DATABASE_URL up
```

Desfaca a ultima migration:

```bash
migrate -path migrations -database "$DATABASE_URL" down 1
```

## Producao

Em producao, prefira rodar migrations como etapa separada do deploy, antes de subir a nova versao da API.

Esse fluxo evita que multiplas instancias da aplicacao tentem alterar o schema ao mesmo tempo e deixa mais claro quando uma mudanca de banco falha.
