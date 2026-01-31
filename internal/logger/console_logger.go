package logger

import (
	"context"
	"log"
	"time"
)

type ConsoleLogger struct{}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (l *ConsoleLogger) Debug(ctx context.Context, msg string, fields map[string]interface{}) {
	l.log("DEBUG", msg, fields)
}

func (l *ConsoleLogger) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	l.log("INFO", msg, fields)
}

func (l *ConsoleLogger) Error(ctx context.Context, msg string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	if err != nil {
		fields["error"] = err.Error()
	}
	l.log("ERROR", msg, fields)
}

func (l *ConsoleLogger) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	l.log("WARN", msg, fields)
}

func (l *ConsoleLogger) log(level, msg string, fields map[string]interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("[%s] %s - %s - %v\n", level, timestamp, msg, fields)
}
