package parse

import (
	"fmt"
	"hego/common/config/domain"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Parser struct {
	sheetName string            // 配置文件名
	fileInfo  os.FileInfo       // 文件修改时间
	cfgs      []domain.LoadFunc // 配置解析器
}

func NewParser(name string, cfgs ...domain.LoadFunc) *Parser {
	return &Parser{
		sheetName: name,
		cfgs:      cfgs,
	}
}

func (d *Parser) Register(cfgs ...domain.LoadFunc) {
	d.cfgs = append(d.cfgs, cfgs...)
}

func (d *Parser) Check(dir string) bool {
	filename := filepath.Join(dir, fmt.Sprintf("%s.bytes", d.sheetName))
	st, _ := os.Stat(filename)
	return os.SameFile(st, d.fileInfo)
}

// 加载配置
func (d *Parser) Load(dir string) (err error) {
	// 更新文件信息
	filename := filepath.Join(dir, fmt.Sprintf("%s.bytes", d.sheetName))
	if d.fileInfo, err = os.Stat(filename); err != nil {
		return
	}
	// 加载整个文件
	var buf []byte
	if buf, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	// 解析配置文件
	for _, f := range d.cfgs {
		if err = f(buf); err != nil {
			return
		}
	}
	return
}
