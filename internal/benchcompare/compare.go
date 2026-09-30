// Package benchcompare compares two previously recorded benchmark results.
//
// It is a read-only analysis layer: it consumes decoded result values and
// never runs workloads, coding agents, verifiers, or AgentCap processing. It
// reports signed differences and deliberately makes no judgment about which
// side is better; regression policy belongs to a separate layer.
package benchcompare

import (
	"errors"
	"fmt"

	"github.com/taqu/agentcap/internal/workload"
)

// ComparisonSchemaVersion identifies the JSON comparison contract. It is
// independent of workload.BenchmarkResultSchemaVersion.
//
// Version 2 added raw_bytes, stateless_bytes, and stateful_bytes to single
// result comparisons.
const ComparisonSchemaVersion = 2

// modeDisabled is the coding-agent mode without any AgentCap hook, in which
// AgentCap processing and show/raw recovery cannot occur.
const modeDisabled = "disabled"

// Availability states why a side does or does not carry a value.
type Availability string

const (
	// Measured values were recorded by the benchmark, including measured zero.
	Measured Availability = "measured"
	// Unavailable values were not recorded by this result kind or schema.
	Unavailable Availability = "unavailable"
	// NotApplicable values cannot exist in this configuration, such as
	// AgentCap processing when AgentCap is disabled.
	NotApplicable Availability = "not_applicable"
)

// Unit describes how a numeric metric is measured.
type Unit string

const (
	UnitBytes       Unit = "bytes"
	UnitCount       Unit = "count"
	UnitNanoseconds Unit = "ns"
)

// Side is the categorical metadata of one compared result. It is displayed
// side by side and never subtracted.
type Side struct {
	SchemaVersion       int    `json:"schema_version"`
	Workload            string `json:"workload"`
	Agent               string `json:"agent,omitempty"`
	Mode                string `json:"mode,omitempty"`
	TrialCount          int    `json:"trial_count"`
	RequestedTrialCount int    `json:"requested_trial_count,omitempty"`
	RunStatus           string `json:"run_status,omitempty"`
	ExecutionStatus     string `json:"execution_status,omitempty"`
}

// NumericMetric compares one measured quantity. Delta is candidate - baseline
// and RelativeDelta is (candidate - baseline) / baseline; each is nil when it
// is not well-defined.
type NumericMetric struct {
	Name            string       `json:"name"`
	Unit            Unit         `json:"unit"`
	Baseline        *int64       `json:"baseline"`
	Candidate       *int64       `json:"candidate"`
	BaselineStatus  Availability `json:"baseline_status"`
	CandidateStatus Availability `json:"candidate_status"`
	Delta           *int64       `json:"delta"`
	RelativeDelta   *float64     `json:"relative_delta"`
}

// Fraction is a count of trials out of the trials observed on one side.
type Fraction struct {
	Count int `json:"count"`
	Total int `json:"total"`
}

// OutcomeMetric compares structured trial outcomes side by side. Counts with
// their own denominators are never reduced to a single numeric delta.
type OutcomeMetric struct {
	Name            string       `json:"name"`
	Baseline        *Fraction    `json:"baseline"`
	Candidate       *Fraction    `json:"candidate"`
	BaselineStatus  Availability `json:"baseline_status"`
	CandidateStatus Availability `json:"candidate_status"`
}

// Comparison is the canonical result of comparing a baseline with a
// candidate. Human and JSON output are both rendered from this model.
type Comparison struct {
	SchemaVersion int                 `json:"schema_version"`
	Kind          workload.ResultKind `json:"kind"`
	Baseline      Side                `json:"baseline"`
	Candidate     Side                `json:"candidate"`
	Outcomes      []OutcomeMetric     `json:"outcomes"`
	Metrics       []NumericMetric     `json:"metrics"`
}

// Metric returns the numeric metric with name, or nil.
func (c *Comparison) Metric(name string) *NumericMetric {
	for i := range c.Metrics {
		if c.Metrics[i].Name == name {
			return &c.Metrics[i]
		}
	}
	return nil
}

// Outcome returns the outcome metric with name, or nil.
func (c *Comparison) Outcome(name string) *OutcomeMetric {
	for i := range c.Outcomes {
		if c.Outcomes[i].Name == name {
			return &c.Outcomes[i]
		}
	}
	return nil
}

// IncompatibleError reports results that must not be compared as one
// configuration change.
type IncompatibleError struct {
	Reason string
}

func (e *IncompatibleError) Error() string { return "incompatible results: " + e.Reason }

// Compare builds the canonical comparison of baseline and candidate. It is a
// pure function of its inputs and does not modify them.
func Compare(baseline, candidate *workload.StoredResult) (*Comparison, error) {
	if baseline == nil || candidate == nil {
		return nil, errors.New("baseline and candidate results are required")
	}
	if err := checkCompatible(baseline, candidate); err != nil {
		return nil, err
	}
	c := &Comparison{SchemaVersion: ComparisonSchemaVersion, Kind: baseline.Kind}
	if baseline.Kind == workload.ResultKindRepeated {
		compareRepeated(c, baseline.Repeated, candidate.Repeated)
	} else {
		compareSingle(c, baseline.Single, candidate.Single)
	}
	return c, nil
}

