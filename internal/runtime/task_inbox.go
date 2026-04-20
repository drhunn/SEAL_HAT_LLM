package runtime

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/drhunn/SEAL_HAT_LLM/internal/modality"
)

type taskFileEnvelope struct {
	ID                          string   `json:"id"`
	Summary                     string   `json:"summary"`
	Class                       string   `json:"class"`
	PrimaryModality             string   `json:"primary_modality"`
	SecondaryModalities         []string `json:"secondary_modalities"`
	CrossModalGroundingRequired bool     `json:"cross_modal_grounding_required"`
	AllowTextOnlyFallback       bool     `json:"allow_text_only_fallback"`
	PreferredExecutor           string   `json:"preferred_executor"`
	AssetRefs                   []string `json:"asset_refs"`
	Prompt                      string   `json:"prompt"`
}

type queuedTask struct {
	Task        Task
	WorkingPath string
	ArchiveName string
}

type TaskInbox struct {
	rootDir      string
	processedDir string
	failedDir    string
	logger       *slog.Logger
}

func NewTaskInbox(rootDir string, logger *slog.Logger) *TaskInbox {
	rootDir = strings.TrimSpace(rootDir)
	if rootDir == "" {
		return nil
	}
	return &TaskInbox{
		rootDir:      rootDir,
		processedDir: filepath.Join(rootDir, "processed"),
		failedDir:    filepath.Join(rootDir, "failed"),
		logger:       logger,
	}
}

func (i *TaskInbox) EnsureDirs() error {
	for _, dir := range []string{i.rootDir, i.processedDir, i.failedDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("ensure task inbox dir %s: %w", dir, err)
		}
	}
	return nil
}

func (i *TaskInbox) ClaimNext() (*queuedTask, error) {
	entries, err := os.ReadDir(i.rootDir)
	if err != nil {
		return nil, fmt.Errorf("read task inbox: %w", err)
	}

	names := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".json") {
			continue
		}
		names = append(names, name)
	}
	if len(names) == 0 {
		return nil, nil
	}
	sort.Strings(names)

	name := names[0]
	src := filepath.Join(i.rootDir, name)
	working := src + ".working"
	if err := os.Rename(src, working); err != nil {
		return nil, fmt.Errorf("claim task file %s: %w", name, err)
	}

	queued := &queuedTask{
		Task: Task{
			ID:      strings.TrimSuffix(name, filepath.Ext(name)),
			Summary: strings.TrimSuffix(name, filepath.Ext(name)),
			Class:   "analysis",
		},
		WorkingPath: working,
		ArchiveName: name,
	}

	data, err := os.ReadFile(working)
	if err != nil {
		return queued, fmt.Errorf("read claimed task file %s: %w", name, err)
	}

	var envelope taskFileEnvelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return queued, fmt.Errorf("decode task file %s: %w", name, err)
	}

	queued.Task = Task{
		ID:                          firstNonEmpty(strings.TrimSpace(envelope.ID), queued.Task.ID),
		Summary:                     firstNonEmpty(strings.TrimSpace(envelope.Summary), queued.Task.Summary),
		Class:                       firstNonEmpty(strings.TrimSpace(envelope.Class), queued.Task.Class),
		PrimaryModality:             modality.Normalize(envelope.PrimaryModality),
		SecondaryModalities:         normalizeSecondaryModalities(envelope.SecondaryModalities),
		CrossModalGroundingRequired: envelope.CrossModalGroundingRequired,
		AllowTextOnlyFallback:       envelope.AllowTextOnlyFallback,
		PreferredExecutor:           strings.TrimSpace(envelope.PreferredExecutor),
		AssetRefs:                   compactStrings(envelope.AssetRefs),
		Prompt:                      strings.TrimSpace(envelope.Prompt),
	}

	return queued, nil
}

func (i *TaskInbox) MarkProcessed(queued *queuedTask) error {
	if queued == nil {
		return nil
	}
	dst := i.uniqueArchivePath(i.processedDir, queued.ArchiveName)
	if err := os.Rename(queued.WorkingPath, dst); err != nil {
		return fmt.Errorf("archive processed task %s: %w", queued.ArchiveName, err)
	}
	if i.logger != nil {
		i.logger.Info("task archived as processed", "task_id", queued.Task.ID, "path", dst)
	}
	return nil
}

func (i *TaskInbox) MarkFailed(queued *queuedTask, taskErr error) error {
	if queued == nil {
		return nil
	}
	dst := i.uniqueArchivePath(i.failedDir, queued.ArchiveName)
	if err := os.Rename(queued.WorkingPath, dst); err != nil {
		return fmt.Errorf("archive failed task %s: %w", queued.ArchiveName, err)
	}
	message := "task failed"
	if taskErr != nil {
		message = strings.TrimSpace(taskErr.Error())
	}
	if err := os.WriteFile(dst+".error.txt", []byte(message+"\n"), 0o644); err != nil {
		return fmt.Errorf("write task failure note: %w", err)
	}
	if i.logger != nil {
		i.logger.Warn("task archived as failed", "task_id", queued.Task.ID, "path", dst, "err", taskErr)
	}
	return nil
}

func (i *TaskInbox) uniqueArchivePath(dir, name string) string {
	candidate := filepath.Join(dir, name)
	if _, err := os.Stat(candidate); os.IsNotExist(err) {
		return candidate
	}
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	return filepath.Join(dir, fmt.Sprintf("%s-%d%s", base, time.Now().UTC().UnixNano(), ext))
}

func normalizeSecondaryModalities(values []string) []modality.Type {
	out := make([]modality.Type, 0, len(values))
	for _, value := range values {
		item := modality.Normalize(value)
		if !item.IsKnown() {
			continue
		}
		out = append(out, item)
	}
	return out
}

func compactStrings(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		out = append(out, value)
	}
	return out
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
