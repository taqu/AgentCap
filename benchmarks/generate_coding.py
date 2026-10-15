"""Regenerate frozen v1 standard-library Go coding fixtures (not run by agents)."""
import json
from pathlib import Path

ROOT = Path(__file__).resolve().parent
catalog = []

def add(name, level, task, files, characteristics):
    fixture = ROOT / 'fixtures' / name
    fixture.mkdir(parents=True, exist_ok=True)
    files = {'go.mod': 'module bench.local/work\n\ngo 1.22\n', **files}
    for path, content in files.items():
        target = fixture / path
        target.parent.mkdir(parents=True, exist_ok=True)
        target.write_text(content.strip() + '\n', encoding='utf-8')
    path = ROOT / 'workloads' / 'coding' / (name + '.yaml')
    path.parent.mkdir(parents=True, exist_ok=True)
    task += '\nPreserve existing tests and go.mod. Use only the Go standard library. Complete and validate the task in this workspace.'
    path.write_text(f'version: 1\nname: coding/{name}\nfixture: ../../fixtures/{name}\ngit: {{init: true}}\ntimeout: {5 + level * 3}m\ntask: |\n' + ''.join('  '+line+'\n' for line in task.splitlines()) + '''verify:
  - run:
      argv: ["git", "diff", "--exit-code", "HEAD", "--", "*_test.go", "go.mod"]
      expect: {exit: 0}
  - run:
      argv: ["go", "test", "-count=1", "./..."]
      expect: {exit: 0}
  - run:
      argv: ["go", "vet", "./..."]
      expect: {exit: 0}
''', encoding='utf-8')
    catalog.append(dict(id='coding/'+name, level=level, workload='workloads/coding/'+name+'.yaml', fixture=name, revision='v1', expected_behavior=task.splitlines()[0], characteristics=characteristics, timeout_minutes=5+level*3))

add('l1-clamp',1,'Repair Clamp so it clamps integers to the inclusive [low, high] interval without changing its signature.', {'clamp.go': 'package work\nfunc Clamp(v, low, high int) int { if v < low { return high }; if v > high { return low }; return v }',
'clamp_test.go': '''package work
import "testing"
func TestClamp(t *testing.T) { for _, c := range [][4]int{{-2,0,10,0},{12,0,10,10},{5,0,10,5},{0,0,10,0},{10,0,10,10},{-9,-5,-1,-5}} { if got:=Clamp(c[0],c[1],c[2]); got!=c[3] {t.Errorf("Clamp%v=%d",c,got)} } }
'''}, 'Single implementation file, boundary conditional repair.')
add('l1-duration',1,'Repair ParseSeconds: accept nonnegative integer seconds, reject negative or malformed input, and return the corresponding time.Duration.', {
'duration.go': '''package work
import ("strconv"; "time"; "fmt")
func ParseSeconds(s string) (time.Duration,error) { n,e:=strconv.Atoi(s); if e!=nil || n<0 {return 0,fmt.Errorf("invalid seconds")}; return time.Duration(n)*time.Millisecond,nil }
''',
'duration_test.go': '''package work
import ("testing";"time")
func TestSeconds(t *testing.T){ for _,n:=range []string{"0","1","60"}{ d,e:=ParseSeconds(n); if e!=nil {t.Fatal(e)}; if n=="60" && d!=time.Minute {t.Errorf("60 = %s",d)}; if n=="1" && d!=time.Second {t.Errorf("1 = %s",d)} }; for _,s:=range []string{"-1","bad","1.5"}{if _,e:=ParseSeconds(s);e==nil{t.Errorf("accepted %q",s)}} }
'''}, 'Single implementation file, unit conversion repair.')
add('l2-config',2,'Fix configuration loading and client construction. An absent PORT defaults to 8080; an explicit port must be 1..65535. Address must include the supplied hostname and selected port.', {
'config/config.go': '''package config
import "strconv"
type Config struct{ Host string; Port int }
func Load(env map[string]string)(Config,error){p,e:=strconv.Atoi(env["PORT"]); return Config{Host:env["HOST"],Port:p},e}
''',
'client/client.go': '''package client
import("bench.local/work/config";"strconv")
func Address(c config.Config)string{return c.Host+strconv.Itoa(c.Port)}
''',
'client/client_test.go': '''package client
import("testing";"bench.local/work/config")
func TestConfig(t *testing.T){for _,s:=range []string{"","9000"}{c,e:=config.Load(map[string]string{"HOST":"localhost","PORT":s});if e!=nil{t.Fatal(e)};want:="localhost:9000";if s==""{want="localhost:8080"};if Address(c)!=want{t.Errorf("got %s want %s",Address(c),want)}};for _,s:=range []string{"0","-1","65536","oops"}{if _,e:=config.Load(map[string]string{"PORT":s});e==nil{t.Errorf("accepted %s",s)}}}
'''}, 'Two packages, parsing defaults and address construction.')
add('l2-pagination',2,'Implement zero-based pagination. Page(items, offset, limit) must return at most limit items, reject negative arguments, return an empty slice beyond the end, and avoid aliasing the input. List must propagate pagination errors.', {
'paging/page.go': '''package paging
func Page(items []string,offset,limit int)([]string,error){return items,nil}
''',
'service/list.go': '''package service
import "bench.local/work/paging"
func List(items []string,offset,limit int)([]string,error){v,_:=paging.Page(items,offset,limit);return v,nil}
''',
'service/list_test.go': '''package service
import("testing";"reflect")
func TestList(t *testing.T){src:=[]string{"a","b","c"};v,e:=List(src,1,1);if e!=nil||!reflect.DeepEqual(v,[]string{"b"}){t.Fatalf("page %v %v",v,e)};v[0]="z";if src[1]!="b"{t.Fatal("aliased")};for _,c:=range [][2]int{{3,1},{5,2},{0,0}}{v,e:=List(src,c[0],c[1]);if e!=nil||len(v)!=0{t.Errorf("empty %v %v",v,e)}};for _,c:=range [][2]int{{-1,1},{0,-1}}{if _,e:=List(src,c[0],c[1]);e==nil{t.Fatal("missing error")}}}
'''}, 'Two packages, new behavior plus error propagation.')

