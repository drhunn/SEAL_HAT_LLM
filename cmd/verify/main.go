package main

import (
	"context"
	"flag"
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

type verifyMode string

const (
	verifyModeSoft   verifyMode = "soft"
	verifyModeStrict verifyMode = "strict"
)

func parseVerifyMode(raw string) (verifyMode, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", string(verifyModeSoft):
		return verifyModeSoft, nil
	case string(verifyModeStrict):
		return verifyModeStrict, nil
	default:
		return "", fmt.Errorf("unsupported verify mode %q (expected soft or strict)", raw)
	}
}

func handleOptionalFailure(mode verifyMode, logger *slog.Logger, message string, err error, attrs ...any) error {
	if err == nil {
		return nil
	}
	attrs = append(attrs, "verify_mode", string(mode), "err", err)
	if mode == verifyModeStrict {
		logger.Error(message, attrs...)
		return fmt.Errorf("%s: %w", message, err)
	}
	logger.Warn(message, attrs...)
	return nil
}

type verifyProposalStore struct {
	proposals []seal.AdaptationProposal
}

func (s *verifyProposalStore) WriteProposal(ctx context.Context, proposal seal.AdaptationProposal) error {
	s.proposals = append(s.proposals, proposal)
	return nil
}

func (s *verifyProposalStore) Count() int {
	return len(s.proposals)
}

type verifyGrowthPlanStore struct {
	plans []den.GrowthPlan
}

func (s *verifyGrowthPlanStore) WriteGrowthPlan(ctx context.Context, plan den.GrowthPlan) error {
	s.plans = append(s.plans, plan)
	return nil
}

func (s *verifyGrowthPlanStore) Count() int {
	return len(s.plans)
}

