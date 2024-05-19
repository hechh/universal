package domain

// token类型
const (
	IDENT   = 0x01
	POINTER = 0x02
	ARRAY   = 0x04
	MAP     = 0x08
)

const (
	PackageAction = "package"
	PackageTpl    = "package.tpl"
)

type IMaker interface {
	GetHelp(string) string                 // help信息
	OpenTpl(*CmdLine) error                // 打开tpl文件
	ParseFile(*CmdLine, interface{}) error // 解析文件
	Gen(*CmdLine) error                    // 生成文件
}

type CmdLine struct {
	Action string
	Param  string
	Tpl    string
	Src    string
	Dst    string
}
