# LEI Explorer

LEI Explorer is a local, self-updating index of the 3.4M companies that hold a Legal
Entity Identifier (LEI). It answers two questions fast: find a company by (fuzzy) name
or LEI, and trace who ultimately owns whom.

The underlying data comes from the Global Legal Entity Identifier Foundation (GLEIF).
Their public API works but is slow for exploration; since GLEIF publishes full snapshots
and deltas as fresh as 8 hours, a local copy serves the same answers in milliseconds.

It is also my toy project for learning Go in its full breadth. The docs below
record what I picked and why.

---

## Relevant Reads

- **Architecture Decisions:** `./docs/ADR.md`
- **Architecture Design:** `./docs/ARCHITECTURE.md`
- **Todos per Version:** `./docs/TODOS.md`
- **Roadmap:** `./docs/ROADMAP.md`
- **Benchmarks:** `./docs/BENCHMARKS.md`

---

## What I Want to Learn

- **Go concurrency**: goroutines, channels, worker pools, `errgroup`, `context`
  propagation. Mostly in the chunked downloader, the ingest ended up sequential
- **Go project structure**: `cmd/`, `internal/`, interfaces, separation of concerns
- **gRPC in Go**: protobuf definitions, generated code, server-side streaming, interceptors
- **TDD in Go**: table-driven tests, `testcontainers-go` for integration tests
- **Boring architecture**: how far a single binary and Postgres can go before anything
  fancier is justified, and proving it with numbers

---

## Principles

These apply to every version:

- **Idempotency**: All writes are upserts keyed on LEI. Ingesting the same data twice
  changes nothing.
- **Resumability**: Progress lives in Postgres (`import_runs` plus idempotent writes).
  A crashed ingest is just restarted.
- **Rejected rows go to a table, not logs**: `rejected_records` keeps the raw line and
  the error so they can be queried and re-driven later. Only component level failures
  go to `stderr`. Needs a TTL at 3M+ records per run.
- **Benchmark scale assumptions**: anything that depends on data size (`pg_trgm`
  search, recursive CTEs, single process ingest) gets measured, numbers in
  `./docs/BENCHMARKS.md`.

## Datasets

Three files, all freely available:

| File                      | Description                                                                                                                                                                                                | Count                                |
| ------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------ |
| LEI-CDF v3.1 Golden Copy  | Level 1. LEI records and related reference data that provide information on who is who.                                                                                                                    | (2026-07-07 00:00) 3,364,801 Records |
| RR-CDF v2.1 Golden Copy   | Level 2. 'Who owns whom' , allows the identification of the direct and ultimate parents of a legal entity whose direct and ultimate parents have an LEI.                                                   | (2026-07-07 00:00) 478,853 Records   |
| Reporting Exceptions v2.1 | Level 2. LEI registrants that have no parent entity, or the child legal entity opts out of reporting for exceptional reasons, or the direct and ultimate parents of the LEI registrant do not have an LEI. | (2026-07-07 00:00) 6,219,311 Records |

- Golden copy: <https://www.gleif.org/en/lei-data/gleif-golden-copy/download-the-golden-copy>
- GLEIF export links:
  <https://gist.githubusercontent.com/m8d3/6a188311e5f0c667854d2e64dd567046/raw/d66c8b7e4eb67c2007c869ca3968a7945e09b124/download-the-golden-copy.md>
- Level 1 LEI-CDF (LEI-CDF v3.1): <https://leidata-preview.gleif.org/api/v2/golden-copies/publishes/lei2/latest.csv>
- Level 2 Relationship Record (RR-CDF v2.1):
  <https://leidata-preview.gleif.org/api/v2/golden-copies/publishes/rr/latest.csv>
- Level 2 Reporting Exceptions (Reporting Exceptions v2.1):
  <https://leidata-preview.gleif.org/api/v2/golden-copies/publishes/repex/latest.csv>

Note: GLEIF records carry no industry or sector classification. Any "companies that do
X" feature needs enrichment from an external dataset first.

---

## Tech Stack

| Concern     | Tool                               | Why                                                                                     |
| ----------- | ---------------------------------- | --------------------------------------------------------------------------------------- |
| Database    | **PostgreSQL**                     | Entities, relationships, fuzzy search (`pg_trgm`) and run tracking all fit in one store |
| Language    | **Go**                             | Fast, typesafe concurrent language that compiles to a single binary                     |
| API         | **gRPC**                           | Query API for the CLI and the MCP server, with server-side streaming                    |
| Testing     | Go `testing` + `testcontainers-go` | TDD with real infra, no mocks for integration tests                                     |
| Local infra | **Docker Compose**                 | Single `docker compose up` to start everything                                          |

## Go Conventions

- Errors propagate up: `log.Fatal` only in `main`, never in library code
- Interfaces at the boundary, defined where they are consumed, kept small. Extract them
  when a second implementation or a test seam needs one, not before
- `context.Context` is always the first argument on I/O functions
- Table-driven tests for all unit tests
- `//go:build integration` tag on any test requiring Docker
- `golangci-lint` wired into the Makefile

## Readings

- Pipelines and cancellation (Go blog)
- `errgroup` pattern
- Download files in resumable fashion in Golang:
  <https://transloadit.com/devtips/build-a-resumable-file-downloader-in-go-with-concurrent-chunks/>
- Serve large files via HTTP: <https://oneuptime.com/blog/post/2026-01-30-how-to-handle-large-file-downloads-in-go/view>
