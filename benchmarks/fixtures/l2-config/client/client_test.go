package client
import("testing";"bench.local/work/config")
func TestConfig(t *testing.T){for _,s:=range []string{"","9000"}{c,e:=config.Load(map[string]string{"HOST":"localhost","PORT":s});if e!=nil{t.Fatal(e)};want:="localhost:9000";if s==""{want="localhost:8080"};if Address(c)!=want{t.Errorf("got %s want %s",Address(c),want)}};for _,s:=range []string{"0","-1","65536","oops"}{if _,e:=config.Load(map[string]string{"PORT":s});e==nil{t.Errorf("accepted %s",s)}}}
