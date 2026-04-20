package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
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
	"github.com/drhunn/SEAL_HAT_LLM/internal/seal"
	"github.com/drhunn/SEAL_HAT_LLM/internal/slots"
	"github.com/drhunn/SEAL_HAT_LLM/internal/telemetry"
)

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
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel()}))
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

	collector := telemetry.NewCollector(logger)
	signalStore := telemetry.NewMemoryStore()
	proposalStore := &verifyProposalStore{}
	growthPlanStore := &verifyGrowthPlanStore{}

	if _, err := store.HealthSnapshot(ctx, cfg.Runtime.SpecialistID); err != nil {
		logger.Warn("health snapshot unavailable", "err", err)
	} else {
		logger.Info("health snapshot call ok")
	}

	retrievalResults, retrievalErr := store.RunCoarseToFineSearch(ctx, cfg.Runtime.Namespace, cfg.Runtime.SpecialistID, memory.ZeroVector(1536), 3, 5, 5)
	if retrievalErr != nil {
		logger.Warn("retrieval smoke test unavailable", "err", retrievalErr)
	} else {
		logger.Info("retrieval smoke test ok", "result_count", len(retrievalResults))
	}
	if err := collector.Write(ctx, signalStore, memory.SignalsForRetrieval(cfg.Runtime.SpecialistID, "analysis", retrievalResults, retrievalErr, collector)...); err != nil {
		logger.Error("retrieval telemetry write failed", "err", err)
		os.Exit(1)
	}

	routingService := routing.NewService(logger)
	hostRegistry := modelhost.NewRegistry()
	hostRegistry.Register("Parent-Generalist-30B", modelhost.NewPromptHost("verify-parent-host"))
	hostRegistry.Register("Image-Analysis-Specialist-01", modelhost.NewAssetSummaryHost("verify-image-host", "image"))
	hostRegistry.Register("Audio-Transcription-Specialist-01", modelhost.NewAssetSummaryHost("verify-audio-host", "audio"))
	hostRegistry.Register("Video-Understanding-Specialist-01", modelhost.NewAssetSummaryHost("verify-video-host", "video"))
	hostRegistry.Register("Document-Layout-OCR-Specialist-01", modelhost.NewAssetSummaryHost("verify-document-host", "document"))
	hostRegistry.Register("Multimodal-Evidence-Fusion-Specialist-01", modelhost.NewFusionHost("verify-fusion-host"))
	executionService := execution.NewService(logger, hostRegistry)
	growthService := growth.NewService(store, logger)
	primary := modality.Normalize(cfg.Runtime.DefaultPrimaryModality)
	routingInput := routing.Input{
		TaskSummary:     "verify multimodal routing",
		TaskClass:       "analysis",
		PrimaryModality: primary,
	}
	routingDecision := routingService.DecideTask(ctx, routingInput)
	logger.Info("routing verify ok",
		"chosen_target", routingDecision.ChosenTarget,
		"primary_modality", routingDecision.PrimaryModality,
	)
	if err := collector.Write(ctx, signalStore, routing.SignalsForDecision(cfg.Runtime.SpecialistID, routingInput, routingDecision, collector)...); err != nil {
		logger.Error("routing telemetry write failed", "err", err)
		os.Exit(1)
	}

	var executionResult execution.Result
	var executionErr error
	if cfg.Runtime.EnableMultimodalSmokeTest {
		executionReq := execution.Request{
			TaskSummary:                 "verify multimodal execution",
			TaskClass:                   "evidence_fusion",
			PrimaryModality:             modality.Image,
			SecondaryModalities:         []modality.Type{primary},
			CrossModalGroundingRequired: true,
			AllowTextOnlyFallback:       cfg.Runtime.AllowTextOnlyFallback,
			AssetRefs:                   []string{"sandbox://verify/image-1"},
			Prompt:                      "Verify image and text fusion.",
		}
		executionResult, executionErr = executionService.Execute(ctx, executionReq)
		if executionErr != nil {
			logger.Error("execution verify failed", "err", executionErr)
			os.Exit(1)
		}
		logger.Info("execution verify ok",
			"executor", executionResult.Plan.ChosenExecutor,
			"execution_mode", executionResult.Plan.ExecutionMode,
			"host", executionResult.HostResult.HostName,
			"handled", executionResult.HostResult.Handled,
		)
		if err := collector.Write(ctx, signalStore, execution.SignalsForExecution(cfg.Runtime.SpecialistID, executionReq, executionResult, executionErr, collector)...); err != nil {
			logger.Error("execution telemetry write failed", "err", err)
			os.Exit(1)
		}
	}

	sealService := seal.NewService(signalStore, proposalStore, logger)
	proposals, err := sealService.ReviewSpecialist(ctx, cfg.Runtime.SpecialistID, 100)
	if err != nil {
		logger.Error("seal review failed", "err", err)
		os.Exit(1)
	}
	logger.Info("seal verify ok", "signal_count", signalStore.Count(), "proposal_count", len(proposals), "persisted_proposal_count", proposalStore.Count())

	if len(proposals) > 0 {
		denService := den.NewService(growthPlanStore, logger)
		plan, err := denService.PlanFromProposal(ctx, proposals[0])
		if err != nil {
			logger.Error("den growth planning failed", "err", err)
			os.Exit(1)
		}
		logger.Info("den verify ok", "proposal_id", proposals[0].ID, "growth_surface", plan.Surface, "growth_plan_count", growthPlanStore.Count())
	}

	growthResult, err := growthService.StageExperiment(ctx, cfg.Runtime.Namespace, cfg.Runtime.SpecialistID, growth.Assessment{
		AbilityName:       "multimodal_grounding",
		GapSummary:        "Verify governed ability-growth storage and staging.",
		EvidenceSummary:   "Verification path confirms current multimodal ability remains scaffold-level.",
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

	logger.Info("verify complete")
}
