package runtime

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/drhunn/LLM-plus-harness/internal/config"
	"github.com/drhunn/LLM-plus-harness/internal/harness"
	"github.com/drhunn/LLM-plus-harness/internal/memory"
	"github.com/drhunn/LLM-plus-harness/internal/slots"
)

type Service struct {
	cfg     *config.AppConfig
	loader  *slots.FilesystemLoader
	store   *memory.PostgresStore
	harness *harness.Service
	logger  *slog.Logger
}

func NewService(cfg *config.AppConfig, loader *slots.FilesystemLoader, store *memory.PostgresStore, harnessService *harness.Service, logger *slog.Logger) *Service {
	return &Service{cfg: cfg, loader: loader, store: store, harness: harnessService, logger: logger}
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

	s.logger.Info("runtime initialized",
		"specialist_id", s.cfg.Runtime.SpecialistID,
		"namespace", s.cfg.Runtime.Namespace,
		"slot_count", len(slotFiles),
	)

	<-ctx.Done()
	s.logger.Info("shutdown requested")
	return nil
}
