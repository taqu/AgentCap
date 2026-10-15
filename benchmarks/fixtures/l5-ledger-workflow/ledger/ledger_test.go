package ledger
import("testing";"bench.local/work/model")
func TestApply(t *testing.T){l:=New();for _,e:=range []model.Event{{ID:"a",Kind:"credit",Amount:9},{ID:"b",Kind:"debit",Amount:4},{ID:"a",Kind:"credit",Amount:9}}{if err:=l.Apply(e);err!=nil{t.Fatal(err)}};if l.Total()!=5{t.Fatalf("total %d",l.Total())};if e:=l.Apply(model.Event{ID:"z",Kind:"other",Amount:1});e==nil{t.Fatal("missing error")};if l.Total()!=5{t.Fatal("error changed total")}}
