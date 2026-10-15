package batch
import("bench.local/work/cache")
func Build(inputs [][]byte)([]string,error){c:=&cache.Cache{};out:=[]string{};for _,in:=range inputs{v,e:=c.Build(in);if e!=nil{continue};out=append(out,v)};return out,nil}
