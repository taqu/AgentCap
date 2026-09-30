package errcatalog

import "testing"

func TestKnownKeys(t *testing.T) {
	for _, k := range []string{"billing.missing", "tax.invalid", "audit.locked"} {
		if Lookup(k) == "internal" {
			t.Errorf("%s unmapped", k)
		}
	}
	if Status("nope.nothing") != 500 {
		t.Error("unknown keys must be 500")
	}
}
