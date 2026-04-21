SET search_path TO agent_core, public;

CREATE TABLE IF NOT EXISTS agent_core.route_episodes (
  id text PRIMARY KEY,
  namespace text NOT NULL,
  specialist_id text NOT NULL REFERENCES agent_core.specialists(specialist_id) ON DELETE CASCADE,
  task_id text NOT NULL,
  task_summary text NOT NULL,
  task_class text NOT NULL,
  primary_modality text NOT NULL,
  chosen_target text,
  target_unit_id text,
  target_role text,
  target_model_ref text,
  chosen_executor text,
  execution_mode text,
  confidence double precision,
  was_fallback boolean NOT NULL DEFAULT false,
  fallback_reason text,
  needs_parent_review boolean NOT NULL DEFAULT false,
  requires_fusion boolean NOT NULL DEFAULT false,
  execution_handled boolean,
  execution_host text,
  status text NOT NULL DEFAULT 'recorded',
  error_text text,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS route_episodes_specialist_idx
  ON agent_core.route_episodes (specialist_id, created_at DESC);

CREATE INDEX IF NOT EXISTS route_episodes_task_idx
  ON agent_core.route_episodes (task_id, created_at DESC);
