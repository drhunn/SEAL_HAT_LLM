package slotsync

import (
	"context"
	"log/slog"

	"github.com/drhunn/SEAL_HAT_LLM/internal/slots"
)

type Service struct {
	loader *slots.FilesystemLoader
	logger *slog.Logger
}

func NewService(loader *slots.FilesystemLoader, logger *slog.Logger) *Service {
	return &Service{loader: loader, logger: logger}
}

func (s *Service) SyncFilesystemView(ctx context.Context, specialistID string) error {
	files, err := s.loader.LoadSpecialistSlots(specialistID)
	if err != nil {
		return err
	}

	s.logger.InfoContext(ctx, "slot sync snapshot", "specialist_id", specialistID, "slot_count", len(files))
	return nil
}
