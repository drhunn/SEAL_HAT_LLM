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
	"github.com/drhunn/SEAL_HAT_LLM/internal/memory"
	"github.com/drhunn/SEAL_HAT_LLM/internal/slots"
)

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

	if _, err := store.HealthSnapshot(ctx, cfg.Runtime.SpecialistID); err != nil {
		logger.Warn("health snapshot unavailable", "err", err)
	} else {
		logger.Info("health snapshot call ok")
	}

	if _, err := store.RunCoarseToFineSearch(ctx, cfg.Runtime.Namespace, cfg.Runtime.SpecialistID, memory.ZeroVector(1536), 3, 5, 5); err != nil {
		logger.Warn("retrieval smoke test unavailable", "err", err)
	} else {
		logger.Info("retrieval smoke test ok")
	}

	logger.Info("verify complete")
}
