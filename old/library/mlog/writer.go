package mlog

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"
	"universal/library/builder"
	"universal/library/queue"
)

type meta struct {
	time time.Time
	buff []byte
}

type LogWriter struct {
	sync.WaitGroup
	lpath   string
	lname   string
	logs    *queue.Queue[*meta]
	builder *builder.Builder
	exit    chan struct{}
	notify  chan struct{}
}

func NewLogWriter(lpath, lname string, size int) *LogWriter {
	w := &LogWriter{
		lpath:   lpath,
		lname:   lname,
		logs:    queue.NewQueue[*meta](),
		notify:  make(chan struct{}, 1),
		exit:    make(chan struct{}),
		builder: builder.NewBuilder(size),
	}
	go w.run()
	return w
}

func (l *LogWriter) Write(tt time.Time, buf []byte) error {
	l.logs.Push(&meta{time: tt, buff: buf})
	select {
	case l.notify <- struct{}{}:
	default:
	}
	return nil
}

func (l *LogWriter) Close() error {
	l.exit <- struct{}{}
	l.Wait()
	return l.builder.Close()
}

func (m *meta) GetFileName(lname string) string {
	return fmt.Sprintf("%s_%04d%02d%02d_%02d.log", lname, m.time.Year(), m.time.Month(), m.time.Day(), m.time.Hour())
}

func (l *LogWriter) run() {
	l.Add(1)
	tt := time.NewTicker(1 * time.Second)
	defer func() {
		tt.Stop()
		l.handle()
		l.Done()
	}()
	for {
		select {
		case <-l.notify:
			l.handle()
		case <-tt.C:
			l.builder.Flush()
		case <-l.exit:
			return
		}
	}
}

func (l *LogWriter) handle() {
	for mm := l.logs.Pop(); mm != nil; mm = l.logs.Pop() {
		l.builder.Set(path.Join(l.lpath, mm.GetFileName(l.lname)))
		l.builder.Write(mm.buff)
	}
}

type StdWriter struct{}

func (d *StdWriter) Write(t time.Time, buf []byte) error {
	_, err := fmt.Fprintln(os.Stdout, t.Format("2006-01-02 15:04:05")+"\t"+string(buf))
	return err
}

func (d *StdWriter) Close() error {
	return nil
}
