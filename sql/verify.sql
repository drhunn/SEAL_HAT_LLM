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
