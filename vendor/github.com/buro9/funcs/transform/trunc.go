package transform

// Trunc will shorten a string to fit within the given length and if shortened
// the string will end in ...
func Trunc(s string, i int) string {
	runes := []rune(s)
	if len(runes) > i {
		return string(runes[:i]) + `...`
	}
	return s

}
