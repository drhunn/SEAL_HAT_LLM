package lifecycle

import "testing"

func TestValidateArtifactTransitionAllowsLegalTransitions(t *testing.T) {
	tests := []struct {
		from ArtifactState
		to   ArtifactState
		want ArtifactEventType
	}{
		{ArtifactStateDefined, ArtifactStateCandidate, ArtifactEventCandidate},
		{ArtifactStateCandidate, ArtifactStateCurrent, ArtifactEventPromoted},
		{ArtifactStateCandidate, ArtifactStateRejected, ArtifactEventRejected},
		{ArtifactStateCandidate, ArtifactStateRolledBack, ArtifactEventRolledBack},
		{ArtifactStateCurrent, ArtifactStateRolledBack, ArtifactEventRolledBack},
		{ArtifactStateRejected, ArtifactStateArchived, ArtifactEventArchived},
	}
	for _, tt := range tests {
		t.Run(string(tt.from)+"->"+string(tt.to), func(t *testing.T) {
			got, err := ValidateArtifactTransition(tt.from, tt.to)
			if err != nil {
				t.Fatalf("ValidateArtifactTransition returned error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("event type = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValidateArtifactTransitionRejectsIllegalTransitions(t *testing.T) {
	tests := []struct {
		from ArtifactState
		to   ArtifactState
	}{
		{ArtifactStateDefined, ArtifactStateCurrent},
		{ArtifactStateCurrent, ArtifactStateCandidate},
		{ArtifactStateRolledBack, ArtifactStateCurrent},
		{ArtifactStateArchived, ArtifactStateCurrent},
		{ArtifactStateCandidate, ArtifactStateCandidate},
	}
	for _, tt := range tests {
		t.Run(string(tt.from)+"->"+string(tt.to), func(t *testing.T) {
			if _, err := ValidateArtifactTransition(tt.from, tt.to); err == nil {
				t.Fatalf("expected illegal transition error")
			}
		})
	}
}
