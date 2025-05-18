package base

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"
)

var ROOT = AbsPath("../../../../")

// 获取绝对路径
func AbsPath(file string) string {
	_, filename, _, _ := runtime.Caller(1)
	datapath := filepath.Join(path.Dir(filename), file)
	abs, _ := filepath.Abs(datapath)
	return abs
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func StringIndex(str string, ch byte) (result []int) {
	result = []int{-1}
	for i := 0; i < len(str); i++ {
		if str[i] == ch {
			result = append(result, i)
		}
	}
	result = append(result, len(str))
	return
}

func StringSplit(str string, ch byte, f func(str string)) {
	j := -1
	for i := 0; i <= len(str); i++ {
		if i == len(str) || ch == str[i] {
			f(str[j+1 : i])
			j = i
		}
	}
}

func RuleSplit(str string) (index []int, rule, desc string) {
	if pos := strings.LastIndex(str, "|#"); pos > 0 {
		desc = str[pos+2:]
		str = str[:pos]
	}
	index = StringIndex(str, '|')
	rule = str
	return
}

// 将驼峰命名分割
func splitWord(pbname []byte, data *[]string) {
	var index int
	// 连续大写匹配
	for ; index < len(pbname) && unicode.IsUpper(rune(pbname[index])); index++ {
	}
	// 连续小写匹配
	if index <= 1 {
		for ; index < len(pbname) && !unicode.IsUpper(rune(pbname[index])); index++ {
		}
	} else if index+1 < len(pbname) {
		index--
	}
	*data = append(*data, strings.ToLower(string(pbname[:index])))
	if index < len(pbname) {
		splitWord(pbname[index:], data)
	}
}

func ToCmd(x string) string {
	tmp := []string{}
	if !strings.Contains(strings.ToLower(x), "cmd") {
		tmp = append(tmp, "CMD")
	}
	splitWord([]byte(x), &tmp)
	return strings.ToUpper(strings.Join(tmp, "_"))
}

func ToEvent(x string) string {
	tmp := []string{}
	splitWord([]byte(x), &tmp)
	return strings.ToUpper(strings.Join(tmp, "_"))
}

func FirstToBig(x string) string {
	if len(x) <= 0 {
		return x
	}
	ret := []byte(x)
	ret[0] = byte(unicode.ToUpper(rune(x[0])))
	return string(ret)
}

func FirstToLow(x string) string {
	if len(x) <= 0 {
		return x
	}
	ret := []byte(x)
	ret[0] = byte(unicode.ToLower(rune(x[0])))
	return string(ret)
}

// 小写 + 下划线 格式
func ToUnderline(pbname string) string {
	buf := []string{}
	splitWord([]byte(pbname), &buf)
	return strings.Join(buf, "_")
}

func ToBigHump(name string) string {
	if len(name) <= 0 {
		return ""
	}
	result := []byte(strings.ToLower(name))
	result[0] = byte(unicode.ToUpper(rune(name[0])))
	return string(result)
}

func isFile(filename string) bool {
	_, err := os.Lstat(filename)
	return !os.IsNotExist(err)
}

func GenProto(buf *bytes.Buffer, gfile string, flag bool) {
	if !flag && isFile(gfile) {
		return
	}
	// 生成文档
	if err := os.MkdirAll(filepath.Dir(gfile), os.FileMode(0777)); err != nil {
		panic(err)
		return
	}
	if err := ioutil.WriteFile(gfile, buf.Bytes(), os.FileMode(0666)); err != nil {
		panic(err)
	}
}

func GenGo(buf *bytes.Buffer, gfile string, flag bool) {
	if !flag && isFile(gfile) {
		return
	}
	// 格式化
	result, err := format.Source(buf.Bytes())
	if err != nil {
		ioutil.WriteFile("./gen.go", buf.Bytes(), os.FileMode(0644))
		panic(err)
		return
	}
	// 生成文档
	if err := os.MkdirAll(filepath.Dir(gfile), os.FileMode(0777)); err != nil {
		panic(err)
		return
	}
	if err := ioutil.WriteFile(gfile, result, os.FileMode(0666)); err != nil {
		panic(err)
	}
}

func TrimSpace(str string) string {
	j := -1
	buf := []byte(str)
	for _, val := range buf {
		if val == ' ' {
			continue
		}
		j++
		buf[j] = val
	}
	buf = buf[:j+1]
	return string(buf)
}
