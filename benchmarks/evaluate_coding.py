"""Orchestrate existing acap bench agent/compare; no alternative agent runner.

python benchmarks/evaluate_coding.py --out benchmarks/results/smoke --smoke
python benchmarks/evaluate_coding.py --out benchmarks/results/main
"""
import argparse
import hashlib
import json
import statistics
import subprocess
import time
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parent
REPO = ROOT.parent

def call(args, **kwargs):
    return subprocess.run(args, capture_output=True, text=True, encoding='utf-8', errors='replace', **kwargs)

def digest(path):
    h = hashlib.sha256()
    for f in sorted(path.rglob('*')):
        if f.is_file():
            h.update(f.relative_to(path).as_posix().encode())
            h.update(f.read_bytes())
    return h.hexdigest()

def summarize(rows):
    result = {'runs': len(rows), 'successful': sum(r['success'] for r in rows)}
    for field in ['tokens', 'seconds', 'visible_bytes', 'commands', 'tools', 'drilldowns','repeated_commands','build_test_invocations','failed_validation_attempts']:
        values = [r[field] for r in rows if r.get(field) is not None]
        result[field] = {'measured':len(values), 'mean':statistics.mean(values) if values else None, 'median':statistics.median(values) if values else None}
    return result

def transcript_metrics(path):
    commands=[]; tools={}; failures=[];names={};bytes_by_tool={};calls_by_tool={};provider_errors=[]
    for line in path.read_text(encoding='utf-8').splitlines():
        if not line.startswith('{'):continue
        try: event=json.loads(line)
        except ValueError:continue
        if event.get('type')=='result' and event.get('is_error'):
            provider_errors.append(event.get('result',event.get('subtype','unknown agent error')))
        for item in event.get('message',{}).get('content',[]):
            if item.get('type')=='tool_use':
                name=item.get('name','unknown');names[item['id']]=name
                calls_by_tool[name]=calls_by_tool.get(name,0)+1
            if item.get('type')=='tool_result':
                content=item.get('content','')
                if isinstance(content,list):content=''.join(b.get('text','') for b in content if b.get('type')=='text')
                if isinstance(content,str):
                    name=names.get(item.get('tool_use_id'),'unknown')
                    bytes_by_tool[name]=bytes_by_tool.get(name,0)+len(content.encode('utf-8'))
            if item.get('type')=='tool_use' and item.get('name')=='Bash':
                command=item.get('input',{}).get('command','')
                tools[item['id']]=command;commands.append(command)
            if item.get('type')=='tool_result' and item.get('tool_use_id') in tools and item.get('is_error'):
                failures.append(tools[item['tool_use_id']])
    validation=lambda s:bool(re.search(r'\bgo\s+(test|build|vet)\b',s))
    stateful=lambda s:bool(re.search(r'(^|[^a-zA-Z0-9_])(cd|pushd|popd|export|unset|source|alias|unalias|umask|set|trap|eval|exec|acap|acap\.exe)([^a-zA-Z0-9_]|$)',s))
    return dict(command_sequence=commands,repeated_command_count=len(commands)-len(set(commands)),build_test_invocations=sum(validation(s) for s in commands),failed_validation_attempts=sum(validation(s) for s in failures),stateful_shell_commands=sum(stateful(s) for s in commands),failed_shell_commands=failures,provider_errors=provider_errors,tool_output_bytes_by_type=bytes_by_tool,tool_calls_by_type=calls_by_tool,repair_iterations=None)

