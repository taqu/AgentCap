package common

import (
	"os"
	"testing"

	"github.com/taqu/agentcap/internal/integration/protocol"
)

func TestIsBypassed(t *testing.T) {
	// Not set → false
	os.Unsetenv(protocol.EnvBypass)
	if IsBypassed() {
		t.Error("IsBypassed: expected false when env not set")
	}

	// Set to "1" → true
	os.Setenv(protocol.EnvBypass, "1")
	defer os.Unsetenv(protocol.EnvBypass)
	if !IsBypassed() {
		t.Error("IsBypassed: expected true when ACAP_BYPASS=1")
	}

	// Set to "0" → false
	os.Setenv(protocol.EnvBypass, "0")
	if IsBypassed() {
		t.Error("IsBypassed: expected false when ACAP_BYPASS=0")
	}
}

func TestIsRecursive(t *testing.T) {
	os.Unsetenv(protocol.EnvDepth)
	if IsRecursive() {
		t.Error("IsRecursive: expected false when env not set")
	}

	os.Setenv(protocol.EnvDepth, "1")
	defer os.Unsetenv(protocol.EnvDepth)
	if !IsRecursive() {
		t.Error("IsRecursive: expected true when ACAP_INTERCEPT_DEPTH=1")
	}

	os.Setenv(protocol.EnvDepth, "0")
	if IsRecursive() {
		t.Error("IsRecursive: expected false when ACAP_INTERCEPT_DEPTH=0")
	}
}

func TestIsAcapCommand(t *testing.T) {
	cases := []struct {
		args []string
		want bool
	}{
		{nil, false},
		{[]string{}, false},
		{[]string{"acap"}, true},
		{[]string{"acap.exe"}, true},
		{[]string{"/usr/local/bin/acap"}, true},
		{[]string{"C:\\tools\\acap.exe"}, true},
		{[]string{"echo"}, false},
		{[]string{"git"}, false},
		{[]string{"acapty"}, false},
	}
	for _, c := range cases {
		got := IsAcapCommand(c.args)
		if got != c.want {
			t.Errorf("IsAcapCommand(%v) = %v, want %v", c.args, got, c.want)
		}
	}
}

func TestEnvWithDepth(t *testing.T) {
	// No existing ACAP_INTERCEPT_DEPTH → appends "ACAP_INTERCEPT_DEPTH=1"
	env := []string{"FOO=bar", "BAZ=qux"}
	result := EnvWithDepth(env)
	key := protocol.EnvDepth + "=1"
	found := false
	for _, e := range result {
		if e == key {
			found = true
		}
	}
	if !found {
		t.Errorf("EnvWithDepth: expected %q in result %v", key, result)
	}
	if len(result) != len(env)+1 {
		t.Errorf("EnvWithDepth: expected len %d, got %d", len(env)+1, len(result))
	}

	// Existing ACAP_INTERCEPT_DEPTH=2 → becomes 3
	env2 := []string{"FOO=bar", protocol.EnvDepth + "=2"}
	result2 := EnvWithDepth(env2)
	key2 := protocol.EnvDepth + "=3"
	found2 := false
	for _, e := range result2 {
		if e == key2 {
			found2 = true
		}
	}
	if !found2 {
		t.Errorf("EnvWithDepth: expected %q in result %v", key2, result2)
	}
	if len(result2) != len(env2) {
		t.Errorf("EnvWithDepth: expected same length %d, got %d", len(env2), len(result2))
	}
}
