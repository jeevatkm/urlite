package util

func Contains(a []string, s string) bool {
	for _, v := range a {
		if v == s {
			return true
		}
	}

	return false
}
