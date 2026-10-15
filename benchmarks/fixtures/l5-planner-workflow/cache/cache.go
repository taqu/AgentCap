package cache
import("bench.local/work/plan")
type Cache struct {last string}
func(c *Cache)Build(data []byte)(string,error){if c.last!=""{return c.last,nil};v,e:=plan.Build(data);c.last=v;return v,e}
