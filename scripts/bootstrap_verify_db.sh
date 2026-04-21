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
  sql/20260421_route_episodes.sql
 do
  psql "$DSN" -v ON_ERROR_STOP=1 -f "$f"
 done
