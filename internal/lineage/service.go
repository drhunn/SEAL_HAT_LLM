package lineage

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

type NodeType string

type EdgeType string

const (
	NodeParentModel NodeType = "parent_model"
	NodeSpecialist  NodeType = "specialist"
	NodeAdapter     NodeType = "adapter"
	NodeBranch      NodeType = "branch"
	NodeExpertBank  NodeType = "expert_bank"
)

const (
	EdgeSpawned     EdgeType = "spawned"
	EdgeSplitInto   EdgeType = "split_into"
	EdgeTunedInto   EdgeType = "tuned_into"
	EdgeDerivedFrom EdgeType = "derived_from"
	EdgeRetiredInto EdgeType = "retired_into"
)

type Node struct {
	ID         string
	NodeType   NodeType
	ExternalID string
	Status     string
	Metadata   map[string]string
	CreatedAt  time.Time
}

type Edge struct {
	ID        string
	FromNode  string
	ToNode    string
	EdgeType  EdgeType
	Metadata  map[string]string
	CreatedAt time.Time
}

type Store interface {
	UpsertNode(ctx context.Context, node Node) error
	CreateEdge(ctx context.Context, edge Edge) error
}

type Service struct {
	store  Store
	logger *slog.Logger
}

func NewService(store Store, logger *slog.Logger) *Service {
	return &Service{store: store, logger: logger}
}

func (s *Service) RecordDerivedSpecialist(ctx context.Context, sourceExternalID, targetExternalID string, metadata map[string]string) error {
	if s.store == nil {
		return fmt.Errorf("lineage store is required")
	}
	sourceExternalID = strings.TrimSpace(sourceExternalID)
	targetExternalID = strings.TrimSpace(targetExternalID)
	if sourceExternalID == "" || targetExternalID == "" {
		return fmt.Errorf("source and target external ids are required")
	}
	now := time.Now().UTC()
	from := Node{ID: "node|" + sourceExternalID, NodeType: NodeSpecialist, ExternalID: sourceExternalID, Status: "active", CreatedAt: now}
	to := Node{ID: "node|" + targetExternalID, NodeType: NodeSpecialist, ExternalID: targetExternalID, Status: "active", Metadata: metadata, CreatedAt: now}
	if err := s.store.UpsertNode(ctx, from); err != nil {
		return err
	}
	if err := s.store.UpsertNode(ctx, to); err != nil {
		return err
	}
	edge := Edge{ID: fmt.Sprintf("edge|%s|%s|%s", sourceExternalID, targetExternalID, EdgeDerivedFrom), FromNode: from.ID, ToNode: to.ID, EdgeType: EdgeDerivedFrom, Metadata: metadata, CreatedAt: now}
	if err := s.store.CreateEdge(ctx, edge); err != nil {
		return err
	}
	s.logger.InfoContext(ctx, "lineage edge created", "from", sourceExternalID, "to", targetExternalID, "edge_type", EdgeDerivedFrom)
	return nil
}

func (s *Service) RecordSpecialistSplit(ctx context.Context, sourceExternalID, leftExternalID, rightExternalID string, metadata map[string]string) error {
	if err := s.RecordDerivedSpecialist(ctx, sourceExternalID, leftExternalID, metadata); err != nil {
		return err
	}
	return s.RecordDerivedSpecialist(ctx, sourceExternalID, rightExternalID, metadata)
}
