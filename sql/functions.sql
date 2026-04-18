SET search_path TO agent_core, public;

CREATE OR REPLACE FUNCTION agent_core.fn_status_multiplier(p_status agent_core.memory_status)
RETURNS numeric
LANGUAGE sql
IMMUTABLE
AS $$
  SELECT CASE p_status
    WHEN 'active' THEN 1.00
    WHEN 'staged' THEN 0.85
    WHEN 'deprecated' THEN 0.75
    WHEN 'contradicted' THEN 0.60
    ELSE 0.50
  END::numeric;
$$;

CREATE OR REPLACE FUNCTION agent_core.fn_stage_memory_record(
  p_namespace text,
  p_specialist_id text,
  p_record_kind agent_core.memory_record_kind,
  p_title text,
  p_summary text,
  p_body text,
  p_importance numeric,
  p_confidence numeric,
  p_tags text[],
  p_embedding_text text,
  p_embedding_model text,
  p_embedding vector(1536),
  p_created_by text
) RETURNS uuid
LANGUAGE plpgsql
AS $$
DECLARE v_record_id uuid;
BEGIN
  INSERT INTO agent_core.memory_records (
    namespace, specialist_id, record_kind, title, summary, body, status,
    importance, confidence, tags, embedding_text, created_by
  ) VALUES (
    p_namespace, p_specialist_id, p_record_kind, p_title, p_summary, p_body, 'staged',
    COALESCE(p_importance, 0.5), COALESCE(p_confidence, 0.5), COALESCE(p_tags, '{}'), p_embedding_text, p_created_by
  ) RETURNING record_id INTO v_record_id;

  INSERT INTO agent_core.memory_embeddings (
    record_id, namespace, specialist_id, embedding_model, embedding
  ) VALUES (
    v_record_id, p_namespace, p_specialist_id, p_embedding_model, p_embedding
  );

  RETURN v_record_id;
END;
$$;

CREATE OR REPLACE FUNCTION agent_core.fn_promote_memory_record(
  p_record_id uuid
) RETURNS void
LANGUAGE plpgsql
AS $$
BEGIN
  UPDATE agent_core.memory_records
  SET status = 'active', updated_at = now()
  WHERE record_id = p_record_id AND status = 'staged';
END;
$$;

CREATE OR REPLACE FUNCTION agent_core.fn_create_postmortem(
  p_namespace text,
  p_specialist_id text,
  p_task_summary text,
  p_expected_behavior text,
  p_actual_behavior text,
  p_what_went_wrong text,
  p_failure_classification text[],
  p_root_cause text,
  p_preventable boolean,
  p_created_by text
) RETURNS uuid
LANGUAGE plpgsql
AS $$
DECLARE v_pm_id uuid;
BEGIN
  INSERT INTO agent_core.memory_postmortems (
    namespace, specialist_id, task_summary, expected_behavior, actual_behavior,
    what_went_wrong, failure_classification, root_cause, preventable, created_by
  ) VALUES (
    p_namespace, p_specialist_id, p_task_summary, p_expected_behavior, p_actual_behavior,
    p_what_went_wrong, COALESCE(p_failure_classification, '{}'), p_root_cause, p_preventable, p_created_by
  ) RETURNING postmortem_id INTO v_pm_id;
  RETURN v_pm_id;
END;
$$;

CREATE OR REPLACE FUNCTION agent_core.fn_update_specialist_health_stats(
  p_specialist_id text
) RETURNS TABLE (specialist_id text)
LANGUAGE plpgsql
AS $$
DECLARE
  v_recent_count integer;
  v_health numeric(5,4);
BEGIN
  SELECT COUNT(*) INTO v_recent_count
  FROM agent_core.memory_postmortems
  WHERE specialist_id = p_specialist_id
    AND created_at >= now() - interval '30 days';

  v_health := GREATEST(0.0, LEAST(1.0, 1.0 - (COALESCE(v_recent_count, 0) * 0.10)));

  UPDATE agent_core.specialists
  SET health_score = v_health,
      updated_at = now()
  WHERE specialists.specialist_id = p_specialist_id;

  RETURN QUERY SELECT p_specialist_id;
END;
$$;

CREATE OR REPLACE FUNCTION agent_core.fn_get_specialist_health_snapshot(
  p_specialist_id text
) RETURNS TABLE (
  specialist_id text,
  status agent_core.specialist_status,
  health_score numeric
)
LANGUAGE sql
STABLE
AS $$
  SELECT s.specialist_id, s.status, s.health_score
  FROM agent_core.specialists s
  WHERE s.specialist_id = p_specialist_id;
$$;

