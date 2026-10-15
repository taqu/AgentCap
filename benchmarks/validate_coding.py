"""Prove each frozen fixture fails initially and passes with a reference repair.

References never enter the copied agent workspace. No network dependencies.
"""
import json
import shutil
import subprocess
import tempfile
from pathlib import Path

ROOT=Path(__file__).resolve().parent
repairs={
'clamp.go': 'package work\nfunc Clamp(v,lo,hi int)int{if v<lo{return lo};if v>hi{return hi};return v}',
'duration.go': '''package work
import("strconv";"time";"fmt")
func ParseSeconds(s string)(time.Duration,error){n,e:=strconv.Atoi(s);if e!=nil||n<0{return 0,fmt.Errorf("invalid")};return time.Duration(n)*time.Second,nil}
''',
'config/config.go': '''package config
import("strconv";"fmt")
type Config struct{Host string;Port int}
func Load(env map[string]string)(Config,error){p:=8080;var e error;if env["PORT"]!=""{p,e=strconv.Atoi(env["PORT"])};if e!=nil||p<1||p>65535{return Config{},fmt.Errorf("port")};return Config{env["HOST"],p},nil}
''',
'client/client.go': '''package client
import("bench.local/work/config";"net";"strconv")
func Address(c config.Config)string{return net.JoinHostPort(c.Host,strconv.Itoa(c.Port))}
''',
'paging/page.go': '''package paging
import "fmt"
func Page(items []string,offset,limit int)([]string,error){if offset<0||limit<0{return nil,fmt.Errorf("negative")};if offset>=len(items)||limit==0{return []string{},nil};end:=len(items);if limit<end-offset{end=offset+limit};return append([]string{},items[offset:end]...),nil}
''',
'service/list.go': '''package service
import "bench.local/work/paging"
func List(items []string,offset,limit int)([]string,error){return paging.Page(items,offset,limit)}
''',
'ledger/ledger.go': '''package ledger
import("bench.local/work/model";"fmt")
type Ledger struct{total int;seen map[string]bool}
func New()*Ledger{return &Ledger{seen:map[string]bool{}}}
func(l *Ledger)Apply(e model.Event)error{if e.ID==""||e.Amount<0||(e.Kind!="credit"&&e.Kind!="debit"){return fmt.Errorf("invalid")};if l.seen[e.ID]{return nil};l.seen[e.ID]=true;if e.Kind=="credit"{l.total+=e.Amount}else{l.total-=e.Amount};return nil}
func(l *Ledger)Total()int{return l.total}
''',
'report/report.go': '''package report
import "fmt"
func Format(total int)string{return fmt.Sprintf("total=%d\\n",total)}
''',
'graph/graph.go': '''package graph
import("fmt";"sort")
type Node struct{Name string;Deps []string}
type Graph map[string]Node
func(g Graph)Order()([]string,error){done:=map[string]bool{};out:=[]string{};for len(done)<len(g){ready:=[]string{};for name,node:=range g{if done[name]{continue};ok:=true;for _,dep:=range node.Deps{if _,exists:=g[dep];!exists{return nil,fmt.Errorf("missing %s",dep)};if !done[dep]{ok=false}};if ok{ready=append(ready,name)}};if len(ready)==0{return nil,fmt.Errorf("cycle")};sort.Strings(ready);name:=ready[0];done[name]=true;out=append(out,name)};return out,nil}
''',
'parse/parse.go': '''package parse
import("encoding/json";"fmt";"bench.local/work/graph")
func Load(data []byte)(graph.Graph,error){var nodes []graph.Node;if e:=json.Unmarshal(data,&nodes);e!=nil{return nil,e};g:=graph.Graph{};for _,n:=range nodes{if n.Name==""{return nil,fmt.Errorf("empty")};if _,ok:=g[n.Name];ok{return nil,fmt.Errorf("duplicate")};g[n.Name]=n};return g,nil}
''',
'cache/cache.go': '''package cache
import "bench.local/work/plan"
type Cache struct{input string;last string;valid bool}
func(c *Cache)Build(data []byte)(string,error){if c.valid&&c.input==string(data){return c.last,nil};v,e:=plan.Build(data);if e!=nil{return "",e};c.input=string(data);c.last=v;c.valid=true;return v,nil}
''',
}

