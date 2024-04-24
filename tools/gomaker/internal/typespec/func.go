package typespec

import "strings"

type BaseFunc struct{}

func (d *BaseFunc) TrimPrefix(str, prefix string) string {
	return strings.TrimPrefix(str, prefix)
}

func (d *BaseFunc) TrimSuffix(str, prefix string) string {
	return strings.TrimSuffix(str, prefix)
}

func (d *BaseFunc) Split(str, sp string) []string {
	return strings.Split(str, sp)
}

func (d *BaseFunc) Join(a, b string) string {
	return a + b
}
