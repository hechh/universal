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
	POINTER = 0x01
	ARRAY   = 0x02
	MAP     = 0x04
	STRUCT  = 0x08
	ALIAS   = 0x10
	ENUM    = 0x20
)
