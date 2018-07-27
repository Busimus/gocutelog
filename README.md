# gocutelog – bridge between Go logging libraries and cutelog
[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/busimus/gocutelog) [![license](https://img.shields.io/badge/license-MIT-green.svg?style=flat-square)](https://raw.githubusercontent.com/busimus/gocutelog/master/LICENSE)

This Go package makes it possible to send log records from Go logging libraries
to a [cutelog](https://github.com/busimus/cutelog) instance without having to
manually manage a socket connection.

Function NewWriter returns a struct that implements io.Writer interface so it
can be used as output by libraries like zerolog, zap, onelog, logrus, etc.

Just like cutelog itself, this package is meant to be used only during
development, so performance or reliability are not the focus here.

## Usage
### zerolog
```go
package main

import (
    "github.com/busimus/gocutelog"
    "github.com/rs/zerolog"
)

func main() {
	w := gocutelog.NewWriter("localhost:19996", "json")
	l := zerolog.New(w)
	l.Info().Msg("Hello world from zerolog!")
}
```

### onelog
```go
package main

import (
    "github.com/busimus/gocutelog"
    "github.com/francoispqt/onelog"
)

func main() {
	w := gocutelog.NewWriter("localhost:19996", "json")
	l := onelog.New(w, onelog.ALL)
	l.Info("Hello world from onelog!")
}
```

### logrus
```go
package main

import (
    "github.com/busimus/gocutelog"
    "github.com/sirupsen/logrus"
)

func main() {
	w := gocutelog.NewWriter("localhost:19996", "json")
	l := logrus.New()
	l.Out = w
	l.Formatter = new(logrus.JSONFormatter)
	l.Info("Hello world from logrus!")
}
```

### zap
```go
package main

import (
    "github.com/busimus/gocutelog"
	"go.uber.org/zap"
   	"go.uber.org/zap/zapcore"
)

func main() {
	w := gocutelog.NewWriter("localhost:19996", "json")
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
```

## License
Released under the [MIT license](https://raw.githubusercontent.com/busimus/gocutelog/master/LICENSE).
