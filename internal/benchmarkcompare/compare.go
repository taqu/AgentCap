package benchmarkcompare

import (
	"errors"
	"fmt"
)

type normalized struct {
	side Side

	visible, commands, wall, processing *int64
	trialsWithShow, trialsWithRaw       *int64
	totalShow, totalRaw                 *int64
}

// Compare creates a deterministic comparison without modifying either source.
func Compare(baseline, candidate *Result) (*Comparison, error) {
	if baseline == nil || candidate == nil {
		return nil, errors.New("baseline and candidate results are required")
	}
	if baseline.Kind != candidate.Kind {
		return nil, fmt.Errorf("incompatible result kinds: baseline %s, candidate %s", baseline.Kind, candidate.Kind)
	}
	b, err := normalize(baseline)
	if err != nil {
		return nil, fmt.Errorf("baseline: %w", err)
	}
	c, err := normalize(candidate)
	if err != nil {
		return nil, fmt.Errorf("candidate: %w", err)
	}
	if b.side.Workload != c.side.Workload {
		return nil, fmt.Errorf("incompatible workloads: baseline %q, candidate %q", b.side.Workload, c.side.Workload)
	}
	if b.side.Agent != c.side.Agent {
		return nil, fmt.Errorf("incompatible agents: baseline %q, candidate %q", b.side.Agent, c.side.Agent)
	}
	statistic := "observation"
	if baseline.Kind == KindRepeated {
		statistic = "median"
	}
	return &Comparison{
		SchemaVersion: SchemaVersion,
		Statistic:     statistic,
		Baseline:      b.side,
		Candidate:     c.side,
		Metrics: Metrics{
			TotalVisibleBytes:      numeric(b.visible, c.visible),
			CommandCount:           numeric(b.commands, c.commands),
			WallTimeNS:             numeric(b.wall, c.wall),
			ProcessingNS:           numeric(b.processing, c.processing),
			TrialsWithShow:         numeric(b.trialsWithShow, c.trialsWithShow),
			TrialsWithRawRetrieval: numeric(b.trialsWithRaw, c.trialsWithRaw),
			TotalShowCount:         numeric(b.totalShow, c.totalShow),
			TotalRawRetrievalCount: numeric(b.totalRaw, c.totalRaw),
		},
	}, nil
}

func normalize(result *Result) (*normalized, error) {
	switch result.Kind {
	case KindRepeated:
		if result.Repeated == nil {
			return nil, errors.New("missing repeated result")
		}
		r := result.Repeated
		if r.SchemaVersion != 4 {
			return nil, fmt.Errorf("unsupported repeated benchmark schema version %d", r.SchemaVersion)
		}
		return &normalized{
			side: Side{
				Kind: KindRepeated, ResultSchemaVersion: r.SchemaVersion,
				Workload: r.Workload, Agent: r.Agent, Mode: r.Mode,
				TrialCount: r.TrialCount, RequestedTrialCount: r.RequestedTrialCount, RunStatus: r.RunStatus,
				SuccessCount: intPointer(r.SuccessCount), TaskFailureCount: intPointer(r.TaskFailureCount), TimeoutCount: intPointer(r.TimeoutCount),
				AgentErrorCount: intPointer(r.AgentErrorCount), CanceledCount: intPointer(r.CanceledCount),
			},
			visible: r.Aggregate.MedianTotalVisibleBytes, commands: r.Aggregate.MedianCommandCount,
			wall: r.Aggregate.MedianWallTimeNS, processing: r.Aggregate.MedianProcessingNS,
			trialsWithShow: integer(int64(r.Aggregate.TrialsWithShow)), trialsWithRaw: integer(int64(r.Aggregate.TrialsWithRawRetrieval)),
			totalShow: integer(int64(r.Aggregate.TotalShowCount)), totalRaw: integer(int64(r.Aggregate.TotalRawRetrievalCount)),
		}, nil
	case KindSingle:
		if result.Single == nil {
			return nil, errors.New("missing single result")
		}
		r := result.Single
		if r.SchemaVersion != 3 && r.SchemaVersion != 4 {
			return nil, fmt.Errorf("unsupported single benchmark schema version %d", r.SchemaVersion)
		}
		side := Side{Kind: KindSingle, ResultSchemaVersion: r.SchemaVersion, Workload: r.Workload, Agent: r.Agent, Mode: r.Mode, TrialCount: 1, RequestedTrialCount: 1, RunStatus: "completed"}
		if r.TaskSuccess != nil {
			side.TimeoutCount = intPointer(0)
			side.AgentErrorCount = intPointer(0)
			side.CanceledCount = intPointer(0)
			if *r.TaskSuccess {
				side.SuccessCount = intPointer(1)
				side.TaskFailureCount = intPointer(0)
			} else {
				side.SuccessCount = intPointer(0)
				side.TaskFailureCount = intPointer(1)
			}
		}
		switch r.ExecutionStatus {
		case "timeout":
			side.TimeoutCount = intPointer(1)
		case "agent_error":
			side.AgentErrorCount = intPointer(1)
		case "canceled":
			side.CanceledCount = intPointer(1)
		}
		var wall *int64
		if result.singleWallTimeAvailable {
			wall = integer(r.WallTimeNS)
		}
		return &normalized{
			side:    side,
			visible: integer(r.TotalVisibleBytes), commands: integer(int64(r.Commands)),
			wall: wall, processing: integer(r.ProcessingNS),
			trialsWithShow: integer(boolInt64(r.ShowCount > 0)), trialsWithRaw: integer(boolInt64(r.RawRetrievalCount > 0)),
			totalShow: integer(int64(r.ShowCount)), totalRaw: integer(int64(r.RawRetrievalCount)),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported result kind %q", result.Kind)
	}
}

func numeric(baseline, candidate *int64) NumericComparison {
	result := NumericComparison{Baseline: copyInteger(baseline), Candidate: copyInteger(candidate)}
	if baseline == nil || candidate == nil {
		return result
	}
	delta := *candidate - *baseline
	result.Delta = &delta
	if *baseline != 0 {
		relative := float64(delta) / float64(*baseline)
		result.RelativeDelta = &relative
	}
	return result
}

func integer(value int64) *int64 { return &value }

func intPointer(value int) *int { return &value }

func copyInteger(value *int64) *int64 {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func boolInt64(value bool) int64 {
	if value {
		return 1
	}
	return 0
}
