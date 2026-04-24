package taskdispatch

import (
	"context"
	"fmt"
	"strings"

	"github.com/drhunn/SEAL_HAT_LLM/internal/runtime"
	"github.com/drhunn/SEAL_HAT_LLM/internal/taskrpc"
)

type SocketResolver interface {
	ResolveSocket(ctx context.Context, unitID string) (string, bool)
}

type StaticSocketResolver map[string]string

func (r StaticSocketResolver) ResolveSocket(_ context.Context, unitID string) (string, bool) {
	socketPath, ok := r[strings.TrimSpace(unitID)]
	return strings.TrimSpace(socketPath), ok && strings.TrimSpace(socketPath) != ""
}

type Client interface {
	RunTask(context.Context, taskrpc.RunTaskRequest) (taskrpc.RunTaskResponse, error)
}

type ClientFactory func(socketPath string) Client

type Dispatcher struct {
	parentUnitID  string
	resolver      SocketResolver
	clientFactory ClientFactory
}

type Options struct {
	ParentUnitID  string
	Resolver      SocketResolver
	ClientFactory ClientFactory
	ClientOptions taskrpc.ClientOptions
}

type Result struct {
	TargetUnitID string
	SocketPath   string
	Request      taskrpc.RunTaskRequest
	Response     taskrpc.RunTaskResponse
}

func New(opts Options) (*Dispatcher, error) {
	parentUnitID := strings.TrimSpace(opts.ParentUnitID)
	if parentUnitID == "" {
		return nil, fmt.Errorf("parent unit id is required")
	}
	if opts.Resolver == nil {
		return nil, fmt.Errorf("task dispatch socket resolver is required")
	}
	clientFactory := opts.ClientFactory
	if clientFactory == nil {
		clientOptions := opts.ClientOptions
		clientFactory = func(socketPath string) Client {
			return taskrpc.NewClientWithOptions(socketPath, clientOptions)
		}
	}
	return &Dispatcher{parentUnitID: parentUnitID, resolver: opts.Resolver, clientFactory: clientFactory}, nil
}

func (d *Dispatcher) Dispatch(ctx context.Context, task runtime.Task) (*Result, error) {
	if d == nil {
		return nil, fmt.Errorf("task dispatcher is required")
	}
	targetUnitID := strings.TrimSpace(task.PreferredUnitID)
	if targetUnitID == "" {
		return nil, fmt.Errorf("task preferred unit id is required for rpc dispatch")
	}
	socketPath, ok := d.resolver.ResolveSocket(ctx, targetUnitID)
	if !ok {
		return nil, fmt.Errorf("no task rpc socket registered for unit %q", targetUnitID)
	}
	client := d.clientFactory(socketPath)
	if client == nil {
		return nil, fmt.Errorf("task rpc client factory returned nil client for socket %q", socketPath)
	}
	req := RequestFromTask(d.parentUnitID, task)
	resp, err := client.RunTask(ctx, req)
	if err != nil {
		return &Result{TargetUnitID: targetUnitID, SocketPath: socketPath, Request: req, Response: resp}, fmt.Errorf("dispatch task %q to unit %q: %w", task.ID, targetUnitID, err)
	}
	return &Result{TargetUnitID: targetUnitID, SocketPath: socketPath, Request: req, Response: resp}, nil
}

func RequestFromTask(parentUnitID string, task runtime.Task) taskrpc.RunTaskRequest {
	secondary := make([]string, 0, len(task.SecondaryModalities))
	for _, item := range task.SecondaryModalities {
		secondary = append(secondary, item.String())
	}
	return taskrpc.RunTaskRequest{
		ProtocolVersion:             taskrpc.ProtocolVersion,
		TaskID:                      strings.TrimSpace(task.ID),
		ParentUnitID:                strings.TrimSpace(parentUnitID),
		TargetUnitID:                strings.TrimSpace(task.PreferredUnitID),
		TaskClass:                   strings.TrimSpace(task.Class),
		Summary:                     strings.TrimSpace(task.Summary),
		Prompt:                      strings.TrimSpace(task.Prompt),
		PrimaryModality:             task.PrimaryModality.String(),
		SecondaryModalities:         secondary,
		CrossModalGroundingRequired: task.CrossModalGroundingRequired,
		AllowTextOnlyFallback:       task.AllowTextOnlyFallback,
		AssetRefs:                   compactStrings(task.AssetRefs),
	}
}

func compactStrings(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		out = append(out, value)
	}
	return out
}
