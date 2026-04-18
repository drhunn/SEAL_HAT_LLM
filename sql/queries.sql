SET search_path TO agent_core, public;

-- Tier 1: regions
WITH scored_regions AS (
  SELECT r.region_id, r.region_name, r.region_kind, r.region_summary,
         (1 - (re.embedding <=> :query_embedding)) AS semantic_score,
         ((1 - (re.embedding <=> :query_embedding)) * 0.70 + r.importance * 0.15 + r.confidence * 0.15) AS final_score
  FROM memory_regions r
  JOIN memory_region_embeddings re ON re.region_id = r.region_id
  WHERE r.namespace = :namespace
    AND r.specialist_id = :specialist_id
    AND r.status = 'active'
)
SELECT * FROM scored_regions ORDER BY final_score DESC LIMIT :top_regions;

-- Tier 2: clusters
WITH scored_clusters AS (
  SELECT c.cluster_id, c.region_id, c.cluster_name, c.cluster_kind, c.cluster_summary,
         (1 - (ce.embedding <=> :query_embedding)) AS semantic_score,
         ((1 - (ce.embedding <=> :query_embedding)) * 0.70 + c.importance * 0.15 + c.confidence * 0.15) AS final_score
  FROM memory_clusters c
  JOIN memory_cluster_embeddings ce ON ce.cluster_id = c.cluster_id
  WHERE c.namespace = :namespace
    AND c.specialist_id = :specialist_id
    AND c.region_id = ANY(:region_ids)
    AND c.status = 'active'
)
SELECT * FROM scored_clusters ORDER BY final_score DESC LIMIT :top_clusters;

-- Tier 3: exact records
WITH candidate_records AS (
  SELECT mr.record_id, mr.record_kind, mr.title, mr.summary, mr.body, mr.status,
         mr.importance, mr.confidence, mcm.membership_weight,
         (1 - (me.embedding <=> :query_embedding)) AS semantic_score
  FROM memory_cluster_members mcm
  JOIN memory_records mr ON mr.record_id = mcm.record_id
  JOIN memory_embeddings me ON me.record_id = mr.record_id
  WHERE mcm.namespace = :namespace
    AND mcm.specialist_id = :specialist_id
    AND mcm.cluster_id = ANY(:cluster_ids)
    AND mr.is_soft_deleted = false
    AND mr.status IN ('active', 'staged', 'deprecated', 'contradicted')
), ranked_records AS (
  SELECT *,
         ((semantic_score * 0.55 + importance * 0.15 + confidence * 0.15 + membership_weight * 0.15) *
          CASE status WHEN 'active' THEN 1.00 WHEN 'staged' THEN 0.85 WHEN 'deprecated' THEN 0.75 WHEN 'contradicted' THEN 0.60 ELSE 0.50 END) AS final_score
  FROM candidate_records
)
SELECT * FROM ranked_records ORDER BY final_score DESC LIMIT :top_records;
