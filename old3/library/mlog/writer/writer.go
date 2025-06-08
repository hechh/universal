package writer

import (
	"fmt"
	"path"
	"sync"
	"time"
	"universal/library/async"
)

type Data struct {
	t   time.Time
	buf []byte
}

type Writer struct {
	sync.WaitGroup
	logPath string
	logName string
	logs    *async.Queue[*Data]
	builder *FileBuilder
	exit    chan struct{}
	notify  chan struct{}
}

func NewWriter(logPath, logName string, size int) *Writer {
	w := &Writer{
		logPath: logPath,
		logName: logName,
		logs:    async.NewQueue[*Data](),
		exit:    make(chan struct{}, 1),
		notify:  make(chan struct{}, 1),
		builder: NewFileBuilder(size),
	}
	go w.run()
	return w
}

func (l *Writer) Write(tt time.Time, buf []byte) error {
	l.logs.Push(&Data{t: tt, buf: buf})
	select {
	case l.notify <- struct{}{}:
	default:
	}
	return nil
}

func (l *Writer) Close() error {
	l.Add(1)
	l.exit <- struct{}{}
	l.Wait()
	return l.builder.Close()
}

func (l *Writer) run() {
	tt := time.NewTicker(1 * time.Second)
	defer func() {
		tt.Stop()
		for mm := l.logs.Pop(); mm != nil; mm = l.logs.Pop() {
			filename := fmt.Sprintf("%s_%04d%02d%02d_%02d.log", l.logName, mm.t.Year(), mm.t.Month(), mm.t.Day(), mm.t.Hour())
			l.builder.SetFile(path.Join(l.logPath, filename))
			l.builder.Write(mm.buf)
		}
		l.builder.Close()
		l.Done()
	}()

	for {
		select {
		case <-l.notify:
			for mm := l.logs.Pop(); mm != nil; mm = l.logs.Pop() {
				filename := fmt.Sprintf("%s_%04d%02d%02d_%02d.log", l.logName, mm.t.Year(), mm.t.Month(), mm.t.Day(), mm.t.Hour())
				l.builder.SetFile(path.Join(l.logPath, filename))
				l.builder.Write(mm.buf)
			}
		case <-tt.C:
			l.builder.Flush()
		case <-l.exit:
			return
		}
	}
}
