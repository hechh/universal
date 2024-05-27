package domain

const (
	PackageAction = "package"
	PackageTpl    = "package.tpl"
)

type CmdLine struct {
	Action string
	Param  string
	Tpl    string
	Src    string
	Dst    string
}

type IMaker interface {
	GetHelp(string) string                 // help信息
	OpenTpl(*CmdLine) error                // 打开tpl文件
	ParseFile(*CmdLine, interface{}) error // 解析文件
	Gen(*CmdLine) error                    // 生成文件
}

// token类型
const (
	ENUM    = 0x01
	STRUCT  = 0x02
	ALIAS   = 0x04
	POINTER = 0x08
	ARRAY   = 0x10
	MAP     = 0x20
)
