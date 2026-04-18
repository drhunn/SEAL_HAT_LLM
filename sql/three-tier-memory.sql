SET search_path TO agent_core, public;

CREATE TABLE IF NOT EXISTS memory_regions (
  region_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  namespace text NOT NULL,
  specialist_id text NOT NULL REFERENCES specialists(specialist_id) ON DELETE CASCADE,
  region_name text NOT NULL,
  region_kind text NOT NULL,
  region_summary text NOT NULL,
  region_tags text[] NOT NULL DEFAULT '{}',
  status memory_status NOT NULL DEFAULT 'active',
  importance numeric(5,4) NOT NULL DEFAULT 0.5,
  confidence numeric(5,4) NOT NULL DEFAULT 0.5,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  created_by text NOT NULL,
  UNIQUE (namespace, region_name)
);

CREATE TABLE IF NOT EXISTS memory_region_embeddings (
  region_id uuid PRIMARY KEY REFERENCES memory_regions(region_id) ON DELETE CASCADE,
  namespace text NOT NULL,
  specialist_id text NOT NULL REFERENCES specialists(specialist_id) ON DELETE CASCADE,
  embedding_model text NOT NULL,
  embedding vector(1536) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_memory_region_embeddings_hnsw ON memory_region_embeddings USING hnsw (embedding vector_cosine_ops);

CREATE TABLE IF NOT EXISTS memory_clusters (
  cluster_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  namespace text NOT NULL,
  specialist_id text NOT NULL REFERENCES specialists(specialist_id) ON DELETE CASCADE,
  region_id uuid NOT NULL REFERENCES memory_regions(region_id) ON DELETE CASCADE,
  cluster_name text NOT NULL,
  cluster_kind text NOT NULL,
  cluster_summary text NOT NULL,
  cluster_tags text[] NOT NULL DEFAULT '{}',
  status memory_status NOT NULL DEFAULT 'active',
  importance numeric(5,4) NOT NULL DEFAULT 0.5,
  confidence numeric(5,4) NOT NULL DEFAULT 0.5,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  created_by text NOT NULL,
  UNIQUE (region_id, cluster_name)
);

CREATE TABLE IF NOT EXISTS memory_cluster_embeddings (
  cluster_id uuid PRIMARY KEY REFERENCES memory_clusters(cluster_id) ON DELETE CASCADE,
  namespace text NOT NULL,
  specialist_id text NOT NULL REFERENCES specialists(specialist_id) ON DELETE CASCADE,
  embedding_model text NOT NULL,
  embedding vector(1536) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_memory_cluster_embeddings_hnsw ON memory_cluster_embeddings USING hnsw (embedding vector_cosine_ops);

CREATE TABLE IF NOT EXISTS memory_cluster_members (
  cluster_member_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  namespace text NOT NULL,
  specialist_id text NOT NULL REFERENCES specialists(specialist_id) ON DELETE CASCADE,
  cluster_id uuid NOT NULL REFERENCES memory_clusters(cluster_id) ON DELETE CASCADE,
  record_id uuid NOT NULL REFERENCES memory_records(record_id) ON DELETE CASCADE,
  membership_weight numeric(5,4) NOT NULL DEFAULT 1.0,
  is_representative boolean NOT NULL DEFAULT false,
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (cluster_id, record_id)
);