# A real data path, with independently tested layers and integration behavior.
pipeline = {
'model/event.go': '''package model
type Event struct { ID string `json:"id"`; Kind string `json:"kind"`; Amount int `json:"amount"` }
''',
'decode/decode.go': '''package decode
import("encoding/json";"bench.local/work/model")
func Parse(line string)(model.Event,error){var e model.Event;err:=json.Unmarshal([]byte(line),&e);return e,err}
''',
'ledger/ledger.go': '''package ledger
import "bench.local/work/model"
type Ledger struct { total int; seen map[string]bool }
func New()*Ledger{return &Ledger{seen:map[string]bool{}}}
func(l *Ledger)Apply(e model.Event)error{l.total+=e.Amount;return nil}
func(l *Ledger)Total()int{return l.total}
''',
'report/report.go': '''package report
import "strconv"
func Format(total int)string{return strconv.Itoa(total)}
''',
'service/service.go': '''package service
import("bench.local/work/decode";"bench.local/work/ledger";"bench.local/work/report")
func Run(lines []string)(string,error){l:=ledger.New();for _,s:=range lines{e,err:=decode.Parse(s);if err!=nil{return "",err};if err=l.Apply(e);err!=nil{return "",err}};return report.Format(l.Total()),nil}
''',
'service/service_test.go': '''package service
import "testing"
func TestRun(t *testing.T){lines:=[]string{`{"id":"a","kind":"credit","amount":9}`,`{"id":"b","kind":"debit","amount":4}`,`{"id":"a","kind":"credit","amount":9}`};s,e:=Run(lines);if e!=nil||s!="total=5\\n"{t.Fatalf("result %q %v",s,e)};for _,bad:=range []string{`broken`,`{"id":"","kind":"credit","amount":1}`,`{"id":"x","kind":"other","amount":1}`,`{"id":"x","kind":"credit","amount":-1}`}{if _,e:=Run([]string{bad});e==nil{t.Errorf("accepted %s",bad)}}}
''',
'ledger/ledger_test.go': '''package ledger
import("testing";"bench.local/work/model")
func TestApply(t *testing.T){l:=New();for _,e:=range []model.Event{{ID:"a",Kind:"credit",Amount:9},{ID:"b",Kind:"debit",Amount:4},{ID:"a",Kind:"credit",Amount:9}}{if err:=l.Apply(e);err!=nil{t.Fatal(err)}};if l.Total()!=5{t.Fatalf("total %d",l.Total())};if e:=l.Apply(model.Event{ID:"z",Kind:"other",Amount:1});e==nil{t.Fatal("missing error")};if l.Total()!=5{t.Fatal("error changed total")}}
'''}
add('l3-ledger',3,'Repair the event processing pipeline. Credits add and debits subtract nonnegative amounts; duplicate event IDs are ignored. Reject empty IDs and unknown kinds without changing ledger state. Run returns total=<integer> followed by a newline. Preserve package APIs.',pipeline,'Five interacting packages; decoding, idempotency, transaction state and reporting.')

