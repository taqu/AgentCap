package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/taqu/agentcap/internal/project"
	"github.com/taqu/agentcap/internal/regression"
	"github.com/taqu/agentcap/internal/workload"
)

// Exit statuses of "acap bench check". A failed policy and an evaluation
// that could not be completed are both non-zero but remain distinguishable.
const (
	exitRegression      = 1
	exitEvaluationError = 2
)

const benchSuiteUsage = `Usage: acap bench suite --out <dir> <regression.yaml>

Runs every deterministic workload listed in a regression suite, exactly as
"acap bench run" does (fresh temporary workspace and benchmark session per
workload), and writes one JSON result per workload into <dir>. <dir> must not
exist or must be empty so stale results are never evaluated.

Coding-agent workloads are rejected: regression suites are deterministic.

Flags:
  --out <dir>  Directory for the per-workload result files (required).
`

const benchCheckUsage = `Usage: acap bench check --baseline <path> --candidate <path> --policy <regression.yaml> [--json]

Applies an explicit regression policy to previously recorded results. <path>
is a result file or a directory of result files such as "acap bench suite"
writes. Nothing is executed and no baseline, policy, or result is modified.

Every workload listed in the policy is compared with "acap bench compare"
semantics (delta = candidate - baseline) and every metric rule is checked for
that workload separately. A check fails when the increase exceeds every
configured limit it can be evaluated against; an increase equal to a limit
passes. A missing workload, a missing measurement, or an undefined relative
change without an absolute limit is "not evaluable" and also fails.

Exit status:
  0  every check passed
  1  at least one regression or not-evaluable check
  2  evaluation could not be completed (invalid policy or results)

Flags:
  --baseline <path>   Baseline result file or directory (required).
  --candidate <path>  Candidate result file or directory (required).
  --policy <file>     Regression suite/policy YAML (required).
  --json              Print the versioned evaluation as JSON only.
`

func benchSuiteCmd(args []string) {
	fs := flag.NewFlagSet("bench suite", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	outFlag := fs.String("out", "", "result directory")
	if len(args) > 0 && (args[0] == "--help" || args[0] == "-h") {
		fmt.Fprint(os.Stdout, benchSuiteUsage)
		os.Exit(0)
	}
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "acap: bench suite: %v\n\n%s", err, benchSuiteUsage)
		os.Exit(1)
	}
	if *outFlag == "" || len(fs.Args()) != 1 {
		fmt.Fprint(os.Stderr, benchSuiteUsage)
		os.Exit(1)
	}
	suite, err := regression.LoadSuite(fs.Args()[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: bench suite: %v\n", err)
		os.Exit(1)
	}
	if err := prepareEmptyDir(*outFlag); err != nil {
		fmt.Fprintf(os.Stderr, "acap: bench suite: %v\n", err)
		os.Exit(1)
	}
	cwd, _ := os.Getwd()
	for _, w := range suite.Workloads {
		run, err := workload.Run(context.Background(), w.Definition, workload.Options{StoreRoot: project.FindRoot(cwd)})
		if err != nil {
			fmt.Fprintf(os.Stderr, "acap: bench suite: %s: %v\n", w.Name, err)
			os.Exit(1)
		}
		path := filepath.Join(*outFlag, workload.ResultFileName(w.Name))
		if err := writeResultFile(path, workload.NewBenchmarkResult(run)); err != nil {
			fmt.Fprintf(os.Stderr, "acap: bench suite: %s: %v\n", w.Name, err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "%s: %s\n", w.Name, path)
	}
}

func prepareEmptyDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return os.MkdirAll(dir, 0o755)
	}
	if err != nil {
		return err
	}
	if len(entries) != 0 {
		return fmt.Errorf("output directory %s is not empty", dir)
	}
	return nil
}

func writeResultFile(path string, result *workload.BenchmarkResult) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	if err := workload.WriteJSON(f, result); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func benchCheckCmd(args []string) {
	fs := flag.NewFlagSet("bench check", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	baselineFlag := fs.String("baseline", "", "baseline results")
	candidateFlag := fs.String("candidate", "", "candidate results")
	policyFlag := fs.String("policy", "", "regression policy")
	jsonFlag := fs.Bool("json", false, "print JSON evaluation")
	if len(args) > 0 && (args[0] == "--help" || args[0] == "-h") {
		fmt.Fprint(os.Stdout, benchCheckUsage)
		os.Exit(0)
	}
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "acap: bench check: %v\n\n%s", err, benchCheckUsage)
		os.Exit(exitEvaluationError)
	}
	if *baselineFlag == "" || *candidateFlag == "" || *policyFlag == "" || len(fs.Args()) != 0 {
		fmt.Fprint(os.Stderr, benchCheckUsage)
		os.Exit(exitEvaluationError)
	}
	evaluation, err := runCheck(*baselineFlag, *candidateFlag, *policyFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: bench check: %v\n", err)
		evaluation = regression.ErrorEvaluation(err)
	}
	if *jsonFlag {
		err = regression.WriteJSON(os.Stdout, evaluation)
	} else {
		err = regression.WriteHuman(os.Stdout, evaluation)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: bench check: %v\n", err)
		os.Exit(exitEvaluationError)
	}
	switch evaluation.Status {
	case regression.StatusPassed:
	case regression.StatusFailed:
		os.Exit(exitRegression)
	default:
		os.Exit(exitEvaluationError)
	}
}

func runCheck(baselinePath, candidatePath, policyPath string) (*regression.Evaluation, error) {
	suite, err := regression.LoadSuite(policyPath)
	if err != nil {
		return nil, err
	}
	baseline, err := workload.LoadResultSet(baselinePath)
	if err != nil {
		return nil, fmt.Errorf("baseline: %w", err)
	}
	candidate, err := workload.LoadResultSet(candidatePath)
	if err != nil {
		return nil, fmt.Errorf("candidate: %w", err)
	}
	names := make([]string, len(suite.Workloads))
	for i, w := range suite.Workloads {
		names[i] = w.Name
	}
	return regression.EvaluateSuite(names, &suite.Policy, baseline, candidate), nil
}
