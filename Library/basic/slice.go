package basic

func GetPrefix(str string, pos int) string {
	if pos < 0 || pos >= len(str) {
		return ""
	}
	return str[:pos]
}

func GetSuffix(str string, pos int) string {
	if pos < 0 || pos >= len(str) {
		return ""
	}
	return str[pos:]
}

func Prefix[T any](str []T, pos int) []T {
	if pos < 0 || pos >= len(str) {
		return nil
	}
	return str[:pos]
}

func Suffix[T any](str []T, pos int) []T {
	if pos < 0 || pos >= len(str) {
		return nil
	}
	return str[pos:]
}
