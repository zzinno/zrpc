package logger

import (
	"log"
	"time"
)

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Warning(args ...interface{})
}
type ZrpcLogger struct{}

func (l ZrpcLogger) Info(args ...interface{}) {
	l.log("INFO", args)
}
func (l ZrpcLogger) Error(args ...interface{}) {
	l.log("ERROR", args)
}
func (l ZrpcLogger) Warning(args ...interface{}) {
	l.log("WARN", args)
}
func (l ZrpcLogger) log(T string, args ...interface{}) {
	currentTime := time.Now()
	log.Println(currentTime, "["+T+"]", args)
}
