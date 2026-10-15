package batch
import("testing";"reflect")
func TestBatch(t *testing.T){v,e:=Build([][]byte{[]byte(`[{"Name":"a"}]`),[]byte(`[{"Name":"b"}]`),[]byte(`[{"Name":"a"}]`)});if e!=nil||!reflect.DeepEqual(v,[]string{"a","b","a"}){t.Fatal(v,e)};v,e=Build([][]byte{[]byte(`[{"Name":"a"}]`),[]byte(`broken`)});if e==nil||v!=nil{t.Fatal("partial success",v,e)}}
