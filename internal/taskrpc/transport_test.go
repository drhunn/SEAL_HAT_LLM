package taskrpc

import (
	"context"
	"path/filepath"
	"testing"
	"time"
)

type testHandler struct {
	resp RunTaskResponse
	fn   func(context.Context, RunTaskRequest) (RunTaskResponse, error)
}

func (h testHandler) RunTask(ctx context.Context, req RunTaskRequest) (RunTaskResponse, error) {
	if h.fn != nil {
		return h.fn(ctx, req)
	}
	return h.resp, nil
}

func TestServerClientRoundTrip(t *testing.T) {
	socketPath := filepath.Join(t.TempDir(), "taskrpc.sock")
	server := NewServer(socketPath, testHandler{resp: RunTaskResponse{TaskID: "task-1", SpecialistUnitID: "spec-1", Status: "ok", ResultSummary: "done", Confidence: 0.9}})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Serve(ctx)
	}()
	time.Sleep(100 * time.Millisecond)

	client := NewClient(socketPath)
	resp, err := client.RunTask(context.Background(), RunTaskRequest{TaskID: "task-1", TargetUnitID: "spec-1", Summary: "do thing"})
	if err != nil {
		t.Fatalf("RunTask returned error: %v", err)
	}
	if resp.TaskID != "task-1" || resp.SpecialistUnitID != "spec-1" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestServerClientPropagatesHandlerError(t *testing.T) {
	socketPath := filepath.Join(t.TempDir(), "taskrpc.sock")
	server := NewServer(socketPath, testHandler{fn: func(ctx context.Context, req RunTaskRequest) (RunTaskResponse, error) {
		return RunTaskResponse{TaskID: req.TaskID, SpecialistUnitID: req.TargetUnitID}, context.DeadlineExceeded
	}})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = server.Serve(ctx)
	}()
	time.Sleep(100 * time.Millisecond)

	client := NewClient(socketPath)
	_, err := client.RunTask(context.Background(), RunTaskRequest{TaskID: "task-2", TargetUnitID: "spec-2", Summary: "fail thing"})
	if err == nil {
		t.Fatalf("expected handler error")
	}
}
