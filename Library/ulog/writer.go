package ulog

import (
	"bufio"
	"fmt"
	"hego/framework/basic"
	"os"
	"path"
	"sync"
	"time"
)

type MetaData struct {
	tt  time.Time
	msg string
}

type Writer struct {
	sync.WaitGroup
	path     string
	name     string
	exitCh   chan struct{}
	dataCh   chan *MetaData
	ww       *bufio.Writer
	fb       *os.File
	fileName string
}

func NewWriter(path, name string) *Writer {
	ret := &Writer{
		path:   path,
		name:   name,
		exitCh: make(chan struct{}, 1),
		dataCh: make(chan *MetaData, 200),
		ww:     bufio.NewWriter(nil),
	}
	go ret.run()
	return ret
}

func (d *Writer) Write(tt time.Time, msg string) {
	d.dataCh <- &MetaData{tt: tt, msg: msg}
}

func (d *Writer) Close() {
	d.exitCh <- struct{}{}
	d.Wait()
}

func (d *Writer) Push(data *MetaData) {
	tt := data.tt
	str := fmt.Sprintf("%s/%04d%02d%02d/%s%02d.log", d.path, tt.Year(), tt.Month(), tt.Day(), d.name, tt.Hour())
	fileName := path.Clean(str)
	// 判断是否切换文件
	if fileName != d.fileName {
		// 关闭之前的文件
		d.closeFb()
		// 切换文件
		if tmpFb, err := basic.NewFile(fileName); err != nil {
			return
		} else {
			d.fileName = fileName
			d.fb = tmpFb
			d.ww.Reset(tmpFb)
		}
	}
	// 写入缓存
	if d.fb != nil {
		d.ww.WriteString(data.msg)
	}
}

func (d *Writer) Flush() {
	// 判断文件是否已经被删除
	if d.fb != nil && !basic.SameFile(d.fb, d.fileName) {
		d.fb.Close()
		d.fb = nil

		if tmpFb, err := basic.NewFile(d.fileName); err != nil {
			return
		} else {
			d.fb = tmpFb
			d.ww.Reset(tmpFb)
		}
	}
	// 把缓存写入文件
	if d.fb != nil {
		d.ww.Flush()
	}
}

func (d *Writer) run() {
	tt := time.NewTicker(1 * time.Second)
	defer func() {
		for {
			select {
			case item := <-d.dataCh:
				d.Push(item)
			default:
				tt.Stop()
				d.closeFb()
				return
			}
		}
	}()

	for {
		select {
		case <-d.exitCh:
			return
		case <-tt.C:
			if len(d.fileName) > 0 {
				d.Flush()
			}
		case item := <-d.dataCh:
			d.Push(item)
		}
	}
}

func (d *Writer) closeFb() {
	if d.fb != nil {
		d.ww.Flush()
		d.fb.Close()
		d.fb = nil
	}
}
