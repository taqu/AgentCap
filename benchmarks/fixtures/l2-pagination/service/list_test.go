package service
import("testing";"reflect")
func TestList(t *testing.T){src:=[]string{"a","b","c"};v,e:=List(src,1,1);if e!=nil||!reflect.DeepEqual(v,[]string{"b"}){t.Fatalf("page %v %v",v,e)};v[0]="z";if src[1]!="b"{t.Fatal("aliased")};for _,c:=range [][2]int{{3,1},{5,2},{0,0}}{v,e:=List(src,c[0],c[1]);if e!=nil||len(v)!=0{t.Errorf("empty %v %v",v,e)}};for _,c:=range [][2]int{{-1,1},{0,-1}}{if _,e:=List(src,c[0],c[1]);e==nil{t.Fatal("missing error")}}}
