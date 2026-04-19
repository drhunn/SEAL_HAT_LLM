package slotsync

import (
	"context"
	"log/slog"

	"github.com/drhunn/SEAL_HAT_LLM/internal/slots"
)

type Service struct {
	loader   *slots.FilesystemLoader
	compiler *slots.Compiler
	logger   *slog.Logger
}

func NewService(loader *slots.FilesystemLoader, logger *slog.Logger) *Service {
	return &Service{loader: loader, compiler: slots.NewCompiler(), logger: logger}
}

func (s *Service) CompileBundle(ctx context.Context, specialistID string) (*slots.Bundle, []slots.File, error) {
	files, err := s.loader.LoadSpecialistSlots(specialistID)
	if err != nil {
		return nil, nil, err
	}

	bundle, err := s.compiler.CompileSpecialist(specialistID, files)
	if err != nil {
		return nil, files, err
	}

	s.logger.DebugContext(ctx, "compiled slot bundle", "specialist_id", specialistID, "schema", bundle.Schema, "source_map_count", len(bundle.SourceMap))
	return bundle, files, nil
}

func (s *Service) SyncFilesystemView(ctx context.Context, specialistID string) error {
	bundle, files, err := s.CompileBundle(ctx, specialistID)
	if err != nil {
		return err
	}

	summary := bundle.Summary()
	s.logger.InfoContext(ctx, "slot sync snapshot",
		"specialist_id", specialistID,
		"slot_count", len(files),
		"core_skill_count", summary.CoreSkillCount,
		"playbook_count", summary.PlaybookCount,
		"allowed_tool_count", summary.AllowedToolCount,
		"memory_pointer_count", summary.MemoryPointerCount,
		"source_map_count", summary.SourceMapCount,
	)
	return nil
}
