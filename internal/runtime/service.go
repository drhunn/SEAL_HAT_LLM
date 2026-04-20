package runtime

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/drhunn/SEAL_HAT_LLM/internal/config"
	"github.com/drhunn/SEAL_HAT_LLM/internal/execution"
	"github.com/drhunn/SEAL_HAT_LLM/internal/growth"
	"github.com/drhunn/SEAL_HAT_LLM/internal/harness"
	"github.com/drhunn/SEAL_HAT_LLM/internal/memory"
	"github.com/drhunn/SEAL_HAT_LLM/internal/modality"
	"github.com/drhunn/SEAL_HAT_LLM/internal/routing"
	"github.com/drhunn/SEAL_HAT_LLM/internal/slotpacket"
	"github.com/drhunn/SEAL_HAT_LLM/internal/slots"
	"github.com/drhunn/SEAL_HAT_LLM/internal/slotsync"
)

type Service struct {
	cfg       *config.AppConfig
	loader    *slots.FilesystemLoader
	store     *memory.PostgresStore
	harness   *harness.Service
	routing   *routing.Service
	execution *execution.Service
	growth    *growth.Service
	slotSync  *slotsync.Service
	logger    *slog.Logger
}

func NewService(cfg *config.AppConfig, loader *slots.FilesystemLoader, store *memory.PostgresStore, harnessService *harness.Service, routingService *routing.Service, executionService *execution.Service, growthService *growth.Service, slotSyncService *slotsync.Service, logger *slog.Logger) *Service {
	return &Service{
		cfg:       cfg,
		loader:    loader,
		store:     store,
		harness:   harnessService,
		routing:   routingService,
		execution: executionService,
		growth:    growthService,
		slotSync:  slotSyncService,
		logger:    logger,
	}
}

