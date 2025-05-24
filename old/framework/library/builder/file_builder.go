package builder

import (
	"os"
	"universal/framework/library/uerror"
	"universal/framework/library/util"
)

type FileBuilder struct {
	n        int
	size     int
	buf      []byte
	filename string
	fp       *os.File
}

func NewFileBuilder(size int) *FileBuilder {
	return &FileBuilder{buf: make([]byte, size), size: size}
}

func (c *FileBuilder) SetWriter(fileName string) (err error) {
	if len(c.filename) <= 0 {
		c.filename = fileName
		c.fp, err = util.NewOrOpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR)
		if err != nil {
			return
		}
	}
	if c.filename != fileName {
		// 保存数据
		if err = c.Flush(); err != nil {
			return
		}
		if err = c.Close(); err != nil {
			return
		}
		c.fp = nil
		c.fp, err = util.NewOrOpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR)
	}
	return
}

func (c *FileBuilder) Write(items []byte) (err error) {
	diff := c.size - c.n
	if len(items) <= diff {
		copy(c.buf[c.n:], items)
		c.n += len(items)
		return
	}
	if diff > 0 {
		copy(c.buf[c.n:], items[:diff])
		c.n += diff
	}
	if err = c.Flush(); err != nil {
		return
	}
	return c.Write(items[diff:])
}

func (c *FileBuilder) Flush() (err error) {
	if c.n <= 0 {
		return
	}
	if c.fp == nil {
		if len(c.filename) <= 0 {
			return uerror.New(1, -1, "写入文件为空")
		}
		c.fp, err = util.NewOrOpenFile(c.filename, os.O_CREATE|os.O_APPEND|os.O_RDWR)
		if err != nil {
			return
		}
	}
	if _, err = c.fp.Write(c.buf[:c.n]); err != nil {
		return
	}
	c.n = 0
	return
}

func (c *FileBuilder) Close() error {
	if c.fp != nil {
		return c.fp.Close()
	}
	return nil
}
