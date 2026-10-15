package service
import "testing"
func TestRun(t *testing.T){lines:=[]string{`{"id":"a","kind":"credit","amount":9}`,`{"id":"b","kind":"debit","amount":4}`,`{"id":"a","kind":"credit","amount":9}`};s,e:=Run(lines);if e!=nil||s!="total=5\n"{t.Fatalf("result %q %v",s,e)};for _,bad:=range []string{`broken`,`{"id":"","kind":"credit","amount":1}`,`{"id":"x","kind":"other","amount":1}`,`{"id":"x","kind":"credit","amount":-1}`}{if _,e:=Run([]string{bad});e==nil{t.Errorf("accepted %s",bad)}}}
