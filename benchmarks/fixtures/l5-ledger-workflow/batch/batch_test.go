package batch
import "testing"
func TestBatch(t *testing.T){v,e:=Run("\n{\"id\":\"a\",\"kind\":\"credit\",\"amount\":9}\n\n{\"id\":\"b\",\"kind\":\"debit\",\"amount\":4}\n");if e!=nil||v!="total=5\n"{t.Fatal(v,e)};v,e=Run("");if e!=nil||v!="total=0\n"{t.Fatal(v,e)}}
