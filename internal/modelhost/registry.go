package modelhost

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

type Request struct {
	TaskSummary   string
	TaskClass     string
	ExecutionMode string
	Executor      string
	AssetRefs     []string
	Prompt        string
}

type Result struct {
	HostName     string
	Executor     string
	Output       string
	Handled      bool
	Metadata     map[string]string
}

type Host interface {
	Name() string
	Execute(ctx context.Context, req Request) (Result, error)
}

type Registry struct {
	mu    sync.RWMutex
	hosts map[string]Host
}

func NewRegistry() *Registry {
	return &Registry{hosts: make(map[string]Host)}
}

func (r *Registry) Register(executor string, host Host) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.hosts[executor] = host
}

func (r *Registry) Resolve(executor string) (Host, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	h, ok := r.hosts[executor]
	return h, ok
}

func (r *Registry) Execute(ctx context.Context, executor string, req Request) (Result, error) {
	host, ok := r.Resolve(executor)
	if !ok {
		return Result{}, fmt.Errorf("no model host registered for executor %q", executor)
	}
	return host.Execute(ctx, req)
}

func (r *Registry) Executors() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]string, 0, len(r.hosts))
	for executor := range r.hosts {
		out = append(out, executor)
	}
	sort.Strings(out)
	return out
}

type StaticHost struct {
	hostName string
	output   string
}

func NewStaticHost(hostName, output string) *StaticHost {
	return &StaticHost{hostName: hostName, output: output}
}

func (h *StaticHost) Name() string {
	return h.hostName
}

func (h *StaticHost) Execute(_ context.Context, req Request) (Result, error) {
	return Result{
		HostName: h.hostName,
		Executor: req.Executor,
		Output:   h.output,
		Handled:  true,
		Metadata: map[string]string{
			"task_class": req.TaskClass,
			"execution_mode": req.ExecutionMode,
		},
	}, nil
}
