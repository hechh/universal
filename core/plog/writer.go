package plog

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"sync"
	"time"
)

type MetaData struct {
	tt  time.Time
	msg []byte
}

type Writer struct {
	sync.WaitGroup
	path   string
	name   string
	exitCh chan struct{}
	metaCh chan *MetaData
}

func NewWriter(path, file string) *Writer {
	ret := &Writer{
		path:   path,
		name:   file,
		exitCh: make(chan struct{}, 0),
		metaCh: make(chan *MetaData, 50),
	}
	ret.Add(1)
	go ret.run()
	return ret
}

func (d *Writer) Write(tt time.Time, data []byte) {
	d.metaCh <- &MetaData{tt: tt, msg: data}
}

func (d *Writer) Close() {
	d.exitCh <- struct{}{}
	d.Wait()
}

func (d *Writer) walk(f func(time.Time, []byte)) {
	select {
	case item := <-d.metaCh:
		f(item.tt, item.msg)
		d.walk(f)
	default:
	}
}

func (d *Writer) run() {
	var fp *FPointer
	writeFunc := func(tt time.Time, msg []byte) {
		fileName := getFileName(tt, d.path, d.name)
		if fp == nil {
			fp = newFPointer(fileName)
		}
		if fp != nil {
			if fp.fileName == fileName {
				if fp.keep() {
					if fb, err := newFile(fp.fileName); err == nil {
						fp.fb = fb
						fp.ww.Reset(fb)
					}
				}
			} else {
				fp.Close()
				fp = nil
				fp = newFPointer(fileName)
			}
		}
		if fp != nil {
			fp.ww.Write(msg)
		}
	}
	defer func() {
		d.walk(writeFunc)
		if fp != nil {
			fp.Close()
		}
		d.Done()
	}()
	for tt := time.NewTimer(1 * time.Second); ; {
		select {
		case <-d.exitCh:
			return
		case <-tt.C:
			if fp != nil {
				fp.ww.Flush()
			}
		case item := <-d.metaCh:
			writeFunc(item.tt, item.msg)
		}
	}
}

type FPointer struct {
	ww       *bufio.Writer
	fb       *os.File
	fileName string
}

func (d *FPointer) Close() {
	d.ww.Flush()
	d.fb.Close()
}

func (d *FPointer) keep() bool {
	st1, err := os.Stat(d.fileName)
	st2, _ := d.fb.Stat()
	return err != nil || !os.SameFile(st1, st2)
}

func newFPointer(fileName string) *FPointer {
	fb, err := newFile(fileName)
	if err != nil {
		return nil
	}
	// 返回
	return &FPointer{ww: bufio.NewWriter(fb), fb: fb, fileName: fileName}
}

func newFile(fileName string) (fb *os.File, err error) {
	// 判断路径是否存在
	pp := path.Dir(fileName)
	if _, err = os.Stat(pp); err != nil {
		if err := os.MkdirAll(pp, os.FileMode(0755)); err != nil {
			return nil, err
		}
	}
	// 创建文件
	if fb, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644); err != nil {
		return nil, err
	}
	return
}

func getFileName(tt time.Time, fpath, name string) string {
	if len(name) > 0 {
		return path.Clean(fmt.Sprintf("%s/%04d%02d%02d/%s%02d.log", fpath, tt.Year(), tt.Month(), tt.Day(), name, tt.Hour()))
	}
	return path.Clean(fmt.Sprintf("%s/%04d%02d%02d/%02d.log", fpath, tt.Year(), tt.Month(), tt.Day(), tt.Hour()))
}
