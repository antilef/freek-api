## Context

Freek es una API en Go para gestión de gastos familiares en Chile. El usuario necesita registrar gastos (efectivo/débito y crédito) y ver proyecciones de pago de tarjetas de crédito a 12 meses. No existe código previo - es un proyecto nuevo.

El dominio tiene dos partes claras:
1. **Ingesta**: escritura de gastos, tarjetas, miembros, categorías
2. **Cálculo**: lectura y proyección de liquidaciones futuras e histórico

El usuario es un hogar (una persona o familia), con miembros como etiquetas para filtrar, no para separar cuentas.

## Goals / Non-Goals

**Goals:**
- API REST funcional para registrar y consultar gastos
- Proyección de liquidaciones de tarjeta a 12 meses
- Soporte CLP (pesos chilenos) sin decimales
- Períodos flexibles configurables por el usuario
- Categorías predefinidas editables por hogar
- Notas libres en gastos para contexto

**Non-Goals:**
- Autenticación / multi-usuario (futuro)
- Importación de PDFs de liquidaciones (big bang futuro)
- Reconciliación con estados de cuenta reales
- Simulaciones de compras
- Interfaz web (solo API)
- Mobile app

## Decisions

### 1. Gorilla/mux como router HTTP

**Decisión**: Usar `gorilla/mux` para el router HTTP.

**Por qué sobre chi o gin**:
- Gorilla es el estándar histórico de Go, ampliamente documentado
- El usuario lo pidió explícitamente
- Manejo limpio de subrouters para versionar la API (/api/v1/...)

**Alternativa descartada**: chi (más moderno pero el usuario prefiere gorilla)

### 2. PostgreSQL como base de datos

**Decisión**: PostgreSQL para almacenamiento.

**Por qué sobre SQLite**:
- Soporte nativo de fechas y rangos de fechas
- Mejor para queries de proyección (CTEs, window functions)
- Escalable si la familia crece o se comparte
- JSONB para datos flexibles (notas, metadata futura)

**Por qué sobre MySQL**:
- Mejor soporte de tipos date/range
- Mais expressivo para queries recursivas (proyección 12 meses)

### 3. CLP como enteros (int64)

**Decisión**: Almacenar montos en pesos chilenos como `int64`, sin decimales.

**Razón**: CLP no tiene centavos. Usar float causa problemas de precisión. Enteros son exactos.

```go
type CLP int64

func (m CLP) String() string {
    return "$" + formatWithDots(int64(m))
}
// $120.000 en vez de 120000.00
```

### 4. Dos dominios separados (ingesta / cálculo)

**Decisión**: Separar la lógica en `internal/ingestion` (escritura) y `internal/calculation` (lectura/proyección).

**Razón**:
- Separación de responsabilidades clara
- El cálculo de proyecciones es complejo y no debe mezclarse con CRUD
- Facilita testing independiente
- Futuro: el cálculo podría ser un microservicio separado

```
internal/
├── ingestion/       ← CRUD: crear, actualizar, borrar
│   ├── domain/
│   ├── repository/
│   ├── service/
│   └── handler/
├── calculation/     ← Queries: proyección, histórico, liquidaciones
│   ├── domain/
│   ├── repository/
│   ├── service/
│   └── handler/
└── shared/          ← Compartido: money, dates, household context
```

### 5. Modelo de hogar con miembros como filtros

**Decisión**: Un solo nivel de cuentas (hogar). Los miembros son etiquetas para filtrar gastos, no para separar presupuestos.

**Razón**:
- Simplifica el modelo enormemente
- Los gastos son del hogar, el "quién compró" es informativo
- Filtrar por titular es la uso principal, no separación contable
- Si en el futuro se necesita separar, se agrega `household_id` por usuario

### 6. Períodos flexibles

**Decisión**: El usuario define su propio "mes fiscal" (ej: del 5 al 4).

**Razón**:
- Los días de pago en Chile varían (1, 5, 10, 15, 20...)
- Un mes calendario (1-30) no refleja la realidad del usuario
- La tabla `period_config` almacena start_day y end_day por hogar

### 7. Migraciones con golang-migrate

**Decisión**: Usar `golang-migrate` para manejar esquema de base de datos.

**Razón**:
- Estándar en proyectos Go
- Archivos SQL versionados y repetibles
- Fácil rollback si algo sale mal
- Funciona bien con CI/CD

### 8. Sin autenticación inicial

**Decisión**: La API no tiene auth por ahora. Un solo hogar, acceso directo.

**Razón**:
- MVP rápido, sin complejidad innecesaria
- Auth se agrega después con JWT o sesiones
- El `household_id` se pasa como header o query param por ahora

## Risks / Trade-offs

**[Riesgo] Proyección 12 meses puede ser lenta con muchos gastos** → Mitigación: Usar CTEs recursivas en PostgreSQL, cachear resultados por período. Si es necesario, materializar vistas.

**[Riesgo] Sin auth, cualquier request puede ver/modificar todo** → Mitigación: Aceptar para MVP local. Auth se agrega antes de producción.

**[Riesgo] CLP como int64 no maneja fracciones** → Mitigación: CLP no tiene centavos, esto es correcto. Si se necesita división (ej: dividir $100.000 en 3 cuotas), se redondea hacia arriba.

**[Riesgo] Gorilla/mux está en maintenance mode** → Mitigación: Es estable y funciona. Si se necesita más features, se migra a chi o stdlib.

**[Trade-off] Dos dominios = más archivos** → Aceptado: La claridad vale la complejidad extra de carpetas.

**[Trade-off] Sin simulaciones** → Aceptado: El usuario pidió proyección real, no simulación. Se puede agregar después.