def report(out, manifest):
    rows = []
    for run in manifest['runs']:
        path = out / (run['key']+'.json')
        try:
            value = json.loads(path.read_text(encoding='utf-8'))
        except (ValueError, FileNotFoundError):
            continue
        if 'task_success' not in value:
            continue
        tokens = sum(value.get(k,0) for k in ['input_tokens','output_tokens','cache_creation_tokens','cache_read_tokens']) if value.get('usage_available') else None
        workflow=transcript_metrics(out/(run['key']+'.log'))
        rows.append(dict(key=run['key'],id=run['id'],level=run['level'],mode=run['mode'], repetition=run['repetition'], success=value['task_success'] and value['execution_status']=='completed',status=value['execution_status'],tokens=tokens,seconds=value['wall_time_ns']/1e9,visible_bytes=value.get('tool_visible_bytes',0),shell_bytes=value['total_visible_bytes'],commands=value['commands'],tools=value.get('tool_count'),drilldowns=value['show_count']+value['raw_retrieval_count'],repeated_commands=workflow['repeated_command_count'],build_test_invocations=workflow['build_test_invocations'],failed_validation_attempts=workflow['failed_validation_attempts'],model=value.get('model'),canonical_result=value,workflow=workflow))
    modes = manifest.get('analysis_modes', manifest['config']['modes'])
    groups = []
    for level in range(1,6):
        for mode in modes:
            selected = [r for r in rows if r['level']==level and r['mode']==mode]
            groups.append(dict(level=level,mode=mode,unconditional=summarize(selected),successful_only=summarize([r for r in selected if r['success']])))
    pairs=[]
    for base in rows:
        if base['mode']!='disabled': continue
        for mode in modes:
            if mode=='disabled': continue
            candidate=next((r for r in rows if r['mode']==mode and r['id']==base['id'] and r['repetition']==base['repetition']),None)
            if candidate is None: continue
            def reduction(field):
                a,b=base.get(field),candidate.get(field)
                return 100*(a-b)/a if a and b is not None else None
            pairs.append(dict(id=base['id'],level=base['level'],mode=mode,repetition=base['repetition'],both_successful=base['success'] and candidate['success'],token_reduction_percent=reduction('tokens'),output_reduction_percent=reduction('visible_bytes'),execution_overhead_percent=-reduction('seconds')))
    paired_levels=[]
    for level in range(1,6):
        for mode in modes:
            selected=[p for p in pairs if p['level']==level and p['mode']==mode]
            stats={'level':level,'mode':mode,'pairs':len(selected),'joint_successes':sum(p['both_successful'] for p in selected)}
            for scope in ['unconditional','successful_only']:
                items=selected if scope=='unconditional' else [p for p in selected if p['both_successful']]
                stats[scope]={}
                for field in ['token_reduction_percent','output_reduction_percent','execution_overhead_percent']:
                    values=[p[field] for p in items if p[field] is not None]
                    stats[scope][field]=statistics.mean(values) if values else None
            paired_levels.append(stats)
    aggregate=dict(rows=rows,groups=groups,pairs=pairs,paired_levels=paired_levels)
    (out/'aggregate.json').write_text(json.dumps(aggregate,indent=2)+'\n',encoding='utf-8')
    def fmt(value): return 'unavailable' if value is None else f'{value:.1f}'
    good=[p for p in pairs if p['both_successful']]
    reductions=[p['token_reduction_percent'] for p in good if p['token_reduction_percent'] is not None]
    overheads=[p['execution_overhead_percent'] for p in good]
    lines=['# Claude Code coding workflow evaluation','',f"Attempted {len(manifest['runs'])} conditions; recorded {len(rows)} trials across {len(set(r['id'] for r in rows))} workloads. Successful: {sum(r['success'] for r in rows)}; unsuccessful: {sum(not r['success'] for r in rows)}.",'',f"Mean paired token reduction among jointly successful enabled-mode pairs: {fmt(statistics.mean(reductions) if reductions else None)}%. Mean paired time overhead: {fmt(statistics.mean(overheads) if overheads else None)}%.",'',f"Compared modes: {', '.join(modes)}. Each enabled run is paired with the matching disabled baseline by workload and repetition.",'','Tokens include uncached input, cache creation, cache reads, and output, using provider-reported final usage. Bytes are UTF-8 textual tool-result content; they are never converted to tokens. AgentCap shell bytes include drill-down. Resource aggregates include failures; successful-only aggregates and matched successful pairs are separate in aggregate.json.','', '| Level | Mode | Success | Mean / median tokens | Mean / median seconds | Mean tools / Bash | Show + raw | Repeated Bash | Build/test | Failed validations |','|---|---|---|---|---|---|---|---|---|---|']
    for g in groups:
        s=g['unconditional'];lines.append(f"| L{g['level']} | {g['mode']} | {s['successful']}/{s['runs']} | {fmt(s['tokens']['mean'])} / {fmt(s['tokens']['median'])} | {fmt(s['seconds']['mean'])} / {fmt(s['seconds']['median'])} | {fmt(s['tools']['mean'])} / {fmt(s['commands']['mean'])} | {fmt(s['drilldowns']['mean'])} | {fmt(s['repeated_commands']['mean'])} | {fmt(s['build_test_invocations']['mean'])} | {fmt(s['failed_validation_attempts']['mean'])} |")
    lines += ['','Paired means below use only jointly successful trials. Negative token/output reduction means increased consumption.','', '| Level | Mode | Successful pairs | Token reduction % | Tool-output reduction % | Time overhead % |','|---|---|---|---|---|---|']
    for s in paired_levels:
        v=s['successful_only'];lines.append(f"| L{s['level']} | {s['mode']} | {s['joint_successes']}/{s['pairs']} | {fmt(v['token_reduction_percent'])} | {fmt(v['output_reduction_percent'])} | {fmt(v['execution_overhead_percent'])} |")
    lines += ['', '| Workload | Rep | Mode | Success / status | Model tokens | All-tool B | Bash calls | AgentCap raw → stateful B | Processing ms | Δ / same | Show / raw | Seconds |','|---|---|---|---|---|---|---|---|---|---|---|---|']
    for r in rows:
        cap=r['canonical_result']
        if r['mode']=='disabled':
            bash=f"— / {r['commands']}";raw_state='n/a';processing='n/a';presentations='n/a'
        else:
            bash=f"{cap.get('intercepted_commands',0)} / {r['commands']}"
            raw_state=f"{cap['raw_bytes']} → {cap['stateful_bytes']}"
            processing=f"{cap['processing_ns']/1e6:.2f}"
            presentations=f"{cap.get('delta_count',0)} / {cap.get('unchanged_count',0)}"
        lines.append(f"| {r['id']} | {r['repetition']} | {r['mode']} | {r['success']} / {r['status']} | {fmt(r['tokens'])} | {r['visible_bytes']} | {bash} | {raw_state} | {processing} | {presentations} | {cap['show_count']} / {cap['raw_retrieval_count']} | {r['seconds']:.1f} |")
    enabled=[r for r in rows if r['mode']!='disabled']
    intercepted=sum(r['canonical_result'].get('intercepted_commands',0) for r in enabled)
    shell=sum(r['commands'] for r in enabled)
    lines += ['',f'AgentCap intercepted {intercepted}/{shell} enabled-mode Bash invocations. Raw/reduced/stateful byte fields cover intercepted commands only. Unintercepted native outputs are included in the all-tool byte comparison above. A zero intercepted byte count is not evidence of complete output elimination.']
    lines += ['','## Regression evidence','']
    regressions=[p for p in pairs if p['both_successful'] and p['token_reduction_percent'] is not None and p['token_reduction_percent']<0]
    for pair in regressions:
        r=next(r for r in rows if r['id']==pair['id'] and r['repetition']==pair['repetition'] and r['mode']==pair['mode'])
        lines.append(f"- {pair['id']} {pair['mode']} repetition {pair['repetition']}: token usage increased {-pair['token_reduction_percent']:.1f}%; time overhead {pair['execution_overhead_percent']:.1f}%; {r['drilldowns']} drill-downs. [Transcript]({r['key']}.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.")
    if not regressions:lines.append('No token regressions observed among available jointly successful pairs.')
    for r in rows:
        if not r['success']:
            lines.append(f"- {r['key']}: verification/execution did not establish successful completion ({r['status']}). Provider errors: {r['workflow']['provider_errors'] or 'none reported'}. [Transcript]({r['key']}.log).")
    lines += ['','## Analysis','']
    changed=sum(r['canonical_result'].get('delta_count',0) for r in enabled)
    unchanged=sum(r['canonical_result'].get('unchanged_count',0) for r in enabled)
    retrievals=sum(r['drilldowns'] for r in enabled)
    native=sum(v for r in enabled for k,v in r['workflow']['tool_output_bytes_by_type'].items() if k!='Bash')
    all_bytes=sum(r['visible_bytes'] for r in enabled)
    lines += [f'Enabled trials recorded {changed} delta presentations, {unchanged} unchanged presentations, and {retrievals} show/raw calls. '+('There is no observed drill-down cost in these trials.' if retrievals==0 else 'Drill-down costs are included in intercepted-shell and total-tool output measurements.'), '', f'Native tools other than Bash contributed {native}/{all_bytes} textual tool-result bytes ({100*native/all_bytes if all_bytes else 0:.1f}%). They are outside the current shell integration. Interception coverage and native-output share limit the workflow benefit achievable by shell reduction.', '', 'The paired level table shows whether reductions increase with intended task complexity. The observed L5 command histories and repeated-command counts are in aggregate.json; intended complexity alone is not evidence of longer actual context history. Stateful effectiveness is limited by the number of observed unchanged/delta presentations, and cannot be inferred from the mode name. Quality is judged by the same protected tests and vet checks in both modes.']
    lines += ['','## Evidence and limitations','', 'Raw JSONL and stderr are retained in each .log file. Canonical acap bench compare reports are retained for each complete pair. Tool-output bytes exclude images and internal system prompts. Cache token accounting is explicit, but billed cost is not inferred. Repair iterations are unavailable. Exact repeated command strings, build/test/vet Bash invocations and failed validation tool results are extracted into aggregate.json; native background notifications are not classified as validation attempts. Shell commands count Bash invocations, not individual commands in compound scripts.','', 'The complexity levels describe intended engineering scope, not empirically calibrated difficulty. L5 requests meaningful layer-by-layer validation, without artificial repetition or output padding. Workflow length and repeated commands remain agent decisions. Three repetitions are the default; a smoke run supports infrastructure conclusions only.','', '## Recommendations','', '1. Use all-tool bytes and provider token usage for release decisions. Shell-only byte reduction can conceal bypassed and native tool costs measured above.', '2. Investigate the recorded shell-state bypass command sequences before broadening adapter coverage. Preserve shell semantics; do not change compression algorithms based on these results alone.', '3. Target follow-up experiments at the levels with token or time regressions in the paired table. Retain the same success criteria and use more repetitions before treating variability as a stable effect.']
    if manifest.get('blockers'):
        lines += ['','## Blockers','',*manifest['blockers']]
    (out/'report.md').write_text('\n'.join(lines)+'\n',encoding='utf-8')

