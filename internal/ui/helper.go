package ui

func isNumeric(str string) bool {
	if str == "" {
		return false
	}
	for _, char := range str {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}
