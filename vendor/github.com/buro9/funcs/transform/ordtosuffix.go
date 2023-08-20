package transform

// OrdToSuffix takes an int, int32 or int64 and returns an nth English suffix
func OrdToSuffix(value interface{}) string {
	var j int64
	switch v := value.(type) {
	case int:
		j = int64(v)
	case int32:
		j = int64(v)
	case int64:
		j = v
	default:
		return ""
	}

	suffix := "th"
	switch j % 10 {
	case 1:
		if j%100 != 11 {
			suffix = "st"
		}
	case 2:
		if j%100 != 12 {
			suffix = "nd"
		}
	case 3:
		if j%100 != 13 {
			suffix = "rd"
		}
	}
	return suffix
}
