package common

import (
	"testing"
)

func TestMapAgentSession(t *testing.T) {
	// Empty input → empty output
	if got := MapAgentSession(""); got != "" {
		t.Errorf("MapAgentSession(\"\") = %q, want \"\"", got)
	}

	// Same input → same output (deterministic)
	id1 := "my-session-token-abc123"
	out1a := MapAgentSession(id1)
	out1b := MapAgentSession(id1)
	if out1a != out1b {
		t.Errorf("MapAgentSession not deterministic: %q != %q", out1a, out1b)
	}

	// Output length is 6
	if len(out1a) != 6 {
		t.Errorf("MapAgentSession length = %d, want 6", len(out1a))
	}

	// Different inputs → different outputs
	id2 := "different-session-xyz999"
	out2 := MapAgentSession(id2)
	if out1a == out2 {
		t.Errorf("MapAgentSession collision: %q and %q both map to %q", id1, id2, out1a)
	}

	// Output is hex (lowercase a-f 0-9)
	for _, c := range out1a {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("MapAgentSession output %q contains non-hex char %c", out1a, c)
		}
	}
}
