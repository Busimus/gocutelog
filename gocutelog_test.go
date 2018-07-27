package gocutelog

import (
	"testing"

	"github.com/francoispqt/onelog"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestZerolog(t *testing.T) {
	w := NewWriter("localhost:19996", "json")
	l := zerolog.New(w)
	l.Info().Msg("Hello world from zerolog!")
}

func TestOnelog(t *testing.T) {
	w := NewWriter("localhost:19996", "json")
	l := onelog.New(w, onelog.ALL)
	l.Info("Hello world from onelog!")
}

func TestLogrus(t *testing.T) {
	w := NewWriter("localhost:19996", "json")
	l := logrus.New()
	l.Out = w
	l.Formatter = new(logrus.JSONFormatter)
	l.Info("Hello world from logrus!")
}

func TestZap(t *testing.T) {
	w := NewWriter("localhost:19996", "json")
	conf := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "exc_info",
		LineEnding:     "",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	enc := zapcore.NewJSONEncoder(conf)
	priority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return true
	})
	core := zapcore.NewCore(enc, w, priority)
	l := zap.New(core)
	l.Info("Hello world from zap!")
}

func TestNotConnected(t *testing.T) {
	w := NewWriter("localhost:19997", "json")
	l := zerolog.New(w)
	l.Info().Msg("Hello?")
	l.Info().Msg("Anyone there?")
}
