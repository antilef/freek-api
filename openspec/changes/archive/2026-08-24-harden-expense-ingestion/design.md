## Context

The ingestion module (households, members, categories, cards, expenses) is implemented and tested, but an audit against `openspec/specs/expense-ingestion/` found four gaps:

1. Interest math: derived values are computed inconsistently — the calculation repo truncates in SQL (`int64(total)/installments`) while `pkg/money.DivideRoundUp` exists unused by that path.
2. `Update` performs no validation and surfaces database errors as 500s.
3. The schema has a `status` column and the read side filters on it, but nothing ever writes a non-default value — soft-delete was anticipated but never built.
4. Omitting `purchase_date` stores Go's zero time (`0001-01-01`) because the handler always binds an explicit value, defeating the DB's `DEFAULT CURRENT_DATE`.

Stakeholder decisions made during exploration:
- **Interest doctrine: derive at read time.** Ingestion stores raw facts; reads compute derived values. (This reverses an earlier draft of this change that proposed persisting calculated columns.)
- Cross-household reference validation is explicitly out of scope: not possible in this deployment model.
- Soft-delete is not supported in this version.

## Goals / Non-Goals

**Goals:**
- Standardize derivation of total-with-interest and monthly amounts on a single round-up formula (`pkg/money.DivideRoundUp`) across all read paths
- Apply create-grade validation to update, with correct HTTP status codes (400 for invalid input, 404 for missing expense)
- Make hard-delete semantics explicit in spec and behavior (404 on missing)
- Default `purchase_date` to today when omitted

**Non-Goals:**
- Persisting derived values (`total_with_interest_clp`, `monthly_amount_clp` columns) — rejected doctrine
- Cross-household referential validation
- Soft-delete, restore endpoints, or status transitions
- Auth
- Fixing the projection service's stub logic and nil dependency (separate follow-up change)

## Decisions

### D1: Store raw facts at ingestion, derive at read time
No schema changes. `expenses` keeps storing `amount_clp`, `installments`, `interest_rate`; every consumer derives totals on demand.

*Why not persist calculated columns?* Considered and rejected (this change's first draft proposed it): duplication of truth, migration + backfill cost, and update paths must remember to recompute. Raw-fact storage means one writer of the formula conceptually lives with the readers, and ingestion stays trivially simple.

*Why not derive inside SQL?* Keeps business math in Go where it is unit-tested via `pkg/money`, not scattered through Postgres dialect. The existing SQL truncation is exactly the bug class this avoids.

### D2: Single validation path shared by Create and Update
Extract validation into one function used by both service methods. Update additionally requires the expense to exist (returns not-found).

*Rationale:* eliminates the drift that caused gap #2; a future field added to Create validation automatically applies to Update.

### D3: Handler maps domain errors to status codes explicitly
Validation errors → 400 with `{"error": "..."}`. Missing expense (GetByID/Update/Delete returning nil) → 404. Everything else remains 500.

*Rationale:* current code returns 500 for all repo errors including not-found; the cheap fix is checking nil before proceeding rather than introducing typed error machinery.

### D4: Default purchase_date in the handler/service layer, not the DB
When `purchase_date.IsZero()`, set it to `time.Now()` before insert.

*Rationale:* the DB default can never fire because Go always sends an explicit value; fixing at the write boundary keeps DB schema untouched.

### D5: Hard delete stays as-is; document it
No code change needed — `DELETE FROM expenses` already matches. Spec makes permanence explicit so future readers don't assume `status` implies soft-delete. The `status` column stays (harmless, read side depends on it) but no behavior claims rest on it.

## Risks / Trade-offs

- [Derivation repeated on every read] → Negligible arithmetic cost; acceptable at this scale. If it ever matters, caching or persisted columns can be revisited without breaking clients.
- [Formula drift risk remains across future read paths] → Mitigated by making `DivideRoundUp` the single mandated helper; new read code should call it, never re-implement division.
- [Breaking behavior: none] → No response shapes change; no schema changes; purely additive validation and status-code corrections.

## Migration Plan

None required — no schema or persisted-data changes. Deploy is a plain rollout; rollback is redeploying the previous binary.

## Open Questions

(none)
