package typespec

import "strings"

type Enum struct {
	ID    string
	Type  string
	Name  string
	Value int32
}

type Field struct {
	Name string
	Type string
	Doc  string
}

func ParseField(val01, val02 []string) (rets []*Field) {
	for i, val := range val01 {
		if len(val02) <= i {
			val02 = append(val02, "")
		}
		val = strings.TrimSpace(val)
		if len(val) <= 0 {
			continue
		}
		pos := strings.Index(val, "/")
		if pos <= 0 {
			rets = append(rets, &Field{Name: val, Type: "string", Doc: val02[i]})
		} else {
			rets = append(rets, &Field{Name: val[:pos], Type: val[pos+1:], Doc: val02[i]})
		}
	}
	return
}