graph = {
'graph/graph.go': '''package graph
type Node struct { Name string; Deps []string }
type Graph map[string]Node
func(g Graph)Order()([]string,error){var names []string;for n:=range g{names=append(names,n)};return names,nil}
''',
'parse/parse.go': '''package parse
import("encoding/json";"bench.local/work/graph")
func Load(data []byte)(graph.Graph,error){var nodes []graph.Node;if e:=json.Unmarshal(data,&nodes);e!=nil{return nil,e};g:=graph.Graph{};for _,n:=range nodes{g[n.Name]=n};return g,nil}
''',
'plan/plan.go': '''package plan
import("bench.local/work/parse";"strings")
func Build(data []byte)(string,error){g,e:=parse.Load(data);if e!=nil{return "",e};order,e:=g.Order();return strings.Join(order,","),e}
''',
'plan/plan_test.go': '''package plan
import "testing"
func TestBuild(t *testing.T){s,e:=Build([]byte(`[{"Name":"app","Deps":["lib"]},{"Name":"lib","Deps":["base"]},{"Name":"base"}]`));if e!=nil||s!="base,lib,app"{t.Fatalf("plan %q %v",s,e)};for _,s:=range []string{`[{"Name":"x","Deps":["missing"]}]`,`[{"Name":"x","Deps":["y"]},{"Name":"y","Deps":["x"]}]`,`[{"Name":"x"},{"Name":"x"}]`,`[{"Name":""}]`}{if _,e:=Build([]byte(s));e==nil{t.Errorf("accepted %s",s)}}}
''',
'graph/graph_test.go': '''package graph
import("testing";"reflect")
func TestOrder(t *testing.T){g:=Graph{"b":Node{Name:"b"},"a":Node{Name:"a"}};for i:=0;i<20;i++{v,e:=g.Order();if e!=nil||!reflect.DeepEqual(v,[]string{"a","b"}){t.Fatalf("unstable %v %v",v,e)}};if _,e:=(Graph{"x":Node{Name:"x",Deps:[]string{"x"}}}).Order();e==nil{t.Fatal("self-cycle accepted")}}
'''}
add('l3-planner',3,'Implement deterministic dependency planning. Dependencies precede dependents, unrelated ready nodes use lexical order, and missing dependencies and cycles are errors. Reject duplicate and empty node names when parsing. Preserve public APIs.',graph,'Graph algorithm, input validation and planner integration across three packages.')