CREATE OR REPLACE FUNCTION agent_core.fn_project_and_write_memory_summary_slot(
  p_namespace text,
  p_specialist_id text,
  p_created_by text,
  p_rationale text
) RETURNS uuid
LANGUAGE plpgsql
AS $$
DECLARE
  v_slot_id uuid;
  v_current_version integer;
  v_next_version integer;
  v_summary text;
BEGIN
  SELECT slot_id, current_version
  INTO v_slot_id, v_current_version
  FROM agent_core.slots
  WHERE specialist_id = p_specialist_id
    AND slot_name = 'MEMORY.md';

  IF v_slot_id IS NULL THEN
    RAISE EXCEPTION 'MEMORY.md slot not found for specialist %', p_specialist_id;
  END IF;

  SELECT COALESCE(string_agg(format('- %s: %s', title, summary), E'\n' ORDER BY importance DESC, created_at DESC), '- No active durable memory available.')
  INTO v_summary
  FROM agent_core.memory_records
  WHERE namespace = p_namespace
    AND specialist_id = p_specialist_id
    AND status = 'active'
    AND is_soft_deleted = false;

  v_next_version := COALESCE(v_current_version, 1) + 1;

  INSERT INTO agent_core.slot_versions (
    slot_id,
    specialist_id,
    slot_name,
    version_number,
    content_markdown,
    created_by,
    approval_level,
    rationale
  ) VALUES (
    v_slot_id,
    p_specialist_id,
    'MEMORY.md',
    v_next_version,
    '# MEMORY\n\n## projected_summary\n' || COALESCE(v_summary, '- No active durable memory available.'),
    p_created_by,
    'harness',
    p_rationale
  );

  UPDATE agent_core.slots
  SET current_version = v_next_version,
      updated_at = now()
  WHERE slot_id = v_slot_id;

  RETURN v_slot_id;
END;
$$;

CREATE OR REPLACE FUNCTION agent_core.fn_run_coarse_to_fine_search(
  p_namespace text,
  p_specialist_id text,
  p_query_embedding vector(1536),
  p_top_regions integer DEFAULT 5,
  p_top_clusters integer DEFAULT 10,
  p_top_records integer DEFAULT 12
) RETURNS TABLE (
  record_id uuid,
  record_kind agent_core.memory_record_kind,
  title text,
  summary text,
  body text,
  status agent_core.memory_status,
  semantic_score numeric,
  cluster_score numeric,
  final_score numeric
)
LANGUAGE sql
STABLE
AS $$
  WITH top_regions AS (
    SELECT r.region_id,
      ((1 - (re.embedding <=> p_query_embedding)) * 0.70 + r.importance * 0.15 + r.confidence * 0.15)::numeric AS region_score
    FROM agent_core.memory_regions r
    JOIN agent_core.memory_region_embeddings re ON re.region_id = r.region_id
    WHERE r.namespace = p_namespace AND r.specialist_id = p_specialist_id AND r.status = 'active'
    ORDER BY region_score DESC
    LIMIT p_top_regions
  ), top_clusters AS (
    SELECT c.cluster_id,
      (((1 - (ce.embedding <=> p_query_embedding)) * 0.70 + c.importance * 0.15 + c.confidence * 0.15) * 0.75 + tr.region_score * 0.25)::numeric AS cluster_score
    FROM agent_core.memory_clusters c
    JOIN agent_core.memory_cluster_embeddings ce ON ce.cluster_id = c.cluster_id
    JOIN top_regions tr ON tr.region_id = c.region_id
    WHERE c.namespace = p_namespace AND c.specialist_id = p_specialist_id AND c.status = 'active'
    ORDER BY cluster_score DESC
    LIMIT p_top_clusters
  )
  SELECT mr.record_id,
         mr.record_kind,
         mr.title,
         mr.summary,
         mr.body,
         mr.status,
         (1 - (me.embedding <=> p_query_embedding))::numeric AS semantic_score,
         tc.cluster_score,
         (((1 - (me.embedding <=> p_query_embedding)) * 0.45 + mr.importance * 0.10 + mr.confidence * 0.10 + mcm.membership_weight * 0.10 + tc.cluster_score * 0.25) * agent_core.fn_status_multiplier(mr.status))::numeric AS final_score
  FROM agent_core.memory_cluster_members mcm
  JOIN agent_core.memory_records mr ON mr.record_id = mcm.record_id
  JOIN agent_core.memory_embeddings me ON me.record_id = mr.record_id
  JOIN top_clusters tc ON tc.cluster_id = mcm.cluster_id
  WHERE mcm.namespace = p_namespace
    AND mcm.specialist_id = p_specialist_id
    AND mr.is_soft_deleted = false
    AND mr.status IN ('active', 'staged', 'deprecated', 'contradicted')
  ORDER BY final_score DESC
  LIMIT p_top_records;
$$;
