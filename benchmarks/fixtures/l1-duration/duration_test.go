package work
import ("testing";"time")
func TestSeconds(t *testing.T){ for _,n:=range []string{"0","1","60"}{ d,e:=ParseSeconds(n); if e!=nil {t.Fatal(e)}; if n=="60" && d!=time.Minute {t.Errorf("60 = %s",d)}; if n=="1" && d!=time.Second {t.Errorf("1 = %s",d)} }; for _,s:=range []string{"-1","bad","1.5"}{if _,e:=ParseSeconds(s);e==nil{t.Errorf("accepted %q",s)}} }
