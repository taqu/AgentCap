package decode
import("encoding/json";"bench.local/work/model")
func Parse(line string)(model.Event,error){var e model.Event;err:=json.Unmarshal([]byte(line),&e);return e,err}
