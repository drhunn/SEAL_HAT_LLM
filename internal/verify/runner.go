package verify

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/drhunn/SEAL_HAT_LLM/internal/config"
	"github.com/drhunn/SEAL_HAT_LLM/internal/db"
	"github.com/drhunn/SEAL_HAT_LLM/internal/den"
	"github.com/drhunn/SEAL_HAT_LLM/internal/execution"
	"github.com/drhunn/SEAL_HAT_LLM/internal/growth"
	"github.com/drhunn/SEAL_HAT_LLM/internal/memory"
	"github.com/drhunn/SEAL_HAT_LLM/internal/modality"
	"github.com/drhunn/SEAL_HAT_LLM/internal/modelhost"
	"github.com/drhunn/SEAL_HAT_LLM/internal/routing"
	rt "github.com/drhunn/SEAL_HAT_LLM/internal/runtime"
	"github.com/drhunn/SEAL_HAT_LLM/internal/seal"
	"github.com/drhunn/SEAL_HAT_LLM/internal/slots"
	"github.com/drhunn/SEAL_HAT_LLM/internal/telemetry"
)

type Mode string

const (
	ModeSoft   Mode = "soft"
	ModeStrict Mode = "strict"
)

func ParseMode(raw string) (Mode, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", string(ModeSoft):
		return ModeSoft, nil
	case string(ModeStrict):
		return ModeStrict, nil
	default:
		return "", fmt.Errorf("unsupported verify mode %q (expected soft or strict)", raw)
	}
}

type proposalStore struct {
	proposals []seal.AdaptationProposal
}

func (s *proposalStore) WriteProposal(_ context.Context, proposal seal.AdaptationProposal) error {
	s.proposals = append(s.proposals, proposal)
	return nil
}

func (s *proposalStore) Count() int {
	return len(s.proposals)
}

type growthPlanStore struct {
	plans []den.GrowthPlan
}

func (s *growthPlanStore) WriteGrowthPlan(_ context.Context, plan den.GrowthPlan) error {
	s.plans = append(s.plans, plan)
	return nil
}

func (s *growthPlanStore) Count() int {
	return len(s.plans)
}

type Runner struct {
	cfg    *config.AppConfig
	mode   Mode
	logger *slog.Logger
	store  *memory.PostgresStore
}

func NewRunner(cfg *config.AppConfig, mode Mode, logger *slog.Logger) *Runner {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	return &Runner{cfg: cfg, mode: mode, logger: logger}
}

func (r *Runner) Run(ctx context.Context) error {
	r.logger.Info("verify mode selected", "verify_mode", string(r.mode))
	if err := r.openStore(ctx); err != nil {
		return err
	}
	if err := r.verifySlotsAndBundle(ctx); err != nil {
		return err
	}
	result, err := r.verifyRuntimeTask(ctx)
	if err != nil {
		return err
	}
	if err := r.verifyAdaptationLoop(ctx, result); err != nil {
		return err
	}
	if err := r.verifyAbilityGrowth(ctx, result); err != nil {
		return err
	}
	r.logger.Info("verify complete", "verify_mode", string(r.mode))
	return nil
}

