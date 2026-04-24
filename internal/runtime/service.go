package runtime

import (
	"context"
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
	PreferredUnitID             string
	AssetRefs                   []string
	Prompt                      string
}

type RemoteDispatchResult struct {
	TargetUnitID string
	SocketPath   string
	Status       string
	ResultSummary string
	OutputJSON   string
	ArtifactRefs []string
	Confidence   float64
	Warnings     []string
	Signals      []string
	ErrorText    string
}

type RemoteDispatcher interface {
	DispatchRemote(context.Context, Task) (*RemoteDispatchResult, error)
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
	TargetUnitID string
	Executor     string
}

type Service struct {
	cfg              *config.AppConfig
	loader           *slots.FilesystemLoader
	store            *memory.PostgresStore
	harness          *harness.Service
	routing          *routing.Service
	execution        *execution.Service
	growth           *growth.Service
	slotSync         *slotsync.Service
	remoteDispatcher RemoteDispatcher
	logger           *slog.Logger
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

func (s *Service) SetRemoteDispatcher(dispatcher RemoteDispatcher) {
	if s == nil {
		return
	}
	s.remoteDispatcher = dispatcher
}
