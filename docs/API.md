# API Documentation

Base URL: `http://localhost:8080`

Formato de respuesta estandar:

```json
{
  "success": true,
  "message": "...",
  "data": {},
  "errors": null
}
```

## Health

### GET /health
- Rol requerido: publico
- Request body: no aplica
- Response ejemplo:
```json
{
  "success": true,
  "message": "Service is healthy",
  "data": { "service": "ok" },
  "errors": null
}
```
- Errores esperados: 500 (inesperado)

### GET /health/db
- Rol requerido: publico
- Request body: no aplica
- Response ejemplo (OK):
```json
{
  "success": true,
  "message": "Database connection is healthy",
  "data": { "database": "ok" },
  "errors": null
}
```
- Errores esperados: 503 (DB no disponible), 500 (inesperado)

## Auth

### POST /api/auth/login
- Rol requerido: publico
- Request body:
```json
{ "email": "admin@test.com", "password": "Admin123" }
```
- Response ejemplo:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "jwt-token",
    "user": { "id": 1, "name": "Administrator", "email": "admin@test.com", "role": "ADMIN", "must_change_password": false }
  },
  "errors": null
}
```
- Errores esperados: 400, 401, 403, 500

### GET /api/auth/me
- Rol requerido: autenticado
- Request body: no aplica
- Response ejemplo:
```json
{
  "success": true,
  "message": "Authenticated user retrieved successfully",
  "data": { "id": 1, "name": "Administrator", "email": "admin@test.com", "role": "ADMIN", "must_change_password": false },
  "errors": null
}
```
- Errores esperados: 401, 403, 404, 500

### PATCH /api/auth/change-password
- Rol requerido: autenticado
- Request body:
```json
{ "current_password": "Admin123", "new_password": "Admin1234" }
```
- Response ejemplo:
```json
{ "success": true, "message": "Password changed successfully", "data": null, "errors": null }
```
- Errores esperados: 400, 401, 403, 404, 500

### POST /api/auth/logout
- Rol requerido: autenticado
- Request body: no aplica
- Response ejemplo: HTTP 204 (sin body)
- Errores esperados: 401, 500

## Users (ADMIN)

### POST /api/users
- Rol requerido: ADMIN
- Request body:
```json
{ "name": "Executor One", "email": "executor@test.com", "role": "EXECUTOR" }
```
- Response ejemplo:
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "user": { "id": 2, "name": "Executor One", "email": "executor@test.com", "role": "EXECUTOR", "must_change_password": true, "is_active": true },
    "temporary_password": "ChangeMe123!"
  },
  "errors": null
}
```
- Errores esperados: 400, 409, 500

### GET /api/users
- Rol requerido: ADMIN
- Request body: no aplica
- Response ejemplo:
```json
{ "success": true, "message": "Users retrieved successfully", "data": [{ "id": 1, "role": "ADMIN" }], "errors": null }
```
- Errores esperados: 401, 403, 500

### GET /api/users/:id
- Rol requerido: ADMIN
- Request body: no aplica
- Response ejemplo:
```json
{ "success": true, "message": "User retrieved successfully", "data": { "id": 2, "role": "EXECUTOR" }, "errors": null }
```
- Errores esperados: 400, 401, 403, 404, 500

### PUT /api/users/:id
- Rol requerido: ADMIN
- Request body:
```json
{ "name": "Executor Updated", "role": "EXECUTOR", "is_active": true }
```
- Response ejemplo:
```json
{ "success": true, "message": "User updated successfully", "data": { "id": 2, "name": "Executor Updated" }, "errors": null }
```
- Errores esperados: 400, 401, 403, 404, 500

### DELETE /api/users/:id
- Rol requerido: ADMIN
- Request body: no aplica
- Response ejemplo:
```json
{ "success": true, "message": "User deleted successfully", "data": null, "errors": null }
```
- Errores esperados: 400, 401, 403, 404, 500

## Tasks ADMIN

### POST /api/tasks
- Rol requerido: ADMIN
- Request body:
```json
{
  "title": "Task A",
  "description": "Implement API",
  "due_date": "2026-12-20T10:00:00Z",
  "assigned_to": 2
}
```
- Response ejemplo:
```json
{ "success": true, "message": "Task created successfully", "data": { "id": 10, "status": "ASSIGNED" }, "errors": null }
```
- Errores esperados: 400, 401, 403, 409, 500

### GET /api/tasks
- Rol requerido: ADMIN
- Request body: no aplica
- Response ejemplo:
```json
{ "success": true, "message": "Tasks retrieved successfully", "data": [{ "id": 10, "status": "ASSIGNED" }], "errors": null }
```
- Errores esperados: 401, 403, 500

### GET /api/tasks/:id
- Rol requerido: ADMIN
- Request body: no aplica
- Response ejemplo:
```json
{ "success": true, "message": "Task retrieved successfully", "data": { "id": 10, "status": "ASSIGNED" }, "errors": null }
```
- Errores esperados: 400, 401, 403, 404, 500

### PUT /api/tasks/:id
- Rol requerido: ADMIN
- Request body:
```json
{ "title": "Task A updated", "description": "Updated", "due_date": "2026-12-22T10:00:00Z", "assigned_to": 2 }
```
- Response ejemplo:
```json
{ "success": true, "message": "Task updated successfully", "data": { "id": 10 }, "errors": null }
```
- Errores esperados: 400, 401, 403, 404, 500

### DELETE /api/tasks/:id
- Rol requerido: ADMIN
- Request body: no aplica
- Response ejemplo:
```json
{ "success": true, "message": "Task deleted successfully", "data": null, "errors": null }
```
- Errores esperados: 400, 401, 403, 404, 500

## Tasks EXECUTOR

### GET /api/tasks/my
- Rol requerido: EXECUTOR
- Request body: no aplica
- Response ejemplo:
```json
{ "success": true, "message": "My tasks retrieved successfully", "data": [{ "id": 10, "status": "ASSIGNED" }], "errors": null }
```
- Errores esperados: 401, 403, 500

### GET /api/tasks/my/:id
- Rol requerido: EXECUTOR
- Request body: no aplica
- Response ejemplo:
```json
{ "success": true, "message": "Task retrieved successfully", "data": { "id": 10, "status": "ASSIGNED" }, "errors": null }
```
- Errores esperados: 400, 401, 403, 404, 500

### PATCH /api/tasks/my/:id/status
- Rol requerido: EXECUTOR
- Request body:
```json
{ "status": "IN_PROGRESS" }
```
- Response ejemplo:
```json
{ "success": true, "message": "Task status updated successfully", "data": { "id": 10, "status": "IN_PROGRESS" }, "errors": null }
```
- Errores esperados: 400, 401, 403, 404, 500

### POST /api/tasks/my/:id/comments
- Rol requerido: EXECUTOR
- Request body:
```json
{ "comment": "No pude completar por bloqueo externo" }
```
- Response ejemplo:
```json
{ "success": true, "message": "Task comment created successfully", "data": { "id": 3, "task_id": 10, "user_id": 2 }, "errors": null }
```
- Errores esperados: 400, 401, 403, 404, 500

## Audit

### GET /api/audit/tasks
- Rol requerido: AUDITOR
- Request body: no aplica
- Response ejemplo:
```json
{ "success": true, "message": "Tasks retrieved successfully", "data": [{ "id": 10, "status": "IN_PROGRESS" }], "errors": null }
```
- Errores esperados: 401, 403, 500
