package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/drhunn/SEAL_HAT_LLM/internal/config"
	"github.com/drhunn/SEAL_HAT_LLM/internal/db"
	"github.com/drhunn/SEAL_HAT_LLM/internal/evals"
	"github.com/drhunn/SEAL_HAT_LLM/internal/execution"
	"github.com/drhunn/SEAL_HAT_LLM/internal/harness"
	workflow "github.com/drhunn/SEAL_HAT_LLM/internal/harness/workflows"
	"github.com/drhunn/SEAL_HAT_LLM/internal/lifecycle"
	"github.com/drhunn/SEAL_HAT_LLM/internal/memory"
	"github.com/drhunn/SEAL_HAT_LLM/internal/modelhost"
	"github.com/drhunn/SEAL_HAT_LLM/internal/postmortem"
	"github.com/drhunn/SEAL_HAT_LLM/internal/routing"
	"github.com/drhunn/SEAL_HAT_LLM/internal/runtime"
	"github.com/drhunn/SEAL_HAT_LLM/internal/slots"
	"github.com/drhunn/SEAL_HAT_LLM/internal/slotsync"
)

func main() {
	cfgPath := flag.String("config", "config/runtime.example.toml", "path to runtime TOML config")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		slog.Error("load config", "err", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel()}))
	slog.SetDefault(logger)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	database, err := db.Open(ctx, cfg.Database.DSN)
	if err != nil {
		logger.Error("open database", "err", err)
		os.Exit(1)
	}
	defer database.Close()

	// If this looks like paperwork, that's because governance is paperwork with better logging.
	store := memory.NewPostgresStore(database, logger)
	slotLoader := slots.NewFilesystemLoader(cfg.Runtime.SlotsRoot)
	slotSyncService := slotsync.NewService(slotLoader, logger)
	pmService := postmortem.NewService(store, logger, cfg.Harness.DefaultCreatedBy)
	evalService := evals.NewService(store, logger)
	lifecycleService := lifecycle.NewService(database, logger)
	recoveryPlanner := workflow.NewDefaultRecoveryPlanner()
	harnessService := harness.NewService(store, pmService, evalService, lifecycleService, recoveryPlanner, logger, cfg)
	routingService := routing.NewService(logger)
	hostRegistry := modelhost.NewRegistry()
	hostRegistry.Register("Parent-Generalist-30B", modelhost.NewStaticHost("local-parent-host", "parent text execution stub ok"))
	hostRegistry.Register("Image-Analysis-Specialist-01", modelhost.NewStaticHost("local-image-host", "image execution stub ok"))
	hostRegistry.Register("Audio-Transcription-Specialist-01", modelhost.NewStaticHost("local-audio-host", "audio execution stub ok"))
	hostRegistry.Register("Video-Understanding-Specialist-01", modelhost.NewStaticHost("local-video-host", "video execution stub ok"))
	hostRegistry.Register("Document-Layout-OCR-Specialist-01", modelhost.NewStaticHost("local-document-host", "document execution stub ok"))
	hostRegistry.Register("Multimodal-Evidence-Fusion-Specialist-01", modelhost.NewStaticHost("local-fusion-host", "multimodal fusion stub ok"))
	executionService := execution.NewService(logger, hostRegistry)
	runtimeService := runtime.NewService(cfg, slotLoader, store, harnessService, routingService, executionService, slotSyncService, logger)

	if err := runtimeService.Start(ctx); err != nil {
		logger.Error("runtime stopped with error", "err", err)
		os.Exit(1)
	}

	logger.Info("runtime stopped cleanly")
}
