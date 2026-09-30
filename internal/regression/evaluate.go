package regression

import (
	"github.com/taqu/agentcap/internal/benchcompare"
	"github.com/taqu/agentcap/internal/workload"
)

// EvaluationSchemaVersion identifies the JSON regression-evaluation contract.
// It is independent of the benchmark result and comparison schemas.
const EvaluationSchemaVersion = 1

// Status is the outcome of one policy check.
type Status string

const (
	Pass         Status = "pass"
	Regression   Status = "regression"
	NotEvaluable Status = "not_evaluable"
)

// WorkloadStatus states whether a suite workload could be compared at all.
type WorkloadStatus string

const (
	WorkloadEvaluated    WorkloadStatus = "evaluated"
	WorkloadMissing      WorkloadStatus = "missing"
	WorkloadIncompatible WorkloadStatus = "incompatible"
)

// Overall evaluation states. StatusError is produced only when evaluation
// could not be completed, for example when a result cannot be loaded.
const (
	StatusPassed = "passed"
	StatusFailed = "failed"
	StatusError  = "error"
)

// Check is one rule applied to one B8 comparison metric. The embedded metric
// carries the comparison's own values and candidate - baseline delta.
type Check struct {
	benchcompare.NumericMetric
	Rule
	Status Status `json:"status"`
	Reason string `json:"reason,omitempty"`
}

// WorkloadEvaluation is every check for one workload. Workloads are never
// merged, so one workload's improvement cannot hide another's regression.
type WorkloadEvaluation struct {
	Workload string         `json:"workload"`
	Status   WorkloadStatus `json:"status"`
	Reason   string         `json:"reason,omitempty"`
	// RawBytes is reported for reference only and is never gated.
	RawBytes *benchcompare.NumericMetric `json:"raw_bytes,omitempty"`
	Checks   []Check                     `json:"checks"`
}

// Summary counts check and workload outcomes.
type Summary struct {
	Workloads       int `json:"workloads"`
	FailedWorkloads int `json:"failed_workloads"`
	Checks          int `json:"checks"`
	Passed          int `json:"passed"`
	Regressions     int `json:"regressions"`
	NotEvaluable    int `json:"not_evaluable"`
}

// Evaluation is the canonical regression result. Human output, JSON output,
// and the CLI exit status are all derived from it.
type Evaluation struct {
	SchemaVersion int                  `json:"schema_version"`
	Status        string               `json:"status"`
	Passed        bool                 `json:"passed"`
	Error         string               `json:"error,omitempty"`
	Summary       Summary              `json:"summary"`
	Workloads     []WorkloadEvaluation `json:"workloads"`
}

// ErrorEvaluation represents an evaluation that could not be completed. It is
// never a passing result and reports no checks.
func ErrorEvaluation(err error) *Evaluation {
	return &Evaluation{SchemaVersion: EvaluationSchemaVersion, Status: StatusError, Error: err.Error(), Workloads: []WorkloadEvaluation{}}
}

// EvaluateComparison applies policy to one B8 comparison. It reuses the
// comparison's delta semantics and performs no delta arithmetic of its own.
func EvaluateComparison(c *benchcompare.Comparison, policy *Policy) []Check {
	checks := make([]Check, 0, len(policy.Rules))
	for _, rule := range policy.Rules {
		m := c.Metric(rule.Metric)
		if m == nil {
			checks = append(checks, Check{
				NumericMetric: benchcompare.NumericMetric{Name: rule.Metric},
				Rule:          rule.Rule, Status: NotEvaluable,
				Reason: "metric is not reported for " + string(c.Kind) + " results",
			})
			continue
		}
		status, reason := evaluate(m, rule.Rule)
		checks = append(checks, Check{NumericMetric: *m, Rule: rule.Rule, Status: status, Reason: reason})
	}
	return checks
}

// evaluate defines the policy semantics. An increase exactly equal to a limit
// passes. With both limits configured, an increase within either limit
// passes; a zero baseline leaves only the absolute limit evaluable.
func evaluate(m *benchcompare.NumericMetric, rule Rule) (Status, string) {
	if m.BaselineStatus != benchcompare.Measured {
		return NotEvaluable, "baseline value is " + string(m.BaselineStatus)
	}
	if m.CandidateStatus != benchcompare.Measured {
		return NotEvaluable, "candidate value is " + string(m.CandidateStatus)
	}
	if m.Delta == nil {
		return NotEvaluable, "comparison defines no delta for this metric"
	}
	if *m.Delta <= 0 {
		return Pass, ""
	}
	limitEvaluated := false
	if rule.MaxAbsoluteIncrease != nil {
		limitEvaluated = true
		if *m.Delta <= *rule.MaxAbsoluteIncrease {
			return Pass, ""
		}
	}
	if rule.MaxRelativeIncrease != nil && m.RelativeDelta != nil {
		limitEvaluated = true
		if *m.RelativeDelta <= *rule.MaxRelativeIncrease {
			return Pass, ""
		}
	}
	if !limitEvaluated {
		return NotEvaluable, "relative change is undefined for a zero baseline and no absolute limit is configured"
	}
	return Regression, ""
}

// EvaluateSuite compares and evaluates every suite workload. baseline and
// candidate are keyed by workload name; results for workloads outside the
// suite are ignored. The evaluation fails when any check regresses or cannot
// be evaluated, when a suite workload is missing or incompatible, or when no
// check was evaluated at all.
func EvaluateSuite(workloads []string, policy *Policy, baseline, candidate map[string]*workload.StoredResult) *Evaluation {
	e := &Evaluation{SchemaVersion: EvaluationSchemaVersion, Workloads: make([]WorkloadEvaluation, 0, len(workloads))}
	for _, name := range workloads {
		we := WorkloadEvaluation{Workload: name, Checks: []Check{}}
		b, k := baseline[name], candidate[name]
		switch {
		case b == nil && k == nil:
			we.Status, we.Reason = WorkloadMissing, "no baseline or candidate result"
		case b == nil:
			we.Status, we.Reason = WorkloadMissing, "no baseline result"
		case k == nil:
			we.Status, we.Reason = WorkloadMissing, "no candidate result"
		default:
			c, err := benchcompare.Compare(b, k)
			if err != nil {
				we.Status, we.Reason = WorkloadIncompatible, err.Error()
				break
			}
			we.Status = WorkloadEvaluated
			if raw := c.Metric("raw_bytes"); raw != nil {
				copied := *raw
				we.RawBytes = &copied
			}
			we.Checks = EvaluateComparison(c, policy)
		}
		e.add(we)
	}
	e.Passed = e.Summary.FailedWorkloads == 0 && e.Summary.Passed > 0
	e.Status = StatusFailed
	if e.Passed {
		e.Status = StatusPassed
	}
	return e
}

func (e *Evaluation) add(we WorkloadEvaluation) {
	e.Workloads = append(e.Workloads, we)
	e.Summary.Workloads++
	failed := we.Status != WorkloadEvaluated
	for _, check := range we.Checks {
		e.Summary.Checks++
		switch check.Status {
		case Pass:
			e.Summary.Passed++
		case Regression:
			e.Summary.Regressions++
			failed = true
		default:
			e.Summary.NotEvaluable++
			failed = true
		}
	}
	if failed {
		e.Summary.FailedWorkloads++
	}
}
