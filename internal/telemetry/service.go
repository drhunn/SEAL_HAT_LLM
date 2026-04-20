package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityModerate Severity = "moderate"
	SeverityHigh     Severity = "high"
)

type Signal struct {
	ID           string
	SpecialistID string
	Category     string
	Severity     Severity
	Surface      string
	TaskClass    string
	Summary      string
	EvidenceRefs []string
	OccurredAt   time.Time
}

type SignalStore interface {
	WriteSignals(ctx context.Context, signals []Signal) error
}

type Collector struct {
	logger *slog.Logger
}

func NewCollector(logger *slog.Logger) *Collector {
	return &Collector{logger: logger}
}

func (c *Collector) NewSignal(specialistID, category, surface, taskClass, summary string, severity Severity, evidenceRefs ...string) Signal {
	return Signal{
		ID:           signalID(specialistID, category, surface),
		SpecialistID: strings.TrimSpace(specialistID),
		Category:     normalize(category),
		Severity:     normalizeSeverity(severity),
		Surface:      normalize(surface),
		TaskClass:    normalize(taskClass),
		Summary:      strings.TrimSpace(summary),
		EvidenceRefs: compact(evidenceRefs),
		OccurredAt:   time.Now().UTC(),
	}
}

func (c *Collector) Write(ctx context.Context, store SignalStore, signals ...Signal) error {
	if len(signals) == 0 {
		return nil
	}
	if store == nil {
		return fmt.Errorf("signal store is required")
	}
	if err := store.WriteSignals(ctx, signals); err != nil {
		return err
	}
	c.logger.InfoContext(ctx, "telemetry signals written", "count", len(signals))
	return nil
}

func signalID(specialistID, category, surface string) string {
	return fmt.Sprintf("%s|%s|%s|%d", normalize(specialistID), normalize(category), normalize(surface), time.Now().UTC().UnixNano())
}

func normalizeSeverity(severity Severity) Severity {
	switch severity {
	case SeverityLow, SeverityModerate, SeverityHigh:
		return severity
	default:
		return SeverityModerate
	}
}

func normalize(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, " ", "_")
	return value
}

func compact(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}
	return out
}
