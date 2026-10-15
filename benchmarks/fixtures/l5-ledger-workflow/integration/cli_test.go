package integration
import("testing";"os/exec";"strings")
func TestCLI(t *testing.T){c:=exec.Command("go","run","../cmd/ledger");c.Stdin=strings.NewReader(`{"id":"a","kind":"credit","amount":7}`);b,e:=c.Output();if e!=nil||string(b)!="total=7\n"{t.Fatal(string(b),e)};c=exec.Command("go","run","../cmd/ledger");c.Stdin=strings.NewReader("bad");if e:=c.Run();e==nil{t.Fatal("invalid input exited zero")}}
