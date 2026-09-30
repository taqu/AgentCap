package util

import "testing"

func TestClamp(t *testing.T) {
	for _, c := range [][4]int{{5, 0, 10, 5}, {-1, 0, 10, 0}, {11, 0, 10, 10}} {
		if got := Clamp(c[0], c[1], c[2]); got != c[3] {
			t.Errorf("Clamp(%d,%d,%d) = %d", c[0], c[1], c[2], got)
		}
	}
}

func TestAbs(t *testing.T) {
	if Abs(-3) != 3 || Abs(4) != 4 {
		t.Fatal("Abs")
	}
}

func TestStrings(t *testing.T) {
	if PadRight("ab", 4) != "ab  " || Truncate("abcdefgh", 6) != "abc..." {
		t.Fatal("strings helpers")
	}
}
