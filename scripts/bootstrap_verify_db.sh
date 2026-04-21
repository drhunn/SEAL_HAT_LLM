#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "usage: $0 <postgres-dsn>" >&2
  exit 1
fi

DSN="$1"

for f in \
  sql/postgres-ddl.sql \
  sql/three-tier-memory.sql \
  sql/functions.sql \
  sql/seed.sql \
  sql/verify.sql \
  sql/multimodal-memory.sql \
  sql/20260420_seal_den_growth.sql \
  sql/20260421_route_episodes.sql \
  sql/20260421_specialist_artifacts.sql \
  sql/20260421_ability_growth_artifact_refs.sql \
  sql/20260421_specialist_artifact_events.sql
 do
  psql "$DSN" -v ON_ERROR_STOP=1 -f "$f"
 done
