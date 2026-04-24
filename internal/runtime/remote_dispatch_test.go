package runtime

import (
	"context"
	"testing"

	"github.com/drhunn/SEAL_HAT_LLM/internal/config"
	"github.com/drhunn/SEAL_HAT_LLM/internal/routing"
	"github.com/drhunn/SEAL_HAT_LLM/internal/unitref"
)

func TestShouldDispatchRemoteRequiresConfiguredNonLocalTarget(t *testing.T) {
	cfg := &config.AppConfig{}
	cfg.Runtime.SpecialistID = "local-unit"
	cfg.TaskDispatch.RemoteUnitSockets = map[string]string{"remote-unit": "/tmp/remote.sock"}

	if shouldDispatchRemote(cfg, nil, routing.Decision{TargetUnitID: "remote-unit"}) {
		t.Fatalf("expected nil dispatcher to disable remote dispatch")
	}
	if shouldDispatchRemote(cfg, fakeRemoteDispatcher{}, routing.Decision{TargetUnitID: "local-unit"}) {
		t.Fatalf("expected local target to avoid remote dispatch")
	}
	if shouldDispatchRemote(cfg, fakeRemoteDispatcher{}, routing.Decision{TargetUnitID: "missing-unit"}) {
		t.Fatalf("expected unmapped target to avoid remote dispatch")
	}
	if !shouldDispatchRemote(cfg, fakeRemoteDispatcher{}, routing.Decision{TargetUnitID: "remote-unit"}) {
		t.Fatalf("expected mapped non-local target to use remote dispatch")
	}
}

func TestExecutionResultForRemoteDispatchMapsRemoteResult(t *testing.T) {
	result := executionResultForRemoteDispatch(routing.Decision{
		ChosenTarget:   "remote-executor",
		TargetUnitID:   "remote-unit",
		TargetRole:     unitref.RoleSpecialist,
		TargetModelRef: "remote-model",
	}, &RemoteDispatchResult{
		SocketPath:    "/tmp/remote.sock",
		Status:        "ok",
		ResultSummary: "remote completed",
		Confidence:    0.8,
	})

	if result.Plan.ExecutionMode != "remote_rpc" {
		t.Fatalf("unexpected execution mode: %q", result.Plan.ExecutionMode)
	}
	if result.Plan.TargetUnitID != "remote-unit" || result.Plan.ChosenExecutor != "remote-executor" {
		t.Fatalf("unexpected plan target: %+v", result.Plan)
	}
	if !result.HostResult.Handled {
		t.Fatalf("expected ok remote status to mark host result handled")
	}
	if result.HostResult.Output != "remote completed" {
		t.Fatalf("unexpected remote output: %q", result.HostResult.Output)
	}
	if result.HostResult.Metadata["socket_path"] != "/tmp/remote.sock" {
		t.Fatalf("missing socket metadata: %+v", result.HostResult.Metadata)
	}
}

type fakeRemoteDispatcher struct{}

func (fakeRemoteDispatcher) DispatchRemote(context.Context, Task) (*RemoteDispatchResult, error) {
	return nil, nil
}
