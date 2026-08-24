# Glossary / Glosario

Domain terms of **freek** (personal finance / household expense tracker).
Canonical names are in **English** (as used in code, DB columns and API), with
Spanish translations. Use these exact terms when writing code, prompts or specs.

> Términos de dominio de **freek** (app de finanzas personales / gastos del hogar).
> Los nombres canónicos están en **inglés** (como se usan en código, columnas de BD y API),
> con su traducción al español. Usa estos términos exactos al escribir código, prompts o specs.

---

## Core entities / Entidades principales

### Household — Hogar
The root aggregate and multi-tenancy boundary. Every other entity is scoped by
`household_id`. Enforced by middleware (`internal/shared/household.go`).

La raíz de todo el dominio y el límite de multi-tenancy. Todas las demás entidades
están acotadas por `household_id`, forzado por middleware.

- `id` int64 · `name` string

### Member — Miembro / Integrante
A named person belonging to a household. Has **no roles** (no owner/admin).
Optional holder of expenses and credit cards.

Una persona con nombre dentro de un hogar. **No tiene roles** (ni owner/admin).
Titular opcional de gastos y tarjetas de crédito.

- `id` · `household_id` · `name` (unique per household)

### Category — Categoría
Classification of expenses. `household_id NULL` means a **global predefined**
category (`is_predefined = true`); predefined ones are immutable
(cannot be renamed or deleted).

Clasificación de los gastos. `household_id NULL` significa categoría
**predefinida global** (`is_predefined = true`); las predefinidas son inmutables
(no se renombran ni eliminan).

- `id` · `household_id *` · `name` · `icon *` · `is_predefined`
- Seeded (Spanish in DB): Alimentación→Food/Groceries, Transporte→Transportation,
  Salud→Health, Educación→Education, Entretenimiento→Entertainment,
  Hogar→Home, Ropa→Clothing, Servicios→Utilities/Bills, Otros→Others/Misc.
- Fallback label in summaries: `'Sin categoría'` → "Uncategorized".

### CreditCard — Tarjeta de crédito
A credit card of the household, with a monthly statement closing day.

Tarjeta de crédito del hogar, con día de cierre mensual del estado de cuenta.

- `id` · `household_id` · `name` · `closing_day` (**1–28**, avoids invalid dates)
- `member_id *` (optional cardholder) · `credit_limit *` (integer CLP, nullable)

### Expense — Gasto
A money outflow. Amounts are **always integer CLP** (no decimals; CLP has no cents).

Una salida de dinero. Los montos son **siempre enteros en CLP** (sin decimales;
el peso chileno no tiene centavos).

- `id` · `household_id` · `amount_clp` (> 0, BIGINT)
- `description` (required) · `notes *`
- `category_id *` (optional, SET NULL on delete)
- `payment_type` enum: `cash | debit | credit`
  - `credit` ⇒ requires `credit_card_id` (validated in service)
- `member_id *` (who incurred it; fallback `'Sin titular'` → "No member")
- `purchase_date` DATE (defaults to today if zero)
- `installments` ≥ 1 ("cuotas"; coerced to 1 minimum)
- `interest_rate` fraction ≥ 0, e.g. `0.02` (applied only to credit)
- `status` VARCHAR(20); only `"active"` used today; projections filter on it

---

## Value objects & domain concepts / Objetos de valor y conceptos de dominio

### Money (CLP) — Dinero / Peso chileno
`pkg/money`: `type CLP int64`. Formatting/parsing uses dot thousands separators
(`$1.250.000`). Helpers:

- `TotalWithInterest(amount, rate)` → `round(amount × (1 + rate))`
- `DivideRoundUp(amount, installments)` → integer division always rounded up
- `MonthlyInstallment(total, n)` → per-cuota amount (*cuota mensual*)

### Installments — Cuotas
Number of monthly payments of a credit purchase (`installments ≥ 1`).
Monthly charge = round-up(total-with-interest ÷ installments).

Número de pagos mensuales de una compra con crédito (`installments ≥ 1`).
Cargo mensual = total con interés ÷ cuotas, redondeado hacia arriba.

### Closing day — Día de cierre
Day of month (1–28) when a card's statement closes.
First closing after purchase: same month if purchase day ≤ closing_day,
otherwise next month (`pkg/dates.NextClosingDate`).

Día del mes (1–28) en que cierra el estado de cuenta de la tarjeta.
Primer cierre tras la compra: mismo mes si día de compra ≤ closing_day,
si no, mes siguiente.

### Settlement / Settlement period — Liquidación / Período de liquidación
A card's billing cycle: from previous month's `closing_day + 1` to this month's
`closing_day`. Identified by `"YYYY-MM"` string
(`pkg/dates.SettlementForMonth`, read model `SettlementPeriod`).

Ciclo de facturación de una tarjeta: desde `closing_day + 1` del mes anterior
hasta `closing_day` del mes actual. Identificado como `"YYYY-MM"`.

### Period — Período (configurable)
`pkg/dates.PeriodConfig{StartDay, EndDay}`: custom monthly period that may span
month boundaries when `StartDay > EndDay`.

Período mensual personalizable que puede cruzar límites de mes si `StartDay > EndDay`.

### Projection — Proyección
Read-model simulation of future months' outstanding credit charges from active
installment expenses (default horizon: 12 months).
`ProjectionResult{Months, Total, MonthlyTotals}` keyed by `"2006-01"`.

Modelo de lectura que simula los cargos futuros pendientes de gastos con cuotas
activas (horizonte por defecto: 12 meses).

### Summary / Compare periods — Resumen / Comparación de períodos
Aggregations by category or member over date ranges; comparing two ranges
returns a `difference`.

Agregaciones por categoría o miembro sobre rangos de fechas; comparar dos
rangos devuelve una `difference` (diferencia).

---

## Relationships / Relaciones

```
Household 1───N Member
Household 1───N Category        (NULL household_id => global/predefined)
Household 1───N CreditCard ───N→1 Member (optional cardholder)
Household 1───N Expense ───N→1 Category   (optional, SET NULL)
                        ├───N→1 CreditCard (required iff payment_type='credit')
                        └───N→1 Member     (optional, SET NULL)
```

## Not yet in the domain / Aún no existen en el dominio
No budget, income, savings, user accounts or authentication entities exist yet;
members are plain named persons identified by integer IDs.

No existen aún entidades de presupuesto, ingresos, ahorro, cuentas de usuario
ni autenticación; los miembros son personas con nombre identificadas por IDs enteros.
