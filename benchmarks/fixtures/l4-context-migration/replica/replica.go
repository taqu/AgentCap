package replica
import("bench.local/work/storage";"bench.local/work/record")
func Copy(s *storage.Store)map[string]record.Record{return s.Snapshot()}
