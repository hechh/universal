package base

func Prefix[T any](vals []T, pos int) []T {
	if pos < 0 || pos >= len(vals) {
		return nil
	}
	return vals[:pos]
}

func Suffix[T any](vals []T, pos int) []T {
	if pos < 0 || pos >= len(vals) {
		return nil
	}
	return vals[pos:]
}

func Ifelse[T any](flag bool, a, b T) T {
	if flag {
		return a
	}
	return b
}

func Sub(a, b int) int {
	return a - b
}

func Add(a, b int) int {
	return a + b
}

func GetProtoName(name string) string {
	return name + ".proto"
}
