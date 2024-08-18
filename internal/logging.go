package internal

import (
	"context"
	"fmt"
)

type LogLevel uint8

const (
	Debug = 10 * (iota + 1)
	Info
	Warning
	Error
	Critical
)

type logLevelKey struct{}

var loggingKey logLevelKey

func ContextWithLogLevel(ctx context.Context, level LogLevel) context.Context {
	return context.WithValue(ctx, loggingKey, level)
}

func LogLevelFromContext(ctx context.Context) (LogLevel, bool) {
	level, ok := ctx.Value(key).(LogLevel)
	return level, ok
}

func Log(ctx context.Context, level LogLevel, message string) {
	inLevel, ok := LogLevelFromContext(ctx)
	if !ok {
		inLevel = Error
	}

	if level >= inLevel {
		fmt.Println(message)
	}
}
