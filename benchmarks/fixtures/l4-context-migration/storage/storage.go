package storage
import ("bench.local/work/record";"context")
type Store struct { data map[string]record.Record }
func New()*Store{return &Store{data:map[string]record.Record{}}}
func(s *Store)Put(ctx context.Context,key,value string,expected int)(record.Record,error){r:=record.Record{Key:key,Value:value,Revision:expected+1};s.data[key]=r;return r,nil}
func(s *Store)Get(ctx context.Context,key string)(record.Record,bool){r,ok:=s.data[key];return r,ok}
func(s *Store)Snapshot(ctx context.Context)map[string]record.Record{return s.data}
