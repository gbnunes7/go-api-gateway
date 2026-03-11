package logger

import "context"

type Logger interface {
	WithContext(ctx context.Context) Logger
	Info() Event
	Error() Event
	Debug() Event
	Warn() Event
}

type Event interface {
	Str(key string, val string) Event
	Int(key string, val int) Event
	Err(err error) Event
	Msg(msg string)
	Msgf(format string, args ...interface{})
}
