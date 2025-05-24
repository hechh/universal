package mlog

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"
	"universal/framework/library/async"
	"universal/framework/library/builder"
	"universal/framework/library/util"
)

type StdWriter struct{}

func (d *StdWriter) Write(f IFormat) error {
	_, err := fmt.Fprintln(os.Stdout, f.GetString())
	return err
}

func (d *StdWriter) Close() error {
	return nil
}

type Writer struct {
	sync.WaitGroup
	logPath string
	logName string
	cache   *async.Queue
	exit    chan struct{}
	notify  chan struct{}
	w       *builder.FileBuilder
}

func NewWriter(logPath, logName string, size int) *Writer {
	w := &Writer{
		logPath: logPath,
		logName: logName,
		cache:   async.NewQueue(),
		exit:    make(chan struct{}, 1),
		notify:  make(chan struct{}, 1),
		w:       builder.NewFileBuilder(size),
	}
	go w.run()
	return w
}

func (l *Writer) Write(f IFormat) error {
	l.cache.Push(f)
	select {
	case l.notify <- struct{}{}:
	default:
	}
	return nil
}

func (l *Writer) Close() error {
	l.exit <- struct{}{}
	l.Add(1)
	l.Wait()
	return nil
}

func (l *Writer) run() {
	tt := time.NewTicker(1 * time.Second)
	defer func() {
		tt.Stop()
		l.do()
		l.w.Flush()
		l.w.Close()
		l.Done()
	}()

	for {
		select {
		case <-l.exit:
			return
		case <-tt.C:
			l.w.Flush()
		case <-l.notify:
			l.do()
		}
	}
}

func (l *Writer) do() {
	for vv := l.cache.Pop(); vv != nil; vv = l.cache.Pop() {
		item, ok := vv.(IFormat)
		if !ok || item == nil {
			continue
		}

		tt := item.GetTime()
		filename := fmt.Sprintf("%s_%04d%02d%02d_%02d.log", l.logName, tt.Year(), tt.Month(), tt.Day(), tt.Hour())
		l.w.SetWriter(path.Join(l.logPath, filename))
		l.w.Write(util.StringToBytes(item.GetString()))
	}
}
