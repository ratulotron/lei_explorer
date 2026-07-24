# Benchmarks

Numbers backing the claims in `ADR.md`. Each entry notes how it was measured so it
can be re-run after schema or code changes.

Environment for all runs unless noted: (machine, CPU, RAM, Postgres version,
dataset publish date).

## Ingest (V0)

Claim: single process, native Go concurrency handles the full golden copy (ADR-003).

| Metric | Value | How measured |
| ------ | ----- | ------------ |
| Full golden copy ingest, wall clock | TBD | `time lei ingest <file>` |
| Throughput (rows/s) | TBD | rows read / wall clock |
| Peak memory | TBD | (fill in tool) |
| Second run (idempotency), wall clock | TBD | same command, expect zero row changes |

## Fuzzy search (V1)

Claim: `pg_trgm` covers fuzzy search over 3.4M rows (ADR-002).

| Metric | Value | How measured |
| ------ | ----- | ------------ |
| p50 / p95 latency, misspelled name | TBD | `EXPLAIN ANALYZE`, N queries |
| GIN index size | TBD | `pg_relation_size` |
| Baseline: same query via `ILIKE`, no index | TBD | for contrast |

## Ownership chain (V1)

Claim: recursive CTEs cover traversal on a sparse, shallow graph (ADR-002).

| Metric | Value | How measured |
| ------ | ----- | ------------ |
| Chain latency by depth (1, 3, 5, max found) | TBD | `EXPLAIN ANALYZE` |
| Deepest chain in dataset | TBD | resolution pass output |
| Ultimate parent resolution pass, wall clock | TBD | part of ingest report |

## Neo4j comparison (V2)

Same queries as above against Neo4j fed from the same pipeline. Filled in when the
V2 comparison runs.
