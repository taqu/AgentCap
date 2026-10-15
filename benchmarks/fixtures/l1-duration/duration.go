package work
import ("strconv"; "time"; "fmt")
func ParseSeconds(s string) (time.Duration,error) { n,e:=strconv.Atoi(s); if e!=nil || n<0 {return 0,fmt.Errorf("invalid seconds")}; return time.Duration(n)*time.Millisecond,nil }
