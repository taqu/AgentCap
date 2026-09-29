// Package workload parses and runs deterministic benchmark workloads.
package workload

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	osexec "os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/taqu/agentcap/internal/bench"
	"gopkg.in/yaml.v3"
)

const SchemaVersion = 1

var namePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._/-]*$`)

// Definition is version 1 of the workload file format.
type Definition struct {
	Version int      `yaml:"version"`
	Name    string   `yaml:"name"`
	Fixture string   `yaml:"fixture"`
	Git     GitSetup `yaml:"git,omitempty"`
	Steps   []Step   `yaml:"steps"`

	fixtureDir string
}

// GitSetup requests a deterministic local repository and baseline commit.
type GitSetup struct {
	Init bool `yaml:"init"`
}

// Step has exactly one operation.
type Step struct {
	Run    *RunStep    `yaml:"run,omitempty"`
	Copy   *CopyStep   `yaml:"copy,omitempty"`
	Write  *WriteStep  `yaml:"write,omitempty"`
	Remove *RemoveStep `yaml:"remove,omitempty"`
	Mkdir  *MkdirStep  `yaml:"mkdir,omitempty"`
}

type RunStep struct {
	Argv   []string `yaml:"argv"`
	Cwd    string   `yaml:"cwd,omitempty"`
	Expect *Expect  `yaml:"expect,omitempty"`
}

type Expect struct {
	Exit *int `yaml:"exit,omitempty"`
}

type CopyStep struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

type WriteStep struct {
	Path    string `yaml:"path"`
	Content string `yaml:"content"`
}

type RemoveStep struct {
	Path string `yaml:"path"`
}

type MkdirStep struct {
	Path string `yaml:"path"`
}

// Load parses and fully validates a workload before any target command runs.
// Fixtures must live below the benchmarks/fixtures directory paired with the
// benchmarks/workloads directory containing path.
func Load(path string) (*Definition, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read workload: %w", err)
	}
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	var d Definition
	if err := dec.Decode(&d); err != nil {
		return nil, fmt.Errorf("parse workload: %w", err)
	}
	var extra any
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return nil, errors.New("parse workload: multiple YAML documents are not supported")
		}
		return nil, fmt.Errorf("parse workload: %w", err)
	}
	if err := d.validate(path); err != nil {
		return nil, err
	}
	return &d, nil
}

func (d *Definition) validate(path string) error {
	if d.Version != SchemaVersion {
		return fmt.Errorf("unsupported workload schema version %d (supported: %d)", d.Version, SchemaVersion)
	}
	if d.Name == "" {
		return errors.New("workload name is required")
	}
	if !namePattern.MatchString(d.Name) || strings.Contains(d.Name, "//") || hasDotDot(d.Name) {
		return fmt.Errorf("invalid workload name %q", d.Name)
	}
	if d.Fixture == "" {
		return errors.New("workload fixture is required")
	}
	benchRoot, err := benchmarkRoot(path)
	if err != nil {
		return err
	}
	fixtureRoot, err := filepath.Abs(filepath.Join(benchRoot, "fixtures"))
	if err != nil {
		return err
	}
	fixture, err := filepath.Abs(filepath.Join(filepath.Dir(path), d.Fixture))
	if err != nil {
		return err
	}
	if !within(fixtureRoot, fixture) {
		return fmt.Errorf("fixture must be inside %s", fixtureRoot)
	}
	fixtureRel, err := filepath.Rel(fixtureRoot, fixture)
	if err != nil {
		return fmt.Errorf("fixture: %w", err)
	}
	fixtureReal, err := securePath(fixtureRoot, fixtureRel, true)
	if err != nil {
		return fmt.Errorf("fixture: %w", err)
	}
	info, err := os.Stat(fixtureReal)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("fixture is not a directory: %s", d.Fixture)
	}
	d.fixtureDir = fixtureReal
	if len(d.Steps) == 0 {
		return errors.New("workload must contain at least one step")
	}
	for i := range d.Steps {
		if err := d.validateStep(i+1, &d.Steps[i]); err != nil {
			return fmt.Errorf("step %d: %w", i+1, err)
		}
	}
	return nil
}

func (d *Definition) validateStep(n int, s *Step) error {
	count := 0
	for _, present := range []bool{s.Run != nil, s.Copy != nil, s.Write != nil, s.Remove != nil, s.Mkdir != nil} {
		if present {
			count++
		}
	}
	if count != 1 {
		return errors.New("exactly one of run, copy, write, remove, or mkdir is required")
	}
	switch {
	case s.Run != nil:
		if len(s.Run.Argv) == 0 {
			return errors.New("run argv must not be empty")
		}
		for _, arg := range s.Run.Argv {
			if strings.IndexByte(arg, 0) >= 0 {
				return errors.New("run argv contains a NUL byte")
			}
		}
		if s.Run.Cwd != "" {
			if err := validateRelative(s.Run.Cwd, true); err != nil {
				return fmt.Errorf("invalid cwd: %w", err)
			}
		}
	case s.Copy != nil:
		if err := validateRelative(s.Copy.From, false); err != nil {
			return fmt.Errorf("invalid copy source: %w", err)
		}
		if err := validateRelative(s.Copy.To, false); err != nil {
			return fmt.Errorf("invalid copy destination: %w", err)
		}
		if _, err := securePath(d.fixtureDir, s.Copy.From, true); err != nil {
			return fmt.Errorf("copy source: %w", err)
		}
	case s.Write != nil:
		if err := validateRelative(s.Write.Path, false); err != nil {
			return fmt.Errorf("invalid write path: %w", err)
		}
	case s.Remove != nil:
		if err := validateRelative(s.Remove.Path, false); err != nil {
			return fmt.Errorf("invalid remove path: %w", err)
		}
	case s.Mkdir != nil:
		if err := validateRelative(s.Mkdir.Path, false); err != nil {
			return fmt.Errorf("invalid mkdir path: %w", err)
		}
	}
	_ = n
	return nil
}

func benchmarkRoot(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(abs)
	for {
		if filepath.Base(dir) == "workloads" {
			return filepath.Dir(dir), nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", errors.New("workload must be located under a benchmarks/workloads directory")
}

func validateRelative(path string, allowDot bool) error {
	if path == "" {
		return errors.New("path is required")
	}
	if filepath.IsAbs(path) || filepath.VolumeName(path) != "" || !filepath.IsLocal(path) {
		return fmt.Errorf("path must be relative and remain within its root: %q", path)
	}
	clean := filepath.Clean(path)
	if !allowDot && clean == "." {
		return errors.New("path must identify a child of the workspace")
	}
	return nil
}

func hasDotDot(path string) bool {
	for _, part := range strings.FieldsFunc(path, func(r rune) bool { return r == '/' || r == '\\' }) {
		if part == ".." {
			return true
		}
	}
	return false
}

func within(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && !filepath.IsAbs(rel)
}

// securePath rejects symlinks in every existing path component. requireFile
// requires the final path to exist; callers still decide whether it is a file
// or directory.
func securePath(root, rel string, requireFile bool) (string, error) {
	if err := validateRelative(rel, false); err != nil {
		return "", err
	}
	target := filepath.Join(root, filepath.Clean(rel))
	if !within(root, target) {
		return "", errors.New("path escapes its root")
	}
	cur := root
	parts := strings.Split(filepath.Clean(rel), string(filepath.Separator))
	for i, part := range parts {
		cur = filepath.Join(cur, part)
		info, err := os.Lstat(cur)
		if err != nil {
			if os.IsNotExist(err) && !requireFile {
				return target, nil
			}
			return "", err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("symlink path component is not allowed: %s", strings.Join(parts[:i+1], string(filepath.Separator)))
		}
	}
	return target, nil
}

// Options controls workspace retention and result persistence.
type Options struct {
	StoreRoot     string
	TempRoot      string
	KeepWorkspace bool
}

type Result struct {
	Name              string
	SessionID         string
	TotalSteps        int
	MutationCount     int
	Aggregate         *bench.SessionMeasurement
	WallDuration      time.Duration
	Steps             []StepResult
	RetainedWorkspace string
}

type StepResult struct {
	Number      int
	Type        string
	Measurement *bench.Measurement
}

// StepError identifies infrastructure or assertion failure separately from a
// target command's ordinary non-zero exit.
type StepError struct {
	Workload string
	Step     int
	Type     string
	Err      error
}

func (e *StepError) Error() string {
	return fmt.Sprintf("workload %q failed at step %d (%s): %v", e.Workload, e.Step, e.Type, e.Err)
}
func (e *StepError) Unwrap() error { return e.Err }

// Run creates a fresh workspace and B2 benchmark session, then executes each
// run step exactly once through bench.MeasureSessionCommandAt.
func Run(ctx context.Context, d *Definition, opts Options) (result *Result, err error) {
	started := time.Now()
	result = &Result{Name: d.Name, TotalSteps: len(d.Steps)}
	tempRoot := opts.TempRoot
	if tempRoot == "" {
		tempRoot = os.TempDir()
	}
	workspace, err := os.MkdirTemp(tempRoot, "acap-bench-")
	if err != nil {
		return result, fmt.Errorf("create temporary workspace: %w", err)
	}
	defer func() {
		result.WallDuration = time.Since(started)
		if opts.KeepWorkspace {
			result.RetainedWorkspace = workspace
			return
		}
		if removeErr := os.RemoveAll(workspace); err == nil && removeErr != nil {
			err = fmt.Errorf("clean temporary workspace: %w", removeErr)
		}
	}()
	if err = copyTree(d.fixtureDir, workspace); err != nil {
		return result, fmt.Errorf("copy fixture: %w", err)
	}
	if d.Git.Init {
		if err = initGit(ctx, workspace); err != nil {
			return result, fmt.Errorf("initialize Git fixture: %w", err)
		}
	}
	if opts.StoreRoot == "" {
		return result, errors.New("persistent store root is required")
	}
	result.SessionID, err = bench.StartSession(opts.StoreRoot)
	if err != nil {
		return result, fmt.Errorf("start benchmark session: %w", err)
	}
	for i := range d.Steps {
		step := &d.Steps[i]
		typeName := stepType(step)
		sr := StepResult{Number: i + 1, Type: typeName}
		result.Steps = append(result.Steps, sr)
		if step.Run != nil {
			cwd := workspace
			if step.Run.Cwd != "" && filepath.Clean(step.Run.Cwd) != "." {
				cwd, err = securePath(workspace, step.Run.Cwd, true)
				if err == nil {
					var info os.FileInfo
					info, err = os.Stat(cwd)
					if err == nil && !info.IsDir() {
						err = errors.New("cwd is not a directory")
					}
				}
			}
			if err == nil {
				var m *bench.Measurement
				m, err = bench.MeasureSessionCommandAt(ctx, step.Run.Argv, cwd, opts.StoreRoot, result.SessionID)
				result.Steps[len(result.Steps)-1].Measurement = m
				if err == nil && step.Run.Expect != nil && step.Run.Expect.Exit != nil && m.ExitCode != *step.Run.Expect.Exit {
					err = fmt.Errorf("exit status %d, expected %d", m.ExitCode, *step.Run.Expect.Exit)
				}
			}
		} else {
			result.MutationCount++
			err = mutate(step, d.fixtureDir, workspace)
		}
		if err != nil {
			result.Aggregate, _ = bench.LoadSession(opts.StoreRoot, result.SessionID)
			return result, &StepError{Workload: d.Name, Step: i + 1, Type: typeName, Err: err}
		}
	}
	result.Aggregate, err = bench.LoadSession(opts.StoreRoot, result.SessionID)
	if err != nil {
		return result, fmt.Errorf("load benchmark session: %w", err)
	}
	return result, nil
}

func stepType(s *Step) string {
	switch {
	case s.Run != nil:
		return "run"
	case s.Copy != nil:
		return "copy"
	case s.Write != nil:
		return "write"
	case s.Remove != nil:
		return "remove"
	default:
		return "mkdir"
	}
}

func mutate(s *Step, fixtureRoot, workspace string) error {
	switch {
	case s.Copy != nil:
		src, err := securePath(fixtureRoot, s.Copy.From, true)
		if err != nil {
			return fmt.Errorf("source fixture file: %w", err)
		}
		dst, err := securePath(workspace, s.Copy.To, false)
		if err != nil {
			return err
		}
		return copyTree(src, dst)
	case s.Write != nil:
		dst, err := securePath(workspace, s.Write.Path, false)
		if err != nil {
			return err
		}
		return os.WriteFile(dst, []byte(s.Write.Content), 0o644)
	case s.Remove != nil:
		dst, err := securePath(workspace, s.Remove.Path, false)
		if err != nil {
			return err
		}
		return os.RemoveAll(dst)
	case s.Mkdir != nil:
		dst, err := securePath(workspace, s.Mkdir.Path, false)
		if err != nil {
			return err
		}
		return os.MkdirAll(dst, 0o755)
	}
	return errors.New("unknown mutation")
}

func copyTree(src, dst string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("symlink is not allowed: %s", src)
	}
	if !info.IsDir() {
		data, err := os.ReadFile(src)
		if err != nil {
			return err
		}
		return os.WriteFile(dst, data, info.Mode().Perm())
	}
	if err := os.MkdirAll(dst, info.Mode().Perm()); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if err := copyTree(filepath.Join(src, entry.Name()), filepath.Join(dst, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

func initGit(ctx context.Context, workspace string) error {
	emptyConfig := filepath.Join(workspace, ".agentcap-empty-gitconfig")
	env := append(os.Environ(),
		"GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL="+emptyConfig,
		"GIT_TERMINAL_PROMPT=0", "GIT_AUTHOR_DATE=2000-01-01T00:00:00Z",
		"GIT_COMMITTER_DATE=2000-01-01T00:00:00Z")
	commands := [][]string{
		{"git", "init", "--quiet"},
		{"git", "config", "user.name", "AgentCap Benchmark"},
		{"git", "config", "user.email", "benchmark@agentcap.invalid"},
		{"git", "config", "commit.gpgSign", "false"},
		{"git", "config", "core.hooksPath", ".git/disabled-hooks"},
		{"git", "add", "--all"},
		{"git", "commit", "--quiet", "-m", "benchmark baseline"},
	}
	for _, argv := range commands {
		cmd := osexec.CommandContext(ctx, argv[0], argv[1:]...)
		cmd.Dir = workspace
		cmd.Env = env
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("%s: %w: %s", strings.Join(argv, " "), err, strings.TrimSpace(string(output)))
		}
	}
	return nil
}
