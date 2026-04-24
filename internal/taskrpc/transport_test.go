package taskrpc

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"strings"
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
	startTestServer(t, server, socketPath)

	client := NewClient(socketPath)
	resp, err := client.RunTask(context.Background(), RunTaskRequest{TaskID: "task-1", TargetUnitID: "spec-1", Summary: "do thing"})
	if err != nil {
		t.Fatalf("RunTask returned error: %v", err)
	}
	if resp.ProtocolVersion != ProtocolVersion {
		t.Fatalf("unexpected protocol version: %q", resp.ProtocolVersion)
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
	startTestServer(t, server, socketPath)

	client := NewClient(socketPath)
	_, err := client.RunTask(context.Background(), RunTaskRequest{TaskID: "task-2", TargetUnitID: "spec-2", Summary: "fail thing"})
	if err == nil || !strings.Contains(err.Error(), context.DeadlineExceeded.Error()) {
		t.Fatalf("expected propagated handler error, got %v", err)
	}
}

func TestClientRejectsOversizedRequest(t *testing.T) {
	client := NewClientWithOptions(filepath.Join(t.TempDir(), "missing.sock"), ClientOptions{MaxRequestBytes: 32})
	_, err := client.RunTask(context.Background(), RunTaskRequest{TaskID: "task-oversized", Summary: strings.Repeat("x", 128)})
	if err == nil || !strings.Contains(err.Error(), "exceeds max bytes") {
		t.Fatalf("expected oversized request error, got %v", err)
	}
}

func TestServerRejectsMalformedJSON(t *testing.T) {
	socketPath := filepath.Join(t.TempDir(), "taskrpc.sock")
	server := NewServer(socketPath, testHandler{resp: RunTaskResponse{Status: "ok"}})
	startTestServer(t, server, socketPath)

	resp := writeRawRequest(t, socketPath, []byte("{bad json\n"))
	if resp.Status != "error" || !strings.Contains(resp.ErrorText, "decode request") {
		t.Fatalf("expected decode error response, got %+v", resp)
	}
}

func TestServerRejectsUnsupportedProtocolVersion(t *testing.T) {
	socketPath := filepath.Join(t.TempDir(), "taskrpc.sock")
	server := NewServer(socketPath, testHandler{resp: RunTaskResponse{Status: "ok"}})
	startTestServer(t, server, socketPath)

	payload, err := json.Marshal(RunTaskRequest{ProtocolVersion: "taskrpc.v0", TaskID: "task-old"})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	resp := writeRawRequest(t, socketPath, append(payload, '\n'))
	if resp.Status != "error" || !strings.Contains(resp.ErrorText, "unsupported task rpc protocol version") {
		t.Fatalf("expected protocol error response, got %+v", resp)
	}
}

func TestServerRejectsOversizedRequest(t *testing.T) {
	socketPath := filepath.Join(t.TempDir(), "taskrpc.sock")
	server := NewServerWithOptions(socketPath, testHandler{resp: RunTaskResponse{Status: "ok"}}, ServerOptions{MaxRequestBytes: 64})
	startTestServer(t, server, socketPath)

	payload, err := json.Marshal(RunTaskRequest{ProtocolVersion: ProtocolVersion, TaskID: "task-big", Summary: strings.Repeat("x", 256)})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	resp := writeRawRequest(t, socketPath, append(payload, '\n'))
	if resp.Status != "error" || !strings.Contains(resp.ErrorText, "exceeds max bytes") {
		t.Fatalf("expected oversized request error response, got %+v", resp)
	}
}

func TestServerRemovesSocketOnShutdown(t *testing.T) {
	socketPath := filepath.Join(t.TempDir(), "taskrpc.sock")
	server := NewServer(socketPath, testHandler{resp: RunTaskResponse{Status: "ok"}})
	ctx, cancel := context.WithCancel(context.Background())
	serveErr := make(chan error, 1)
	go func() { serveErr <- server.Serve(ctx) }()
	waitForSocket(t, socketPath)

	cancel()
	if err := <-serveErr; err != nil {
		t.Fatalf("server returned error: %v", err)
	}
	if _, err := os.Stat(socketPath); !os.IsNotExist(err) {
		t.Fatalf("expected socket to be removed after shutdown, got %v", err)
	}
}

func startTestServer(t *testing.T, server *Server, socketPath string) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	serveErr := make(chan error, 1)
	go func() { serveErr <- server.Serve(ctx) }()
	waitForSocket(t, socketPath)
	t.Cleanup(func() {
		cancel()
		select {
		case err := <-serveErr:
			if err != nil {
				t.Fatalf("server returned error during cleanup: %v", err)
			}
		case <-time.After(time.Second):
			t.Fatalf("timed out waiting for server shutdown")
		}
	})
}

func waitForSocket(t *testing.T, socketPath string) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("unix", socketPath, 50*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for socket %s", socketPath)
}

func writeRawRequest(t *testing.T, socketPath string, payload []byte) RunTaskResponse {
	t.Helper()
	conn, err := net.DialTimeout("unix", socketPath, time.Second)
	if err != nil {
		t.Fatalf("dial socket: %v", err)
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatalf("set deadline: %v", err)
	}
	if _, err := conn.Write(payload); err != nil {
		t.Fatalf("write request: %v", err)
	}
	var resp RunTaskResponse
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return resp
}
