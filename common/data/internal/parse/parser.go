package parse

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"universal/library/baselib/uerror"
)

type Parser struct {
	fileName string             // 配置文件名
	fileInfo os.FileInfo        // 文件修改时间
	loadFun  func([]byte) error // 配置解析器
}

func NewParser(name string, f func([]byte) error) *Parser {
	return &Parser{
		fileName: name,
		loadFun:  f,
	}
}

func (d *Parser) IsChange(dir, ext string) bool {
	filename := filepath.Join(dir, fmt.Sprintf("%s.%s", d.fileName, ext))
	st, _ := os.Stat(filename)
	return os.SameFile(st, d.fileInfo)
}

// 加载配置
func (d *Parser) Load(dir, ext string) error {
	// 更新文件信息
	filename := filepath.Join(dir, fmt.Sprintf("%s.%s", d.fileName, ext))

	fileInfo, err := os.Stat(filename)
	if err != nil {
		return uerror.New(1, -1, "文件不存在: %s", filename)
	}
	d.fileInfo = fileInfo

	// 加载整个文件
	var buf []byte
	if buf, err = ioutil.ReadFile(filename); err != nil {
		return uerror.New(1, -1, err.Error())
	}

	// 解析配置文件
	return d.loadFun(buf)
}
