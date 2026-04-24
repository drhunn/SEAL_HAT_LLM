package lifecycle

import (
	"fmt"
	"strings"
)

type ArtifactState string

const (
	ArtifactStateDefined  ArtifactState = "defined"
	ArtifactStateCandidate ArtifactState = "candidate"
	ArtifactStateCurrent   ArtifactState = "current"
	ArtifactStateRolledBack ArtifactState = "rolled_back"
	ArtifactStateRejected  ArtifactState = "rejected"
	ArtifactStateArchived  ArtifactState = "archived"
)

type ArtifactEventType string

const (
	ArtifactEventDefined  ArtifactEventType = "defined"
	ArtifactEventCandidate ArtifactEventType = "candidate"
	ArtifactEventPromoted ArtifactEventType = "promoted"
	ArtifactEventRolledBack ArtifactEventType = "rolled_back"
	ArtifactEventRejected ArtifactEventType = "rejected"
	ArtifactEventArchived ArtifactEventType = "archived"
)

var allowedArtifactTransitions = map[ArtifactState]map[ArtifactState]ArtifactEventType{
	ArtifactStateDefined: {
		ArtifactStateCandidate: ArtifactEventCandidate,
		ArtifactStateArchived:  ArtifactEventArchived,
	},
	ArtifactStateCandidate: {
		ArtifactStateCurrent:    ArtifactEventPromoted,
		ArtifactStateRejected:   ArtifactEventRejected,
		ArtifactStateRolledBack: ArtifactEventRolledBack,
		ArtifactStateArchived:   ArtifactEventArchived,
	},
	ArtifactStateCurrent: {
		ArtifactStateRolledBack: ArtifactEventRolledBack,
		ArtifactStateArchived:   ArtifactEventArchived,
	},
	ArtifactStateRolledBack: {
		ArtifactStateArchived: ArtifactEventArchived,
	},
	ArtifactStateRejected: {
		ArtifactStateArchived: ArtifactEventArchived,
	},
}

func NormalizeArtifactState(value string) (ArtifactState, error) {
	state := ArtifactState(strings.TrimSpace(value))
	switch state {
	case ArtifactStateDefined, ArtifactStateCandidate, ArtifactStateCurrent, ArtifactStateRolledBack, ArtifactStateRejected, ArtifactStateArchived:
		return state, nil
	default:
		return "", fmt.Errorf("unknown artifact state %q", value)
	}
}

func ValidateArtifactTransition(from, to ArtifactState) (ArtifactEventType, error) {
	if from == to {
		return "", fmt.Errorf("artifact already in state %q", to)
	}
	allowed, ok := allowedArtifactTransitions[from]
	if !ok {
		return "", fmt.Errorf("artifact state %q has no outgoing transitions", from)
	}
	eventType, ok := allowed[to]
	if !ok {
		return "", fmt.Errorf("illegal artifact transition %q -> %q", from, to)
	}
	return eventType, nil
}
