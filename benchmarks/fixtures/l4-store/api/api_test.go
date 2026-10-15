package api
import("testing";"bench.local/work/storage";"bench.local/work/export";"bench.local/work/replica")
func TestClients(t *testing.T){s:=storage.New();a:=API{s};if _,e:=a.Update("x","value",0);e!=nil{t.Fatal(e)};if _,e:=a.Update("x","bad",0);e==nil{t.Fatal("conflict")};b,e:=export.JSON(s);if e!=nil||len(b)==0{t.Fatal(e)};r:=replica.Copy(s);delete(r,"x");if _,ok:=s.Get("x");!ok{t.Fatal("replica aliases store")}}
