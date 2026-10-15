package work
import "testing"
func TestClamp(t *testing.T) { for _, c := range [][4]int{{-2,0,10,0},{12,0,10,10},{5,0,10,5},{0,0,10,0},{10,0,10,10},{-9,-5,-1,-5}} { if got:=Clamp(c[0],c[1],c[2]); got!=c[3] {t.Errorf("Clamp%v=%d",c,got)} } }