func checkCompatible(baseline, candidate *workload.StoredResult) error {
	if baseline.Kind != candidate.Kind {
		return &IncompatibleError{Reason: fmt.Sprintf("result kinds differ (baseline %s, candidate %s); a single observation is not comparable to a repeated-trial median", baseline.Kind, candidate.Kind)}
	}
	if (baseline.Kind == workload.ResultKindSingle && (baseline.Single == nil || candidate.Single == nil)) ||
		(baseline.Kind == workload.ResultKindRepeated && (baseline.Repeated == nil || candidate.Repeated == nil)) {
		return fmt.Errorf("malformed %s result", baseline.Kind)
	}
	bw, ba := identity(baseline)
	cw, ca := identity(candidate)
	if bw != cw {
		return &IncompatibleError{Reason: fmt.Sprintf("workloads differ (baseline %q, candidate %q)", bw, cw)}
	}
	if ba != ca {
		return &IncompatibleError{Reason: fmt.Sprintf("agents differ (baseline %q, candidate %q)", ba, ca)}
	}
	return nil
}

func identity(r *workload.StoredResult) (workloadName, agent string) {
	if r.Repeated != nil {
		return r.Repeated.Workload, r.Repeated.Agent
	}
	return r.Single.Workload, r.Single.Agent
}

func compareRepeated(c *Comparison, b, k *workload.RepeatedBenchmarkResult) {
	c.Baseline = repeatedSide(b)
	c.Candidate = repeatedSide(k)
	bNA, kNA := b.Mode == modeDisabled, k.Mode == modeDisabled

	c.Outcomes = []OutcomeMetric{
		outcome("task_success", fraction(b.SuccessCount, b.TrialCount), fraction(k.SuccessCount, k.TrialCount)),
		outcome("task_failure", fraction(b.TaskFailureCount, b.TrialCount), fraction(k.TaskFailureCount, k.TrialCount)),
		outcome("timeout", fraction(b.TimeoutCount, b.TrialCount), fraction(k.TimeoutCount, k.TrialCount)),
		outcome("agent_error", fraction(b.AgentErrorCount, b.TrialCount), fraction(k.AgentErrorCount, k.TrialCount)),
		outcome("canceled", fraction(b.CanceledCount, b.TrialCount), fraction(k.CanceledCount, k.TrialCount)),
		outcome("trials_with_show",
			recoveryFraction(bNA, b.Aggregate.TrialsWithShow, b.TrialCount),
			recoveryFraction(kNA, k.Aggregate.TrialsWithShow, k.TrialCount)),
		outcome("trials_with_raw_retrieval",
			recoveryFraction(bNA, b.Aggregate.TrialsWithRawRetrieval, b.TrialCount),
			recoveryFraction(kNA, k.Aggregate.TrialsWithRawRetrieval, k.TrialCount)),
	}

	ba, ka := b.Aggregate, k.Aggregate
	sameTrials := b.TrialCount == k.TrialCount
	c.Metrics = []NumericMetric{
		numeric("median_total_visible_bytes", UnitBytes, optional(ba.MedianTotalVisibleBytes), optional(ka.MedianTotalVisibleBytes), true),
		numeric("median_command_count", UnitCount, optional(ba.MedianCommandCount), optional(ka.MedianCommandCount), true),
		numeric("median_wall_time_ns", UnitNanoseconds, optional(ba.MedianWallTimeNS), optional(ka.MedianWallTimeNS), true),
		numeric("median_processing_ns", UnitNanoseconds,
			agentCap(bNA, optional(ba.MedianProcessingNS)), agentCap(kNA, optional(ka.MedianProcessingNS)), true),
		// Totals are sums over each side's own trials, so a difference is only
		// meaningful when both sides observed the same number of trials.
		numeric("total_show_count", UnitCount,
			agentCap(bNA, count(b.TrialCount, int64(ba.TotalShowCount))), agentCap(kNA, count(k.TrialCount, int64(ka.TotalShowCount))), sameTrials),
		numeric("total_raw_retrieval_count", UnitCount,
			agentCap(bNA, count(b.TrialCount, int64(ba.TotalRawRetrievalCount))), agentCap(kNA, count(k.TrialCount, int64(ka.TotalRawRetrievalCount))), sameTrials),
	}
}

