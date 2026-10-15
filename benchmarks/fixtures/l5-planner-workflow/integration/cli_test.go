package integration
import("testing";"os/exec";"strings")
func TestCLI(t *testing.T){c:=exec.Command("go","run","../cmd/planner");c.Stdin=strings.NewReader(`[{"Name":"b","Deps":["a"]},{"Name":"a"}]`);b,e:=c.Output();if e!=nil||string(b)!="a,b\n"{t.Fatal(string(b),e)};c=exec.Command("go","run","../cmd/planner");c.Stdin=strings.NewReader(`[{"Name":"a","Deps":["a"]}]`);if e:=c.Run();e==nil{t.Fatal("cycle exited zero")}}
