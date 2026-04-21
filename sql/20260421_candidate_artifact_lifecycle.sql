SET search_path TO agent_core, public;

UPDATE agent_core.specialist_artifacts
SET activation_status = 'current'
WHERE activation_status = 'active';

DROP INDEX IF EXISTS agent_core.specialist_artifacts_specialist_unit_idx;
DROP INDEX IF EXISTS specialist_artifacts_specialist_unit_idx;

CREATE UNIQUE INDEX IF NOT EXISTS specialist_artifacts_current_idx
  ON agent_core.specialist_artifacts (specialist_id)
  WHERE activation_status = 'current';
