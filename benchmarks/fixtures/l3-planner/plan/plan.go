package plan
import("bench.local/work/parse";"strings")
func Build(data []byte)(string,error){g,e:=parse.Load(data);if e!=nil{return "",e};order,e:=g.Order();return strings.Join(order,","),e}
