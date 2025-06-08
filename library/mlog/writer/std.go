package writer

import (
	"fmt"
	"os"
	"time"
)

type StdWriter struct{}

func NewStdWriter() *StdWriter {
	return &StdWriter{}
}

func (d *StdWriter) Write(t time.Time, buf []byte) error {
	_, err := fmt.Fprintln(os.Stdout, t.Format("2006-01-02 15:04:05")+string(buf))
	return err
}

func (d *StdWriter) Close() error {
	return nil
}
