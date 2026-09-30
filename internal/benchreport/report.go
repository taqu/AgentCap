// Package benchreport aggregates Phase 7 coding-agent benchmark records into
// batch summaries, paired OFF/FULL comparisons, and single-run inspections.
// It never executes anything and never adjusts or discards records silently:
// excluded runs are listed with their reason.
package benchreport

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/taqu/agentcap/internal/workload"
)

// Load reads every single-trial result JSON carrying workflow metrics below dir.
func Load(dir string) ([]*workload.BenchmarkResult, error) {
	var out []*workload.BenchmarkResult
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".json") || filepath.Base(path) == "result.json" {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var r workload.BenchmarkResult
		if json.Unmarshal(data, &r) != nil || r.Workflow == nil {
			return nil
		}
		out = append(out, &r)
		return nil
	})
	sort.Slice(out, func(i, j int) bool { return out[i].Workflow.RunID < out[j].Workflow.RunID })
	return out, err
}

// Excluded reports why a record must not enter performance aggregates (§43, §44).
func Excluded(r *workload.BenchmarkResult) string {
	if r.Workflow.Outcome == "INFRA_ERROR" {
		return "infrastructure error"
	}
	for _, a := range r.Workflow.Anomalies {
		if strings.HasPrefix(a, "ADAPTER:") {
			return strings.TrimSpace(strings.TrimPrefix(a, "ADAPTER:"))
		}
	}
	return ""
}

func modeLabel(m string) string {
	switch m {
	case "disabled":
		return "OFF"
	case "integrated":
		return "FULL"
	}
	return strings.ToUpper(m)
}

// ---------------------------------------------------------------- statistics

type dist struct{ v []float64 }

func (d *dist) add(x float64) { d.v = append(d.v, x) }
func (d *dist) pct(p float64) float64 {
	if len(d.v) == 0 {
		return 0
	}
	s := append([]float64(nil), d.v...)
	sort.Float64s(s)
	pos := p * float64(len(s)-1)
	lo := int(pos)
	if lo+1 >= len(s) {
		return s[lo]
	}
	return s[lo] + (s[lo+1]-s[lo])*(pos-float64(lo))
}
func (d *dist) med() float64 { return d.pct(0.5) }
func (d *dist) sum() float64 {
	t := 0.0
	for _, x := range d.v {
		t += x
	}
	return t
}

// group accumulates one agent x mode (x optional task/category) cell.
type group struct {
	runs, pass, fail, timeout            int
	bytes, bytesOK, wall, cmds, repeated dist
	hostIn, native                       dist
	show, raw, immRaw, results           int
	showBytes, rawBytes                  int64
	trialsRaw, trialsShow                int
	intercepted, bypassed, adapterFail   int
	adapterNS, coreNS                    int64
	full, delta, unchanged               int
	families                             map[string]int64
	familyCount                          map[string]int
}

func newGroup() *group { return &group{families: map[string]int64{}, familyCount: map[string]int{}} }

func (g *group) add(r *workload.BenchmarkResult) {
	w := r.Workflow
	g.runs++
	switch w.Outcome {
	case "PASS":
		g.pass++
		g.bytesOK.add(float64(w.TotalAcapCostBytes))
	case "TIMEOUT":
		g.timeout++
	default:
		g.fail++
	}
	g.bytes.add(float64(w.TotalAcapCostBytes))
	g.wall.add(float64(r.WallTimeNS) / 1e9)
	g.cmds.add(float64(w.ShellCommands))
	g.repeated.add(float64(w.RepeatedCommands))
	g.native.add(float64(w.NativeToolBytes))
	if w.HostInputTokens+w.HostCacheReadTokens+w.HostCacheWriteTokens > 0 {
		g.hostIn.add(float64(w.HostInputTokens + w.HostCacheReadTokens + w.HostCacheWriteTokens))
	}
	g.show += w.ShowCalls
	g.raw += w.RawCalls
	g.immRaw += w.ImmediateRawCalls
	g.showBytes += w.ShowVisibleBytes
	g.rawBytes += w.RawVisibleBytes
	g.results += w.AcapResults
	if w.RawCalls > 0 {
		g.trialsRaw++
	}
	if w.ShowCalls > 0 {
		g.trialsShow++
	}
	g.intercepted += w.Intercepted
	g.bypassed += w.Bypassed
	g.adapterFail += w.AdapterFailures
	g.adapterNS += w.AdapterLatencyNS
	g.coreNS += w.CoreLatencyNS
	g.full += w.FullResults
	g.delta += w.DeltaResults
	g.unchanged += w.UnchangedResults
	for k, v := range w.FamilyBytes {
		g.families[k] += v
	}
	for k, v := range w.FamilyCount {
		g.familyCount[k] += v
	}
}

