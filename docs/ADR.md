# Architecture Decision Records (ADRs)

## ADR Template

### Status

- Proposed / Accepted / Deprecated / Superseded

### Context

Describe the issue or problem statement that motivated this decision.

### Decision

State the architectural decision that was made.

### Consequences

List the resulting context after applying the decision:

- ✅ Positive consequences
- ⚠️ Neutral consequences
- ❌ Negative consequences

---

## ADR Log

### ADR-001: Single binary CLI, sub-commands to access features for v0

**Status:** Accepted  
**Date:** 02.04.2026
**Author:** Ratul Minhaz

**Context:**

Establishing method of running application. Generally Go applications expose functionalities through the
CLI, the `cmd/<feature>` pattern. We have distinct actions a user can take such as initiating a sync, searching entities
using partial text, or fetching an ownership chain.

**Decision:**

Use a single binary built with sub commands to access features like sync, fuzzy search. Eg.

- `lei sync`: triggers sync with available GLEIF datasets since last sync.
- `lei search <name>`: fetches and prints out companies that roughly match with the search phrase.
- `lei chain <lei>`: see full chain of ownership given any LEI.

**Consequences:**

- ✅ Simple, zero fuss use
- ✅ Consistent with Golang monolithic applications
- ⚠️ Initial setup overhead

---

### ADR-002: PostgreSQL for search and relationships in v1

**Status:** Accepted
**Date:** 02.04.2026
**Author:** Ratul Minhaz

**Context:**

Although the hierarchical ownership data naturally models as a graph in Neo4j, a second store means dual writes and
keeping both consistent on every sync. Most of the data is plainly relational anyway: entity reference data (names,
addresses, registration status, legal form) plus operational tables like `import_runs` and `rejected_records`. The
GLEIF relationship graph is also sparse and shallow (~479K relationship records over ~3.4M entities, mostly short
parent chains), so traversals can be implemented in PostgreSQL through recursive CTEs, and later with the Apache AGE
extension if needed. Similarly Elasticsearch brings resource usage overhead for fuzzy search that can be accomplished
with the `pg_trgm` extension in PostgreSQL at this scale.

**Decision:**

Intentionally using PostgreSQL for fuzzy search and relationship storage.

**Consequences:**

- ✅ Single storage to maintain
- ⚠️ Revisited in V2 as a measured Neo4j comparison (see ROADMAP)
- ⚠️ Claim not measured yet, numbers go to `BENCHMARKS.md` (search and chain sections)
- ❌ Graph based community and similarity not available

---

### ADR-003: Single process, no message broker

**Status:** Accepted
**Date:** 06.07.2026
**Author:** Ratul Minhaz

**Context:**

Initial ingestion of the data might be 3.4M+ entities across many tables, however consequent syncs are far smaller
(below 10k for daily.) This level of load can be processed using Golang's native concurrency. One benefit broker based
ingestion has is the DLQ, we can implement that via logging and dumping failed data in a dedicated table.

**Decision:**

Using single process, native Go concurrency. No message broker in v1.

**Consequences:**

- ✅ Native stack, no infra dependency
- ⚠️ Claim not measured yet, numbers go to `BENCHMARKS.md` (ingest section)
- ⚠️ Not elegantly clear cut separation
- ❌ Might not scale with drastically large load

---

### ADR-004: Ultimate parent via post-ingest resolution pass

**Status:** Proposed
**Date:** 12.07.2026
**Author:** Ratul Minhaz

**Context:**

Group rollups and chain queries need `ultimate_parent_lei` (and depth) on every entity.
It cannot be computed while streaming: a parent may appear after its children in the
snapshot, and relationships arrive in a separate file (RR-CDF) from the entities
(LEI-CDF). The data may also contain cycles, which must be detected rather than
looped over.

**Decision:**

Materialize `ultimate_parent_lei` and depth in a dedicated resolution pass that runs
after all files of an import run are merged. One recursive CTE walks every entity to
its root, guards against cycles by tracking the visited path, and writes the result
back to the entities table. The pass reports unresolved parents, dangling relationship
references, and cycles found, which becomes the data quality report.

**Consequences:**

- ✅ Ingest stays a simple stream, resolution is set based SQL
- ✅ Re-running the pass is idempotent, same inputs give the same columns
- ✅ Cycle and dangling reference detection falls out for free
- ⚠️ Chain data is stale between merge and resolution, so the pass is part of every run
- ❌ Full-table pass on every sync, acceptable until measured otherwise (`BENCHMARKS.md`)

---

## Field Guide

### Status Field

- **Proposed:** Decision is under consideration
- **Accepted:** Decision has been approved and implemented
- **Deprecated:** Decision is no longer used
- **Superseded:** Decision has been replaced by another ADR

### Context Field

Should answer "What is the issue we're trying to solve?" Include:

- Background information
- Constraints
- Requirements

### Decision Field

Should answer "What is our decision?" Include:

- Clear statement of the chosen approach
- Why this option was selected

### Consequences Field

Should answer "What becomes easier, harder, or impossible?" Include:

- Positive outcomes
- Risks or drawbacks
- Dependencies created
