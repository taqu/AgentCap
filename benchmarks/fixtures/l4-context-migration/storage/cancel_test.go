package storage
import("testing";"context")
func TestCancel(t *testing.T){s:=New();ctx,c:=context.WithCancel(context.Background());c();if _,e:=s.Put(ctx,"x","v",0);e==nil{t.Fatal("canceled write")};if _,ok:=s.Get(context.Background(),"x");ok{t.Fatal("canceled mutation")}}
