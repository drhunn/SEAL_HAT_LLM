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
	"github.com/drhunn/SEAL_HAT_LLM/internal/telemetry"
)

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

	var taskFile TaskFile
	if err := json.Unmarshal(data, &taskFile); err != nil {
		return queued, fmt.Errorf("decode task file %s: %w", name, err)
	}
	if err := validateTaskFile(taskFile); err != nil {
		return queued, fmt.Errorf("validate task file %s: %w", name, err)
	}
	var taskHints struct {
		PreferredUnitID string `json:"preferred_unit_id"`
	}
	if err := json.Unmarshal(data, &taskHints); err != nil {
		return queued, fmt.Errorf("decode task unit hints %s: %w", name, err)
	}

	queued.Task = Task{
		ID:                          firstNonEmpty(strings.TrimSpace(taskFile.ID), queued.Task.ID),
		Summary:                     firstNonEmpty(strings.TrimSpace(taskFile.Summary), queued.Task.Summary),
		Class:                       firstNonEmpty(strings.TrimSpace(taskFile.Class), queued.Task.Class),
		PrimaryModality:             modality.Normalize(taskFile.PrimaryModality),
		SecondaryModalities:         normalizeSecondaryModalities(taskFile.SecondaryModalities),
		CrossModalGroundingRequired: taskFile.CrossModalGroundingRequired,
		AllowTextOnlyFallback:       taskFile.AllowTextOnlyFallback,
		PreferredExecutor:           strings.TrimSpace(taskFile.PreferredExecutor),
		PreferredUnitID:             strings.TrimSpace(taskHints.PreferredUnitID),
		AssetRefs:                   compactStrings(taskFile.AssetRefs),
		Prompt:                      strings.TrimSpace(taskFile.Prompt),
	}

	return queued, nil
}

func (i *TaskInbox) MarkProcessed(queued *queuedTask, result *TaskResult) error {
	if queued == nil {
		return nil
	}
	dst := i.uniqueArchivePath(i.processedDir, queued.ArchiveName)
	if err := os.Rename(queued.WorkingPath, dst); err != nil {
		return fmt.Errorf("archive processed task %s: %w", queued.ArchiveName, err)
	}
	artifact := TaskArtifact{
		Version:         TaskArtifactSchemaVersion,
		TaskID:          queued.Task.ID,
		Status:          "processed",
		Summary:         queued.Task.Summary,
		Class:           queued.Task.Class,
		PrimaryModality: queued.Task.PrimaryModality.String(),
		ArchivedPath:    dst,
		CompletedAt:     time.Now().UTC(),
	}
	if result != nil {
		artifact.Executor = result.ExecutionResult.Plan.ChosenExecutor
		artifact.HostName = result.ExecutionResult.HostResult.HostName
		artifact.SignalSummaries = signalSummaries(result.Signals)
	}
	if err := i.writeArtifact(dst, artifact); err != nil {
		return err
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
	artifact := TaskArtifact{
		Version:         TaskArtifactSchemaVersion,
		TaskID:          queued.Task.ID,
		Status:          "failed",
		Summary:         queued.Task.Summary,
		Class:           queued.Task.Class,
		PrimaryModality: queued.Task.PrimaryModality.String(),
		ArchivedPath:    dst,
		Error:           message,
		CompletedAt:     time.Now().UTC(),
	}
	if err := i.writeArtifact(dst, artifact); err != nil {
		return err
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

func (i *TaskInbox) writeArtifact(archivePath string, artifact TaskArtifact) error {
	payload, err := json.MarshalIndent(artifact, "", "  ")
	if err != nil {
		return fmt.Errorf("encode task artifact: %w", err)
	}
	if err := os.WriteFile(archivePath+".result.json", payload, 0o644); err != nil {
		return fmt.Errorf("write task artifact: %w", err)
	}
	return nil
}

func validateTaskFile(taskFile TaskFile) error {
	if strings.TrimSpace(taskFile.Version) != TaskInboxSchemaVersion {
		return fmt.Errorf("version must be %q", TaskInboxSchemaVersion)
	}
	if strings.TrimSpace(taskFile.Summary) == "" {
		return fmt.Errorf("summary is required")
	}
	return nil
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

func signalSummaries(signals []telemetry.Signal) []string {
	out := make([]string, 0, len(signals))
	for _, signal := range signals {
		if strings.TrimSpace(signal.Summary) == "" {
			continue
		}
		out = append(out, signal.Summary)
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
