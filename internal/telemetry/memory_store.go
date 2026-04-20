package telemetry

import (
	"context"
	"strings"
	"sync"
)

type MemoryStore struct {
	mu      sync.Mutex
	signals []Signal
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (s *MemoryStore) WriteSignals(ctx context.Context, signals []Signal) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.signals = append(s.signals, signals...)
	return nil
}

func (s *MemoryStore) ListSignals(ctx context.Context, specialistID string, limit int) ([]Signal, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Signal, 0)
	for i := len(s.signals) - 1; i >= 0; i-- {
		signal := s.signals[i]
		if strings.TrimSpace(specialistID) != "" && signal.SpecialistID != specialistID {
			continue
		}
		out = append(out, signal)
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out, nil
}

func (s *MemoryStore) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.signals)
}
