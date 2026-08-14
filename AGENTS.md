# AGENTS.md

Repository-specific instructions for AI development agents working on SudokuSolver.

## Project Identity

- **Project name:** SudokuSolver
- **GitHub repository:** `ImtehQ/SudokuSolver`
- **Main/release branch:** `main`
- **Distribution environment:** GitHub Releases
- **Production server:** none; this is a standalone CLI

## Product

SudokuSolver analyzes a Sudoku grid rather than simply printing a solution. Its primary score is the percentage of randomized, constraint-respecting completion trials that reach a valid full Sudoku. This is an empirical analytical score, not a mathematical probability that a deterministic puzzle is solvable. The optional `--verify` mode checks exact solution existence without returning the solved grid.

## Production Safety

There is currently no production database, upload store, persistent volume, or server-side application state. GitHub source history and published release artifacts must still be preserved. Never rewrite `main`, delete releases/tags as a routine fix, or expose credentials.

## Application

- **Language/runtime:** Go 1.23+
- **Package manager:** Go modules
- **CLI entry point:** `cmd/sudokusolver`
- **Core analysis package:** `internal/sudoku`
- **Database:** none
- **Authentication:** none
- **External runtime services:** none

## Development

Install dependencies/build metadata:

```bash
go mod download
```

Run locally:

```bash
go run ./cmd/sudokusolver --help
```

Validation:

```bash
test -z "$(gofmt -l .)"
go vet ./...
go test ./...
go build ./cmd/sudokusolver
```

There are no integration, E2E, database migration, or formatting-tool commands beyond the checks above because the application has no external services or database.

## Git Workflow

Every logical change uses a new branch and PR. Never develop directly on `main`.

Preferred branch pattern:

```text
agent/short-description
```

PRs target `main`. Prefer squash merge when repository requirements allow it. Preserve unrelated changes and keep diffs focused.

## CI

- `.github/workflows/ci.yml` runs on pull requests and pushes to `main`.
- It checks `gofmt`, `go vet`, `go test`, and `go build`.
- Do not claim CI succeeded until the workflow reaches a successful conclusion.

## Releases / Deployment

- `.github/workflows/release.yml` is the normal production distribution process.
- It runs automatically on pushes to `main`.
- The version comes from `VERSION`; every new release-producing change must bump this semantic version before merge.
- The workflow refuses to reuse an existing tag.
- It creates standalone Windows, macOS, and Linux archives plus SHA-256 checksums.
- It smoke-tests the Linux x86-64 build before creating the GitHub Release.
- Do not create a second manual deployment when the merge-triggered release workflow is already running.

A deployment is verified only after the Release workflow succeeds and the GitHub Release exists for the intended `main` commit with the expected assets.

## Persistent Data / Backups

- **Database:** none
- **Uploads:** none
- **Persistent application state:** none
- **Release artifacts:** persistent distribution artifacts; do not delete routinely
- **Source history:** Git/GitHub history is persistent project history

No application-data backup or restore mechanism is required because the CLI stores no production data. Do not describe source control as a tested disaster-recovery backup for user data.

## Rollback

There is no mutable server deployment to roll back. A software rollback means directing users to a previously published immutable GitHub Release or publishing a new corrected version through a new branch/PR. Never delete or rewrite current repository history as rollback behavior.

## Security-sensitive areas

Take special care with:

- GitHub Actions permissions;
- release credentials/tokens;
- parsing untrusted puzzle/file input;
- future network or update-check functionality if introduced.

Never commit secrets.

## Definition of Done

For an applicable development request:

```text
inspect -> branch -> implement -> test -> validate -> commit -> push -> PR -> CI -> merge -> release -> release verification
```

If permissions, mandatory review, Actions policy, or another external gate blocks a stage, stop at that exact stage and report what completed, what is blocked, production impact, next action, and whether a human must act.