func compareSingle(c *Comparison, b, k *workload.BenchmarkResult) {
	c.Baseline = singleSide(b)
	c.Candidate = singleSide(k)
	bNA, kNA := b.Mode == modeDisabled, k.Mode == modeDisabled

	c.Outcomes = []OutcomeMetric{
		outcome("task_success", taskSuccess(b), taskSuccess(k)),
	}
	c.Metrics = []NumericMetric{
		numeric("raw_bytes", UnitBytes, measured(b.RawBytes), measured(k.RawBytes), true),
		numeric("stateless_bytes", UnitBytes, measured(b.StatelessBytes), measured(k.StatelessBytes), true),
		numeric("stateful_bytes", UnitBytes, measured(b.StatefulBytes), measured(k.StatefulBytes), true),
		numeric("total_visible_bytes", UnitBytes, measured(b.TotalVisibleBytes), measured(k.TotalVisibleBytes), true),
		numeric("command_count", UnitCount, measured(int64(b.Commands)), measured(int64(k.Commands)), true),
		numeric("wall_time_ns", UnitNanoseconds, wallTime(b), wallTime(k), true),
		numeric("processing_ns", UnitNanoseconds, agentCap(bNA, measured(b.ProcessingNS)), agentCap(kNA, measured(k.ProcessingNS)), true),
		numeric("show_count", UnitCount, agentCap(bNA, measured(int64(b.ShowCount))), agentCap(kNA, measured(int64(k.ShowCount))), true),
		numeric("raw_retrieval_count", UnitCount, agentCap(bNA, measured(int64(b.RawRetrievalCount))), agentCap(kNA, measured(int64(k.RawRetrievalCount))), true),
		numeric("show_bytes", UnitBytes, agentCap(bNA, measured(b.ShowBytes)), agentCap(kNA, measured(k.ShowBytes)), true),
		numeric("raw_retrieval_bytes", UnitBytes, agentCap(bNA, measured(b.RawRetrievalBytes)), agentCap(kNA, measured(k.RawRetrievalBytes)), true),
	}
}

func repeatedSide(r *workload.RepeatedBenchmarkResult) Side {
	return Side{
		SchemaVersion: r.SchemaVersion, Workload: r.Workload, Agent: r.Agent, Mode: r.Mode,
		TrialCount: r.TrialCount, RequestedTrialCount: r.RequestedTrialCount, RunStatus: r.RunStatus,
	}
}

func singleSide(r *workload.BenchmarkResult) Side {
	return Side{
		SchemaVersion: r.SchemaVersion, Workload: r.Workload, Agent: r.Agent, Mode: r.Mode,
		TrialCount: 1, ExecutionStatus: r.ExecutionStatus,
	}
}

// value is one side of a metric before it is placed in the comparison.
type value struct {
	v      *int64
	status Availability
}

func measured(v int64) value { return value{v: &v, status: Measured} }

// optional maps a nil B7 median (no participating trials) to Unavailable.
func optional(v *int64) value {
	if v == nil {
		return value{status: Unavailable}
	}
	return measured(*v)
}

// count is Unavailable when no trials were observed.
func count(trials int, v int64) value {
	if trials == 0 {
		return value{status: Unavailable}
	}
	return measured(v)
}

// agentCap marks AgentCap-produced measurements as NotApplicable when the
// configuration has no AgentCap hook, rather than reporting a zero.
func agentCap(disabled bool, v value) value {
	if disabled {
		return value{status: NotApplicable}
	}
	return v
}

// wallTime is only part of the coding-agent result schema. Deterministic
// workload results do not record it.
func wallTime(r *workload.BenchmarkResult) value {
	if r.Agent == "" {
		return value{status: Unavailable}
	}
	return measured(r.WallTimeNS)
}

// numeric applies the single, central delta definition: candidate - baseline.
func numeric(name string, unit Unit, b, k value, withDelta bool) NumericMetric {
	m := NumericMetric{
		Name: name, Unit: unit,
		Baseline: b.v, Candidate: k.v,
		BaselineStatus: b.status, CandidateStatus: k.status,
	}
	if !withDelta || b.status != Measured || k.status != Measured {
		return m
	}
	delta := *k.v - *b.v
	m.Delta = &delta
	if *b.v != 0 {
		relative := float64(delta) / float64(*b.v)
		m.RelativeDelta = &relative
	}
	return m
}

type fractionValue struct {
	f      *Fraction
	status Availability
}

func fraction(n, total int) fractionValue {
	if total == 0 {
		return fractionValue{status: Unavailable}
	}
	return fractionValue{f: &Fraction{Count: n, Total: total}, status: Measured}
}

func recoveryFraction(disabled bool, n, total int) fractionValue {
	if disabled {
		return fractionValue{status: NotApplicable}
	}
	return fraction(n, total)
}

func taskSuccess(r *workload.BenchmarkResult) fractionValue {
	if r.TaskSuccess == nil {
		return fractionValue{status: Unavailable}
	}
	n := 0
	if *r.TaskSuccess {
		n = 1
	}
	return fraction(n, 1)
}

func outcome(name string, b, k fractionValue) OutcomeMetric {
	return OutcomeMetric{
		Name: name, Baseline: b.f, Candidate: k.f,
		BaselineStatus: b.status, CandidateStatus: k.status,
	}
}
