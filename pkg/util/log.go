package util

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"net/http"
	"os"
)

type EcsLogger struct {
	log.Logger
}

// LogHTTP implements the Logger interface of the Scaleway API.
func (l *EcsLogger) LogHTTP(r *http.Request) {
	_ = level.Debug(l).Log("msg", "HTTP request", "method", r.Method, "url", r.URL.String())
}

// Fatalf implements the Logger interface of the Scaleway API.
func (l *EcsLogger) Fatalf(format string, v ...interface{}) {
	_ = level.Error(l).Log("msg", fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Debugf implements the Logger interface of the Scaleway API.
func (l *EcsLogger) Debugf(format string, v ...interface{}) {
	_ = level.Debug(l).Log("msg", fmt.Sprintf(format, v...))
}

// Infof implements the Logger interface of the Scaleway API.
func (l *EcsLogger) Infof(format string, v ...interface{}) {
	_ = level.Info(l).Log("msg", fmt.Sprintf(format, v...))
}

// Warnf implements the Logger interface of the Scaleway API.
func (l *EcsLogger) Warnf(format string, v ...interface{}) {
	_ = level.Warn(l).Log("msg", fmt.Sprintf(format, v...))
}

// Warnf implements the Logger interface of the promhttp package.
func (l *EcsLogger) Println(v ...interface{}) {
	_ = level.Error(l).Log("msg", fmt.Sprintln(v...))
}
