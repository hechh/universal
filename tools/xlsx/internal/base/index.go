package base

import (
	"fmt"
	"hego/tools/xlsx/domain"
	"strings"
)

func (d *Index) Arg(pkg, split string) string {
	strs := []string{}
	for _, field := range d.List {
		strs = append(strs, field.Name+" "+field.Type.GetType(pkg))
	}
	return strings.Join(strs, split)
}

func (d *Index) Value(ref, split string) string {
	strs := []string{}
	for _, field := range d.List {
		if len(ref) > 0 {
			strs = append(strs, ref+"."+field.Name)
		} else {
			strs = append(strs, field.Name)
		}
	}
	return strings.Join(strs, split)
}

func (d *Index) IndexValue(val string) string {
	switch d.Type.TypeOf {
	case domain.TypeOfBase:
		return val
	}
	return d.Name + "{" + val + "}"
}

func (d *Index) Templ(pkg, split string) string {
	strs := []string{}
	for _, field := range d.List {
		strs = append(strs, field.Type.GetType(pkg))
	}
	return strings.Join(strs, split)
}

func (d *Index) GetType(pkg, cfg string) string {
	switch d.Type.TypeOf {
	case domain.TypeOfBase:
		switch d.Type.ValueOf {
		case domain.ValueOfMap:
			return "map[" + d.Type.Name + "]*" + cfg
		case domain.ValueOfGroup:
			return "map[" + d.Type.Name + "][]*" + cfg
		}
	case domain.TypeOfEnum:
		switch d.Type.ValueOf {
		case domain.ValueOfMap:
			return "map[" + d.Type.GetType(pkg) + "]*" + cfg
		case domain.ValueOfGroup:
			return "map[" + d.Type.GetType(pkg) + "][]*" + cfg
		}
	case domain.TypeOfStruct:
		switch d.Type.ValueOf {
		case domain.ValueOfMap:
			return fmt.Sprintf("map[%s]*%s", d.Name, cfg)
		case domain.ValueOfGroup:
			return fmt.Sprintf("map[%s][]*%s", d.Name, cfg)
		}
	}
	return "[]*" + cfg
}
