package cache
import "testing"
func TestInvalidation(t *testing.T){c:=&Cache{};a,e:=c.Build([]byte(`[{"Name":"a"}]`));if e!=nil||a!="a"{t.Fatal(a,e)};b,e:=c.Build([]byte(`[{"Name":"b"}]`));if e!=nil||b!="b"{t.Fatal(b,e)};if _,e=c.Build([]byte(`broken`));e==nil{t.Fatal("cached error hidden")};a,e=c.Build([]byte(`[{"Name":"a"}]`));if e!=nil||a!="a"{t.Fatal(a,e)}}
