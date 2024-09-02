package local

import (
	"fmt"
	"strings"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/typespec"

	"github.com/spf13/cast"
)

type XlsxParser struct {
	enums  map[string]*typespec.Value // 中文--->所有代对
	tables map[string]string          // 中文--->英文
}

func (d *XlsxParser) Visit(row int, rowData []string) domain.Visitor {
	for _, val := range rowData {
		if !strings.Contains(val, ":") {
			continue
		}

		ss := strings.Split(strings.TrimSpace(val), ":")
		switch ss[0] {
		case "C", "c":
		case "CS", "cs", "Cs", "cS", "SC", "Sc", "sC", "sc", "s", "S":

		case "E", "e":
			val := &typespec.Value{
				Name: fmt.Sprintf("%s_%s", ss[2], ss[3]),
				Type: manager.GetOrAddType(&typespec.Type{
					Kind:     domain.ENUM,
					Selector: "pb",
					Name:     ss[2],
				}),
				Value:   cast.ToInt32(ss[4]),
				Comment: ss[1],
			}
			// 存储类型
			d.enums[val.Comment] = val
			manager.AddValue(val)
		}
	}
	return d
}