def refresh_comparisons(out, manifest, binary):
    recorded=set()
    for run in manifest['runs']:
        if run['mode']=='disabled':
            continue
        fixture=run['key'].rsplit('-',2)[0]
        baseline=out/f"{fixture}-{run['repetition']}-disabled.json"
        candidate=out/(run['key']+'.json')
        result=call([str(binary.resolve()),'bench','compare','--json',str(baseline),str(candidate)],cwd=REPO)
        if result.returncode:
            raise RuntimeError(f"compare {run['key']}: {result.stderr.strip()}")
        target=out/f"{fixture}-{run['repetition']}-{run['mode']}-comparison.json"
        target.write_text(result.stdout,encoding='utf-8')
        recorded.add(target.name)
    # Replace only earlier comparison artifacts from this run that lacked a
    # mode in the filename and could therefore be overwritten.
    for stale in out.glob('*-comparison.json'):
        if stale.name not in recorded:
            if stale.parent.resolve()!=out.resolve():
                raise RuntimeError(f"unsafe comparison path {stale}")
            stale.unlink()

def main():
    p=argparse.ArgumentParser()
    p.add_argument('--out',type=Path,required=True)
    p.add_argument('--binary',type=Path,default=REPO/'.tmp/acap-bench.exe')
    p.add_argument('--smoke',action='store_true')
    p.add_argument('--report-only',action='store_true')
    p.add_argument('--refresh-comparisons',action='store_true',help='rebuild one mode-specific acap bench compare report per enabled trial')
    p.add_argument('--append',action='store_true',help='append additional comparison modes to an existing main run')
    p.add_argument('--additional-modes',nargs='+',choices=['stateless','stateful'])
    args=p.parse_args();out=args.out.resolve()
    if args.report_only or args.refresh_comparisons:
        manifest=json.loads((out/'manifest.json').read_text(encoding='utf-8'))
        if args.refresh_comparisons:
            refresh_comparisons(out,manifest,args.binary)
        if args.report_only:
            report(out,manifest)
        return
    out.mkdir(parents=True,exist_ok=True)
    if args.additional_modes and not args.append: raise SystemExit('--additional-modes requires --append')
    if not args.append and any(out.iterdir()): raise SystemExit('Output directory must be empty (prevent accidental result mixing).')
    config=json.loads((ROOT/'coding-matrix.json').read_text())
    workloads=config['workloads'][:1] if args.smoke else config['workloads']
    reps=1 if args.smoke else config['repetitions']
    if args.append:
        manifest=json.loads((out/'manifest.json').read_text(encoding='utf-8'))
        if manifest.get('run_status')!='completed': raise SystemExit('Cannot append to an incomplete experiment.')
        if manifest.get('binary_sha256') and manifest['binary_sha256']!=hashlib.sha256(args.binary.read_bytes()).hexdigest(): raise SystemExit('Benchmark executable differs from the completed baseline/Integrated experiment.')
        if args.additional_modes is None or not args.additional_modes: raise SystemExit('--append requires --additional-modes')
        added=list(args.additional_modes)
        existing=manifest.get('analysis_modes',manifest['config']['modes'])
        if len(set(added))!=len(added) or any(mode in existing for mode in added): raise SystemExit('Additional modes must be new and unique.')
        if manifest['config']['repetitions']!=reps or {w['id']:w for w in manifest['config']['workloads']}!={w['id']:w for w in config['workloads']}: raise SystemExit('Coding workload matrix changed since baseline execution.')
        manifest['analysis_modes']=existing+added
        manifest['optional_mode_experiment']=dict(modes=added,repetitions=reps,baseline_source='matching disabled trial by workload and repetition',order='alternating per workload and repetition')
        manifest['expected_conditions']+=len(workloads)*reps*len(added)
    else:
        manifest=dict(config=config,smoke=args.smoke,expected_conditions=len(workloads)*reps*len(config['modes']),run_status='running',binary_sha256=hashlib.sha256(args.binary.read_bytes()).hexdigest(),revision=call(['git','rev-parse','HEAD'],cwd=REPO).stdout.strip(),versions={tool:call([tool,'--version'] if tool!='go' else [tool,'version']).stdout.strip() for tool in ['claude','git','go']},fixture_hashes={w['id']:digest(ROOT/'fixtures'/w['fixture']) for w in workloads},definition_hashes={w['id']:hashlib.sha256((ROOT/w['workload']).read_bytes()).hexdigest() for w in workloads},runs=[],blockers=[])
    manifest['run_status']='running'
    def save():
        (out/'manifest.json').write_text(json.dumps(manifest,indent=2)+'\n',encoding='utf-8')
        report(out,manifest)
    save()
    for rep in range(1,reps+1):
        for index,w in enumerate(workloads):
            mode_set=args.additional_modes or config['modes']
            modes=mode_set if (rep+index)%2 else list(reversed(mode_set))
            for mode in modes:
                key=f"{w['fixture']}-{rep}-{mode}"
                cmd=[str(args.binary.resolve()),'bench','agent','--workload',str(ROOT/w['workload']),'--agent','claude','--mode',mode,'--model',config['model'],'--json','--verbose']
                print(f'Running {key}',flush=True);started=time.time()
                result=call(cmd,cwd=REPO)
                (out/(key+'.json')).write_text(result.stdout,encoding='utf-8')
                (out/(key+'.log')).write_text(result.stderr,encoding='utf-8')
                manifest['runs'].append(dict(key=key,id=w['id'],level=w['level'],mode=mode,repetition=rep,exit=result.returncode,elapsed=time.time()-started,argv=cmd))
                save()
                print(f'Finished {key}: exit {result.returncode}',flush=True)
                if mode!='disabled':
                    base=out/f"{w['fixture']}-{rep}-disabled.json";cand=out/(key+'.json')
                    comparison=call([str(args.binary.resolve()),'bench','compare','--json',str(base),str(cand)],cwd=REPO)
                    (out/f"{w['fixture']}-{rep}-{mode}-comparison.json").write_text(comparison.stdout,encoding='utf-8')
                    if comparison.returncode: manifest['blockers'].append(comparison.stderr.strip())
                    save()
    manifest['run_status']='completed'
    save()

if __name__=='__main__': main()
