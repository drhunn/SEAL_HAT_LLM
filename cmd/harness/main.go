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
	"github.com/drhunn/SEAL_HAT_LLM/internal/unit"
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

	unitSpec, err := unit.SpecFromConfig(cfg)
	if err != nil {
		logger.Error("build unit spec", "err", err)
		os.Exit(1)
	}

	localRuntime, err := unit.NewLocalRuntime(unitSpec, cfg, database, logger)
	if err != nil {
		logger.Error("bootstrap local runtime unit", "unit_id", unitSpec.UnitID, "err", err)
		os.Exit(1)
	}

	logger.Info("starting local model unit",
		"unit_id", unitSpec.UnitID,
		"role", unitSpec.Role,
		"store_mode", unitSpec.StoreMode,
		"tool_plane_mode", unitSpec.ToolPlaneMode,
	)

	if err := localRuntime.Start(ctx); err != nil {
		logger.Error("runtime stopped with error", "unit_id", unitSpec.UnitID, "err", err)
		os.Exit(1)
	}

	logger.Info("runtime stopped cleanly", "unit_id", unitSpec.UnitID)
}
