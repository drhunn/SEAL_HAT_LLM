package memory

import "testing"

func TestZeroVector(t *testing.T) {
	got := ZeroVector(4)
	want := "[0,0,0,0]"
	if got != want {
		t.Fatalf("ZeroVector(4) = %q, want %q", got, want)
	}
}

func TestEncodeVector(t *testing.T) {
	got := EncodeVector([]float32{1.0, 2.5, 0.0})
	want := "[1,2.5,0]"
	if got != want {
		t.Fatalf("EncodeVector mismatch: got %q want %q", got, want)
	}
}
