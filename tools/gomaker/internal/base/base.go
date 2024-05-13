package base

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"universal/common/pb"
	"universal/framework/fbasic"
	"universal/tools/gomaker/domain"

	"github.com/xuri/excelize/v2"
)

func GetAbsPath(src string, root string) string {
	if len(src) <= 0 {
		return root
	}
	if !filepath.IsAbs(src) {
		return filepath.Join(root, src)
	}
	return src
}

func GetPath(dst string, defaultEnv string) string {
	if len(dst) > 0 {
		return dst
	}
	return defaultEnv
}

func GetFilePathBase(dst string) string {
	ext := filepath.Ext(dst)
	if len(ext) <= 0 {
		return filepath.Base(dst)
	}
	return filepath.Base(filepath.Dir(dst))
}

func ParseXlsx(file string, f func(sheet string, row int, cols []string) bool) error {
	fb, err := excelize.OpenFile(file)
	if err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_OpenFile, err)
	}
	defer fb.Close()
	for _, sheet := range fb.GetSheetList() {
		rows, err := fb.GetRows(sheet)
		if err != nil {
			return err
		}
		for row, cols := range rows {
			if !f(sheet, row, cols) {
				break
			}
		}
	}
	return nil
}

func InitCmdLine(cwdPath string, cmdLine *domain.CmdLine) error {
	flag.StringVar(&cmdLine.Action, "action", "", "操作模式")
	flag.StringVar(&cmdLine.Param, "param", "", "不同action的param使用方式不同")
	flag.StringVar(&cmdLine.Tpl, "tpl", "", "模板文件路径, 默认从TPL_GO环境变量中读取, -tpl={tpl文件目录}")
	flag.StringVar(&cmdLine.Src, "src", "", "解析文件或目录, -param={go、xlsx文件目录 或 *.xlsx、*.go}")
	flag.StringVar(&cmdLine.Dst, "dst", "", "生成文件或目录, -dst={生成文件目录}")
	flag.Parse()
	// action模式
	if len(cmdLine.Action) <= 0 {
		return fbasic.NewUError(1, pb.ErrorCode_Parameter, "-action")
	}
	// 模版文件
	cmdLine.Tpl = GetAbsPath(GetPath(cmdLine.Tpl, os.Getenv("TPL_GO")), cwdPath)
	// 地址
	cmdLine.Dst = GetAbsPath(GetPath(cmdLine.Dst, cwdPath), cwdPath)
	return nil
}

type MapTpl map[string]*template.Template

func (d MapTpl) GenPackage(dst string, buf *bytes.Buffer) error {
	tpl, ok := d[domain.PACKAGE]
	if !ok {
		return fbasic.NewUError(1, -1, fmt.Sprintf("The tpl of %s.tpl is not supported", domain.PACKAGE))
	}
	if err := tpl.ExecuteTemplate(buf, domain.PACKAGE+".tpl", GetFilePathBase(dst)); err != nil {
		return fbasic.NewUError(1, -1, err)
	}
	return nil
}
