# Arquitectura

## Arquitectura por Capas

Esta API sigue una arquitectura por capas para mantener responsabilidades aisladas y el codigo mantenible.

- handler: recibe solicitudes HTTP, valida aspectos de transporte y las mapea a casos de uso de la aplicacion.
- service: contiene reglas de negocio y orquesta el comportamiento del dominio.
- repository: encapsula el acceso a datos y los detalles de persistencia.
- database: gestiona la configuracion tecnica de base de datos, pooling y migraciones SQL.
- domain: define entidades centrales e invariantes usadas por los servicios.
- dto: separa los modelos de entrada/salida de la API respecto de las entidades de dominio.

## Flujo de Solicitud

El flujo de solicitud es:

request -> handler -> service -> repository -> database

El flujo de respuesta regresa en direccion inversa hasta escribir la respuesta HTTP.

## Por que Migraciones SQL en lugar de AutoMigrate

Se usan migraciones SQL versionadas para garantizar una evolucion de esquema explicita y reproducible.

- Cada cambio de esquema es trazable por version y script.
- Los cambios son deterministas en entornos locales, CI y produccion.
- Restricciones, indices y claves foraneas se revisan como artefactos de primera clase.
- La estrategia de rollback puede planificarse por migracion.

Este proyecto evita intencionalmente AutoMigrate para mantener control explicito del esquema.

## Por que Docker Compose

Docker Compose proporciona un entorno local consistente para la API y PostgreSQL.

- Arranque reproducible para todo el equipo.
- Dependencias entre servicios explicitamente definidas.
- Mismas variables de entorno y modelo de red que en entornos de integracion.

## Por que Gin

Gin se usa como framework HTTP liviano y de alto rendimiento.

- Router rapido y soporte de middlewares.
- Binding de requests y renderizado JSON claros.
- Boilerplate minimo para endpoints REST.

## Por que las Reglas de Negocio Viven en Servicios

Los servicios aislan la logica de negocio de las preocupaciones de transporte y persistencia.

- Los handlers se enfocan en preocupaciones HTTP.
- Los repositories se enfocan en acceso a datos.
- Las reglas de negocio se mantienen testeables y reutilizables entre interfaces.

## Por que los DTOs se Separan del Dominio

Los DTOs desacoplan contratos externos del modelado interno de dominio.

- Los contratos de API pueden evolucionar sin filtrar detalles de persistencia.
- Las entidades de dominio se mantienen enfocadas en la semantica central del modelo.
- Los modelos de entrada y salida pueden ajustarse por caso de uso de cada endpoint.
