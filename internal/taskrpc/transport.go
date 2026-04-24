package taskrpc

import (
	"context"
	"errors"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	ProtocolVersion        = "taskrpc.v1"
	DefaultMaxMessageBytes = int64(1 << 20)
	DefaultDeadline        = 30 * time.Second
)

type RunTaskRequest struct {
	ProtocolVersion            string            `json:"protocol_version"`
	TaskID                     string            `json:"task_id"`
	ParentUnitID               string            `json:"parent_unit_id"`
	TargetUnitID               string            `json:"target_unit_id"`
	TaskClass                  string            `json:"task_class"`
	Summary                    string            `json:"summary"`
	Prompt                     string            `json:"prompt"`
	PrimaryModality            string            `json:"primary_modality"`
	SecondaryModalities        []string          `json:"secondary_modalities"`
	CrossModalGroundingRequired bool             `json:"cross_modal_grounding_required"`
	AllowTextOnlyFallback      bool              `json:"allow_text_only_fallback"`
	AssetRefs                  []string          `json:"asset_refs"`
	Constraints                map[string]string `json:"constraints"`
	RouteEpisodeID             string            `json:"route_episode_id"`
}

type RunTaskResponse struct {
	ProtocolVersion string   `json:"protocol_version"`
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

type ServerOptions struct {
	MaxRequestBytes int64
	Deadline        time.Duration
}

type Server struct {
	socketPath      string
	handler         TaskHandler
	maxRequestBytes int64
	deadline        time.Duration
}

func NewServer(socketPath string, handler TaskHandler) *Server {
	return NewServerWithOptions(socketPath, handler, ServerOptions{})
}

func NewServerWithOptions(socketPath string, handler TaskHandler, opts ServerOptions) *Server {
	maxRequestBytes := opts.MaxRequestBytes
	if maxRequestBytes <= 0 {
		maxRequestBytes = DefaultMaxMessageBytes
	}
	deadline := opts.Deadline
	if deadline <= 0 {
		deadline = DefaultDeadline
	}
	return &Server{socketPath: socketPath, handler: handler, maxRequestBytes: maxRequestBytes, deadline: deadline}
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
	if s.deadline > 0 {
		_ = conn.SetDeadline(time.Now().Add(s.deadline))
	}

	var req RunTaskRequest
	if err := decodeLimited(conn, s.maxRequestBytes, &req); err != nil {
		_ = encodeResponse(conn, RunTaskResponse{Status: "error", ErrorText: fmt.Sprintf("decode request: %v", err)})
		return
	}
	if req.ProtocolVersion != ProtocolVersion {
		_ = encodeResponse(conn, RunTaskResponse{TaskID: req.TaskID, Status: "error", ErrorText: fmt.Sprintf("unsupported task rpc protocol version %q", req.ProtocolVersion)})
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
	_ = encodeResponse(conn, resp)
}

type ClientOptions struct {
	MaxRequestBytes  int64
	MaxResponseBytes int64
	Deadline         time.Duration
}

type Client struct {
	socketPath       string
	maxRequestBytes  int64
	maxResponseBytes int64
	deadline         time.Duration
}

func NewClient(socketPath string) *Client {
	return NewClientWithOptions(socketPath, ClientOptions{})
}

func NewClientWithOptions(socketPath string, opts ClientOptions) *Client {
	maxRequestBytes := opts.MaxRequestBytes
	if maxRequestBytes <= 0 {
		maxRequestBytes = DefaultMaxMessageBytes
	}
	maxResponseBytes := opts.MaxResponseBytes
	if maxResponseBytes <= 0 {
		maxResponseBytes = DefaultMaxMessageBytes
	}
	deadline := opts.Deadline
	if deadline <= 0 {
		deadline = DefaultDeadline
	}
	return &Client{socketPath: socketPath, maxRequestBytes: maxRequestBytes, maxResponseBytes: maxResponseBytes, deadline: deadline}
}

func (c *Client) RunTask(ctx context.Context, req RunTaskRequest) (RunTaskResponse, error) {
	if c == nil || c.socketPath == "" {
		return RunTaskResponse{}, fmt.Errorf("task rpc socket path is required")
	}
	if req.ProtocolVersion == "" {
		req.ProtocolVersion = ProtocolVersion
	}
	payload, err := json.Marshal(req)
	if err != nil {
		return RunTaskResponse{}, fmt.Errorf("encode task rpc request: %w", err)
	}
	if int64(len(payload)) > c.maxRequestBytes {
		return RunTaskResponse{}, fmt.Errorf("task rpc request exceeds max bytes %d", c.maxRequestBytes)
	}

	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "unix", c.socketPath)
	if err != nil {
		return RunTaskResponse{}, fmt.Errorf("dial task rpc socket: %w", err)
	}
	defer conn.Close()
	if c.deadline > 0 {
		_ = conn.SetDeadline(time.Now().Add(c.deadline))
	}
	if _, err := conn.Write(append(payload, '\n')); err != nil {
		return RunTaskResponse{}, fmt.Errorf("write task rpc request: %w", err)
	}
	var resp RunTaskResponse
	if err := decodeLimited(conn, c.maxResponseBytes, &resp); err != nil {
		return RunTaskResponse{}, fmt.Errorf("decode task rpc response: %w", err)
	}
	if resp.ProtocolVersion != ProtocolVersion {
		return resp, fmt.Errorf("unsupported task rpc response protocol version %q", resp.ProtocolVersion)
	}
	if resp.Status == "error" && resp.ErrorText != "" {
		return resp, errors.New(resp.ErrorText)
	}
	return resp, nil
}

func encodeResponse(w io.Writer, resp RunTaskResponse) error {
	resp.ProtocolVersion = ProtocolVersion
	return json.NewEncoder(w).Encode(resp)
}

func decodeLimited(r io.Reader, maxBytes int64, out any) error {
	if maxBytes <= 0 {
		maxBytes = DefaultMaxMessageBytes
	}
	limited := &io.LimitedReader{R: r, N: maxBytes + 1}
	if err := json.NewDecoder(limited).Decode(out); err != nil {
		if limited.N <= 0 {
			return fmt.Errorf("message exceeds max bytes %d", maxBytes)
		}
		return err
	}
	if limited.N <= 0 {
		return fmt.Errorf("message exceeds max bytes %d", maxBytes)
	}
	return nil
}
