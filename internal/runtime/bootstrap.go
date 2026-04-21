package runtime

import (
	"context"
	"fmt"
	"time"

	"github.com/drhunn/SEAL_HAT_LLM/internal/config"
	"github.com/drhunn/SEAL_HAT_LLM/internal/growth"
	"github.com/drhunn/SEAL_HAT_LLM/internal/modality"
	"github.com/drhunn/SEAL_HAT_LLM/internal/slotpacket"
)

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
	s.persistSpecialistArtifact(runCtx, "artifacts/slot_packet.json", packet.VersionHash)

	s.logger.Info("runtime initialized",
		"specialist_id", s.cfg.Runtime.SpecialistID,
		"namespace", s.cfg.Runtime.Namespace,
		"slot_count", len(slotFiles),
	)

	if s.cfg.Runtime.EnableMultimodalSmokeTest && s.routing != nil && s.execution != nil {
		if err := s.runStartupTask(runCtx); err != nil {
			return err
		}
	}

	if s.cfg.Runtime.EnableTaskInbox {
		inbox := NewTaskInbox(s.cfg.Runtime.TaskInboxDir, s.logger)
		if inbox == nil {
			return fmt.Errorf("task inbox enabled but runtime.task_inbox_dir is empty")
		}
		if err := inbox.EnsureDirs(); err != nil {
			return fmt.Errorf("ensure task inbox: %w", err)
		}
		s.logger.Info("task inbox enabled",
			"dir", s.cfg.Runtime.TaskInboxDir,
			"poll_interval_seconds", s.cfg.Runtime.TaskPollIntervalSeconds,
		)
		return s.runTaskInbox(ctx, inbox)
	}

	<-ctx.Done()
	s.logger.Info("shutdown requested")
	return nil
}

func (s *Service) runStartupTask(ctx context.Context) error {
	startTask := defaultStartupTask(s.cfg)
	result, err := s.ProcessTask(ctx, startTask)
	if err != nil {
		s.persistFailedRouteEpisode(ctx, startTask, err)
		return fmt.Errorf("process startup task: %w", err)
	}
	s.persistRouteEpisode(ctx, result)
	if s.growth == nil {
		return nil
	}
	growthResult, err := s.growth.StageExperiment(ctx, s.cfg.Runtime.Namespace, s.cfg.Runtime.SpecialistID, growth.Assessment{
		AbilityName:       "multimodal_grounding",
		GapSummary:        "Bounded startup task still relies on simulated multimodal execution rather than live grounded backends.",
		EvidenceSummary:   fmt.Sprintf("task=%s executor=%s host=%s retrieval_results=%d signals=%d", result.Task.Summary, result.ExecutionResult.Plan.ChosenExecutor, result.ExecutionResult.HostResult.HostName, len(result.RetrievalResults), len(result.Signals)),
		TriedMemoryFix:    true,
		TriedRoutingFix:   true,
		TriedPromptFix:    true,
		PreferredSurface:  "modality_branch",
		RequestedBy:       s.cfg.Harness.DefaultCreatedBy,
		ParentApprovedBy:  "parent:startup-task",
		HarnessVerifiedBy: s.cfg.Harness.DefaultCreatedBy,
		Notes:             "startup bounded task path",
	})
	if err != nil {
		s.logger.Warn("ability growth task follow-up failed", "err", err)
		return nil
	}
	s.logger.Info("ability growth task follow-up ok",
		"status", growthResult.Status,
		"experiment_id", growthResult.ExperimentID,
	)
	return nil
}

func defaultStartupTask(cfg *config.AppConfig) Task {
	return normalizeTask(Task{
		ID:                          "startup-task",
		Summary:                     "startup multimodal execution task",
		Class:                       "evidence_fusion",
		PrimaryModality:             modality.Image,
		SecondaryModalities:         []modality.Type{modality.Text},
		CrossModalGroundingRequired: true,
		AllowTextOnlyFallback:       cfg.Runtime.AllowTextOnlyFallback,
		PreferredUnitID:             cfg.Runtime.SpecialistID,
		AssetRefs:                   []string{"sandbox://startup-smoke/image-1"},
		Prompt:                      "Compare image evidence with text context.",
	}, cfg)
}

func (s *Service) runTaskInbox(ctx context.Context, inbox *TaskInbox) error {
	if err := s.processAvailableTasks(ctx, inbox); err != nil {
		return err
	}

	ticker := time.NewTicker(time.Duration(s.cfg.Runtime.TaskPollIntervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("task inbox shutdown requested")
			return nil
		case <-ticker.C:
			if err := s.processAvailableTasks(ctx, inbox); err != nil {
				return err
			}
		}
	}
}

func (s *Service) processAvailableTasks(ctx context.Context, inbox *TaskInbox) error {
	for {
		queued, err := inbox.ClaimNext()
		if err != nil {
			if queued != nil {
				if markErr := inbox.MarkFailed(queued, err); markErr != nil {
					return fmt.Errorf("mark failed inbox task: %w", markErr)
				}
				continue
			}
			return err
		}
		if queued == nil {
			return nil
		}

		taskCtx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.Runtime.RequestTimeoutSeconds)*time.Second)
		result, taskErr := s.ProcessTask(taskCtx, queued.Task)
		cancel()
		if taskErr != nil {
			s.persistFailedRouteEpisode(ctx, queued.Task, taskErr)
			if markErr := inbox.MarkFailed(queued, taskErr); markErr != nil {
				return fmt.Errorf("mark task as failed: %w", markErr)
			}
			continue
		}
		s.persistRouteEpisode(ctx, result)
		if err := inbox.MarkProcessed(queued, result); err != nil {
			return fmt.Errorf("mark task as processed: %w", err)
		}
	}
}
