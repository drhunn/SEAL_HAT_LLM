package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/BurntSushi/toml"
)

type AppConfig struct {
	App struct {
		Name        string `toml:"name"`
		Environment string `toml:"environment"`
		LogLevel    string `toml:"log_level"`
	} `toml:"app"`
	Database struct {
		DSN string `toml:"dsn"`
	} `toml:"database"`
	Runtime struct {
		SpecialistID              string `toml:"specialist_id"`
		Namespace                 string `toml:"namespace"`
		SlotsRoot                 string `toml:"slots_root"`
		ConfigRoot                string `toml:"config_root"`
		RequestTimeoutSeconds     int    `toml:"request_timeout_seconds"`
		DefaultPrimaryModality    string `toml:"default_primary_modality"`
		AllowTextOnlyFallback     bool   `toml:"allow_text_only_fallback"`
		EnableMultimodalSmokeTest bool   `toml:"enable_multimodal_smoke_test"`
		EnableTaskInbox           bool   `toml:"enable_task_inbox"`
		TaskInboxDir              string `toml:"task_inbox_dir"`
		TaskPollIntervalSeconds   int    `toml:"task_poll_interval_seconds"`
	} `toml:"runtime"`
	Harness struct {
		AutoCreatePostmortems bool   `toml:"auto_create_postmortems"`
		AutoUpdateHealth      bool   `toml:"auto_update_health"`
		DefaultCreatedBy      string `toml:"default_created_by"`
	} `toml:"harness"`
}

func Load(path string) (*AppConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config: %w", err)
	}
	defer f.Close()

	var cfg AppConfig
	if _, err := toml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	if cfg.Database.DSN == "" {
		return nil, fmt.Errorf("database.dsn is required")
	}
	if cfg.Runtime.SpecialistID == "" {
		return nil, fmt.Errorf("runtime.specialist_id is required")
	}
	if cfg.Runtime.Namespace == "" {
		return nil, fmt.Errorf("runtime.namespace is required")
	}
	if cfg.Runtime.SlotsRoot == "" {
		cfg.Runtime.SlotsRoot = "./specialists"
	}
	if cfg.Harness.DefaultCreatedBy == "" {
		cfg.Harness.DefaultCreatedBy = "harness:runtime"
	}
	if cfg.Runtime.RequestTimeoutSeconds <= 0 {
		cfg.Runtime.RequestTimeoutSeconds = 30
	}
	if cfg.Runtime.DefaultPrimaryModality == "" {
		cfg.Runtime.DefaultPrimaryModality = "text"
	}
	if cfg.Runtime.TaskPollIntervalSeconds <= 0 {
		cfg.Runtime.TaskPollIntervalSeconds = 5
	}
	if cfg.Runtime.EnableTaskInbox && cfg.Runtime.TaskInboxDir == "" {
		cfg.Runtime.TaskInboxDir = "./artifacts/task_inbox"
	}

	return &cfg, nil
}

func (c *AppConfig) LogLevel() slog.Level {
	switch c.App.LogLevel {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
