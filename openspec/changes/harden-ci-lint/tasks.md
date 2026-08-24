## 1. Config

- [ ] 1.1 Expand `.golangci.yml`: enable `revive`, `errcheck`, `noctx`, `rowserrcheck`, `nolintlint` (require-explanation, require-specific); add `gofumpt` formatter; add commented exclude-rule for `errcheck` in `_test.go`; keep current linters
- [ ] 1.2 Add inline exclusions at the 4 handler `Encode` sites with reason `response already committed`

## 2. Correctness fixes

- [ ] 2.1 Add `rows.Err()` check after iteration loops in the 3 read-repo query methods (`expense_read_repo.go`)
- [ ] 2.2 Replace non-context calls flagged by `noctx` (`db.Ping` in `main.go` and `repo_test.go`; `httptest.NewRequest` in handler tests)

## 3. Hygiene

- [ ] 3.1 Run `golangci-lint fmt` to format all files (8 expected)
- [ ] 3.2 Add doc comments to all exported symbols flagged by `revive`

## 4. Policy documentation

- [ ] 4.1 Add "Linter exclusions" section to AGENTS.md: last resort, always reasoned, new recurring patterns require an OpenSpec change; reference Constitution #11

## 5. Verification

- [ ] 5.1 `golangci-lint run ./...` and `golangci-lint fmt --diff ./...` both clean; full test suite green
- [ ] 5.2 Negative test: temporarily add a bare `//nolint` and confirm lint fails; remove it
