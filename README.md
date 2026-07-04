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
- CRUD de usuarios para perfil ADMIN:
	- Crear usuarios EXECUTOR/AUDITOR con contrasena temporal.
	- Listar usuarios.
	- Obtener usuario por ID.
	- Actualizar usuario (sin tocar password_hash ni must_change_password).
	- Soft delete (is_active=false).
- CRUD de tareas para perfil ADMIN:
	- Crear tareas con estado inicial ASSIGNED.
	- Listar tareas.
	- Obtener tarea por ID.
	- Actualizar tareas solo si estan en ASSIGNED.
	- Soft delete de tareas solo si estan en ASSIGNED.
- Middleware de autorizacion por rol (ADMIN para `/api/users`).
- Middleware de autorizacion por rol (ADMIN para `/api/tasks`).
- Tests unitarios para AuthService y UserService con repositorios fake.
- Tests unitarios para TaskService (reglas de asignacion, vencimiento y estado).

Pendiente (segun enunciado):

- Funcionalidad de perfil Ejecutor.
- Funcionalidad de perfil Auditor.
- Comentarios de tareas vencidas (flujo Ejecutor).

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
	- Auth middleware: valida Bearer JWT.
	- Role middleware: controla permisos por perfil.
- `security`: JWT y password hashing.
- `response`: formato estandar de salida API.

Flujo:

`request -> handler -> service -> repository -> database`

Mas detalle en [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Seguridad aplicada

- Password hashing con bcrypt.
- JWT firmado (HMAC) con `sub`, `role`, `iat`, `exp`.
- Endpoints protegidos con middleware Bearer.
- Endpoints de usuarios protegidos por rol ADMIN.
- Endpoints de tareas protegidos por rol ADMIN.
- Seeder de admin con password hasheado.
- No se expone `password_hash` en respuestas.

## Migraciones SQL

Las migraciones viven en `migrations/` y se ejecutan al iniciar la API.

Reglas:

- No se usa AutoMigrate.
- Evolucion de esquema explicita, versionada y reproducible.
- Soporte de soft delete en tareas via `deleted_at` (migracion incremental).

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

### Users (solo ADMIN)

- `POST /api/users` (protegido + rol ADMIN)
- `GET /api/users` (protegido + rol ADMIN)
- `GET /api/users/:id` (protegido + rol ADMIN)
- `PUT /api/users/:id` (protegido + rol ADMIN)
- `DELETE /api/users/:id` (protegido + rol ADMIN)

### Tasks (solo ADMIN)

- `POST /api/tasks` (protegido + rol ADMIN)
- `GET /api/tasks` (protegido + rol ADMIN)
- `GET /api/tasks/:id` (protegido + rol ADMIN)
- `PUT /api/tasks/:id` (protegido + rol ADMIN)
- `DELETE /api/tasks/:id` (protegido + rol ADMIN)

## Reglas de negocio de usuarios (Sprint 2)

- El ADMIN solo puede crear usuarios con rol `EXECUTOR` o `AUDITOR`.
- No se permite crear ni actualizar usuarios al rol `ADMIN` desde el CRUD.
- El usuario creado parte con:
	- `must_change_password=true`
	- `is_active=true`
	- contrasena temporal aleatoria segura (solo visible en respuesta inicial para esta prueba tecnica).
- El delete de usuarios es logico (`is_active=false`), no fisico.
- Nunca se expone `password_hash` en respuestas.

## Reglas de negocio de tareas (Sprint 3A)

- Solo ADMIN puede administrar tareas.
- En create:
	- `assigned_to` debe existir.
	- `assigned_to` debe tener rol `EXECUTOR`.
	- `due_date` debe ser futura.
	- el estado inicial siempre es `ASSIGNED`.
- En update:
	- solo se puede actualizar si la tarea esta en `ASSIGNED`.
	- campos permitidos: `title`, `description`, `due_date`, `assigned_to`.
	- no se permite asignar a roles distintos de `EXECUTOR`.
- En delete:
	- solo se permite eliminar si esta en `ASSIGNED`.
	- se aplica soft delete (`deleted_at`), no borrado fisico.

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

Incluye tambien tests unitarios de UserService para:

- creacion de usuarios EXECUTOR y AUDITOR;
- rechazo de rol ADMIN;
- rechazo de email duplicado;
- validacion de `must_change_password` e `is_active` en alta;
- restriccion de update a rol ADMIN;
- soft delete.

Incluye tests unitarios de TaskService para:

- create exitoso con estado inicial `ASSIGNED`;
- rechazo de asignacion a `AUDITOR`;
- rechazo de asignacion a `ADMIN`;
- rechazo por `due_date` vencida;
- update solo permitido en `ASSIGNED`;
- delete solo permitido en `ASSIGNED`.

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

# Crear usuario EXECUTOR (solo ADMIN)
curl -X POST "http://localhost:8080/api/users" \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer TOKEN" \
	-d '{"name":"Juan Ejecutor","email":"juan@test.com","role":"EXECUTOR"}'

# Listar usuarios (solo ADMIN)
curl -X GET "http://localhost:8080/api/users" \
	-H "Authorization: Bearer TOKEN"

# Crear tarea (solo ADMIN)
curl -X POST "http://localhost:8080/api/tasks" \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer TOKEN" \
	-d '{"title":"Tarea Sprint 3A","description":"Implementacion inicial","due_date":"2026-07-20T10:00:00Z","assigned_to":2}'

# Listar tareas (solo ADMIN)
curl -X GET "http://localhost:8080/api/tasks" \
	-H "Authorization: Bearer TOKEN"
```

## Mapa de cumplimiento del enunciado (Ejercicio 1)

- Login con perfiles: implementado (token con role).
- Cambio de contrasena: implementado.
- Logout: implementado.
- CRUD de usuarios (ADMIN): implementado.
- CRUD de tareas (ADMIN): implementado.
- Reglas por perfil Administrador: implementado para auth + usuarios + tareas.
- Reglas por perfil Ejecutor/Auditor: pendiente.
- Diagrama y decisiones de arquitectura: implementado en [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Uso de IA en este proyecto

Se utilizo IA como apoyo para:

- Acelerar scaffolding de capas y estructura base.
- Proponer mejoras de consistencia en respuestas y auth.
- Revisar redaccion tecnica de documentacion.

Las decisiones de arquitectura, alcance, reglas de negocio aplicadas y validaciones finales fueron ajustadas y verificadas manualmente en el proyecto.

