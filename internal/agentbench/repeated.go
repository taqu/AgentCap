package agentbench

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/taqu/agentcap/internal/workload"
)

// RepeatedRun contains the stable repeated result and the internal trial data
// needed for verbose logs and retained-workspace reporting.
type RepeatedRun struct {
	Benchmark *workload.RepeatedBenchmarkResult
	Trials    []*Trial
}

// RunRepeated sequentially invokes the B6 single-trial runner with an isolated
// workspace, store, and session scope for every repetition.
func RunRepeated(ctx context.Context, definition *workload.Definition, opts Options, repeat int) (*RepeatedRun, error) {
	if repeat < 1 {
		return nil, fmt.Errorf("repeat must be at least 1 (got %d)", repeat)
	}
	if definition == nil {
		return nil, errors.New("workload definition is required")
	}
	run := &RepeatedRun{}
	for i := 1; i <= repeat; i++ {
		trialRoot, err := os.MkdirTemp(opts.TempRoot, "agentcap-agent-trial-")
		if err != nil {
			run.Benchmark, _ = AggregateTrials(definition.Name, adapterName(opts.Adapter), opts.Mode, repeat, run.Trials, "incomplete")
			return run, fmt.Errorf("trial %d workspace root: %w", i, err)
		}
		stateRoot := filepath.Join(trialRoot, "state")
		if err := os.Mkdir(stateRoot, 0o700); err != nil {
			_ = os.RemoveAll(trialRoot)
			run.Benchmark, _ = AggregateTrials(definition.Name, adapterName(opts.Adapter), opts.Mode, repeat, run.Trials, "incomplete")
			return run, fmt.Errorf("trial %d state root: %w", i, err)
		}

		trialOpts := opts
		trialOpts.TempRoot = trialRoot
		trialOpts.StoreRoot = stateRoot
		trial, trialErr := Run(ctx, definition, trialOpts)
		if trial != nil {
			if trial.Benchmark != nil {
				trial.Benchmark.Trial = i
			}
			run.Trials = append(run.Trials, trial)
		}
		if opts.KeepWorkspace && trial != nil {
			trial.RetainedStateRoot = stateRoot
		} else if removeErr := os.RemoveAll(trialRoot); trialErr == nil && removeErr != nil {
			trialErr = fmt.Errorf("clean trial %d state: %w", i, removeErr)
		}
		if trialErr != nil {
			run.Benchmark, _ = AggregateTrials(definition.Name, adapterName(opts.Adapter), opts.Mode, repeat, run.Trials, "incomplete")
			return run, fmt.Errorf("trial %d: %w", i, trialErr)
		}
	}

	var err error
	run.Benchmark, err = AggregateTrials(definition.Name, adapterName(opts.Adapter), opts.Mode, repeat, run.Trials, "completed")
	if err != nil {
		return run, err
	}
	return run, nil
}

func adapterName(adapter Adapter) string {
	if adapter == nil {
		return ""
	}
	return adapter.Name()
}

// AggregateTrials is the pure canonical aggregation layer. All non-nil valid
// trials participate, including task failures, timeouts, and agent errors.
func AggregateTrials(workloadName, agent string, mode Mode, requested int, trials []*Trial, status string) (*workload.RepeatedBenchmarkResult, error) {
	result := &workload.RepeatedBenchmarkResult{
		SchemaVersion: workload.BenchmarkResultSchemaVersion,
		Workload:      workloadName, Agent: agent, Mode: string(mode),
		RequestedTrialCount: requested, RunStatus: status,
		Trials: make([]*workload.BenchmarkResult, 0, len(trials)),
	}
	var visible, commands, wall, processing []int64
	for _, trial := range trials {
		if trial == nil || trial.Benchmark == nil {
			continue
		}
		b := trial.Benchmark
		if b.Workload != workloadName || b.Agent != agent || b.Mode != string(mode) {
			return nil, fmt.Errorf("incompatible trial %d configuration", b.Trial)
		}
		result.Trials = append(result.Trials, b)
		visible = append(visible, b.TotalVisibleBytes)
		commands = append(commands, int64(b.Commands))
		wall = append(wall, b.WallTimeNS)
		processing = append(processing, b.ProcessingNS)
		if b.TaskSuccess != nil {
			if *b.TaskSuccess {
				result.SuccessCount++
			} else {
				result.TaskFailureCount++
			}
		}
		switch b.ExecutionStatus {
		case "timeout":
			result.TimeoutCount++
		case "agent_error":
			result.AgentErrorCount++
		case "canceled":
			result.CanceledCount++
		}
		if b.ShowCount > 0 {
			result.Aggregate.TrialsWithShow++
		}
		if b.RawRetrievalCount > 0 {
			result.Aggregate.TrialsWithRawRetrieval++
		}
		result.Aggregate.TotalShowCount += b.ShowCount
		result.Aggregate.TotalRawRetrievalCount += b.RawRetrievalCount
	}
	result.TrialCount = len(result.Trials)
	result.Aggregate.MedianTotalVisibleBytes = medianValue(visible)
	result.Aggregate.MedianCommandCount = medianValue(commands)
	result.Aggregate.MedianWallTimeNS = medianValue(wall)
	result.Aggregate.MedianProcessingNS = medianValue(processing)
	return result, nil
}

func medianValue(values []int64) *int64 {
	if len(values) == 0 {
		return nil
	}
	value := median(values)
	return &value
}

func median(values []int64) int64 {
	if len(values) == 0 {
		return 0
	}
	ordered := append([]int64(nil), values...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i] < ordered[j] })
	middle := len(ordered) / 2
	if len(ordered)%2 == 1 {
		return ordered[middle]
	}
	lower, upper := ordered[middle-1], ordered[middle]
	return lower + (upper-lower)/2
}
