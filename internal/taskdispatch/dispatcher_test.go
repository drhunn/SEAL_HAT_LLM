package taskdispatch

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/drhunn/SEAL_HAT_LLM/internal/modality"
	"github.com/drhunn/SEAL_HAT_LLM/internal/runtime"
	"github.com/drhunn/SEAL_HAT_LLM/internal/taskrpc"
)

type recordingClient struct {
	resp taskrpc.RunTaskResponse
	err  error
	req  taskrpc.RunTaskRequest
}

func (c *recordingClient) RunTask(ctx context.Context, req taskrpc.RunTaskRequest) (taskrpc.RunTaskResponse, error) {
	c.req = req
	return c.resp, c.err
}

func TestRequestFromTaskMapsRuntimeTask(t *testing.T) {
	req := RequestFromTask("parent-1", runtime.Task{
		ID:                          " task-1 ",
		Summary:                     " do thing ",
		Class:                       " analysis ",
		PrimaryModality:             modality.Text,
		SecondaryModalities:         []modality.Type{modality.Image},
		CrossModalGroundingRequired: true,
		AllowTextOnlyFallback:       true,
		PreferredUnitID:             " spec-1 ",
		AssetRefs:                   []string{" sandbox://a ", ""},
		Prompt:                      " solve it ",
	})
	if req.ProtocolVersion != taskrpc.ProtocolVersion {
		t.Fatalf("protocol mismatch: %q", req.ProtocolVersion)
	}
	if req.ParentUnitID != "parent-1" || req.TargetUnitID != "spec-1" || req.TaskID != "task-1" {
		t.Fatalf("unexpected mapped ids: %+v", req)
	}
	if req.PrimaryModality != "text" || len(req.SecondaryModalities) != 1 || req.SecondaryModalities[0] != "image" {
		t.Fatalf("unexpected modality mapping: %+v", req)
	}
	if len(req.AssetRefs) != 1 || req.AssetRefs[0] != "sandbox://a" {
		t.Fatalf("unexpected asset refs: %+v", req.AssetRefs)
	}
}

func TestDispatchRequiresPreferredUnitID(t *testing.T) {
	dispatcher, err := New(Options{ParentUnitID: "parent-1", Resolver: StaticSocketResolver{"spec-1": "/tmp/spec-1.sock"}})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	_, err = dispatcher.Dispatch(context.Background(), runtime.Task{ID: "task-1"})
	if err == nil || !strings.Contains(err.Error(), "preferred unit id") {
		t.Fatalf("expected preferred unit id error, got %v", err)
	}
}

func TestDispatchRequiresRegisteredSocket(t *testing.T) {
	dispatcher, err := New(Options{ParentUnitID: "parent-1", Resolver: StaticSocketResolver{"other": "/tmp/other.sock"}})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	_, err = dispatcher.Dispatch(context.Background(), runtime.Task{ID: "task-1", PreferredUnitID: "spec-1"})
	if err == nil || !strings.Contains(err.Error(), "no task rpc socket") {
		t.Fatalf("expected socket resolution error, got %v", err)
	}
}

func TestDispatchCallsResolvedClient(t *testing.T) {
	client := &recordingClient{resp: taskrpc.RunTaskResponse{ProtocolVersion: taskrpc.ProtocolVersion, TaskID: "task-1", SpecialistUnitID: "spec-1", Status: "ok"}}
	dispatcher, err := New(Options{
		ParentUnitID: "parent-1",
		Resolver:     StaticSocketResolver{"spec-1": "/tmp/spec-1.sock"},
		ClientFactory: func(socketPath string) Client {
			if socketPath != "/tmp/spec-1.sock" {
				t.Fatalf("unexpected socket path: %s", socketPath)
			}
			return client
		},
	})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	result, err := dispatcher.Dispatch(context.Background(), runtime.Task{ID: "task-1", Summary: "do thing", Class: "analysis", PreferredUnitID: "spec-1"})
	if err != nil {
		t.Fatalf("Dispatch returned error: %v", err)
	}
	if result == nil || result.TargetUnitID != "spec-1" || result.SocketPath != "/tmp/spec-1.sock" {
		t.Fatalf("unexpected dispatch result: %+v", result)
	}
	if client.req.ParentUnitID != "parent-1" || client.req.TargetUnitID != "spec-1" || client.req.TaskID != "task-1" {
		t.Fatalf("unexpected request sent to client: %+v", client.req)
	}
}

func TestDispatchReturnsPartialResultOnClientError(t *testing.T) {
	client := &recordingClient{resp: taskrpc.RunTaskResponse{ProtocolVersion: taskrpc.ProtocolVersion, TaskID: "task-1", Status: "error", ErrorText: "boom"}, err: errors.New("boom")}
	dispatcher, err := New(Options{
		ParentUnitID:  "parent-1",
		Resolver:      StaticSocketResolver{"spec-1": "/tmp/spec-1.sock"},
		ClientFactory: func(socketPath string) Client { return client },
	})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	result, err := dispatcher.Dispatch(context.Background(), runtime.Task{ID: "task-1", PreferredUnitID: "spec-1"})
	if err == nil || !strings.Contains(err.Error(), "dispatch task") {
		t.Fatalf("expected dispatch error, got %v", err)
	}
	if result == nil || result.Response.ErrorText != "boom" {
		t.Fatalf("expected partial result with response, got %+v", result)
	}
}
