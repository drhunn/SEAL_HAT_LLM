package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	rt "github.com/drhunn/SEAL_HAT_LLM/internal/runtime"
)

func main() {
	inboxDir := flag.String("inbox-dir", "./artifacts/task_inbox", "path to the runtime task inbox directory")
	id := flag.String("id", "", "task id; defaults to a slug + timestamp")
	summary := flag.String("summary", "", "required task summary")
	class := flag.String("class", "analysis", "task class")
	primaryModality := flag.String("primary-modality", "text", "primary modality")
	secondaryModalities := flag.String("secondary-modalities", "", "comma-separated secondary modalities")
	crossModalGrounding := flag.Bool("cross-modal-grounding", false, "whether the task requires cross-modal grounding")
	allowTextOnlyFallback := flag.Bool("allow-text-only-fallback", false, "whether text-only fallback is permitted")
	preferredExecutor := flag.String("preferred-executor", "", "optional preferred executor")
	assetRefs := flag.String("asset-refs", "", "comma-separated asset refs")
	prompt := flag.String("prompt", "", "optional prompt")
	flag.Parse()

	trimmedSummary := strings.TrimSpace(*summary)
	if trimmedSummary == "" {
		fmt.Fprintln(os.Stderr, "summary is required")
		os.Exit(1)
	}

	trimmedID := strings.TrimSpace(*id)
	if trimmedID == "" {
		trimmedID = fmt.Sprintf("%s-%d", slugify(trimmedSummary), time.Now().UTC().UnixNano())
	}
	if err := os.MkdirAll(*inboxDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "create inbox dir: %v\n", err)
		os.Exit(1)
	}

	taskFile := rt.TaskFile{
		Version:                     rt.TaskInboxSchemaVersion,
		ID:                          trimmedID,
		Summary:                     trimmedSummary,
		Class:                       strings.TrimSpace(*class),
		PrimaryModality:             strings.TrimSpace(*primaryModality),
		SecondaryModalities:         splitCSV(*secondaryModalities),
		CrossModalGroundingRequired: *crossModalGrounding,
		AllowTextOnlyFallback:       *allowTextOnlyFallback,
		PreferredExecutor:           strings.TrimSpace(*preferredExecutor),
		AssetRefs:                   splitCSV(*assetRefs),
		Prompt:                      strings.TrimSpace(*prompt),
	}

	payload, err := json.MarshalIndent(taskFile, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "encode task file: %v\n", err)
		os.Exit(1)
	}

	path := filepath.Join(*inboxDir, trimmedID+".json")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create task file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	if _, err := f.Write(payload); err != nil {
		fmt.Fprintf(os.Stderr, "write task file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(path)
}

func splitCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		out = append(out, part)
	}
	return out
}

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, " ", "-")
	builder := strings.Builder{}
	for _, ch := range value {
		switch {
		case ch >= 'a' && ch <= 'z':
			builder.WriteRune(ch)
		case ch >= '0' && ch <= '9':
			builder.WriteRune(ch)
		case ch == '-':
			builder.WriteRune(ch)
		}
	}
	out := strings.Trim(builder.String(), "-")
	if out == "" {
		return "task"
	}
	return out
}
