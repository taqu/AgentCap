package export
import("bench.local/work/storage";"encoding/json")
func JSON(s *storage.Store)([]byte,error){return json.Marshal(s.Snapshot())}
