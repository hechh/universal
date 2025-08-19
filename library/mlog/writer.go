package mlog

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"
	"universal/library/queue"
)

type StdWriter struct{}

func (d *StdWriter) Write(m *meta) error {
	defer put(m)
	_, err := fmt.Fprintln(os.Stdout, m.tt.Format("2006-01-02 15:04:05")+"\t"+m.buf.String())
	return err
}

func (d *StdWriter) Close() error {
	return nil
}

type LogWriter struct {
	sync.WaitGroup
	lpath   string
	lname   string
	logs    *queue.Queue[*meta]
	builder *Builder
	exit    chan struct{}
	notify  chan struct{}
}

func NewLogWriter(lpath, lname string, size int) *LogWriter {
	w := &LogWriter{
		lpath:   lpath,
		lname:   lname,
		logs:    queue.New[*meta](),
		notify:  make(chan struct{}, 1),
		exit:    make(chan struct{}),
		builder: NewBuilder(size),
	}
	w.Add(1)
	go w.run()
	return w
}

func (l *LogWriter) Write(mdata *meta) error {
	l.logs.Push(mdata)
	select {
	case l.notify <- struct{}{}:
	default:
	}
	return nil
}

func (l *LogWriter) Close() error {
	close(l.exit)
	l.Wait()
	l.builder.Flush()
	return l.builder.Close()
}

func (l *LogWriter) GetFileName(m *meta) string {
	return path.Join(l.lpath, fmt.Sprintf("%s_%04d%02d%02d_%02d.log", l.lname, m.tt.Year(), m.tt.Month(), m.tt.Day(), m.tt.Hour()))
}

func (l *LogWriter) run() {
	tt := time.NewTicker(1 * time.Second)
	defer func() {
		tt.Stop()
		for mm := l.logs.Pop(); mm != nil; mm = l.logs.Pop() {
			l.builder.Set(l.GetFileName(mm))
			l.builder.Write(mm.buf.Bytes())
			put(mm)
		}
		l.Done()
	}()
	for {
		select {
		case <-l.notify:
			for mm := l.logs.Pop(); mm != nil; mm = l.logs.Pop() {
				l.builder.Set(l.GetFileName(mm))
				l.builder.Write(mm.buf.Bytes())
				put(mm)
			}
		case <-tt.C:
			l.builder.Flush()
		case <-l.exit:
			return
		}
	}
}