func (s *Service) Start(ctx context.Context) error {
	runCtx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.Runtime.RequestTimeoutSeconds)*time.Second)
	defer cancel()

	if err := s.harness.StartupChecks(runCtx, s.cfg.Runtime.SpecialistID, s.cfg.Runtime.Namespace); err != nil {
		return fmt.Errorf("startup checks: %w", err)
	}

	slotFiles, err := s.loader.LoadSpecialistSlots(s.cfg.Runtime.SpecialistID)
	if err != nil {
		return fmt.Errorf("load specialist slots: %w", err)
	}

	if s.slotSync != nil {
		if err := s.slotSync.SyncFilesystemView(runCtx, s.cfg.Runtime.SpecialistID); err != nil {
			s.logger.Warn("slot sync failed", "specialist_id", s.cfg.Runtime.SpecialistID, "err", err)
		}
	}

	packet := slotpacket.BuildFromUnknown(s.cfg.Runtime.SpecialistID, "active", true, slotFiles)
	if err := slotpacket.WriteJSON("artifacts/slot_packet.json", packet); err != nil {
		s.logger.Warn("slot packet export failed", "specialist_id", s.cfg.Runtime.SpecialistID, "err", err)
	} else {
		s.logger.Info("slot packet exported",
			"specialist_id", s.cfg.Runtime.SpecialistID,
			"path", "artifacts/slot_packet.json",
			"version_hash", packet.VersionHash,
		)
	}

	if s.routing != nil {
		decision := s.routing.DecideTask(runCtx, routing.Input{
			TaskSummary:     "startup routing smoke test",
			TaskClass:       "governance",
			PrimaryModality: modality.Normalize(s.cfg.Runtime.DefaultPrimaryModality),
		})
		if _, err := s.store.CreateRoutingAudit(runCtx, memory.RoutingAuditInput{
			TaskID:                "startup-smoke",
			RoutedBy:              "harness:runtime",
			InitialClassifier:     "startup_probe",
			TaskSummary:           decision.TaskSummary,
			TaskClass:             decision.TaskClass,
			ChosenTarget:          decision.ChosenTarget,
			Confidence:            decision.Confidence,
			Impact:                "low",
			WasFallback:           decision.WasFallback,
			FallbackReason:        decision.FallbackReason,
			WasOverride:           false,
			OverrideBy:            "",
			MultiSpecialistReview: decision.RequiresFusion,
			Notes:                 "startup routing smoke test",
		}); err != nil {
			s.logger.Warn("routing audit write failed", "err", err)
		}
	}

	if s.execution != nil && s.cfg.Runtime.EnableMultimodalSmokeTest {
		req := execution.Request{
			TaskSummary:                 "startup multimodal execution smoke test",
			TaskClass:                   "evidence_fusion",
			PrimaryModality:             modality.Image,
			SecondaryModalities:         []modality.Type{modality.Text},
			CrossModalGroundingRequired: true,
			AllowTextOnlyFallback:       s.cfg.Runtime.AllowTextOnlyFallback,
			AssetRefs:                   []string{"sandbox://startup-smoke/image-1"},
			Prompt:                      "Compare image evidence with text context.",
		}
		result, err := s.execution.Execute(runCtx, req)
		if err != nil {
			s.logger.Warn("multimodal execution smoke test failed", "err", err)
		} else {
			s.logger.Info("multimodal execution smoke test ok",
				"executor", result.Plan.ChosenExecutor,
				"execution_mode", result.Plan.ExecutionMode,
				"host", result.HostResult.HostName,
				"handled", result.HostResult.Handled,
			)
			persistInput := memory.MultimodalExecutionInput{
				Namespace:       s.cfg.Runtime.Namespace,
				SpecialistID:    s.cfg.Runtime.SpecialistID,
				TaskSummary:     req.TaskSummary,
				PrimaryModality: req.PrimaryModality.String(),
				ExecutionMode:   result.Plan.ExecutionMode,
				Executor:        result.Plan.ChosenExecutor,
				HostName:        result.HostResult.HostName,
				Output:          result.HostResult.Output,
				RequiresFusion:  result.Plan.RequiresFusion,
				AssetURIs:       req.AssetRefs,
				CreatedBy:       s.cfg.Harness.DefaultCreatedBy,
			}
			if err := s.store.PersistMultimodalExecution(runCtx, persistInput); err != nil {
				s.logger.Warn("persist multimodal execution failed", "err", err)
			}
		}
	}

	if s.growth != nil {
		growthResult, err := s.growth.StageExperiment(runCtx, s.cfg.Runtime.Namespace, s.cfg.Runtime.SpecialistID, growth.Assessment{
			AbilityName:       "multimodal_grounding",
			GapSummary:        "Persistent need for stronger multimodal grounding beyond current text-first execution scaffolding.",
			EvidenceSummary:   "Multimodal routing/execution exists, but live modality backends and deeper grounded retrieval are still immature.",
			TriedMemoryFix:    true,
			TriedRoutingFix:   true,
			TriedPromptFix:    true,
			PreferredSurface:  "modality_branch",
			RequestedBy:       s.cfg.Harness.DefaultCreatedBy,
			ParentApprovedBy:  "parent:startup-smoke",
			HarnessVerifiedBy: s.cfg.Harness.DefaultCreatedBy,
			Notes:             "startup governed ability-growth smoke path",
		})
		if err != nil {
			s.logger.Warn("ability growth smoke path failed", "err", err)
		} else {
			s.logger.Info("ability growth smoke path ok",
				"status", growthResult.Status,
				"experiment_id", growthResult.ExperimentID,
			)
		}
	}

	results, err := s.store.RunCoarseToFineSearch(runCtx, s.cfg.Runtime.Namespace, s.cfg.Runtime.SpecialistID, memory.ZeroVector(1536), 3, 5, 5)
	if err != nil {
		s.logger.Warn("coarse-to-fine retrieval smoke test failed", "err", err)
	} else {
		s.logger.Info("coarse-to-fine retrieval smoke test ok", "result_count", len(results))
	}

	s.logger.Info("runtime initialized",
		"specialist_id", s.cfg.Runtime.SpecialistID,
		"namespace", s.cfg.Runtime.Namespace,
		"slot_count", len(slotFiles),
	)

	<-ctx.Done()
	s.logger.Info("shutdown requested")
	return nil
}
