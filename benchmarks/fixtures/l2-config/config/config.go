package config
import "strconv"
type Config struct{ Host string; Port int }
func Load(env map[string]string)(Config,error){p,e:=strconv.Atoi(env["PORT"]); return Config{Host:env["HOST"],Port:p},e}
