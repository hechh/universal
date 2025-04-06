package test

import (
	"hego/Library/ulog"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// BenchmarkError-8  928797	      1089 ns/op	     944 B/op	      17 allocs/o
func BenchmarkError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ulog.Error("--%d-------askfja;lksdjf;alkdjf;alkjdf;alsjkdf-------", i)
	}
	b.Log(b.N)
}

// BenchmarkLogrus-8   	   21776	     60060 ns/op	     472 B/op	      15 allocs/op
func BenchmarkLogrus(b *testing.B) {
	logger := logrus.New()
	logger.Out = os.Stdout
	logger.Level = logrus.InfoLevel

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("This is a log message")
	}
}

// BenchmarkZap-8   	 1340110	       846.7 ns/op	       2 B/op	       0 allocs/op
func BenchmarkZap(b *testing.B) {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	logger, _ := config.Build()
	defer logger.Sync()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("This is a log message")
	}
}

// BenchmarkZerolog-8   	   21171	     53735 ns/op	       0 B/op	       0 allocs/o
func BenchmarkZerolog(b *testing.B) {
	logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info().Msg("This is a log message")
	}
}
