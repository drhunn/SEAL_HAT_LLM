package modelhost

import "testing"

func TestNewSimulatedRegistryRegistersExpectedExecutors(t *testing.T) {
	registry := NewSimulatedRegistry("verify")

	expected := []string{
		ExecutorAudioTranscription,
		ExecutorDocumentLayoutOCR,
		ExecutorImageAnalysis,
		ExecutorMultimodalEvidenceFusion,
		ExecutorParentGeneralist,
		ExecutorVideoUnderstanding,
	}

	executors := registry.Executors()
	if len(executors) != len(expected) {
		t.Fatalf("expected %d executors, got %d: %+v", len(expected), len(executors), executors)
	}
	for i := range expected {
		if executors[i] != expected[i] {
			t.Fatalf("unexpected executor order at %d: got %q want %q", i, executors[i], expected[i])
		}
	}

	host, ok := registry.Resolve(ExecutorParentGeneralist)
	if !ok {
		t.Fatalf("expected parent executor to be registered")
	}
	if host.Name() != "verify-parent-host" {
		t.Fatalf("unexpected parent host name: %q", host.Name())
	}

	defaultRegistry := NewSimulatedRegistry("")
	defaultHost, ok := defaultRegistry.Resolve(ExecutorMultimodalEvidenceFusion)
	if !ok {
		t.Fatalf("expected fusion executor to be registered")
	}
	if defaultHost.Name() != "local-fusion-host" {
		t.Fatalf("unexpected default fusion host name: %q", defaultHost.Name())
	}
}
