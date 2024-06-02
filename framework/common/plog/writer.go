package plog

import (
	"fmt"
	"os"
	"path"
	"sync"
	"sync/atomic"
	"time"
)

type Writer struct {
	sync.RWMutex
	path  string
	name  string
	value atomic.Value
}

func NewWriter(path, file string) *Writer {
	return &Writer{path: path, name: file}
}

func (d *Writer) fileName() string {
	tt := time.Now()
	if len(d.name) > 0 {
		return path.Join(d.path, fmt.Sprintf("%s-%04d%02d%02d.log", d.name, tt.Year(), tt.Month(), tt.Day()))
	}
	return path.Join(d.path, fmt.Sprintf("%04d%02d%02d.log", tt.Year(), tt.Month(), tt.Day()))
}

func (d *Writer) newFile() (fb *os.File, err error) {
	// 判断路径是否存在
	if _, err = os.Stat(d.path); os.IsNotExist(err) {
		if err := os.MkdirAll(d.path, os.FileMode(0755)); err != nil {
			return nil, err
		}
	}
	// 创建文件
	fileName := d.fileName()
	if fb, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644); err != nil {
		return nil, err
	}
	return
}

func (d *Writer) openFile() (*os.File, error) {
	d.Lock()
	defer d.Unlock()
	// 查看是否已经有了
	if fb, ok := d.value.Load().(*os.File); ok && fb != nil {
		return fb, nil
	}
	// 重新创建
	fb, err := d.newFile()
	if err != nil {
		return nil, err
	}
	d.value.Store(fb)
	return fb, nil
}

func (d *Writer) check(fb *os.File) bool {
	st1, err := fb.Stat()
	if err != nil {
		return false
	}
	st2, err := os.Stat(d.fileName())
	if err != nil {
		return false
	}
	return os.SameFile(st1, st2)
}

func (d *Writer) Write(data []byte) (n int, err error) {
	var fb *os.File
	var ok bool
	// 获取文件句柄
	if fb, ok = d.value.Load().(*os.File); ok && fb != nil {
		// 判断句柄是否有效
		if !d.check(fb) {
			// 重新打开文件
			if fb, err = d.openFile(); err != nil {
				return
			}
		}
	} else {
		if fb, err = d.openFile(); err != nil {
			return
		}
	}
	// 写入数据
	return fb.Write(data)
}

func (d *Writer) Close() error {
	if fb, ok := d.value.Load().(*os.File); ok && fb != nil {
		return fb.Close()
	}
	return nil
}
