package value

import "testing"

func TestValue(t *testing.T) {
	if Value() != 1 {
		t.Fatalf("Value() = %d, want 1", Value())
	}
}