storage='''package storage
import("bench.local/work/record";"sync";"fmt")
type Store struct{mu sync.RWMutex;data map[string]record.Record}
func New()*Store{return &Store{data:map[string]record.Record{}}}
func(s *Store)Put(key,value string,expected int)(record.Record,error){s.mu.Lock();defer s.mu.Unlock();if key==""||expected<0{return record.Record{},fmt.Errorf("invalid")};old:=s.data[key];if old.Revision!=expected{return record.Record{},fmt.Errorf("conflict")};r:=record.Record{Key:key,Value:value,Revision:expected+1};s.data[key]=r;return r,nil}
func(s *Store)Get(key string)(record.Record,bool){s.mu.RLock();defer s.mu.RUnlock();r,ok:=s.data[key];return r,ok}
func(s *Store)Snapshot()map[string]record.Record{s.mu.RLock();defer s.mu.RUnlock();out:=map[string]record.Record{};for k,v:=range s.data{out[k]=v};return out}
'''

def run(args,cwd):
    r=subprocess.run(args,cwd=cwd,capture_output=True,text=True,timeout=120)
    return {'argv':args,'exit':r.returncode,'output':r.stdout+r.stderr}

def main():
    config=json.loads((ROOT/'coding-matrix.json').read_text())
    evidence=[]
    for w in config['workloads']:
        name=w['fixture']
        with tempfile.TemporaryDirectory(prefix='acap-fixture-check-') as temp:
            path=Path(temp)/'work';shutil.copytree(ROOT/'fixtures'/name,path)
            before=run(['go','test','-count=1','./...'],path)
            assert before['exit']!=0, name+' must fail initially'
            assert 'syntax error' not in before['output'], name+' has an unintended syntax error'
            for f in list(path.rglob('*.go')):
                rel=f.relative_to(path).as_posix()
                if rel in repairs:f.write_text(repairs[rel],encoding='utf-8')
            if name.startswith('l4-'):
                code=storage
                if 'context-migration' in name:
                    code=code.replace(';"sync"',';"context";"sync"').replace('Put(key,value string,expected int)','Put(ctx context.Context,key,value string,expected int)').replace('Get(key string)','Get(ctx context.Context,key string)').replace('Snapshot()','Snapshot(ctx context.Context)').replace('s.mu.Lock();defer s.mu.Unlock();','s.mu.Lock();defer s.mu.Unlock();if e:=ctx.Err();e!=nil{return record.Record{},e};')
                    for f in ['api/api.go','export/export.go','replica/replica.go']:
                        target=path/f;target.write_text(target.read_text().replace('import(', 'import("context";').replace('.Put(','.Put(context.Background(),').replace('.Snapshot()', '.Snapshot(context.Background())'))
                (path/'storage/storage.go').write_text(code,encoding='utf-8')
            if name=='l5-ledger-workflow':
                (path/'batch/batch.go').write_text('''package batch
import("bench.local/work/service";"strings")
func Run(input string)(string,error){lines:=[]string{};for _,s:=range strings.Split(input,"\\n"){if strings.TrimSpace(s)!=""{lines=append(lines,s)}};return service.Run(lines)}
''')
                target=path/'cmd/ledger/main.go';target.write_text(target.read_text().replace('fmt.Fprintln(os.Stderr,e);return','fmt.Fprintln(os.Stderr,e);os.Exit(1)'))
            if name=='l5-planner-workflow':
                target=path/'batch/batch.go';target.write_text(target.read_text().replace('if e!=nil{continue}', 'if e!=nil{return nil,e}'))
                target=path/'cmd/planner/main.go';target.write_text(target.read_text().replace('fmt.Fprintln(os.Stderr,e);return','fmt.Fprintln(os.Stderr,e);os.Exit(1)').replace('fmt.Print(s)','fmt.Println(s)'))
            after=run(['go','test','-count=1','./...'],path)
            vet=run(['go','vet','./...'],path)
            assert after['exit']==0,(name,after['output'])
            assert vet['exit']==0,(name,vet['output'])
            evidence.append(dict(id=w['id'],initial=before,repaired=after,vet=vet))
            print(name+': fails initially, reference passes',flush=True)
    (ROOT/'fixture-validation.json').write_text(json.dumps(evidence,indent=2)+'\n')

if __name__=='__main__':main()
