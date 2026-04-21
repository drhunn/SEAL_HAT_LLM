SET search_path TO agent_core, public;

CREATE TABLE IF NOT EXISTS agent_core.specialist_artifact_events (
  event_id text PRIMARY KEY,
  artifact_id text NOT NULL REFERENCES agent_core.specialist_artifacts(id) ON DELETE CASCADE,
  specialist_id text NOT NULL REFERENCES agent_core.specialists(specialist_id) ON DELETE CASCADE,
  event_type text NOT NULL,
  experiment_id text,
  actor text NOT NULL,
  reason text,
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS specialist_artifact_events_artifact_idx
  ON agent_core.specialist_artifact_events (artifact_id, created_at DESC);

CREATE INDEX IF NOT EXISTS specialist_artifact_events_experiment_idx
  ON agent_core.specialist_artifact_events (experiment_id);
