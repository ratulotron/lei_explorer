# Roadmap

## V0: Walking skeleton

Get the golden copy into Postgres end to end. Nothing else.

- Ingest the 3.4M record golden copy from a manually downloaded file
- Upserts keyed on LEI, so re-running the same file changes nothing
- Killing the ingest halfway and restarting should just work, no cleanup

Done when: full golden copy in, a second run changes zero rows, and an integration
test covers the kill + restart case. Ingest numbers go to `BENCHMARKS.md`.

## V1: Query the data

- Ownership chain up and down. If a chain ends, show the reporting exception
  (natural person parent, non-consolidating, opted out)
- Fuzzy search by name, prefix search on partial LEI
- Data quality counts at the end of every ingest: dangling relationship references,
  unresolved ultimate parents, cycles

Done when: `lei chain --dot` draws a real corporate tree, a misspelled name still
finds the right company, and search + chain numbers are in `BENCHMARKS.md`.

## V1.5: Sync + API

- `lei sync` pulls deltas and catches up after missed days. Re-baseline from the
  golden copy if the gap is too big. Scheduling stays outside the binary (cron)
- Chunked concurrent downloads that can resume. Most of the concurrency learning
  is here, the ingest ended up sequential
- gRPC API: `GetEntity`, `SearchByName`, `GetOwnershipChain` (streaming)
- MCP server on top of the gRPC client, so LLM agents can query

Done when: blank machine + `docker compose up` + `lei sync` reaches current data,
Claude answers "who ultimately owns X?" over MCP, and stale data heals on the next
sync.

## V2: Maybe

- Neo4j fed from the same pipeline, compared against the recursive CTEs on
  ergonomics and latency. Checks whether ADR-002 was right. Write up as an ADR
- Dropped: Redpanda (`rejected_records` already covers the DLQ idea, and the
  pipeline is Postgres-bound) and Grafana (the run report and `/metrics` are
  enough for a single-user tool)
