package logger

import (
	"context"
	"fmt"
	"os"

	"api-gateway-go/internal/constants"

	"github.com/rs/zerolog"
)

type ZerologLogger struct {
	logger  zerolog.Logger
	traceID string
}

type zerologEvent struct {
	event *zerolog.Event
}

func (e *zerologEvent) Str(key string, val string) Event {
	e.event.Str(key, val)
	return e
}

func (e *zerologEvent) Int(key string, val int) Event {
	e.event.Int(key, val)
	return e
}

func (e *zerologEvent) Err(err error) Event {
	e.event.Err(err)
	return e
}

func (e *zerologEvent) Msg(msg string) {
	e.event.Msg(msg)
}

func (e *zerologEvent) Msgf(format string, args ...interface{}) {
	e.event.Msg(fmt.Sprintf(format, args...))
}

func NewZerologLogger(z zerolog.Logger) *ZerologLogger {
	return &ZerologLogger{logger: z}
}

func New() Logger {
	z := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return NewZerologLogger(z)
}

func (z *ZerologLogger) WithContext(ctx context.Context) Logger {
	traceID, _ := ctx.Value(constants.TraceIDKey).(string)
	return &ZerologLogger{
		logger:  z.logger,
		traceID: traceID,
	}
}

func (z *ZerologLogger) addTraceID(e *zerolog.Event) *zerolog.Event {
	if z.traceID != "" {
		return e.Str("trace_id", z.traceID)
	}
	return e
}

func (z *ZerologLogger) Info() Event {
	return &zerologEvent{event: z.addTraceID(z.logger.Info())}
}

func (z *ZerologLogger) Error() Event {
	return &zerologEvent{event: z.addTraceID(z.logger.Error())}
}

func (z *ZerologLogger) Debug() Event {
	return &zerologEvent{event: z.addTraceID(z.logger.Debug())}
}

func (z *ZerologLogger) Warn() Event {
	return &zerologEvent{event: z.addTraceID(z.logger.Warn())}
}
