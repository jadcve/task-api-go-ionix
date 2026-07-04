# Task API Go Ionix

API REST desarrollada en Go para la prueba tecnica, con enfoque en arquitectura por capas, seguridad basica, migraciones SQL versionadas y calidad incremental sobre completitud.

## Estado del proyecto

Implementado actualmente:

- Migraciones SQL versionadas (sin AutoMigrate).
- Dominio base: usuarios, tareas y comentarios.
- Autenticacion JWT:
	- Login.
	- Usuario autenticado (`/api/auth/me`).
	- Cambio de contrasena.
	- Logout.
- Middleware de autenticacion Bearer.
- Seeder idempotente de admin inicial.
- Respuesta HTTP estandar para endpoints de negocio.

Pendiente (segun enunciado):

- CRUD completo de usuarios (perfil Administrador).
- CRUD completo de tareas y reglas avanzadas de estado.
- Funcionalidad de perfil Ejecutor.
- Funcionalidad de perfil Auditor.

## Stack tecnico

- Go 1.25+
- Gin (HTTP)
- PostgreSQL
- pgx (acceso a datos)
- golang-migrate (migraciones SQL)
- bcrypt (hash de contrasenas)
- JWT (auth token)
- Docker + Docker Compose

## Arquitectura

Arquitectura por capas:

- `handler`: entrada HTTP y mapeo request/response.
- `service`: reglas de negocio y orquestacion.
- `repository`: persistencia y consultas SQL.
- `database`: conexion, migraciones y seed.
- `domain`: entidades y enums de negocio.
- `dto`: contratos de entrada/salida.
- `middleware`: cross-cutting concerns HTTP.
- `security`: JWT y password hashing.
- `response`: formato estandar de salida API.

Flujo:

`request -> handler -> service -> repository -> database`

Mas detalle en [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Seguridad aplicada

- Password hashing con bcrypt.
- JWT firmado (HMAC) con `sub`, `role`, `iat`, `exp`.
- Endpoints protegidos con middleware Bearer.
- Seeder de admin con password hasheado.
- No se expone `password_hash` en respuestas.

## Migraciones SQL

Las migraciones viven en `migrations/` y se ejecutan al iniciar la API.

Reglas:

- No se usa AutoMigrate.
- Evolucion de esquema explicita, versionada y reproducible.

Comandos:

```bash
make migrate-up
make migrate-down
```

## Variables de entorno

Copiar `.env.example` a `.env`.

```env
APP_ENV=development
APP_PORT=8080

DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=task_api
DB_SSLMODE=disable

JWT_SECRET=change-me
JWT_EXPIRATION_HOURS=24
MIGRATIONS_PATH=file://migrations

ADMIN_EMAIL=admin@test.com
ADMIN_PASSWORD=Admin123
ADMIN_NAME=Administrator
```

## Ejecucion local

```bash
go mod tidy
make run
```

## Ejecucion con Docker Compose

```bash
make up
make logs
```

Detener y limpiar volumen:

```bash
make down
```

## Endpoints disponibles

### Salud

- `GET /health`

### Auth

- `POST /api/auth/login`
- `GET /api/auth/me` (protegido)
- `PATCH /api/auth/change-password` (protegido)
- `POST /api/auth/logout` (protegido)

## Formato estandar de respuesta

Exito:

```json
{
	"success": true,
	"message": "Login successful",
	"data": {
		"token": "...",
		"user": {
			"id": 1,
			"name": "Administrator",
			"email": "admin@test.com",
			"role": "ADMIN",
			"must_change_password": false
		}
	},
	"errors": null
}
```

Error:

```json
{
	"success": false,
	"message": "Invalid credentials",
	"data": null,
	"errors": null
}
```

## Pruebas

Ejecutar:

```bash
make test
```

Incluye tests unitarios del modulo Auth service con repositorio fake (sin dependencia de PostgreSQL real).

## Validacion manual rapida

```bash
# Login
curl -X POST "http://localhost:8080/api/auth/login" \
	-H "Content-Type: application/json" \
	-d '{"email":"admin@test.com","password":"Admin123"}'

# Me (reemplazar TOKEN)
curl -X GET "http://localhost:8080/api/auth/me" \
	-H "Authorization: Bearer TOKEN"

# Change password
curl -X PATCH "http://localhost:8080/api/auth/change-password" \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer TOKEN" \
	-d '{"current_password":"Admin123","new_password":"Admin123"}'

# Logout
curl -i -X POST "http://localhost:8080/api/auth/logout" \
	-H "Authorization: Bearer TOKEN"
```

## Mapa de cumplimiento del enunciado (Ejercicio 1)

- Login con perfiles: implementado (token con role).
- Cambio de contrasena: implementado.
- Logout: implementado.
- CRUD usuarios/tareas: pendiente.
- Reglas por perfil Administrador/Ejecutor/Auditor: parcial.
- Diagrama y decisiones de arquitectura: implementado en [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Uso de IA en este proyecto

Se utilizo IA como apoyo para:

- Acelerar scaffolding de capas y estructura base.
- Proponer mejoras de consistencia en respuestas y auth.
- Revisar redaccion tecnica de documentacion.

Las decisiones de arquitectura, alcance, reglas de negocio aplicadas y validaciones finales fueron ajustadas y verificadas manualmente en el proyecto.

