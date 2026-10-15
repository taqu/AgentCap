package main
import("io";"os";"fmt";"bench.local/work/plan")
func main(){b,_:=io.ReadAll(os.Stdin);s,e:=plan.Build(b);if e!=nil{fmt.Fprintln(os.Stderr,e);return};fmt.Print(s)}
