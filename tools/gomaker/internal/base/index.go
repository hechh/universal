package base

import (
	"strings"
)

// 注释规则中的： field:a@int,b@string,...
type Index struct {
	Field  string
	Values Params
}

func (d *Index) GetRange() string {
	if len(d.Field) <= 0 {
		return d.Values.UniqueID()
	}
	return d.Field
}

func (d *Index) UniqueID() string {
	id := strings.Join(d.Values.Vals(""), "")
	if len(d.Field) <= 0 {
		return id
	}
	return d.Field + id
}

func (d *Index) Index() string {
	return "index" + d.UniqueID()
}

// fmt.Sprintf格式化
func (d *Index) Format() string {
	aa := d.Values.Format()
	if len(d.Field) > 0 && len(aa) > 0 {
		return d.Field + ":" + aa
	}
	if len(aa) > 0 {
		return aa
	}
	return d.Field
}

func (d *Index) Count() int {
	return d.Values.Count()
}

// a@xx,b@yy   a,b@yy
func ParseParams(x string) (rets Params) {
	typ := x[strings.Index(x, "@")+1:]
	StringSplit(x, ',', func(str string) {
		if pos := strings.Index(str, "@"); pos > 0 {
			rets = append(rets, &Param{Name: str[:pos], Type: str[pos+1:]})
		} else {
			rets = append(rets, &Param{Name: str, Type: typ})
		}
	})
	return
}

// ff:a@xx    ff  ff:a,b@xx
func ParseIndex(x string) *Index {
	if pos := strings.Index(x, ":"); pos > 0 {
		return &Index{Field: x[:pos], Values: ParseParams(x[pos+1:])}
	}
	if strings.Contains(x, "@") {
		return &Index{Values: ParseParams(x)}
	}
	return &Index{Field: x}
}
