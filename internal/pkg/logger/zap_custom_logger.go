package logger

import (
	"fmt"
	"runtime"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap/zapcore"
)

type WarehouseZapCore struct {
	zapcore.Core
}

func (w *WarehouseZapCore) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if w.Enabled(entry.Level) {
		return checked.AddCore(entry, w)
	}
	return checked
}

func (w *WarehouseZapCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	ftr := []zapcore.Field{} // fields to return

	switch entry.Level {
	case zapcore.InfoLevel:
		ftr = fields
		entry.Caller = zapcore.EntryCaller{Defined: false}
	case zapcore.ErrorLevel:
		defer sentry.Flush(2 * time.Second)

		method := "unknown method"
		reason := "unknown reason"
		caller := zapcore.EntryCaller{Defined: false}
		allFrames := []runtime.Frame{}
		for _, field := range fields {
			switch field.Key {
			case logMethodField:
				method = field.String
				ftr = append(ftr, field)
			case logReasonField:
				reason = field.String
				ftr = append(ftr, field)
			case logFullStacktraceField:
				frames, ok := field.Interface.([]runtime.Frame)
				if ok {
					allFrames = frames
				}
			case logCallerField:
				c, ok := field.Interface.(zapcore.EntryCaller)
				if ok {
					caller = c
				}
			}
		}

		stacktrace := sentry.NewStacktrace()
		sentryFrames := []sentry.Frame{}
		for _, f := range allFrames {
			sentryFrames = append(sentryFrames, sentry.NewFrame(f))
		}
		stacktrace.Frames = sentryFrames

		typeField := fmt.Sprintf("[%s](%s) %s", entry.LoggerName, method, reason)
		event := sentry.NewEvent()
		event.Level = sentry.LevelError
		event.Exception = append(event.Exception, sentry.Exception{
			Type:       typeField,
			Value:      entry.Message,
			Module:     entry.LoggerName,
			Stacktrace: stacktrace,
		})

		sentry.CaptureEvent(event)

		entry.Caller = caller
	default:
		ftr = fields
	}

	return w.Core.Write(entry, ftr)
}
