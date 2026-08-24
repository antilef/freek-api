# AGENTS.md

## Database

Run `./up-db.sh` to create (or start) a Docker container with PostgreSQL 18 for this project. The script also applies all migrations from `migrations/`. Connection string: `postgres://localhost:5432/freek_test?sslmode=disable`.

The test suite (`go test ./...`) requires the database to be up.
