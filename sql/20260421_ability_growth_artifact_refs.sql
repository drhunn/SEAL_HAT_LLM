SET search_path TO agent_core, public;

ALTER TABLE IF EXISTS agent_core.ability_growth_experiments
  ADD COLUMN IF NOT EXISTS specialist_artifact_id text REFERENCES agent_core.specialist_artifacts(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS ability_growth_experiments_artifact_idx
  ON agent_core.ability_growth_experiments (specialist_artifact_id);
