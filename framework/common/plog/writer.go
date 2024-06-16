package plog

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"sync"
	"time"
	"universal/framework/base"
)

type metaData struct {
	tt  time.Time
	msg []byte
}

type Writer struct {
	sync.WaitGroup
	serverId   int32
	serverName string
	path       string
	fileName   string
	ww         *bufio.Writer
	fb         *os.File
	exitCh     chan struct{}
	metaCh     chan *metaData
}

func NewWriter(id int32, name string, path string) *Writer {
	ret := &Writer{
		serverId:   id,
		serverName: name,
		path:       path,
		ww:         bufio.NewWriter(nil),
		exitCh:     make(chan struct{}, 0),
		metaCh:     make(chan *metaData, 100),
	}
	ret.Add(1)
	go ret.run()
	return ret
}

func (d *Writer) Write(tt time.Time, msg []byte) {
	d.metaCh <- &metaData{tt: tt, msg: msg}
}

func (d *Writer) Close() {
	d.exitCh <- struct{}{}
	d.Wait()
}

func (d *Writer) run() {
	tt := time.NewTicker(1 * time.Second)
	defer func() {
		for {
			select {
			case item := <-d.metaCh:
				d.push(item.tt, item.msg)
			default:
				tt.Stop()     // 关闭定时器
				d.closeFile() // 关闭文件句柄
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
				d.flush()
			}
		case item := <-d.metaCh:
			d.push(item.tt, item.msg)
		}
	}
}

func (d *Writer) closeFile() {
	if d.fb != nil {
		d.ww.Flush()
		d.fb.Close()
		d.fb = nil
	}
	d.Done()
}

func (d *Writer) push(tt time.Time, msg []byte) {
	fileName := getFileName(tt, d.serverId, d.serverName, d.path)
	// 判断是否切换文件
	if d.fileName != fileName {
		// 关闭之前的文件
		if d.fb != nil {
			d.ww.Flush()
			d.fb.Close()
			d.fb = nil
		}
		// 切换文件
		if tmpFb, err := base.NewFile(fileName); err != nil {
			return
		} else {
			d.fileName = fileName
			d.fb = tmpFb
			d.ww.Reset(tmpFb)
		}
	}
	// 写入缓存
	if d.fb != nil {
		d.ww.Write(msg)
	}
}

func (d *Writer) flush() {
	// 判断文件是否已经被删除
	if d.fb != nil && !base.SameFile(d.fb, d.fileName) {
		d.fb.Close()
		d.fb = nil
		if tmpFb, err := base.NewFile(d.fileName); err != nil {
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

func getFileName(tt time.Time, id int32, name, fpath string) string {
	if len(name) > 0 {
		return path.Clean(fmt.Sprintf("%s/%04d%02d%02d/%s-%02d-%02d.log", fpath, tt.Year(), tt.Month(), tt.Day(), name, id, tt.Hour()))
	}
	return path.Clean(fmt.Sprintf("%s/%04d%02d%02d/%02d-%02d.log", fpath, tt.Year(), tt.Month(), tt.Day(), id, tt.Hour()))
}
