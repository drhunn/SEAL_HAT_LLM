CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE SCHEMA IF NOT EXISTS agent_core;
SET search_path TO agent_core, public;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'specialist_status') THEN
    CREATE TYPE specialist_status AS ENUM ('proposed','approved','bootstrapping','inactive','shadow','active','degraded','suspended','retired','rejected');
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'slot_class') THEN
    CREATE TYPE slot_class AS ENUM ('constitutional','operational','summary','mixed','ephemeral_log','distillation','bootstrap');
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'approval_level') THEN
    CREATE TYPE approval_level AS ENUM ('none','harness','parent','human_admin');
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'memory_record_kind') THEN
    CREATE TYPE memory_record_kind AS ENUM ('fact','decision','pointer','evidence','failure_pattern','recovery_pattern','tool_pattern','skill_pattern','prompt_pattern','routing_pattern','contradiction','heartbeat_delta','dream_distillation','postmortem','eval_case','self_edit_candidate','handoff_pattern','artifact_reference');
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'memory_status') THEN
    CREATE TYPE memory_status AS ENUM ('staged','active','contradicted','deprecated','rejected','archived','resolved','approved');
  END IF;
END$$;

CREATE TABLE IF NOT EXISTS specialists (
  specialist_id text PRIMARY KEY,
  name text NOT NULL,
  domain text NOT NULL,
  role text NOT NULL,
  lineage_parent_id text NOT NULL,
  status specialist_status NOT NULL DEFAULT 'proposed',
  priority text NOT NULL DEFAULT 'normal',
  memory_namespace text NOT NULL UNIQUE,
  tool_policy_profile text NOT NULL,
  postmortem_required boolean NOT NULL DEFAULT true,
  seal_enabled boolean NOT NULL DEFAULT true,
  harness_enabled boolean NOT NULL DEFAULT true,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS slots (
  slot_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  specialist_id text NOT NULL REFERENCES specialists(specialist_id) ON DELETE CASCADE,
  slot_name text NOT NULL,
  slot_class slot_class NOT NULL,
  owner text NOT NULL,
  approval_level approval_level NOT NULL,
  is_persistent boolean NOT NULL DEFAULT true,
  is_versioned boolean NOT NULL DEFAULT true,
  current_version integer NOT NULL DEFAULT 1,
  is_locked boolean NOT NULL DEFAULT false,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (specialist_id, slot_name)
);

CREATE TABLE IF NOT EXISTS slot_versions (
  slot_version_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  slot_id uuid NOT NULL REFERENCES slots(slot_id) ON DELETE CASCADE,
  specialist_id text NOT NULL REFERENCES specialists(specialist_id) ON DELETE CASCADE,
  slot_name text NOT NULL,
  version_number integer NOT NULL,
  content_markdown text NOT NULL,
  created_by text NOT NULL,
  approved_by text,
  approval_level approval_level NOT NULL DEFAULT 'none',
  rationale text,
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (slot_id, version_number)
);

CREATE TABLE IF NOT EXISTS memory_records (
  record_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  namespace text NOT NULL,
  specialist_id text NOT NULL REFERENCES specialists(specialist_id) ON DELETE CASCADE,
  record_kind memory_record_kind NOT NULL,
  title text NOT NULL,
  summary text NOT NULL,
  body text NOT NULL,
  status memory_status NOT NULL DEFAULT 'staged',
  importance numeric(5,4) NOT NULL DEFAULT 0.5,
  confidence numeric(5,4) NOT NULL DEFAULT 0.5,
  tags text[] NOT NULL DEFAULT '{}',
  embedding_text text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  created_by text NOT NULL,
  is_soft_deleted boolean NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS memory_embeddings (
  record_id uuid PRIMARY KEY REFERENCES memory_records(record_id) ON DELETE CASCADE,
  namespace text NOT NULL,
  specialist_id text NOT NULL REFERENCES specialists(specialist_id) ON DELETE CASCADE,
  embedding_model text NOT NULL,
  embedding vector(1536) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_memory_embeddings_hnsw ON memory_embeddings USING hnsw (embedding vector_cosine_ops);

CREATE TABLE IF NOT EXISTS memory_postmortems (
  postmortem_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  namespace text NOT NULL,
  specialist_id text NOT NULL REFERENCES specialists(specialist_id) ON DELETE CASCADE,
  task_summary text NOT NULL,
  expected_behavior text NOT NULL,
  actual_behavior text NOT NULL,
  what_went_wrong text NOT NULL,
  failure_classification text[] NOT NULL DEFAULT '{}',
  root_cause text NOT NULL,
  preventable boolean NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  created_by text NOT NULL
);

CREATE TABLE IF NOT EXISTS routing_audit (
  route_audit_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  task_id text,
  routed_by text NOT NULL,
  task_summary text NOT NULL,
  task_class text NOT NULL,
  chosen_target text NOT NULL,
  confidence numeric(5,4) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);
