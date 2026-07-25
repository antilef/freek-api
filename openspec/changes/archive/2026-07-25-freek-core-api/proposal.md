## Why

No existe una herramienta simple y local para que una persona o familia chilena pueda registrar sus gastos del mes (efectivo/débito y crédito) y ver proyectados los pagos de sus tarjetas de crédito a 12 meses. Las apps existentes son muy genéricas o muy complejas. Freek resuelve esto con un enfoque en liquidaciones de tarjeta y proyecciones de pago, usando CLP y períodos flexibles que se adaptan a los días de pago reales del usuario.

## What Changes

- **API REST en Go** con Gorilla/mux para ingesta y consulta de gastos
- **Modelo de hogar**: una cuenta familiar con miembros como etiquetas para filtrar, no para separar cuentas
- **Gastos duales**: cada gasto puede ser efectivo/débito o crédito, con soporte para cuotas (con o sin interés)
- **Tarjetas de crédito**: configurables con día de cierre propio, titulares como filtro
- **Proyección 12 meses**: motor de cálculo que proyecta liquidaciones futuras por tarjeta y total del hogar
- **Histórico**: consulta de gastos pasados con filtros flexibles (por fecha, titular, tarjeta, categoría)
- **Categorías**: tabla predefinida pero editable por cada hogar
- **Notas libres**: campo de texto en gastos para contexto ("compré un vino para la junta de amigos")
- **Períodos flexibles**: el usuario define su propio período de mes (ej: del 1 al 30, o del 5 al 4)
- **PostgreSQL** como base de datos con migraciones versionadas
- **CLP como enteros**: sin decimales, formato chileno con puntos de miles

## Capabilities

### New Capabilities

- `household-management`: CRUD de hogares y miembros (etiquetas para filtrar gastos)
- `credit-card-management`: CRUD de tarjetas de crédito con titular, día de cierre y límite
- `expense-ingestion`: registro de gastos con monto CLP, categoría, notas, tipo de pago, cuotas, tasa de interés
- `category-management`: categorías predefinidas editables por hogar
- `expense-projection`: motor de cálculo que proyecta liquidaciones de tarjeta a 12 meses
- `expense-history`: consulta de gastos pasados con filtros por fecha, titular, tarjeta, categoría y período flexible

### Modified Capabilities

_Ninguna - proyecto nuevo._

## Impact

- **Código nuevo**: API completa en Go, sin código existente
- **Dependencias**: Gorilla/mux, driver PostgreSQL (lib/pq o pgx), golang-migrate
- **Base de datos**: PostgreSQL con esquema nuevo (households, members, credit_cards, categories, expenses)
- **API endpoints**: REST para todos los CRUDs y consultas de proyección
- **Sin autenticación por ahora**: single-user (un hogar), auth se agrega después
