package runtime

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/drhunn/SEAL_HAT_LLM/internal/config"
	"github.com/drhunn/SEAL_HAT_LLM/internal/harness"
	"github.com/drhunn/SEAL_HAT_LLM/internal/memory"
	"github.com/drhunn/SEAL_HAT_LLM/internal/routing"
	"github.com/drhunn/SEAL_HAT_LLM/internal/slots"
	"github.com/drhunn/SEAL_HAT_LLM/internal/slotsync"
)

type Service struct {
	cfg      *config.AppConfig
	loader   *slots.FilesystemLoader
	store    *memory.PostgresStore
	harness  *harness.Service
	routing  *routing.Service
	slotSync *slotsync.Service
	logger   *slog.Logger
}

func NewService(cfg *config.AppConfig, loader *slots.FilesystemLoader, store *memory.PostgresStore, harnessService *harness.Service, routingService *routing.Service, slotSyncService *slotsync.Service, logger *slog.Logger) *Service {
	return &Service{
		cfg:      cfg,
		loader:   loader,
		store:    store,
		harness:  harnessService,
		routing:  routingService,
		slotSync: slotSyncService,
		logger:   logger,
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

	// Startup smoke tests: because "it compiled" is not an availability strategy.
	if s.routing != nil {
		decision := s.routing.Decide(runCtx, "startup routing smoke test", "governance")
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
			MultiSpecialistReview: false,
			Notes:                 "startup routing smoke test",
		}); err != nil {
			s.logger.Warn("routing audit write failed", "err", err)
		}
	}

	results, err := s.store.RunCoarseToFineSearch(runCtx, s.cfg.Runtime.Namespace, s.cfg.Runtime.SpecialistID, memory.ZeroVector(1536), 3, 5, 5)
	if err != nil {
		s.logger.Warn("coarse-to-fine retrieval smoke test failed", "err", err)
	} else {
		// Retrieval-first, transcript-hoarding last.
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
