# AGENTS.md

## Database

Run `./up-db.sh` to create (or start) a Docker container with PostgreSQL 18 for this project. The script also applies all migrations from `migrations/`. Connection string: `postgres://localhost:5432/freek_test?sslmode=disable`.

The test suite (`go test ./...`) requires the database to be up.

## Linter exclusions

The lint gate (`.golangci.yml`) runs with `nolintlint`, so any `//nolint`
directive without a written reason fails CI. When generating code:

- Fix the finding instead of excluding it. Exclusion is a last resort.
- If exclusion is truly needed, scope it to one line and always include a
  machine-readable linter and a reason:
  `json.NewEncoder(w).Encode(result) //nolint:errcheck // response already committed`
- The only path-scoped exclusion is `errcheck` in `_test.go` files (test
  fixture/seed calls). Do not add new global or path-scoped exclusions.
- If you find yourself needing a new *category* of exclusion (a recurring
  pattern), stop: that is a design decision. Propose an OpenSpec change to
  amend `.golangci.yml` per `docs/CONSTITUTION.md` principle #11. Never
  sprinkle bare `//nolint` directives to make lint pass.
