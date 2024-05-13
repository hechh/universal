package domain

import "html/template"

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

type IParser interface {
	GetHelp() string
	GetAction() string
	OpenTpl(string, *CmdLine) error
	ParseFile(string, *CmdLine, interface{}) error
	Gen(string, *CmdLine) error
}

type CmdLine struct {
	Action string
	Param  string
	Tpl    string
	Src    string
	Dst    string
}

type GenFunc func(string, *CmdLine, map[string]*template.Template) error
