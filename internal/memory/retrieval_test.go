package memory

import "testing"

func TestZeroVectorLiteral(t *testing.T) {
	got := ZeroVectorLiteral(4).String()
	want := "[0,0,0,0]"
	if got != want {
		t.Fatalf("ZeroVectorLiteral(4) = %q, want %q", got, want)
	}
}

func TestNewVectorLiteral(t *testing.T) {
	got := NewVectorLiteral([]float32{1.0, 2.5, 0.0}).String()
	want := "[1,2.5,0]"
	if got != want {
		t.Fatalf("NewVectorLiteral mismatch: got %q want %q", got, want)
	}
}

func TestZeroVectorCompatibilityWrapper(t *testing.T) {
	got := ZeroVector(4)
	want := "[0,0,0,0]"
	if got != want {
		t.Fatalf("ZeroVector(4) = %q, want %q", got, want)
	}
}

func TestEncodeVectorCompatibilityWrapper(t *testing.T) {
	got := EncodeVector([]float32{1.0, 2.5, 0.0})
	want := "[1,2.5,0]"
	if got != want {
		t.Fatalf("EncodeVector mismatch: got %q want %q", got, want)
	}
}
