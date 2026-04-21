package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/drhunn/SEAL_HAT_LLM/internal/config"
	iverify "github.com/drhunn/SEAL_HAT_LLM/internal/verify"
)

func main() {
	cfgPath := flag.String("config", "config/runtime.example.toml", "path to runtime TOML config")
	modeFlag := flag.String("mode", string(iverify.ModeSoft), "verify mode: soft or strict")
	flag.Parse()

	mode, err := iverify.ParseMode(*modeFlag)
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Runtime.RequestTimeoutSeconds)*time.Second)
	defer cancel()

	runner := iverify.NewRunner(cfg, mode, logger)
	if err := runner.Run(ctx); err != nil {
		logger.Error("verify failed", "err", err)
		os.Exit(1)
	}
}
