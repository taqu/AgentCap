package main
import("os";"fmt";"io";"bench.local/work/batch")
func main(){data,e:=io.ReadAll(os.Stdin);if e!=nil{panic(e)};v,e:=batch.Run(string(data));if e!=nil{fmt.Fprintln(os.Stderr,e);return};fmt.Print(v)}
