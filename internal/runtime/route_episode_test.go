package runtime

import (
	"errors"
	"testing"

	"github.com/drhunn/SEAL_HAT_LLM/internal/execution"
	"github.com/drhunn/SEAL_HAT_LLM/internal/modelhost"
)

func TestRouteEpisodeStatusForOutcome(t *testing.T) {
	tests := []struct {
		name string
		result execution.Result
		err error
		wantStatus string
		wantError string
	}{
		{
			name: "success",
			result: execution.Result{HostResult: modelhost.Result{Handled: true}},
			wantStatus: "succeeded",
		},
		{
			name: "unhandled",
			result: execution.Result{HostResult: modelhost.Result{Handled: false}},
			wantStatus: "unhandled",
		},
		{
			name: "failure",
			result: execution.Result{HostResult: modelhost.Result{Handled: false}},
			err: errors.New("boom"),
			wantStatus: "failed",
			wantError: "boom",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, errorText := routeEpisodeStatusForOutcome(tt.result, tt.err)
			if status != tt.wantStatus || errorText != tt.wantError {
				t.Fatalf("got status=%q error=%q want status=%q error=%q", status, errorText, tt.wantStatus, tt.wantError)
			}
		})
	}
}
