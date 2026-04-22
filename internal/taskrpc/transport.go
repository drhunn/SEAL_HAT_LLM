package taskrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RunTaskRequest struct {
	TaskID                      string            `json:"task_id"`
	ParentUnitID                string            `json:"parent_unit_id"`
	TargetUnitID                string            `json:"target_unit_id"`
	TaskClass                   string            `json:"task_class"`
	Summary                     string            `json:"summary"`
	Prompt                      string            `json:"prompt"`
	PrimaryModality             string            `json:"primary_modality"`
	SecondaryModalities         []string          `json:"secondary_modalities"`
	CrossModalGroundingRequired bool              `json:"cross_modal_grounding_required"`
	AllowTextOnlyFallback       bool              `json:"allow_text_only_fallback"`
	AssetRefs                   []string          `json:"asset_refs"`
	Constraints                 map[string]string `json:"constraints"`
	RouteEpisodeID              string            `json:"route_episode_id"`
}

type RunTaskResponse struct {
	TaskID           string   `json:"task_id"`
	SpecialistUnitID string   `json:"specialist_unit_id"`
	Status           string   `json:"status"`
	ResultSummary    string   `json:"result_summary"`
	OutputJSON       string   `json:"output_json"`
	ArtifactRefs     []string `json:"artifact_refs"`
	Confidence       float64  `json:"confidence"`
	Warnings         []string `json:"warnings"`
	Signals          []string `json:"signals"`
	ErrorText        string   `json:"error_text,omitempty"`
}

type TaskHandler interface {
	RunTask(context.Context, RunTaskRequest) (RunTaskResponse, error)
}

type Server struct {
	socketPath string
	handler    TaskHandler
}

func NewServer(socketPath string, handler TaskHandler) *Server {
	return &Server{socketPath: socketPath, handler: handler}
}

func (s *Server) Serve(ctx context.Context) error {
	if s == nil || s.handler == nil {
		return fmt.Errorf("task rpc handler is required")
	}
	if s.socketPath == "" {
		return fmt.Errorf("task rpc socket path is required")
	}
	if err := os.MkdirAll(filepath.Dir(s.socketPath), 0o755); err != nil {
		return fmt.Errorf("create task rpc socket dir: %w", err)
	}
	if err := os.RemoveAll(s.socketPath); err != nil {
		return fmt.Errorf("remove stale task rpc socket: %w", err)
	}
	listener, err := net.Listen("unix", s.socketPath)
	if err != nil {
		return fmt.Errorf("listen on task rpc socket: %w", err)
	}
	defer func() {
		_ = listener.Close()
		_ = os.Remove(s.socketPath)
	}()

	var closeOnce sync.Once
	go func() {
		<-ctx.Done()
		closeOnce.Do(func() { _ = listener.Close() })
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				return fmt.Errorf("accept task rpc connection: %w", err)
			}
		}
		go s.handleConn(ctx, conn)
	}
}

func (s *Server) handleConn(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(30 * time.Second))

	var req RunTaskRequest
	if err := json.NewDecoder(conn).Decode(&req); err != nil {
		_ = json.NewEncoder(conn).Encode(RunTaskResponse{Status: "error", ErrorText: fmt.Sprintf("decode request: %v", err)})
		return
	}
	resp, err := s.handler.RunTask(ctx, req)
	if err != nil {
		resp.Status = "error"
		resp.ErrorText = err.Error()
	}
	if resp.Status == "" {
		resp.Status = "ok"
	}
	_ = json.NewEncoder(conn).Encode(resp)
}

type Client struct {
	socketPath string
}

func NewClient(socketPath string) *Client {
	return &Client{socketPath: socketPath}
}

func (c *Client) RunTask(ctx context.Context, req RunTaskRequest) (RunTaskResponse, error) {
	if c == nil || c.socketPath == "" {
		return RunTaskResponse{}, fmt.Errorf("task rpc socket path is required")
	}
	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "unix", c.socketPath)
	if err != nil {
		return RunTaskResponse{}, fmt.Errorf("dial task rpc socket: %w", err)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(30 * time.Second))
	if err := json.NewEncoder(conn).Encode(req); err != nil {
		return RunTaskResponse{}, fmt.Errorf("encode task rpc request: %w", err)
	}
	var resp RunTaskResponse
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		return RunTaskResponse{}, fmt.Errorf("decode task rpc response: %w", err)
	}
	if resp.Status == "error" && resp.ErrorText != "" {
		return resp, fmt.Errorf(resp.ErrorText)
	}
	return resp, nil
}
