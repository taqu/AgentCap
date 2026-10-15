package plan
import "testing"
func TestBuild(t *testing.T){s,e:=Build([]byte(`[{"Name":"app","Deps":["lib"]},{"Name":"lib","Deps":["base"]},{"Name":"base"}]`));if e!=nil||s!="base,lib,app"{t.Fatalf("plan %q %v",s,e)};for _,s:=range []string{`[{"Name":"x","Deps":["missing"]}]`,`[{"Name":"x","Deps":["y"]},{"Name":"y","Deps":["x"]}]`,`[{"Name":"x"},{"Name":"x"}]`,`[{"Name":""}]`}{if _,e:=Build([]byte(s));e==nil{t.Errorf("accepted %s",s)}}}
