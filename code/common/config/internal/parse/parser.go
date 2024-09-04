package parse

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"universal/common/config/domain"
)

type Parser struct {
	sheetName string           // 配置文件名
	fileInfo  os.FileInfo      // 文件修改时间
	cfgs      []domain.IConfig // 配置解析器
}

func NewParser(name string, cfgs ...domain.IConfig) *Parser {
	return &Parser{
		sheetName: name,
		cfgs:      cfgs,
	}
}

func (d *Parser) Register(cfgs ...domain.IConfig) {
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
	for _, val := range d.cfgs {
		if err = val.LoadFile(buf); err != nil {
			return
		}
	}
	return
}
