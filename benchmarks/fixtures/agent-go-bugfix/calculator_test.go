package calculator

import "testing"

func TestAdd(t *testing.T) {
	if got := Add(7, 5); got != 12 {
		t.Fatalf("Add(7, 5) = %d, want 12", got)
	}
}
