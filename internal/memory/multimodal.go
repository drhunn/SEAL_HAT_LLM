package memory

import (
	"context"
	"fmt"
	"strings"
)

type MultimodalAssetInput struct {
	Namespace    string
	SpecialistID string
	Modality     string
	URI          string
	MIMEType     string
	MetadataJSON string
	CreatedBy    string
}

type MultimodalExecutionInput struct {
	Namespace       string
	SpecialistID    string
	TaskSummary     string
	PrimaryModality string
	ExecutionMode   string
	Executor        string
	HostName        string
	Output          string
	RequiresFusion  bool
	AssetURIs       []string
	CreatedBy       string
}

func (s *PostgresStore) RegisterMultimodalAsset(ctx context.Context, in MultimodalAssetInput) (string, error) {
	const q = `
		INSERT INTO agent_core.multimodal_assets (
			namespace,
			specialist_id,
			modality,
			uri,
			mime_type,
			metadata_json,
			created_by
		) VALUES ($1,$2,$3,$4,$5,$6::jsonb,$7)
		RETURNING asset_id::text`
	metadata := in.MetadataJSON
	if strings.TrimSpace(metadata) == "" {
		metadata = `{}`
	}
	var id string
	if err := s.db.QueryRow(ctx, q, in.Namespace, in.SpecialistID, in.Modality, in.URI, in.MIMEType, metadata, in.CreatedBy).Scan(&id); err != nil {
		return "", fmt.Errorf("register multimodal asset: %w", err)
	}
	return id, nil
}

func (s *PostgresStore) CreateExecutionArtifactRecord(ctx context.Context, in MultimodalExecutionInput) (string, error) {
	const q = `
		INSERT INTO agent_core.memory_records (
			namespace,
			specialist_id,
			record_kind,
			title,
			summary,
			body,
			status,
			importance,
			confidence,
			tags,
			embedding_text,
			created_by
		) VALUES (
			$1,$2,'artifact_reference',$3,$4,$5,'staged',$6,$7,$8,$9,$10
		)
		RETURNING record_id::text`

	title := fmt.Sprintf("Execution artifact: %s", in.TaskSummary)
	summary := fmt.Sprintf("executor=%s host=%s mode=%s primary_modality=%s", in.Executor, in.HostName, in.ExecutionMode, in.PrimaryModality)
	body := fmt.Sprintf(
		"task_summary: %s\nprimary_modality: %s\nexecution_mode: %s\nexecutor: %s\nhost_name: %s\nrequires_fusion: %t\nasset_uris: %s\noutput: %s",
		in.TaskSummary,
		in.PrimaryModality,
		in.ExecutionMode,
		in.Executor,
		in.HostName,
		in.RequiresFusion,
		strings.Join(in.AssetURIs, ", "),
		in.Output,
	)
	tags := []string{"multimodal", in.PrimaryModality, in.ExecutionMode, in.Executor}
	var id string
	if err := s.db.QueryRow(ctx, q, in.Namespace, in.SpecialistID, title, summary, body, 0.60, 0.70, tags, in.Output, in.CreatedBy).Scan(&id); err != nil {
		return "", fmt.Errorf("create execution artifact record: %w", err)
	}
	return id, nil
}

func (s *PostgresStore) LinkRecordToAsset(ctx context.Context, recordID, assetID, relationType string) error {
	const q = `
		INSERT INTO agent_core.memory_record_assets (
			record_id,
			asset_id,
			relation_type
		) VALUES ($1::uuid,$2::uuid,$3)
		ON CONFLICT (record_id, asset_id, relation_type) DO NOTHING`
	if relationType == "" {
		relationType = "source"
	}
	if _, err := s.db.Exec(ctx, q, recordID, assetID, relationType); err != nil {
		return fmt.Errorf("link record to asset: %w", err)
	}
	return nil
}

func (s *PostgresStore) PersistMultimodalExecution(ctx context.Context, in MultimodalExecutionInput) error {
	recordID, err := s.CreateExecutionArtifactRecord(ctx, in)
	if err != nil {
		return err
	}

	for _, uri := range in.AssetURIs {
		assetID, err := s.RegisterMultimodalAsset(ctx, MultimodalAssetInput{
			Namespace:    in.Namespace,
			SpecialistID: in.SpecialistID,
			Modality:     in.PrimaryModality,
			URI:          uri,
			MIMEType:     "application/octet-stream",
			MetadataJSON: fmt.Sprintf(`{"executor":%q,"host":%q}`, in.Executor, in.HostName),
			CreatedBy:    in.CreatedBy,
		})
		if err != nil {
			return err
		}
		if err := s.LinkRecordToAsset(ctx, recordID, assetID, "source"); err != nil {
			return err
		}
	}

	return nil
}
