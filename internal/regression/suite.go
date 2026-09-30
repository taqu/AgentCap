// Package regression applies an explicit regression policy to B8 benchmark
// comparisons. It is the only layer that decides whether a change is a
// regression; measurement and comparison stay descriptive. The evaluator is
// pure: it never runs workloads, agents, or AgentCap processing, and it never
// rewrites baselines or thresholds.
package regression

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"

	"github.com/taqu/agentcap/internal/workload"
	"gopkg.in/yaml.v3"
)

// SuiteSchemaVersion is the version of the regression suite file format.
const SuiteSchemaVersion = 1

// Rule is the upper bound on an increase (candidate - baseline) of one
// metric. At least one limit is required. When both are set, the absolute
// limit acts as a noise floor: an increase is a regression only when it
// exceeds every limit that can be evaluated.
type Rule struct {
	MaxRelativeIncrease *float64 `yaml:"max_relative_increase" json:"max_relative_increase,omitempty"`
	MaxAbsoluteIncrease *int64   `yaml:"max_absolute_increase" json:"max_absolute_increase,omitempty"`
}

// MetricRule is a Rule bound to a B8 comparison metric name.
type MetricRule struct {
	Metric string
	Rule
}

// Policy is the ordered set of per-metric rules. Every rule is required.
type Policy struct {
	Rules []MetricRule
}

// SuiteWorkload is one deterministic workload of the regression suite.
type SuiteWorkload struct {
	Path       string // resolved workload file path
	Name       string // stable workload identity used to match results
	Definition *workload.Definition
}

// Suite is a regression suite file: the deterministic workloads that CI runs
// and the policy applied to every one of them.
type Suite struct {
	Path      string
	Workloads []SuiteWorkload
	Policy    Policy
}

// Metrics lists, in report order, the B8 comparison metrics a policy may
// constrain. raw_bytes is deliberately absent: raw output is a reference
// measurement of the fixture and tools, not of AgentCap.
var Metrics = []string{
	// single-result metrics (deterministic workloads and single agent trials)
	"stateless_bytes",
	"stateful_bytes",
	"total_visible_bytes",
	"show_bytes",
	"raw_retrieval_bytes",
	"command_count",
	"show_count",
	"raw_retrieval_count",
	"processing_ns",
	"wall_time_ns",
	// repeated-result aggregate metrics
	"median_total_visible_bytes",
	"median_command_count",
	"median_processing_ns",
	"median_wall_time_ns",
	"total_show_count",
	"total_raw_retrieval_count",
}

type suiteFile struct {
	Version   int             `yaml:"version"`
	Workloads []string        `yaml:"workloads"`
	Metrics   map[string]Rule `yaml:"metrics"`
}

// LoadSuite parses and validates a regression suite. Workload paths are
// relative to the suite file. Every workload must be a deterministic
// (non-agent) workload.
func LoadSuite(path string) (*Suite, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read regression suite: %w", err)
	}
	var file suiteFile
	if err := decodeStrict(data, &file); err != nil {
		return nil, fmt.Errorf("parse regression suite %s: %w", path, err)
	}
	if file.Version != SuiteSchemaVersion {
		return nil, fmt.Errorf("%s: unsupported regression suite version %d (supported: %d)", path, file.Version, SuiteSchemaVersion)
	}
	policy, err := NewPolicy(file.Metrics)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	suite := &Suite{Path: path, Policy: *policy}
	if len(file.Workloads) == 0 {
		return nil, fmt.Errorf("%s: at least one workload is required", path)
	}
	names := map[string]bool{}
	for _, rel := range file.Workloads {
		workloadPath := rel
		if !filepath.IsAbs(workloadPath) {
			workloadPath = filepath.Join(filepath.Dir(path), rel)
		}
		definition, err := workload.Load(workloadPath)
		if err != nil {
			return nil, fmt.Errorf("%s: workload %s: %w", path, rel, err)
		}
		if definition.Task != "" {
			return nil, fmt.Errorf("%s: workload %s is a coding-agent workload; regression suites accept only deterministic workloads", path, rel)
		}
		if names[definition.Name] {
			return nil, fmt.Errorf("%s: workload %q is listed more than once", path, definition.Name)
		}
		names[definition.Name] = true
		suite.Workloads = append(suite.Workloads, SuiteWorkload{Path: workloadPath, Name: definition.Name, Definition: definition})
	}
	return suite, nil
}

// NewPolicy validates metric rules and orders them canonically. Unknown
// metrics are rejected so a typo cannot silently remove a check.
func NewPolicy(rules map[string]Rule) (*Policy, error) {
	if len(rules) == 0 {
		return nil, errors.New("at least one metric rule is required")
	}
	known := map[string]bool{}
	for _, name := range Metrics {
		known[name] = true
	}
	for name, rule := range rules {
		if !known[name] {
			if name == "raw_bytes" {
				return nil, errors.New("raw_bytes is a reference measurement and cannot be a regression metric")
			}
			return nil, fmt.Errorf("unknown metric %q", name)
		}
		if rule.MaxRelativeIncrease == nil && rule.MaxAbsoluteIncrease == nil {
			return nil, fmt.Errorf("metric %s: max_relative_increase or max_absolute_increase is required", name)
		}
		if r := rule.MaxRelativeIncrease; r != nil && (math.IsNaN(*r) || math.IsInf(*r, 0) || *r < 0) {
			return nil, fmt.Errorf("metric %s: max_relative_increase must be a finite fraction >= 0 (0.10 means 10%%)", name)
		}
		if a := rule.MaxAbsoluteIncrease; a != nil && *a < 0 {
			return nil, fmt.Errorf("metric %s: max_absolute_increase must be >= 0", name)
		}
	}
	policy := &Policy{}
	for _, name := range Metrics {
		if rule, ok := rules[name]; ok {
			policy.Rules = append(policy.Rules, MetricRule{Metric: name, Rule: rule})
		}
	}
	return policy, nil
}

func decodeStrict(data []byte, out any) error {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(out); err != nil {
		return err
	}
	var extra any
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("multiple YAML documents are not supported")
		}
		return err
	}
	return nil
}
