SET search_path TO agent_core, public;

CREATE TABLE IF NOT EXISTS agent_core.specialist_artifacts (
  id text PRIMARY KEY,
  specialist_id text NOT NULL REFERENCES agent_core.specialists(specialist_id) ON DELETE CASCADE,
  unit_id text NOT NULL,
  role text NOT NULL,
  model_ref text NOT NULL,
  harness_config_ref text,
  store_mode text NOT NULL,
  store_bootstrap_ref text,
  slot_bundle_ref text,
  slot_version_hash text,
  eval_suite_ref text,
  activation_status text NOT NULL DEFAULT 'defined',
  rollback_ref text,
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS specialist_artifacts_specialist_unit_idx
  ON agent_core.specialist_artifacts (specialist_id, unit_id);

CREATE INDEX IF NOT EXISTS specialist_artifacts_specialist_idx
  ON agent_core.specialist_artifacts (specialist_id, updated_at DESC);
