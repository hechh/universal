package domain

import "text/template"

// token类型
const (
	IDENT   = 0x01
	POINTER = 0x02
	ARRAY   = 0x04
	MAP     = 0x08
)

const (
	PACKAGE    = "package"
	UERRORS    = "uerrors"
	PLAYER_FUN = "playerFun"
	ENTITY     = "entity"
)

type GenFunc func(string, *CmdLine, map[string]*template.Template) error

type IParser interface {
	GetHelp() string                               // help信息
	GetAction() string                             // action类型
	OpenTpl(string, *CmdLine) error                // 打开tpl文件
	ParseFile(string, *CmdLine, interface{}) error // 解析文件
	Gen(string, *CmdLine) error                    // 生成文件
}

type CmdLine struct {
	Action string
	Param  string
	Tpl    string
	Src    string
	Dst    string
}