func main() {
	cfgPath := flag.String("config", "config/runtime.example.toml", "path to runtime TOML config")
	modeFlag := flag.String("mode", string(verifyModeSoft), "verify mode: soft or strict")
	flag.Parse()

	mode, err := parseVerifyMode(*modeFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse verify mode: %v\n", err)
		os.Exit(1)
	}

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel()}))
	logger.Info("verify mode selected", "verify_mode", string(mode))
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Runtime.RequestTimeoutSeconds)*time.Second)
	defer cancel()

	pool, err := db.Open(ctx, cfg.Database.DSN)
	if err != nil {
		logger.Error("database open failed", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	store := memory.NewPostgresStore(pool, logger)
	if err := store.Ping(ctx); err != nil {
		logger.Error("database ping failed", "err", err)
		os.Exit(1)
	}
	logger.Info("database ping ok")

	loader := slots.NewFilesystemLoader(cfg.Runtime.SlotsRoot)
	slotFiles, err := loader.LoadSpecialistSlots(cfg.Runtime.SpecialistID)
	if err != nil {
		logger.Error("slot load failed", "specialist_id", cfg.Runtime.SpecialistID, "err", err)
		os.Exit(1)
	}
	logger.Info("slot load ok", "specialist_id", cfg.Runtime.SpecialistID, "count", len(slotFiles))

	compiler := slots.NewCompiler()
	bundle, err := compiler.CompileSpecialist(cfg.Runtime.SpecialistID, slotFiles)
	if err != nil {
		logger.Error("slot bundle compile failed", "specialist_id", cfg.Runtime.SpecialistID, "err", err)
		os.Exit(1)
	}
	bundleBytes, err := bundle.EncodeTOML()
	if err != nil {
		logger.Error("slot bundle TOML encode failed", "specialist_id", cfg.Runtime.SpecialistID, "err", err)
		os.Exit(1)
	}
	summary := bundle.Summary()
	logger.Info("slot bundle compile ok",
		"specialist_id", cfg.Runtime.SpecialistID,
		"schema", bundle.Schema,
		"core_skill_count", summary.CoreSkillCount,
		"playbook_count", summary.PlaybookCount,
		"allowed_tool_count", summary.AllowedToolCount,
		"memory_pointer_count", summary.MemoryPointerCount,
		"source_map_count", summary.SourceMapCount,
		"encoded_bytes", len(bundleBytes),
	)
	versionLabel := time.Now().UTC().Format("20060102T150405Z")
	if err := store.PersistSlotBundleVersion(ctx, cfg.Runtime.SpecialistID, bundleBytes, bundle.SourceMap, versionLabel, cfg.Harness.DefaultCreatedBy); err != nil {
		if fatalErr := handleOptionalFailure(mode, logger, "slot bundle persistence unavailable", err, "version_label", versionLabel); fatalErr != nil {
			os.Exit(1)
		}
	} else {
		logger.Info("slot bundle persisted", "version_label", versionLabel)
	}

	if _, err := store.HealthSnapshot(ctx, cfg.Runtime.SpecialistID); err != nil {
		if fatalErr := handleOptionalFailure(mode, logger, "health snapshot unavailable", err, "specialist_id", cfg.Runtime.SpecialistID); fatalErr != nil {
			os.Exit(1)
		}
	} else {
		logger.Info("health snapshot call ok")
	}

	routingService := routing.NewService(logger)
	hostRegistry := modelhost.NewSimulatedRegistry("verify")
	executionService := execution.NewService(logger, hostRegistry)
	taskProcessor := rt.NewService(cfg, nil, store, nil, routingService, executionService, nil, nil, logger)

	verifyTask := rt.Task{
		ID:      "verify-task",
		Summary: "verify bounded runtime task",
		Class:   "analysis",
		Prompt:  "Verify runtime task processing.",
	}
	if cfg.Runtime.EnableMultimodalSmokeTest {
		verifyTask = rt.Task{
			ID:                          "verify-task",
			Summary:                     "verify multimodal execution",
			Class:                       "evidence_fusion",
			PrimaryModality:             modality.Image,
			SecondaryModalities:         []modality.Type{modality.Normalize(cfg.Runtime.DefaultPrimaryModality)},
			CrossModalGroundingRequired: true,
			AllowTextOnlyFallback:       cfg.Runtime.AllowTextOnlyFallback,
			AssetRefs:                   []string{"sandbox://verify/image-1"},
			Prompt:                      "Verify image and text fusion.",
		}
	}

	result, err := taskProcessor.ProcessTask(ctx, verifyTask)
	if err != nil {
		logger.Error("runtime task verify failed", "err", err)
		os.Exit(1)
	}
	for _, warning := range result.Warnings {
		if fatalErr := handleOptionalFailure(mode, logger, "runtime task persistence warning", fmt.Errorf(warning), "task_id", result.Task.ID); fatalErr != nil {
			os.Exit(1)
		}
	}
	logger.Info("runtime task verify ok",
		"task_id", result.Task.ID,
		"executor", result.ExecutionResult.Plan.ChosenExecutor,
		"execution_mode", result.ExecutionResult.Plan.ExecutionMode,
		"host", result.ExecutionResult.HostResult.HostName,
		"signal_count", len(result.Signals),
		"warning_count", len(result.Warnings),
	)

	signalStore := telemetry.NewMemoryStore()
	if err := signalStore.WriteSignals(ctx, result.Signals); err != nil {
		logger.Error("verify signal store write failed", "err", err)
		os.Exit(1)
	}

	proposalStore := &verifyProposalStore{}
	growthPlanStore := &verifyGrowthPlanStore{}
	sealService := seal.NewService(signalStore, proposalStore, logger)
	proposals, err := sealService.ReviewSpecialist(ctx, cfg.Runtime.SpecialistID, 100)
	if err != nil {
		logger.Error("seal review failed", "err", err)
		os.Exit(1)
	}
	for _, proposal := range proposals {
		if err := store.WriteProposal(ctx, proposal); err != nil {
			if fatalErr := handleOptionalFailure(mode, logger, "adaptation proposal persistence unavailable", err, "proposal_id", proposal.ID); fatalErr != nil {
				os.Exit(1)
			}
		}
	}
	logger.Info("seal verify ok", "signal_count", len(result.Signals), "proposal_count", len(proposals), "persisted_proposal_count", proposalStore.Count())

	if len(proposals) > 0 {
		denService := den.NewService(growthPlanStore, logger)
		plan, err := denService.PlanFromProposal(ctx, proposals[0])
		if err != nil {
			logger.Error("den growth planning failed", "err", err)
			os.Exit(1)
		}
		if err := store.WriteGrowthPlan(ctx, *plan); err != nil {
			if fatalErr := handleOptionalFailure(mode, logger, "growth plan persistence unavailable", err, "growth_plan_id", plan.ID); fatalErr != nil {
				os.Exit(1)
			}
		}
		logger.Info("den verify ok", "proposal_id", proposals[0].ID, "growth_surface", plan.Surface, "growth_plan_count", growthPlanStore.Count())
	}

	growthService := growth.NewService(store, logger)
	growthResult, err := growthService.StageExperiment(ctx, cfg.Runtime.Namespace, cfg.Runtime.SpecialistID, growth.Assessment{
		AbilityName:       "multimodal_grounding",
		GapSummary:        "Verify governed ability-growth storage and staging.",
		EvidenceSummary:   fmt.Sprintf("bounded verify task executor=%s host=%s signals=%d warnings=%d", result.ExecutionResult.Plan.ChosenExecutor, result.ExecutionResult.HostResult.HostName, len(result.Signals), len(result.Warnings)),
		TriedMemoryFix:    true,
		TriedRoutingFix:   true,
		TriedPromptFix:    true,
		PreferredSurface:  "modality_branch",
		RequestedBy:       cfg.Harness.DefaultCreatedBy,
		ParentApprovedBy:  "parent:verify",
		HarnessVerifiedBy: cfg.Harness.DefaultCreatedBy,
		Notes:             "verify governed ability-growth path",
	})
	if err != nil {
		logger.Error("ability-growth verify failed", "err", err)
		os.Exit(1)
	}
	logger.Info("ability-growth verify ok",
		"status", growthResult.Status,
		"experiment_id", growthResult.ExperimentID,
	)

	logger.Info("verify complete", "verify_mode", string(mode))
}
