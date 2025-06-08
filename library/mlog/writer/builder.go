package writer

import (
	"os"
	"universal/library/fileutil"
)

type FileBuilder struct {
	n        int
	size     int
	buff     []byte
	filename string
	fp       *os.File
}

func NewFileBuilder(size int) *FileBuilder {
	return &FileBuilder{buff: make([]byte, size), size: size}
}

func (d *FileBuilder) SetFile(filename string) (err error) {
	if d.filename != filename {
		if err := d.Close(); err != nil {
			return err
		}
		d.fp = nil
		d.fp, err = fileutil.CreateFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR)
		return
	}

	// 日志文件被删除
	if !fileutil.IsSameFile(d.fp, filename) {
		if d.fp != nil {
			d.fp.Close()
			d.fp = nil
		}
		d.fp, err = fileutil.CreateFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR)
	}
	return
}

func (d *FileBuilder) Write(buf []byte) error {
	diff := d.size - d.n
	if ll := len(buf); ll <= diff {
		copy(d.buff[d.n:], buf)
		d.n += ll
		return nil
	}
	if diff > 0 {
		copy(d.buff[d.n:], buf[:diff])
		d.n += diff
	}
	if err := d.Flush(); err != nil {
		return err
	}
	return d.Write(buf[diff:])
}

func (d *FileBuilder) Flush() error {
	if d.fp == nil || d.n <= 0 {
		return nil
	}
	if _, err := d.fp.Write(d.buff[:d.n]); err != nil {
		return err
	}
	d.n = 0
	return nil
}

func (d *FileBuilder) Close() error {
	if err := d.Flush(); err != nil {
		return err
	}
	if d.fp == nil {
		return nil
	}
	return d.fp.Close()
}
