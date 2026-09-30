package claude

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSettingsRoundTripPreservesUnknownFields(t *testing.T) {
	root := t.TempDir()
	os.Mkdir(filepath.Join(root, ".claude"), 0755)
	path := filepath.Join(root, settingsFile)
	initial := []byte(`{"permissions":{"deny":["Bash(rm *)"]},"env":{"FOO":"bar"},"hooks":{"SessionStart":[{"hooks":[{"type":"command","command":"startup","timeout":11}]}],"PreToolUse":[{"matcher":"Bash","future_matcher":true,"hooks":[{"type":"command","command":"my-agentcap-audit","timeout":19,"async":true,"if":"Bash(git *)","future_handler":42}]}]}}`)
	os.WriteFile(path, initial, 0600)
	if IsInstalled(root) {
		t.Fatal("substring incorrectly treated as ownership")
	}
	if err := Install(root, false); err != nil {
		t.Fatal(err)
	}
	first, _ := os.ReadFile(path)
	if err := Install(root, false); err != nil {
		t.Fatal(err)
	}
	second, _ := os.ReadFile(path)
	if !bytes.Equal(first, second) {
		t.Fatal("idempotent installation rewrote settings")
	}
	if err := Remove(root); err != nil {
		t.Fatal(err)
	}
	final, _ := os.ReadFile(path)
	var before, after any
	json.Unmarshal(initial, &before)
	json.Unmarshal(final, &after)
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("configuration changed: %s", final)
	}
	if err := Remove(root); err != nil {
		t.Fatal(err)
	}
	absent, _ := os.ReadFile(path)
	if !bytes.Equal(final, absent) {
		t.Fatal("absent removal rewrote settings")
	}
	if _, err := os.Stat(filepath.Join(root, "CLAUDE.md")); !os.IsNotExist(err) {
		t.Fatal("installer changed prompt instructions")
	}
}

func TestMalformedSettingsRemainUntouched(t *testing.T) {
	for _, initial := range []string{`not json`, `null`, `[]`, `{"hooks":null}`, `{"hooks":[]}`, `{"hooks":{"PreToolUse":null}}`, `{"hooks":{"PreToolUse":[{"hooks":"broken"}]}}`} {
		root := t.TempDir()
		os.Mkdir(filepath.Join(root, ".claude"), 0755)
		path := filepath.Join(root, settingsFile)
		os.WriteFile(path, []byte(initial), 0600)
		if err := Install(root, false); err == nil {
			t.Fatalf("accepted %s", initial)
		}
		if err := Remove(root); err == nil {
			t.Fatalf("removal accepted %s", initial)
		}
		actual, _ := os.ReadFile(path)
		if string(actual) != initial {
			t.Fatalf("destroyed %s", initial)
		}
	}
}
