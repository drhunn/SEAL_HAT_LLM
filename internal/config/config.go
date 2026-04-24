package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

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
		StoreMode                 string `toml:"store_mode"`
		RequestTimeoutSeconds     int    `toml:"request_timeout_seconds"`
		DefaultPrimaryModality    string `toml:"default_primary_modality"`
		AllowTextOnlyFallback     bool   `toml:"allow_text_only_fallback"`
		EnableMultimodalSmokeTest bool   `toml:"enable_multimodal_smoke_test"`
		EnableTaskInbox           bool   `toml:"enable_task_inbox"`
		TaskInboxDir              string `toml:"task_inbox_dir"`
		TaskPollIntervalSeconds   int    `toml:"task_poll_interval_seconds"`
		EnableTaskRPCServer       bool   `toml:"enable_task_rpc_server"`
		TaskRPCSocketPath         string `toml:"task_rpc_socket_path"`
	} `toml:"runtime"`
	EmbeddedPostgres struct {
		DataDir      string `toml:"data_dir"`
		Port         int    `toml:"port"`
		User         string `toml:"user"`
		DatabaseName string `toml:"database_name"`
		BinDir       string `toml:"bin_dir"`
		SQLRoot      string `toml:"sql_root"`
	} `toml:"embedded_postgres"`
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

	if cfg.Runtime.SpecialistID == "" {
		return nil, fmt.Errorf("runtime.specialist_id is required")
	}
	if cfg.Runtime.Namespace == "" {
		return nil, fmt.Errorf("runtime.namespace is required")
	}
	if cfg.Runtime.SlotsRoot == "" {
		cfg.Runtime.SlotsRoot = "./specialists"
	}
	if cfg.Runtime.ConfigRoot == "" {
		cfg.Runtime.ConfigRoot = "./config"
	}
	if strings.TrimSpace(cfg.Runtime.StoreMode) == "" {
		cfg.Runtime.StoreMode = "shared_dsn"
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
	if cfg.Runtime.EnableTaskRPCServer && strings.TrimSpace(cfg.Runtime.TaskRPCSocketPath) == "" {
		cfg.Runtime.TaskRPCSocketPath = filepath.Join("./artifacts/taskrpc", cfg.Runtime.SpecialistID+".sock")
	}

	switch strings.TrimSpace(cfg.Runtime.StoreMode) {
	case "shared_dsn":
		if cfg.Database.DSN == "" {
			return nil, fmt.Errorf("database.dsn is required for runtime.store_mode=shared_dsn")
		}
	case "embedded_postgres":
		if cfg.EmbeddedPostgres.DataDir == "" {
			cfg.EmbeddedPostgres.DataDir = filepath.Join("./artifacts/embedded_postgres", cfg.Runtime.SpecialistID)
		}
		if cfg.EmbeddedPostgres.Port <= 0 {
			cfg.EmbeddedPostgres.Port = 55432
		}
		if cfg.EmbeddedPostgres.User == "" {
			cfg.EmbeddedPostgres.User = "postgres"
		}
		if cfg.EmbeddedPostgres.DatabaseName == "" {
			cfg.EmbeddedPostgres.DatabaseName = "postgres"
		}
		if strings.TrimSpace(cfg.EmbeddedPostgres.SQLRoot) == "" {
			cfg.EmbeddedPostgres.SQLRoot = "./sql"
		}
	default:
		return nil, fmt.Errorf("unsupported runtime.store_mode %q", cfg.Runtime.StoreMode)
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
