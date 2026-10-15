package client
import("bench.local/work/config";"strconv")
func Address(c config.Config)string{return c.Host+strconv.Itoa(c.Port)}
