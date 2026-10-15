package storage
import("testing";"sync")
func TestCAS(t *testing.T){s:=New();r,e:=s.Put("a","first",0);if e!=nil||r.Revision!=1{t.Fatal(r,e)};if _,e=s.Put("a","stale",0);e==nil{t.Fatal("stale accepted")};r,_=s.Get("a");if r.Value!="first"{t.Fatal("conflict mutated")};if _,e=s.Put("","x",0);e==nil{t.Fatal("empty key")};if _,e=s.Put("b","x",-1);e==nil{t.Fatal("negative revision")};snap:=s.Snapshot();delete(snap,"a");if _,ok:=s.Get("a");!ok{t.Fatal("snapshot aliases")}}
func TestConcurrent(t *testing.T){s:=New();var wg sync.WaitGroup;for i:=0;i<20;i++{wg.Add(1);go func(){defer wg.Done();s.Put("key","v",0);s.Get("key");s.Snapshot()}()};wg.Wait();r,_:=s.Get("key");if r.Revision!=1{t.Fatalf("revision %d",r.Revision)}}
