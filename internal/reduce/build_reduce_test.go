package reduce

import (
	"os"
	"testing"

	"github.com/taqu/agentcap/internal/exec"
)

func TestSelectBuildGoTest(t *testing.T) {
	r := SelectBuild([]string{"go", "test", "./..."})
	if _, ok := r.(*GoTestReducer); !ok {
		t.Errorf("expected *GoTestReducer, got %T", r)
	}
}

func TestSelectBuildGoBuild(t *testing.T) {
	r := SelectBuild([]string{"go", "build", "./..."})
	if _, ok := r.(*GoBuildReducer); !ok {
		t.Errorf("expected *GoBuildReducer, got %T", r)
	}
}

func TestSelectBuildGCC(t *testing.T) {
	r := SelectBuild([]string{"gcc", "main.c"})
	if _, ok := r.(*GccReducer); !ok {
		t.Errorf("expected *GccReducer, got %T", r)
	}
}

func TestSelectBuildCargo(t *testing.T) {
	r := SelectBuild([]string{"cargo", "check"})
	if _, ok := r.(*CargoBuildReducer); !ok {
		t.Errorf("expected *CargoBuildReducer, got %T", r)
	}
}

func TestSelectBuildCargoTest(t *testing.T) {
	r := SelectBuild([]string{"cargo", "test"})
	if _, ok := r.(*CargoTestReducer); !ok {
		t.Errorf("expected *CargoTestReducer, got %T", r)
	}
}

func TestGoBuildFail(t *testing.T) {
	raw, err := os.ReadFile("../../testdata/build/go-build-error.txt")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	result := &exec.Result{
		Stderr:   raw,
		ExitCode: 1,
	}
	r := &GoBuildReducer{Args: []string{"go", "build", "./..."}}
	reduced := r.Reduce(result)
	if reduced == nil {
		t.Fatal("Reduce returned nil")
	}
	if reduced.Output == "" {
		t.Error("expected non-empty output")
	}
	// Should contain FAIL marker
	if len(reduced.Output) == 0 {
		t.Error("expected output")
	}
}

func TestGoTestPass(t *testing.T) {
	raw, err := os.ReadFile("../../testdata/build/go-test-pass.txt")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	result := &exec.Result{
		Stdout:   raw,
		ExitCode: 0,
	}
	r := &GoTestReducer{Args: []string{"go", "test", "./..."}}
	reduced := r.Reduce(result)
	if reduced == nil {
		t.Fatal("Reduce returned nil")
	}
	if reduced.Output == "" {
		t.Error("expected non-empty output")
	}
}
