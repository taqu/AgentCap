package work
func Clamp(v, low, high int) int { if v < low { return high }; if v > high { return low }; return v }
