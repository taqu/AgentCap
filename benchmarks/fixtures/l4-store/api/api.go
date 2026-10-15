package api
import("bench.local/work/storage";"bench.local/work/record")
type API struct{ Store *storage.Store }
func(a API)Update(k,v string,rev int)(record.Record,error){return a.Store.Put(k,v,rev)}
