-- SEAL and DEN growth schema scaffold
-- Adds the first durable tables for telemetry signals, adaptation proposals,
-- growth plans, lineage state, slot bundle versions, and promotion decisions.

create table if not exists adaptation_signals (
  id uuid primary key,
  specialist_id text not null,
  category text not null,
  severity text not null,
  surface text not null,
  task_class text,
  summary text not null,
  evidence_refs jsonb not null default '[]'::jsonb,
  occurred_at timestamptz not null default now(),
  created_at timestamptz not null default now()
);

create index if not exists adaptation_signals_specialist_idx
  on adaptation_signals (specialist_id, occurred_at desc);

create table if not exists gap_clusters (
  id uuid primary key,
  specialist_id text not null,
  category text not null,
  surface text not null,
  count integer not null,
  persistence_score double precision not null,
  severity_score double precision not null,
  reversibility_score double precision not null,
  summaries jsonb not null default '[]'::jsonb,
  evidence_refs jsonb not null default '[]'::jsonb,
  created_at timestamptz not null default now()
);

create table if not exists adaptation_proposals (
  id uuid primary key,
  specialist_id text not null,
  cluster_id uuid,
  surface text not null,
  reason text not null,
  requested_by text not null,
  risk_level text not null,
  requires_harness boolean not null default true,
  requires_parent boolean not null default false,
  rollback_required boolean not null default true,
  status text not null default 'proposed',
  created_at timestamptz not null default now(),
  approved_at timestamptz
);

create index if not exists adaptation_proposals_specialist_idx
  on adaptation_proposals (specialist_id, status, created_at desc);

create table if not exists growth_plans (
  id uuid primary key,
  proposal_id uuid not null,
  specialist_id text not null,
  surface text not null,
  reason text not null,
  experiment_name text not null,
  freeze_plan jsonb not null,
  rollback_plan jsonb not null,
  status text not null default 'planned',
  created_at timestamptz not null default now()
);

create table if not exists growth_experiments (
  id uuid primary key,
  growth_plan_id uuid not null,
  specialist_id text not null,
  surface text not null,
  status text not null default 'staged',
  artifact_refs jsonb not null default '[]'::jsonb,
  metrics jsonb not null default '{}'::jsonb,
  created_at timestamptz not null default now(),
  completed_at timestamptz
);

create table if not exists lineage_nodes (
  id uuid primary key,
  node_type text not null,
  external_id text not null unique,
  status text not null default 'active',
  metadata jsonb not null default '{}'::jsonb,
  created_at timestamptz not null default now()
);

create table if not exists lineage_edges (
  id uuid primary key,
  from_node_id uuid not null,
  to_node_id uuid not null,
  edge_type text not null,
  metadata jsonb not null default '{}'::jsonb,
  created_at timestamptz not null default now()
);

create table if not exists slot_bundle_versions (
  id uuid primary key,
  specialist_id text not null,
  bundle_toml text not null,
  source_map jsonb not null,
  version_label text not null,
  created_by text not null,
  created_at timestamptz not null default now()
);

create index if not exists slot_bundle_versions_specialist_idx
  on slot_bundle_versions (specialist_id, created_at desc);

create table if not exists promotion_decisions (
  id uuid primary key,
  experiment_id uuid not null,
  decision text not null,
  reason text not null,
  approved_by text not null,
  created_at timestamptz not null default now()
);