store = {
'record/record.go': '''package record
type Record struct {Key string; Value string; Revision int}
''',
'storage/storage.go': '''package storage
import "bench.local/work/record"
type Store struct { data map[string]record.Record }
func New()*Store{return &Store{data:map[string]record.Record{}}}
func(s *Store)Put(key,value string,expected int)(record.Record,error){r:=record.Record{Key:key,Value:value,Revision:expected+1};s.data[key]=r;return r,nil}
func(s *Store)Get(key string)(record.Record,bool){r,ok:=s.data[key];return r,ok}
func(s *Store)Snapshot()map[string]record.Record{return s.data}
''',
'api/api.go': '''package api
import("bench.local/work/storage";"bench.local/work/record")
type API struct{ Store *storage.Store }
func(a API)Update(k,v string,rev int)(record.Record,error){return a.Store.Put(k,v,rev)}
''',
'export/export.go': '''package export
import("bench.local/work/storage";"encoding/json")
func JSON(s *storage.Store)([]byte,error){return json.Marshal(s.Snapshot())}
''',
'replica/replica.go': '''package replica
import("bench.local/work/storage";"bench.local/work/record")
func Copy(s *storage.Store)map[string]record.Record{return s.Snapshot()}
''',
'storage/storage_test.go': '''package storage
import("testing";"sync")
func TestCAS(t *testing.T){s:=New();r,e:=s.Put("a","first",0);if e!=nil||r.Revision!=1{t.Fatal(r,e)};if _,e=s.Put("a","stale",0);e==nil{t.Fatal("stale accepted")};r,_=s.Get("a");if r.Value!="first"{t.Fatal("conflict mutated")};if _,e=s.Put("","x",0);e==nil{t.Fatal("empty key")};if _,e=s.Put("b","x",-1);e==nil{t.Fatal("negative revision")};snap:=s.Snapshot();delete(snap,"a");if _,ok:=s.Get("a");!ok{t.Fatal("snapshot aliases")}}
func TestConcurrent(t *testing.T){s:=New();var wg sync.WaitGroup;for i:=0;i<20;i++{wg.Add(1);go func(){defer wg.Done();s.Put("key","v",0);s.Get("key");s.Snapshot()}()};wg.Wait();r,_:=s.Get("key");if r.Revision!=1{t.Fatalf("revision %d",r.Revision)}}
''',
'api/api_test.go': '''package api
import("testing";"bench.local/work/storage";"bench.local/work/export";"bench.local/work/replica")
func TestClients(t *testing.T){s:=storage.New();a:=API{s};if _,e:=a.Update("x","value",0);e!=nil{t.Fatal(e)};if _,e:=a.Update("x","bad",0);e==nil{t.Fatal("conflict")};b,e:=export.JSON(s);if e!=nil||len(b)==0{t.Fatal(e)};r:=replica.Copy(s);delete(r,"x");if _,ok:=s.Get("x");!ok{t.Fatal("replica aliases store")}}
'''}
add('l4-store',4,'Repair the revisioned storage subsystem and all its callers. Expected revision 0 creates a missing key; updates require the current revision and increment it. Reject stale revisions, negative revisions and empty keys atomically. Snapshots and replica copies must not alias storage. Concurrent access must be safe. Preserve public APIs and export behavior.',store,'Five packages; optimistic concurrency, atomic failure, ownership and parallel access.')

migrated=dict(store)
migrated['storage/storage.go']=migrated['storage/storage.go'].replace('"bench.local/work/record"','("bench.local/work/record";"context")').replace('Put(key,value string,expected int)','Put(ctx context.Context,key,value string,expected int)').replace('Get(key string)','Get(ctx context.Context,key string)').replace('Snapshot()','Snapshot(ctx context.Context)')
for p in ['api/api_test.go','storage/storage_test.go']:
    migrated[p]=migrated[p].replace('import(', 'import("context";').replace('.Put(','.Put(context.Background(),').replace('.Get(','.Get(context.Background(),').replace('.Snapshot()', '.Snapshot(context.Background())')
migrated['storage/cancel_test.go']='''package storage
import("testing";"context")
func TestCancel(t *testing.T){s:=New();ctx,c:=context.WithCancel(context.Background());c();if _,e:=s.Put(ctx,"x","v",0);e==nil{t.Fatal("canceled write")};if _,ok:=s.Get(context.Background(),"x");ok{t.Fatal("canceled mutation")}}
'''
add('l4-context-migration',4,'Complete the storage API migration to context.Context across all packages. Put must reject a canceled context without mutation. Also restore revision conflict handling, snapshot independence and concurrency safety. Keep the new storage signatures and existing external API, export and replica signatures.',migrated,'Five packages with compilation failures; API migration plus regression repair.')

extended=dict(pipeline)
extended['batch/batch.go']='''package batch
import("bench.local/work/service";"strings")
func Run(input string)(string,error){return service.Run(strings.Split(input,"\\n"))}
'''
extended['batch/batch_test.go']='''package batch
import "testing"
func TestBatch(t *testing.T){v,e:=Run("\\n{\\"id\\":\\"a\\",\\"kind\\":\\"credit\\",\\"amount\\":9}\\n\\n{\\"id\\":\\"b\\",\\"kind\\":\\"debit\\",\\"amount\\":4}\\n");if e!=nil||v!="total=5\\n"{t.Fatal(v,e)};v,e=Run("");if e!=nil||v!="total=0\\n"{t.Fatal(v,e)}}
'''
extended['cmd/ledger/main.go']='''package main
import("os";"fmt";"io";"bench.local/work/batch")
func main(){data,e:=io.ReadAll(os.Stdin);if e!=nil{panic(e)};v,e:=batch.Run(string(data));if e!=nil{fmt.Fprintln(os.Stderr,e);return};fmt.Print(v)}
'''
extended['integration/cli_test.go']='''package integration
import("testing";"os/exec";"strings")
func TestCLI(t *testing.T){c:=exec.Command("go","run","../cmd/ledger");c.Stdin=strings.NewReader(`{"id":"a","kind":"credit","amount":7}`);b,e:=c.Output();if e!=nil||string(b)!="total=7\\n"{t.Fatal(string(b),e)};c=exec.Command("go","run","../cmd/ledger");c.Stdin=strings.NewReader("bad");if e:=c.Run();e==nil{t.Fatal("invalid input exited zero")}}
'''
add('l5-ledger-workflow',5,'Deliver the ledger batch CLI in stages: first restore ledger unit behavior (credits, debits, ID idempotency, invalid-event rejection); then integrate newline-separated JSON batches, ignoring blank lines; finally make the CLI print total=<integer> with a newline on success and exit nonzero on invalid input. Validate the unit layer before the batch layer, then the subprocess integration. Preserve existing APIs and tests.',extended,'Eight packages; unit, batch and subprocess validation stages with interacting failures.')

