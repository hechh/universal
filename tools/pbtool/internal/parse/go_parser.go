package parse

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"universal/library/uerror"
	"universal/tools/pbtool/domain"
	"universal/tools/pbtool/internal/base"
	"universal/tools/pbtool/internal/manager"

	"github.com/iancoleman/strcase"
)

type Parser struct {
	rule string
}

// 解析go文件
func ParseFiles(v ast.Visitor, files ...string) error {
	fset := token.NewFileSet()
	for _, filename := range files {
		// 解析语法树
		fs, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
		if err != nil {
			return uerror.N(1, -1, "filename: %v, error: %v", filename, err)
		}
		// 遍历语法树
		ast.Walk(v, fs)
	}
	return nil
}

func (d *Parser) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.File:
		return d
	case *ast.GenDecl:
		d.rule = ""
		if n.Doc == nil || len(n.Doc.List) <= 0 || n.Tok != token.TYPE {
			return nil
		}
		for _, str := range n.Doc.List {
			rule := strings.TrimPrefix(str.Text, "//")
			rule = strings.TrimSpace(rule)
			if strings.HasPrefix(rule, domain.RULE_HEAD) {
				rule = strings.TrimPrefix(rule, domain.RULE_HEAD)
				d.rule = rule
			}
		}
		if len(d.rule) <= 0 {
			return nil
		}
		return d
	case *ast.TypeSpec:
		switch vv := n.Type.(type) {
		case *ast.StructType:
			strs := strings.Split(d.rule, "|")
			switch strings.ToLower(strs[0]) {
			case domain.STRING:
				if item := parseString(n.Name.Name, vv, strs...); item != nil {
					manager.AddString(item)
				}
			case domain.HASH:
				if item := parseHash(n.Name.Name, vv, strs...); item != nil {
					manager.AddHash(item)
				}
			}
		}
	}
	return nil
}

func parseHash(name string, vv *ast.StructType, strs ...string) *base.Hash {
	desc := ""
	if len(strs) > 4 {
		desc = strs[4]
	}
	kfmt, keys := parseKey(strs[2])
	ffmt, fields := parseKey(strs[3])
	return &base.Hash{
		Pkg:     strcase.ToSnake(name),
		Name:    name,
		DbName:  strs[1],
		Desc:    desc,
		KFormat: kfmt,
		Keys:    keys,
		FFormat: ffmt,
		Fields:  fields,
	}
}

func parseString(name string, vv *ast.StructType, strs ...string) *base.String {
	desc := ""
	if len(strs) > 3 {
		desc = strs[3]
	}
	format, keys := parseKey(strs[2])

	// 初始化
	return &base.String{
		Pkg:    strcase.ToSnake(name),
		Name:   name,
		DbName: strs[1],
		Desc:   desc,
		Format: format,
		Keys:   keys,
	}
}

func parseKey(str string) (format string, ffs []*base.Field) {
	ffmts := []string{}
	var keys []string
	if index := strings.Index(str, ":"); index > 0 {
		ffmts = strings.Split(str[:index], ",")
		keys = strings.Split(str[index+1:], ",")
	} else {
		keys = strings.Split(str, ",")
	}
	// 解析key类型
	for _, key := range keys {
		lls := strings.Split(key, "@")
		switch strings.ToLower(lls[1]) {
		case "string":
			ffmts = append(ffmts, "%s")
		default:
			ffmts = append(ffmts, "%d")
		}
		ffs = append(ffs, &base.Field{Name: lls[0], Type: lls[1]})
	}
	format = strings.Join(ffmts, ":")
	return
}
