package ledger
import "bench.local/work/model"
type Ledger struct { total int; seen map[string]bool }
func New()*Ledger{return &Ledger{seen:map[string]bool{}}}
func(l *Ledger)Apply(e model.Event)error{l.total+=e.Amount;return nil}
func(l *Ledger)Total()int{return l.total}
