package service
import "bench.local/work/paging"
func List(items []string,offset,limit int)([]string,error){v,_:=paging.Page(items,offset,limit);return v,nil}
