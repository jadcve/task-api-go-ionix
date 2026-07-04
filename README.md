# task-api-go-ionix

Base inicial de API REST en Go para desafio tecnico.

## Stack

- Go
- Gin
- PostgreSQL
- Docker Compose
- Variables de entorno

## Arquitectura base

```text
cmd/api/main.go
internal/config
internal/database
internal/handler
internal/service
internal/repository
internal/domain
internal/dto
internal/middleware
internal/security
```

No hay entidades ni logica de negocio implementadas todavia.

## Endpoint disponible

- `GET /health`

Respuesta:

```json
{
	"status": "ok"
}
```

## Variables de entorno

Copiar `.env.example` a `.env` y ajustar valores si hace falta.

```env
APP_PORT=8080
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=tasks
DB_SSLMODE=disable
```

## Ejecucion con Docker Compose

Levanta los servicios `api` y `postgres`:

```bash
make up
```

Ver logs:

```bash
make logs
```

Bajar servicios:

```bash
make down
```

## Ejecucion local

```bash
make run
```

## Tests

```bash
make test
```