func (r *Runner) openStore(ctx context.Context) error {
	pool, err := db.Open(ctx, r.cfg.Database.DSN)
	if err != nil {
		return fmt.Errorf("database open failed: %w", err)
	}
	r.store = memory.NewPostgresStore(pool, r.logger)
	if err := r.store.Ping(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	r.logger.Info("database ping ok")
	return nil
}

func (r *Runner) verifySlotsAndBundle(ctx context.Context) error {
	loader := slots.NewFilesystemLoader(r.cfg.Runtime.SlotsRoot)
	slotFiles, err := loader.LoadSpecialistSlots(r.cfg.Runtime.SpecialistID)
	if err != nil {
		return fmt.Errorf("slot load failed: %w", err)
	}
	r.logger.Info("slot load ok", "specialist_id", r.cfg.Runtime.SpecialistID, "count", len(slotFiles))

	compiler := slots.NewCompiler()
	bundle, err := compiler.CompileSpecialist(r.cfg.Runtime.SpecialistID, slotFiles)
	if err != nil {
		return fmt.Errorf("slot bundle compile failed: %w", err)
	}
	bundleBytes, err := bundle.EncodeTOML()
	if err != nil {
		return fmt.Errorf("slot bundle TOML encode failed: %w", err)
	}
	summary := bundle.Summary()
	r.logger.Info("slot bundle compile ok",
		"specialist_id", r.cfg.Runtime.SpecialistID,
		"schema", bundle.Schema,
		"core_skill_count", summary.CoreSkillCount,
		"playbook_count", summary.PlaybookCount,
		"allowed_tool_count", summary.AllowedToolCount,
		"memory_pointer_count", summary.MemoryPointerCount,
		"source_map_count", summary.SourceMapCount,
		"encoded_bytes", len(bundleBytes),
	)
	versionLabel := time.Now().UTC().Format("20060102T150405Z")
	if err := r.store.PersistSlotBundleVersion(ctx, r.cfg.Runtime.SpecialistID, bundleBytes, bundle.SourceMap, versionLabel, r.cfg.Harness.DefaultCreatedBy); err != nil {
		if fatalErr := r.handleOptionalFailure("slot bundle persistence unavailable", err, "version_label", versionLabel); fatalErr != nil {
			return fatalErr
		}
	} else {
		r.logger.Info("slot bundle persisted", "version_label", versionLabel)
	}

	if _, err := r.store.HealthSnapshot(ctx, r.cfg.Runtime.SpecialistID); err != nil {
		if fatalErr := r.handleOptionalFailure("health snapshot unavailable", err, "specialist_id", r.cfg.Runtime.SpecialistID); fatalErr != nil {
			return fatalErr
		}
	} else {
		r.logger.Info("health snapshot call ok")
	}
	return nil
}

func (r *Runner) verifyRuntimeTask(ctx context.Context) (*rt.TaskResult, error) {
	routingService := routing.NewService(r.logger)
	hostRegistry := modelhost.NewSimulatedRegistry("verify")
	executionService := execution.NewService(r.logger, hostRegistry)
	taskProcessor := rt.NewService(r.cfg, nil, r.store, nil, routingService, executionService, nil, nil, r.logger)

	verifyTask := rt.Task{
		ID:      "verify-task",
		Summary: "verify bounded runtime task",
		Class:   "analysis",
		Prompt:  "Verify runtime task processing.",
	}
	if r.cfg.Runtime.EnableMultimodalSmokeTest {
		verifyTask = rt.Task{
			ID:                          "verify-task",
			Summary:                     "verify multimodal execution",
			Class:                       "evidence_fusion",
			PrimaryModality:             modality.Image,
			SecondaryModalities:         []modality.Type{modality.Normalize(r.cfg.Runtime.DefaultPrimaryModality)},
			CrossModalGroundingRequired: true,
			AllowTextOnlyFallback:       r.cfg.Runtime.AllowTextOnlyFallback,
			AssetRefs:                   []string{"sandbox://verify/image-1"},
			Prompt:                      "Verify image and text fusion.",
		}
	}

	result, err := taskProcessor.ProcessTask(ctx, verifyTask)
	if err != nil {
		return nil, fmt.Errorf("runtime task verify failed: %w", err)
	}
	for _, warning := range result.Warnings {
		if fatalErr := r.handleOptionalFailure("runtime task persistence warning", fmt.Errorf(warning), "task_id", result.Task.ID); fatalErr != nil {
			return nil, fatalErr
		}
	}
	r.logger.Info("runtime task verify ok",
		"task_id", result.Task.ID,
		"executor", result.ExecutionResult.Plan.ChosenExecutor,
		"execution_mode", result.ExecutionResult.Plan.ExecutionMode,
		"host", result.ExecutionResult.HostResult.HostName,
		"signal_count", len(result.Signals),
		"warning_count", len(result.Warnings),
	)
	return &result, nil
}

func (r *Runner) verifyAdaptationLoop(ctx context.Context, result *rt.TaskResult) error {
	signalStore := telemetry.NewMemoryStore()
	if err := signalStore.WriteSignals(ctx, result.Signals); err != nil {
		return fmt.Errorf("verify signal store write failed: %w", err)
	}

	proposalWriter := &proposalStore{}
	growthWriter := &growthPlanStore{}
	sealService := seal.NewService(signalStore, proposalWriter, r.logger)
	proposals, err := sealService.ReviewSpecialist(ctx, r.cfg.Runtime.SpecialistID, 100)
	if err != nil {
		return fmt.Errorf("seal review failed: %w", err)
	}
	for _, proposal := range proposals {
		if err := r.store.WriteProposal(ctx, proposal); err != nil {
			if fatalErr := r.handleOptionalFailure("adaptation proposal persistence unavailable", err, "proposal_id", proposal.ID); fatalErr != nil {
				return fatalErr
			}
		}
	}
	r.logger.Info("seal verify ok", "signal_count", len(result.Signals), "proposal_count", len(proposals), "persisted_proposal_count", proposalWriter.Count())

	if len(proposals) == 0 {
		return nil
	}

	denService := den.NewService(growthWriter, r.logger)
	plan, err := denService.PlanFromProposal(ctx, proposals[0])
	if err != nil {
		return fmt.Errorf("den growth planning failed: %w", err)
	}
	if err := r.store.WriteGrowthPlan(ctx, *plan); err != nil {
		if fatalErr := r.handleOptionalFailure("growth plan persistence unavailable", err, "growth_plan_id", plan.ID); fatalErr != nil {
			return fatalErr
		}
	}
	r.logger.Info("den verify ok", "proposal_id", proposals[0].ID, "growth_surface", plan.Surface, "growth_plan_count", growthWriter.Count())
	return nil
}

func (r *Runner) verifyAbilityGrowth(ctx context.Context, result *rt.TaskResult) error {
	growthService := growth.NewService(r.store, r.logger)
	growthResult, err := growthService.StageExperiment(ctx, r.cfg.Runtime.Namespace, r.cfg.Runtime.SpecialistID, growth.Assessment{
		AbilityName:       "multimodal_grounding",
		GapSummary:        "Verify governed ability-growth storage and staging.",
		EvidenceSummary:   fmt.Sprintf("bounded verify task executor=%s host=%s signals=%d warnings=%d", result.ExecutionResult.Plan.ChosenExecutor, result.ExecutionResult.HostResult.HostName, len(result.Signals), len(result.Warnings)),
		TriedMemoryFix:    true,
		TriedRoutingFix:   true,
		TriedPromptFix:    true,
		PreferredSurface:  "modality_branch",
		RequestedBy:       r.cfg.Harness.DefaultCreatedBy,
		ParentApprovedBy:  "parent:verify",
		HarnessVerifiedBy: r.cfg.Harness.DefaultCreatedBy,
		Notes:             "verify governed ability-growth path",
	})
	if err != nil {
		return fmt.Errorf("ability-growth verify failed: %w", err)
	}
	r.logger.Info("ability-growth verify ok",
		"status", growthResult.Status,
		"experiment_id", growthResult.ExperimentID,
	)
	return nil
}

func (r *Runner) handleOptionalFailure(message string, err error, attrs ...any) error {
	if err == nil {
		return nil
	}
	attrs = append(attrs, "verify_mode", string(r.mode), "err", err)
	if r.mode == ModeStrict {
		r.logger.Error(message, attrs...)
		return fmt.Errorf("%s: %w", message, err)
	}
	r.logger.Warn(message, attrs...)
	return nil
}
