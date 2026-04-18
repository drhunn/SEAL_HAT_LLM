# BOOTSTRAP RUNBOOK

## order
Run these in order:

1. `sql/postgres-ddl.sql`
2. `sql/three-tier-memory.sql`
3. `sql/functions.sql`
4. `sql/seed.sql`
5. `sql/verify.sql`

## phase 1 — core schema
Create the base schema, enums, specialists registry, slot store, memory tables, audit tables, and routing audit.

Stop if:
- `pgvector` or `pgcrypto` is missing
- core tables fail to create

## phase 2 — hierarchical memory
Create region, cluster, and cluster membership tables for coarse-to-fine retrieval.

Stop if:
- region/cluster tables fail
- indexes fail

## phase 3 — function layer
Install callable functions for:
- stage/promote memory
- create postmortem
- create self-edit candidate
- coarse-to-fine search
- project `MEMORY.md`

## phase 4 — first specialist seed
Seed the first specialist:
- `ComputerScience-SoftwareEngineering-Specialist-01`
- slot rows and v1 slot content
- starter regions and clusters
- starter durable records
- starter `MEMORY.md` projection

## phase 5 — verify
Run `sql/verify.sql` and confirm:
- required slots exist
- constitutional slots are locked
- required regions/clusters exist
- active seed records exist
- cluster membership exists
- `MEMORY.md` was projected
- health snapshot looks sane

## hard acceptance checks
- parent frozen rule present
- mandatory postmortem rule present
- specialist memory namespace exists
- coarse-to-fine search works structurally
- no specialist may self-activate or self-expand scope
