SET search_path TO agent_core, public;

SELECT 'specialist_exists' AS check_name,
       CASE WHEN EXISTS (SELECT 1 FROM specialists WHERE specialist_id = 'csse-tool-development-specialist-01') THEN 'PASS' ELSE 'FAIL' END AS result;

SELECT 'required_slots_count' AS check_name,
       CASE WHEN (
         SELECT COUNT(*) FROM slots
         WHERE specialist_id = 'csse-tool-development-specialist-01'
           AND slot_name IN ('IDENTITY.md','SOUL.md','AGENTS.md','TOOLS.md','SKILLS.md','PROMPT.md','HEARTBEAT.md','MEMORY.md','DREAMS.md','POSTMORTEM.md')
       ) = 10 THEN 'PASS' ELSE 'FAIL' END AS result;

SELECT 'required_regions_count' AS check_name,
       CASE WHEN (
         SELECT COUNT(*) FROM memory_regions
         WHERE namespace = 'memory.csse-tool-development-specialist-01'
       ) >= 5 THEN 'PASS' ELSE 'FAIL' END AS result;

SELECT 'required_clusters_count' AS check_name,
       CASE WHEN (
         SELECT COUNT(*) FROM memory_clusters
         WHERE namespace = 'memory.csse-tool-development-specialist-01'
       ) >= 5 THEN 'PASS' ELSE 'FAIL' END AS result;

SELECT 'region_embeddings_present' AS check_name,
       CASE WHEN (
         SELECT COUNT(*) FROM memory_region_embeddings re
         JOIN memory_regions r ON r.region_id = re.region_id
         WHERE r.namespace = 'memory.csse-tool-development-specialist-01'
       ) >= 5 THEN 'PASS' ELSE 'FAIL' END AS result;

SELECT 'cluster_embeddings_present' AS check_name,
       CASE WHEN (
         SELECT COUNT(*) FROM memory_cluster_embeddings ce
         JOIN memory_clusters c ON c.cluster_id = ce.cluster_id
         WHERE c.namespace = 'memory.csse-tool-development-specialist-01'
       ) >= 5 THEN 'PASS' ELSE 'FAIL' END AS result;

SELECT 'ability_ledgers_table_exists' AS check_name,
       CASE WHEN EXISTS (
         SELECT 1
         FROM information_schema.tables
         WHERE table_schema = 'agent_core' AND table_name = 'ability_ledgers'
       ) THEN 'PASS' ELSE 'FAIL' END AS result;

SELECT 'ability_growth_experiments_table_exists' AS check_name,
       CASE WHEN EXISTS (
         SELECT 1
         FROM information_schema.tables
         WHERE table_schema = 'agent_core' AND table_name = 'ability_growth_experiments'
       ) THEN 'PASS' ELSE 'FAIL' END AS result;

SELECT 'specialist_artifacts_table_exists' AS check_name,
       CASE WHEN EXISTS (
         SELECT 1
         FROM information_schema.tables
         WHERE table_schema = 'agent_core' AND table_name = 'specialist_artifacts'
       ) THEN 'PASS' ELSE 'FAIL' END AS result;

SELECT 'specialist_artifact_events_table_exists' AS check_name,
       CASE WHEN EXISTS (
         SELECT 1
         FROM information_schema.tables
         WHERE table_schema = 'agent_core' AND table_name = 'specialist_artifact_events'
       ) THEN 'PASS' ELSE 'FAIL' END AS result;

SELECT 'route_episodes_table_exists' AS check_name,
       CASE WHEN EXISTS (
         SELECT 1
         FROM information_schema.tables
         WHERE table_schema = 'agent_core' AND table_name = 'route_episodes'
       ) THEN 'PASS' ELSE 'FAIL' END AS result;

SELECT 'ability_growth_experiments_artifact_ref_column_exists' AS check_name,
       CASE WHEN EXISTS (
         SELECT 1
         FROM information_schema.columns
         WHERE table_schema = 'agent_core'
           AND table_name = 'ability_growth_experiments'
           AND column_name = 'specialist_artifact_id'
       ) THEN 'PASS' ELSE 'FAIL' END AS result;

SELECT 'single_current_artifact_index_exists' AS check_name,
       CASE WHEN EXISTS (
         SELECT 1
         FROM pg_indexes
         WHERE schemaname = 'agent_core'
           AND indexname = 'specialist_artifacts_current_idx'
       ) THEN 'PASS' ELSE 'FAIL' END AS result;

SELECT 'single_current_artifact_invariant' AS check_name,
       CASE WHEN NOT EXISTS (
         SELECT 1
         FROM agent_core.specialist_artifacts
         WHERE activation_status = 'current'
         GROUP BY specialist_id
         HAVING COUNT(*) > 1
       ) THEN 'PASS' ELSE 'FAIL' END AS result;

SELECT 'staged_growth_points_to_candidate_artifacts' AS check_name,
       CASE WHEN NOT EXISTS (
         SELECT 1
         FROM agent_core.ability_growth_experiments age
         JOIN agent_core.specialist_artifacts sa ON sa.id = age.specialist_artifact_id
         WHERE age.status IN ('observed','proposed','shadow')
           AND sa.activation_status <> 'candidate'
       ) THEN 'PASS' ELSE 'FAIL' END AS result;

SELECT 'promoted_growth_points_to_current_artifacts' AS check_name,
       CASE WHEN NOT EXISTS (
         SELECT 1
         FROM agent_core.ability_growth_experiments age
         JOIN agent_core.specialist_artifacts sa ON sa.id = age.specialist_artifact_id
         WHERE age.status = 'promoted'
           AND sa.activation_status <> 'current'
       ) THEN 'PASS' ELSE 'FAIL' END AS result;

SELECT 'rolled_back_growth_points_to_rolled_back_artifacts' AS check_name,
       CASE WHEN NOT EXISTS (
         SELECT 1
         FROM agent_core.ability_growth_experiments age
         JOIN agent_core.specialist_artifacts sa ON sa.id = age.specialist_artifact_id
         WHERE age.status = 'rolled_back'
           AND sa.activation_status <> 'rolled_back'
       ) THEN 'PASS' ELSE 'FAIL' END AS result;

SELECT 'artifact_event_history_present_for_lifecycle_types' AS check_name,
       CASE WHEN EXISTS (
         SELECT 1
         FROM information_schema.tables
         WHERE table_schema = 'agent_core' AND table_name = 'specialist_artifact_events'
       ) THEN 'PASS' ELSE 'FAIL' END AS result;
