## 1. Service Layer

- [x] 1.1 Extract shared validation function (amount > 0, description required, credit requires card, installments >= 1, status default) used by Create and Update
- [x] 1.2 In Create: default zero `purchase_date` to today before insert
- [x] 1.3 In Update: run shared validation, return not-found when expense doesn't exist
- [x] 1.4 Add/extend service tests: update invalid amount rejected, update nonexistent returns not-found, purchase_date defaults to today

## 2. Handler Layer

- [x] 2.1 Map validation failures to 400 with JSON error body on Update (currently 500)
- [x] 2.2 Return 404 on Update/Delete when expense doesn't exist
- [x] 2.3 Add handler tests: update invalid amount → 400, update/delete missing id → 404, create without purchase_date → current date in response

## 3. Read-Side Consistency

- [x] 3.1 Replace SQL truncation in `internal/calculation/repository/impl/postgres/expense_read_repo.go` (`ListForProjection`) with derivation using `pkg/money.DivideRoundUp`
- [x] 3.2 Verify history/summary read paths derive identically; refactor any duplicated division to the shared helper
- [x] 3.3 Add tests asserting rounding behavior (e.g., total 115000 / 3 installments → 38334)

## 4. Verification

- [x] 4.1 Run full test suite (`go test ./...`)
- [x] 4.2 Manual smoke: create credit expense with interest → raw values stored as-is; delete then list confirms removal; projection shows rounded-up installments