func rate(a, b int) string {
	if b == 0 {
		return "–"
	}
	return fmt.Sprintf("%.0f%%", 100*float64(a)/float64(b))
}

func savings(on, off float64) string {
	if off == 0 {
		return "–"
	}
	return fmt.Sprintf("%+.0f%%", -100*(1-on/off))
}

func kb(x float64) string {
	if x >= 10240 {
		return fmt.Sprintf("%.1f KiB", x/1024)
	}
	return fmt.Sprintf("%.0f B", x)
}

// ---------------------------------------------------------------- report

type key struct{ agent, mode, sub string }

// Write renders the Markdown batch report (§61–§68, §99).
func Write(w io.Writer, records []*workload.BenchmarkResult) {
	var included []*workload.BenchmarkResult
	var excluded []*workload.BenchmarkResult
	for _, r := range records {
		if Excluded(r) != "" {
			excluded = append(excluded, r)
		} else {
			included = append(included, r)
		}
	}
	agents, modes, tasks, cats := set(), set(), set(), set()
	batches, versions, models := set(), set(), set()
	cells := map[key]*group{}
	get := func(k key) *group {
		if cells[k] == nil {
			cells[k] = newGroup()
		}
		return cells[k]
	}
	for _, r := range included {
		wf := r.Workflow
		agents[r.Agent], modes[r.Mode], tasks[r.Workload], cats[wf.Category] = true, true, true, true
		batches[wf.Batch], versions[wf.AcapCommit], models[r.Agent+": "+wf.Model+" ("+wf.AgentVer+")"] = true, true, true
		get(key{r.Agent, r.Mode, ""}).add(r)
		get(key{r.Agent, r.Mode, "task:" + r.Workload}).add(r)
		get(key{r.Agent, r.Mode, "cat:" + wf.Category}).add(r)
		get(key{"*", r.Mode, "cat:" + wf.Category}).add(r)
	}

	fmt.Fprintf(w, "# AgentCap Phase 7 benchmark report\n\n")
	fmt.Fprintf(w, "- Batches: %s\n- AgentCap versions: %s\n- Agents/models: %s\n", keys(batches), keys(versions), keys(models))
	fmt.Fprintf(w, "- Tasks: %d · Runs: %d included, %d excluded · Modes: %s\n", len(tasks), len(included), len(excluded), keys(modes))
	fmt.Fprintf(w, "- Primary metric: total agent-visible command-result bytes per workflow (capsules + deltas + `acap show` + `acap raw` + instruction overhead), counted once from the agent's own tool-result stream. Tokens marked *est* are bytes/4; host tokens are reported separately.\n\n")

	for _, agent := range sortedKeys(agents) {
		fmt.Fprintf(w, "## %s\n\n", agent)
		fmt.Fprintf(w, "| mode | runs | pass | median context (IQR) | tokens est | median ctx, passing | savings vs OFF (all / passing) | median wall s | median shell cmds | repeated | host input tok (median) | native tool bytes (median) |\n|---|---|---|---|---|---|---|---|---|---|---|---|\n")
		off := cells[key{agent, "disabled", ""}]
		for _, m := range []string{"disabled", "integrated"} {
			g := cells[key{agent, m, ""}]
			if g == nil {
				continue
			}
			sAll, sOK := "–", "–"
			if m != "disabled" && off != nil {
				sAll, sOK = savings(g.bytes.med(), off.bytes.med()), savings(g.bytesOK.med(), off.bytesOK.med())
			}
			fmt.Fprintf(w, "| %s | %d | %d/%d (%s) | %s (%s–%s) | %.0f | %s | %s / %s | %.0f | %.0f | %.0f | %.0f | %s |\n",
				modeLabel(m), g.runs, g.pass, g.runs, rate(g.pass, g.runs), kb(g.bytes.med()), kb(g.bytes.pct(.25)), kb(g.bytes.pct(.75)),
				g.bytes.med()/4, kb(g.bytesOK.med()), sAll, sOK, g.wall.med(), g.cmds.med(), g.repeated.med(), g.hostIn.med(), kb(g.native.med()))
		}
		if g := cells[key{agent, "integrated", ""}]; g != nil {
			fmt.Fprintf(w, "\n**AgentCap behavior (FULL):** %d results (full %d / delta %d / unchanged %d) · intercepted %d, bypassed %d, adapter failures %d · `show` %d calls (%s, in %d runs) · `raw` %d calls (%s, in %d runs; immediate %d) · raw fallback rate %s of results · targeted drill-down rate %s of results\n",
				g.results, g.full, g.delta, g.unchanged, g.intercepted, g.bypassed, g.adapterFail, g.show, kb(float64(g.showBytes)), g.trialsShow, g.raw, kb(float64(g.rawBytes)), g.trialsRaw, g.immRaw, rate(g.raw, g.results), rate(g.show, g.results))
			if g.intercepted > 0 {
				fmt.Fprintf(w, "\n**Latency (FULL, per intercepted command):** adapter %.2f ms · AgentCap core %.2f ms\n", float64(g.adapterNS)/1e6/float64(g.intercepted), float64(g.coreNS)/1e6/float64(g.intercepted))
			}
		}
		if g := cells[key{agent, "disabled", ""}]; g != nil {
			fmt.Fprintf(w, "\n**OFF isolation:** %d hook decisions all bypassed (%d intercepted).\n", g.bypassed, g.intercepted)
		}
		fmt.Fprintf(w, "\n### %s — per task (medians; OFF → FULL)\n\n| task | pass OFF | pass FULL | context OFF | context FULL | Δ | cmds OFF→FULL | wall s OFF→FULL | show/raw FULL |\n|---|---|---|---|---|---|---|---|---|\n", agent)
		for _, t := range sortedKeys(tasks) {
			a, b := cells[key{agent, "disabled", "task:" + t}], cells[key{agent, "integrated", "task:" + t}]
			if a == nil || b == nil {
				continue
			}
			fmt.Fprintf(w, "| %s | %d/%d | %d/%d | %s | %s | %s | %.0f→%.0f | %.0f→%.0f | %d/%d |\n", strings.TrimPrefix(t, "phase7/"),
				a.pass, a.runs, b.pass, b.runs, kb(a.bytes.med()), kb(b.bytes.med()), savings(b.bytes.med(), a.bytes.med()),
				a.cmds.med(), b.cmds.med(), a.wall.med(), b.wall.med(), b.show, b.raw)
		}
		fmt.Fprintf(w, "\n### %s — command families (total bytes over all runs)\n\n| family | OFF bytes | OFF calls | FULL bytes | FULL calls | Δ |\n|---|---|---|---|---|---|\n", agent)
		fam := set()
		for _, m := range []string{"disabled", "integrated"} {
			if g := cells[key{agent, m, ""}]; g != nil {
				for f := range g.families {
					fam[f] = true
				}
			}
		}
		famList := sortedKeys(fam)
		offF, onF := cells[key{agent, "disabled", ""}], cells[key{agent, "integrated", ""}]
		sort.Slice(famList, func(i, j int) bool {
			return famTotal(offF, famList[i])+famTotal(onF, famList[i]) > famTotal(offF, famList[j])+famTotal(onF, famList[j])
		})
		for _, f := range famList {
			fmt.Fprintf(w, "| %s | %s | %d | %s | %d | %s |\n", f, kb(float64(famTotal(offF, f))), famCount(offF, f), kb(float64(famTotal(onF, f))), famCount(onF, f),
				savings(float64(famTotal(onF, f)), float64(famTotal(offF, f))))
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintf(w, "## Per task category (all agents, medians)\n\n| category | runs OFF/FULL | pass OFF | pass FULL | context OFF | context FULL | Δ |\n|---|---|---|---|---|---|---|\n")
	for _, c := range sortedKeys(cats) {
		a, b := cells[key{"*", "disabled", "cat:" + c}], cells[key{"*", "integrated", "cat:" + c}]
		if a == nil || b == nil {
			continue
		}
		fmt.Fprintf(w, "| %s | %d/%d | %s | %s | %s | %s | %s |\n", c, a.runs, b.runs, rate(a.pass, a.runs), rate(b.pass, b.runs), kb(a.bytes.med()), kb(b.bytes.med()), savings(b.bytes.med(), a.bytes.med()))
	}

	fmt.Fprintf(w, "\n## Anomalies and excluded runs (for manual review)\n\n")
	n := 0
	for _, r := range records {
		reason := Excluded(r)
		notes := r.Workflow.Anomalies
		if reason == "" && len(notes) == 0 && !(r.Mode == "integrated" && r.Workflow.Outcome != "PASS") {
			continue
		}
		n++
		fmt.Fprintf(w, "- `%s` %s %s %s — outcome %s", r.Workflow.RunID, r.Agent, modeLabel(r.Mode), strings.TrimPrefix(r.Workload, "phase7/"), r.Workflow.Outcome)
		if reason != "" {
			fmt.Fprintf(w, " — **excluded: %s**", reason)
		}
		if len(notes) > 0 {
			fmt.Fprintf(w, " — %s", strings.Join(notes, "; "))
		}
		fmt.Fprintln(w)
	}
	if n == 0 {
		fmt.Fprintln(w, "None.")
	}
}

func famTotal(g *group, f string) int64 {
	if g == nil {
		return 0
	}
	return g.families[f]
}
func famCount(g *group, f string) int {
	if g == nil {
		return 0
	}
	return g.familyCount[f]
}

func set() map[string]bool { return map[string]bool{} }
func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
func keys(m map[string]bool) string {
	out := sortedKeys(m)
	for i, s := range out {
		if s == "" {
			out[i] = "(unset)"
		}
	}
	return strings.Join(out, ", ")
}

// WriteRun renders one run's metadata, timeline and metrics (§97).
func WriteRun(w io.Writer, r *workload.BenchmarkResult) {
	wf := r.Workflow
	fmt.Fprintf(w, "run %s  batch %s\n", wf.RunID, wf.Batch)
	fmt.Fprintf(w, "task %s (%s, %s)  agent %s %s  model %s  mode %s  acap %s adapter %s  %s  started %s\n",
		r.Workload, wf.Category, wf.Language, r.Agent, wf.AgentVer, wf.Model, modeLabel(r.Mode), wf.AcapCommit, wf.AdapterVer, wf.Platform, wf.StartedAt)
	fmt.Fprintf(w, "outcome %s  wall %.1fs  shell cmds %d (repeated %d)  context %d B (~%d tok est)  show %d/%d B  raw %d/%d B  native tools %d/%d B\n",
		wf.Outcome, float64(r.WallTimeNS)/1e9, wf.ShellCommands, wf.RepeatedCommands, wf.TotalAcapCostBytes, wf.TotalAcapCostTokens,
		wf.ShowCalls, wf.ShowVisibleBytes, wf.RawCalls, wf.RawVisibleBytes, wf.NativeToolCalls, wf.NativeToolBytes)
	fmt.Fprintf(w, "store: raw %d B, stateless %d B, stateful %d B  intercepted %d bypassed %d %v  adapter %.2fms core %.2fms\n",
		r.RawBytes, r.StatelessBytes, r.StatefulBytes, wf.Intercepted, wf.Bypassed, wf.BypassReasons, float64(wf.AdapterLatencyNS)/1e6, float64(wf.CoreLatencyNS)/1e6)
	fmt.Fprintf(w, "host tokens: input %d, cache read %d, cache write %d, output %d\n\ntimeline:\n", wf.HostInputTokens, wf.HostCacheReadTokens, wf.HostCacheWriteTokens, wf.HostOutputTokens)
	for _, e := range wf.Trace {
		meta := ""
		if e.ResultID != "" {
			meta = fmt.Sprintf(" [@acap %s %s]", e.ResultID, e.Kind)
		}
		errMark := ""
		if e.IsError {
			errMark = " (error)"
		}
		label := e.Command
		if e.Tool != "shell" {
			label = "<" + e.Tool + ">"
		}
		fmt.Fprintf(w, "%3d %-10s %7d B%s%s  %s\n", e.Index, e.Family, e.Bytes, meta, errMark, label)
	}
	for _, a := range wf.Anomalies {
		fmt.Fprintf(w, "anomaly: %s\n", a)
	}
}
