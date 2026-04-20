SET search_path TO agent_core, public;

-- SEAL and DEN growth schema scaffold
-- Adds the first durable tables for telemetry signals, adaptation proposals,
-- growth plans, lineage state, slot bundle versions, and promotion decisions.

CREATE TABLE IF NOT EXISTS agent_core.adaptation_signals (
  id text PRIMARY KEY,
  specialist_id text NOT NULL REFERENCES agent_core.specialists(specialist_id) ON DELETE CASCADE,
  category text NOT NULL,
  severity text NOT NULL,
  surface text NOT NULL,
  task_class text,
  summary text NOT NULL,
  evidence_refs jsonb NOT NULL DEFAULT '[]'::jsonb,
  occurred_at timestamptz NOT NULL DEFAULT now(),
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS adaptation_signals_specialist_idx
  ON agent_core.adaptation_signals (specialist_id, occurred_at DESC);

CREATE TABLE IF NOT EXISTS agent_core.gap_clusters (
  id text PRIMARY KEY,
  specialist_id text NOT NULL REFERENCES agent_core.specialists(specialist_id) ON DELETE CASCADE,
  category text NOT NULL,
  surface text NOT NULL,
  count integer NOT NULL,
  persistence_score double precision NOT NULL,
  severity_score double precision NOT NULL,
  reversibility_score double precision NOT NULL,
  summaries jsonb NOT NULL DEFAULT '[]'::jsonb,
  evidence_refs jsonb NOT NULL DEFAULT '[]'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS agent_core.adaptation_proposals (
  id text PRIMARY KEY,
  specialist_id text NOT NULL REFERENCES agent_core.specialists(specialist_id) ON DELETE CASCADE,
  cluster_id text REFERENCES agent_core.gap_clusters(id) ON DELETE SET NULL,
  surface text NOT NULL,
  reason text NOT NULL,
  requested_by text NOT NULL,
  risk_level text NOT NULL,
  requires_harness boolean NOT NULL DEFAULT true,
  requires_parent boolean NOT NULL DEFAULT false,
  rollback_required boolean NOT NULL DEFAULT true,
  status text NOT NULL DEFAULT 'proposed',
  created_at timestamptz NOT NULL DEFAULT now(),
  approved_at timestamptz
);

CREATE INDEX IF NOT EXISTS adaptation_proposals_specialist_idx
  ON agent_core.adaptation_proposals (specialist_id, status, created_at DESC);

CREATE TABLE IF NOT EXISTS agent_core.growth_plans (
  id text PRIMARY KEY,
  proposal_id text NOT NULL REFERENCES agent_core.adaptation_proposals(id) ON DELETE CASCADE,
  specialist_id text NOT NULL REFERENCES agent_core.specialists(specialist_id) ON DELETE CASCADE,
  surface text NOT NULL,
  reason text NOT NULL,
  experiment_name text NOT NULL,
  freeze_plan jsonb NOT NULL,
  rollback_plan jsonb NOT NULL,
  status text NOT NULL DEFAULT 'planned',
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS agent_core.growth_experiments (
  id text PRIMARY KEY,
  growth_plan_id text NOT NULL REFERENCES agent_core.growth_plans(id) ON DELETE CASCADE,
  specialist_id text NOT NULL REFERENCES agent_core.specialists(specialist_id) ON DELETE CASCADE,
  surface text NOT NULL,
  status text NOT NULL DEFAULT 'staged',
  artifact_refs jsonb NOT NULL DEFAULT '[]'::jsonb,
  metrics jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now(),
  completed_at timestamptz
);

CREATE TABLE IF NOT EXISTS agent_core.lineage_nodes (
  id text PRIMARY KEY,
  node_type text NOT NULL,
  external_id text NOT NULL UNIQUE,
  status text NOT NULL DEFAULT 'active',
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS agent_core.lineage_edges (
  id text PRIMARY KEY,
  from_node_id text NOT NULL REFERENCES agent_core.lineage_nodes(id) ON DELETE CASCADE,
  to_node_id text NOT NULL REFERENCES agent_core.lineage_nodes(id) ON DELETE CASCADE,
  edge_type text NOT NULL,
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS agent_core.slot_bundle_versions (
  id text PRIMARY KEY,
  specialist_id text NOT NULL REFERENCES agent_core.specialists(specialist_id) ON DELETE CASCADE,
  bundle_toml text NOT NULL,
  source_map jsonb NOT NULL,
  version_label text NOT NULL,
  created_by text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS slot_bundle_versions_specialist_idx
  ON agent_core.slot_bundle_versions (specialist_id, created_at DESC);

CREATE TABLE IF NOT EXISTS agent_core.promotion_decisions (
  id text PRIMARY KEY,
  experiment_id text NOT NULL REFERENCES agent_core.growth_experiments(id) ON DELETE CASCADE,
  decision text NOT NULL,
  reason text NOT NULL,
  approved_by text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);
