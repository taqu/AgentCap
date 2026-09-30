// Package util holds small shared helpers.
package util

// Clamp limits v to [lo, hi].
func Clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// Abs returns |v|.
func Abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
