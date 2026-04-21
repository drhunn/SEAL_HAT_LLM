package runtime

import (
	"log/slog"

	"github.com/drhunn/SEAL_HAT_LLM/internal/config"
	"github.com/drhunn/SEAL_HAT_LLM/internal/execution"
	"github.com/drhunn/SEAL_HAT_LLM/internal/growth"
	"github.com/drhunn/SEAL_HAT_LLM/internal/harness"
	"github.com/drhunn/SEAL_HAT_LLM/internal/memory"
	"github.com/drhunn/SEAL_HAT_LLM/internal/modality"
	"github.com/drhunn/SEAL_HAT_LLM/internal/routing"
	"github.com/drhunn/SEAL_HAT_LLM/internal/slots"
	"github.com/drhunn/SEAL_HAT_LLM/internal/slotsync"
	"github.com/drhunn/SEAL_HAT_LLM/internal/telemetry"
)

type Task struct {
	ID                          string
	Summary                     string
	Class                       string
	PrimaryModality             modality.Type
	SecondaryModalities         []modality.Type
	CrossModalGroundingRequired bool
	AllowTextOnlyFallback       bool
	PreferredExecutor           string
	AssetRefs                   []string
	Prompt                      string
}

type TaskResult struct {
	Task             Task
	RetrievalResults []memory.RetrievalResult
	RoutingDecision  routing.Decision
	ExecutionResult  execution.Result
	Signals          []telemetry.Signal
	Warnings         []string
}

type TaskReview struct {
	TaskID       string
	TaskSummary  string
	Reason       string
	ChosenTarget string
	Executor     string
}

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
