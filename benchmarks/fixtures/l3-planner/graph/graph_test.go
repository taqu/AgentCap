package graph
import("testing";"reflect")
func TestOrder(t *testing.T){g:=Graph{"b":Node{Name:"b"},"a":Node{Name:"a"}};for i:=0;i<20;i++{v,e:=g.Order();if e!=nil||!reflect.DeepEqual(v,[]string{"a","b"}){t.Fatalf("unstable %v %v",v,e)}};if _,e:=(Graph{"x":Node{Name:"x",Deps:[]string{"x"}}}).Order();e==nil{t.Fatal("self-cycle accepted")}}
