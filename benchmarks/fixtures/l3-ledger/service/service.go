package service
import("bench.local/work/decode";"bench.local/work/ledger";"bench.local/work/report")
func Run(lines []string)(string,error){l:=ledger.New();for _,s:=range lines{e,err:=decode.Parse(s);if err!=nil{return "",err};if err=l.Apply(e);err!=nil{return "",err}};return report.Format(l.Total()),nil}
