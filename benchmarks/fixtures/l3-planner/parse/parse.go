package parse
import("encoding/json";"bench.local/work/graph")
func Load(data []byte)(graph.Graph,error){var nodes []graph.Node;if e:=json.Unmarshal(data,&nodes);e!=nil{return nil,e};g:=graph.Graph{};for _,n:=range nodes{g[n.Name]=n};return g,nil}
