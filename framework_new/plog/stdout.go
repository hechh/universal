package plog

import (
	"bufio"
	"os"
	"sync"
	"time"
)

type Stdout struct {
	sync.WaitGroup
	ww     *bufio.Writer
	exitCh chan struct{}
	byteCh chan []byte
}

func NewStdout() *Stdout {
	ret := &Stdout{
		ww:     bufio.NewWriter(os.Stdout),
		exitCh: make(chan struct{}, 0),
		byteCh: make(chan []byte, 50),
	}
	ret.Add(1)
	go ret.run()
	return ret
}

func (d *Stdout) Write(msg []byte) {
	d.byteCh <- msg
}

func (d *Stdout) Close() {
	d.exitCh <- struct{}{}
	d.Wait()
}

func (d *Stdout) run() {
	tt := time.NewTicker(1 * time.Second)
	defer func() {
		for {
			select {
			case item := <-d.byteCh:
				d.ww.Write(item)
			default:
				tt.Stop()
				d.ww.Flush()
				d.Done()
				return
			}
		}
	}()
	for {
		select {
		case <-d.exitCh:
			return
		case <-tt.C:
			d.ww.Flush()
		case item := <-d.byteCh:
			d.ww.Write(item)
		}
	}
}
