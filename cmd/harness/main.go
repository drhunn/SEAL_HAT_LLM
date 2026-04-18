package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/drhunn/LLM-plus-harness/internal/config"
	"github.com/drhunn/LLM-plus-harness/internal/db"
	"github.com/drhunn/LLM-plus-harness/internal/harness"
	"github.com/drhunn/LLM-plus-harness/internal/memory"
	"github.com/drhunn/LLM-plus-harness/internal/postmortem"
	"github.com/drhunn/LLM-plus-harness/internal/runtime"
	"github.com/drhunn/LLM-plus-harness/internal/slots"
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

	store := memory.NewPostgresStore(database, logger)
	slotLoader := slots.NewFilesystemLoader(cfg.Runtime.SlotsRoot)
	pmService := postmortem.NewService(store, logger, cfg.Harness.DefaultCreatedBy)
	harnessService := harness.NewService(store, pmService, logger, cfg)
	runtimeService := runtime.NewService(cfg, slotLoader, store, harnessService, logger)

	if err := runtimeService.Start(ctx); err != nil {
		logger.Error("runtime stopped with error", "err", err)
		os.Exit(1)
	}

	logger.Info("runtime stopped cleanly")
}
