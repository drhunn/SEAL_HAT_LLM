SET search_path TO agent_core, public;

DO $$
DECLARE
  v_specialist_id text := 'csse-tool-development-specialist-01';
  v_namespace text := 'memory.csse-tool-development-specialist-01';
  v_zero_vec vector(1536);
BEGIN
  v_zero_vec := array_fill(0.0::real, ARRAY[1536])::vector(1536);

  INSERT INTO specialists (
    specialist_id, name, domain, role, lineage_parent_id, status, priority,
    memory_namespace, tool_policy_profile, postmortem_required, seal_enabled, harness_enabled
  ) VALUES (
    v_specialist_id,
    'ComputerScience-SoftwareEngineering-Specialist-01',
    'computer science, software engineering, tool development, runtime design, harness infrastructure, memory integration',
    'foundational_tool_development_specialist',
    'Parent-Generalist-30B',
    'active',
    'critical',
    v_namespace,
    'csse_tooldev_default',
    true,
    true,
    true
  ) ON CONFLICT (specialist_id) DO NOTHING;

  INSERT INTO slots (specialist_id, slot_name, slot_class, owner, approval_level, is_locked)
  VALUES
    (v_specialist_id, 'IDENTITY.md', 'constitutional', 'parent', 'parent', true),
    (v_specialist_id, 'SOUL.md', 'constitutional', 'parent', 'parent', true),
    (v_specialist_id, 'AGENTS.md', 'constitutional', 'parent', 'parent', true),
    (v_specialist_id, 'TOOLS.md', 'operational', 'shared', 'harness', false),
    (v_specialist_id, 'SKILLS.md', 'operational', 'shared', 'harness', false),
    (v_specialist_id, 'PROMPT.md', 'operational', 'shared', 'harness', false),
    (v_specialist_id, 'HEARTBEAT.md', 'mixed', 'shared', 'harness', false),
    (v_specialist_id, 'MEMORY.md', 'summary', 'shared', 'harness', false),
    (v_specialist_id, 'DREAMS.md', 'distillation', 'shared', 'harness', false),
    (v_specialist_id, 'POSTMORTEM.md', 'operational', 'shared', 'harness', false)
  ON CONFLICT (specialist_id, slot_name) DO NOTHING;

  INSERT INTO memory_regions (namespace, specialist_id, region_name, region_kind, region_summary, region_tags, created_by)
  VALUES
    (v_namespace, v_specialist_id, 'tool_patterns', 'tools', 'Durable patterns and anti-patterns for tool design and validation.', ARRAY['tools','patterns'], 'harness:seed'),
    (v_namespace, v_specialist_id, 'runtime_architecture', 'design', 'Core runtime architecture decisions and boundaries.', ARRAY['runtime','architecture'], 'harness:seed'),
    (v_namespace, v_specialist_id, 'slot_governance', 'decisions', 'Constitutional vs operational slot rules and ACL logic.', ARRAY['slots','governance'], 'harness:seed'),
    (v_namespace, v_specialist_id, 'postgres_pgvector_memory', 'design', 'Schema, projection, contradiction, audit, and retrieval design.', ARRAY['postgres','pgvector','memory'], 'harness:seed'),
    (v_namespace, v_specialist_id, 'postmortems', 'postmortems', 'Failure analysis records and durable lessons learned.', ARRAY['postmortem','failure'], 'harness:seed')
  ON CONFLICT (namespace, region_name) DO NOTHING;

  INSERT INTO memory_region_embeddings (region_id, namespace, specialist_id, embedding_model, embedding)
  SELECT region_id, namespace, specialist_id, 'seed-zero-placeholder', v_zero_vec
  FROM memory_regions r
  WHERE r.namespace = v_namespace
    AND NOT EXISTS (SELECT 1 FROM memory_region_embeddings e WHERE e.region_id = r.region_id);

  INSERT INTO memory_clusters (namespace, specialist_id, region_id, cluster_name, cluster_kind, cluster_summary, created_by)
  SELECT v_namespace, v_specialist_id, region_id,
         CASE region_name
           WHEN 'tool_patterns' THEN 'tool_contracts'
           WHEN 'runtime_architecture' THEN 'runtime_boundaries'
           WHEN 'slot_governance' THEN 'constitutional_vs_operational_slots'
           WHEN 'postgres_pgvector_memory' THEN 'pg_schema_design'
           WHEN 'postmortems' THEN 'postmortem_patterns'
         END,
         'topic_cluster',
         region_summary,
         'harness:seed'
  FROM memory_regions
  WHERE namespace = v_namespace
  ON CONFLICT (region_id, cluster_name) DO NOTHING;

  INSERT INTO memory_cluster_embeddings (cluster_id, namespace, specialist_id, embedding_model, embedding)
  SELECT cluster_id, namespace, specialist_id, 'seed-zero-placeholder', v_zero_vec
  FROM memory_clusters c
  WHERE c.namespace = v_namespace
    AND NOT EXISTS (SELECT 1 FROM memory_cluster_embeddings e WHERE e.cluster_id = c.cluster_id);
END $$;
