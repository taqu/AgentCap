package batch
import("bench.local/work/service";"strings")
func Run(input string)(string,error){return service.Run(strings.Split(input,"\n"))}
