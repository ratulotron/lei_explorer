
# TODOs

## V0

- [x] Initialize Go project
    - [x] CLI entry point `cmd/ingestor/main.go` takes file path and prints record count
    - [x] Makefile for triggering pipeline
- [x] Add Dockerized infra
    - [x] PostgreSQL with `pg_trgm` extension
- [ ] `internal/store`: `Store` interface + PostgreSQL implementation
- [ ] `lei ingest <file>`: stream parse of the LEI-CDF CSV, 6 most obvious columns first
    - [ ] Write down the validation contract: what makes a row rejected vs accepted
- [ ] `entities` table, loaded via `pgx.CopyFrom` into staging, merged with `ON CONFLICT`
- [ ] `import_runs` table: run id, source file, status, counts, timestamps
- [ ] `rejected_records` table: raw line + error reason for every bad row
- [ ] End-of-run report: rows read, upserted, rejected, duration
- [ ] Ingest the full 3.4M record golden copy
- [ ] Record ingest throughput and peak memory in `BENCHMARKS.md`

## V1

- [ ] Relationships
    - [ ] Ingest RR-CDF into a `relationships` table, same staging + merge pattern
    - [ ] Ingest Reporting Exceptions into a `reporting_exceptions` table
    - [ ] Resolution pass after merge: materialize `ultimate_parent_lei` and depth
      via recursive CTE, with cycle detection (ADR-004)
- [ ] Queries
    - [ ] Ownership chain up and down via recursive CTE, with cycle guard
    - [ ] Corporate group rollups via `ultimate_parent_lei`
- [ ] Search
    - [ ] Fuzzy name search with `pg_trgm` (GIN index)
    - [ ] Prefix search on partial LEI
- [ ] CLI
    - [ ] `lei search <name>`
    - [ ] `lei chain <lei>`, with a `--dot` flag emitting Graphviz for a visual network
    - [ ] `lei chain` shows the reporting exception reason when a chain ends without a parent
- [ ] Data quality report at end of ingest: dangling relationship references,
  unresolved ultimate parents, cycles found
- [ ] Benchmarks in `BENCHMARKS.md`, linked from ADR-002: search latency
  (p50/p95, index size, vs `ILIKE` baseline), chain latency by depth

## V1.5

- [ ] Sync
    - [ ] `lei sync`: download deltas, verify checksum, apply in order
    - [ ] Compute catch-up from `import_runs` (which snapshot or delta was last applied)
    - [ ] Re-baseline from the golden copy when the gap is too large to replay
    - [ ] Scheduling via cron (document in README), no in-process scheduler
- [ ] Downloader (most of the concurrency learning lives here)
    - [ ] Chunked concurrent download with bounded goroutines, `errgroup`, and
      context cancellation
    - [ ] Resume partial downloads if GLEIF supports HTTP Range requests, otherwise
      fall back to a clean re-download
- [ ] gRPC
    - [ ] Define `lei_search.proto`, `protoc` Makefile target
    - [ ] Server: `GetEntity`, `SearchByName`, `GetOwnershipChain` (server streaming)
    - [ ] Interceptors: logging (`log/slog`), panic recovery. Timeouts are context
      deadlines set by the caller, not an interceptor
    - [ ] Health endpoint and basic `/metrics` counters from the ingest report
- [ ] MCP server: thin wrapper over the gRPC client, so agents can call `SearchByName`
  and `GetOwnershipChain` directly

## V2

- [ ] Neo4j read replica fed by the same pipeline, comparison ADR with measurements
