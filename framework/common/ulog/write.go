package ulog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type IWriter interface {
	Close() error
	Write([]byte) (int, error)
}

type Writer struct {
	sync.RWMutex
	filename string
	value    atomic.Value
}

func NewWriter(file string) *Writer {
	ext := filepath.Ext(file)
	if len(ext) <= 0 {
		file += ".log"
	}
	return &Writer{filename: file}
}

func (d *Writer) getFileName() string {
	ext := filepath.Ext(d.filename)
	bb := strings.TrimSuffix(d.filename, ext)
	tt := time.Now()
	return fmt.Sprintf("%s-%04d%02d%02d.%s", bb, tt.Year(), tt.Month(), tt.Day(), ext)
}

func (d *Writer) newFile() (fb *os.File, err error) {
	// 判断路径是否存在
	fileName := d.getFileName()
	path := filepath.Dir(fileName)
	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.FileMode(0755)); err != nil {
			return nil, err
		}
	}
	// 创建文件
	if fb, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644); err != nil {
		return nil, err
	}
	return
}

func (d *Writer) getFb() (*os.File, error) {
	d.Lock()
	defer d.Unlock()
	// 查看是否已经有了
	fb, ok := d.value.Load().(*os.File)
	if ok && fb != nil {
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
	st2, err := os.Stat(d.getFileName())
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
			if fb, err = d.getFb(); err != nil {
				return
			}
		}
	} else {
		// 重新打开文件
		if fb, err = d.getFb(); err != nil {
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