extended_graph=dict(graph)
extended_graph['cache/cache.go']='''package cache
import("bench.local/work/plan")
type Cache struct {last string}
func(c *Cache)Build(data []byte)(string,error){if c.last!=""{return c.last,nil};v,e:=plan.Build(data);c.last=v;return v,e}
'''
extended_graph['cache/cache_test.go']='''package cache
import "testing"
func TestInvalidation(t *testing.T){c:=&Cache{};a,e:=c.Build([]byte(`[{"Name":"a"}]`));if e!=nil||a!="a"{t.Fatal(a,e)};b,e:=c.Build([]byte(`[{"Name":"b"}]`));if e!=nil||b!="b"{t.Fatal(b,e)};if _,e=c.Build([]byte(`broken`));e==nil{t.Fatal("cached error hidden")};a,e=c.Build([]byte(`[{"Name":"a"}]`));if e!=nil||a!="a"{t.Fatal(a,e)}}
'''
extended_graph['batch/batch.go']='''package batch
import("bench.local/work/cache")
func Build(inputs [][]byte)([]string,error){c:=&cache.Cache{};out:=[]string{};for _,in:=range inputs{v,e:=c.Build(in);if e!=nil{continue};out=append(out,v)};return out,nil}
'''
extended_graph['batch/batch_test.go']='''package batch
import("testing";"reflect")
func TestBatch(t *testing.T){v,e:=Build([][]byte{[]byte(`[{"Name":"a"}]`),[]byte(`[{"Name":"b"}]`),[]byte(`[{"Name":"a"}]`)});if e!=nil||!reflect.DeepEqual(v,[]string{"a","b","a"}){t.Fatal(v,e)};v,e=Build([][]byte{[]byte(`[{"Name":"a"}]`),[]byte(`broken`)});if e==nil||v!=nil{t.Fatal("partial success",v,e)}}
'''
extended_graph['cmd/planner/main.go']='''package main
import("io";"os";"fmt";"bench.local/work/plan")
func main(){b,_:=io.ReadAll(os.Stdin);s,e:=plan.Build(b);if e!=nil{fmt.Fprintln(os.Stderr,e);return};fmt.Print(s)}
'''
extended_graph['integration/cli_test.go']='''package integration
import("testing";"os/exec";"strings")
func TestCLI(t *testing.T){c:=exec.Command("go","run","../cmd/planner");c.Stdin=strings.NewReader(`[{"Name":"b","Deps":["a"]},{"Name":"a"}]`);b,e:=c.Output();if e!=nil||string(b)!="a,b\\n"{t.Fatal(string(b),e)};c=exec.Command("go","run","../cmd/planner");c.Stdin=strings.NewReader(`[{"Name":"a","Deps":["a"]}]`);if e:=c.Run();e==nil{t.Fatal("cycle exited zero")}}
'''
add('l5-planner-workflow',5,'Repair the dependency planner through staged validation: establish deterministic topological order and parsing rejection of duplicates, empty names, missing dependencies and cycles; fix cache invalidation for changed inputs without masking errors; make batch processing fail atomically (nil output on any error); finish the CLI with newline-terminated plans and nonzero error exits. Validate graph/parse/plan first, cache/batch next, subprocess integration last. Preserve public APIs and tests.',extended_graph,'Eight packages; algorithm, cache history, atomic batch and subprocess integration stages.')

(ROOT/'coding-matrix.json').write_text(json.dumps(dict(version=1, model='claude-sonnet-4-6', repetitions=3, modes=['disabled','integrated'], workloads=catalog),indent=2)+'\n',encoding='utf-8')


