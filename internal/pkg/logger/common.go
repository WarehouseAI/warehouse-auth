package logger

import (
	"runtime"
	"strings"
	"unicode"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getMethod(frames []runtime.Frame) string {
	for _, frame := range frames {
		functionName := getFunctionName(frame)
		if unicode.IsUpper(rune(functionName[0])) {
			return functionName
		}
	}

	return getFunctionName(frames[0])
}

func getFunctionName(frame runtime.Frame) string {
	return frame.Function[strings.LastIndex(frame.Function, ".")+1:]
}

func getFunctionsCallTrace(skip int) []runtime.Frame {
	const stackSize = 64
	var (
		callers = make([]uintptr, stackSize)
		n       = runtime.Callers(0, callers)
		frames  = runtime.CallersFrames(callers[:n])
		toSkip  = 2 + skip
	)

	var resultFrames []runtime.Frame
	for frame, hasNext := frames.Next(); hasNext; frame, hasNext = frames.Next() {
		if strings.HasPrefix(frame.Function, "github.com/custom-app/coffee") {
			if strings.Contains(frame.File, "/http/handler.go") {
				break
			}
			resultFrames = append(resultFrames, frame)
		}
	}

	return resultFrames[toSkip:]
}

func setZapMethod(frames []runtime.Frame) zap.Field {
	method := getMethod(frames)
	return zap.String(logMethodField, method)
}
func zapReason(reason string) zap.Field {
	return zap.String(logReasonField, reason)
}
func zapStacktrace(frames []runtime.Frame) zap.Field {
	return zap.Uintptr(logStacktraceField, frames[0].PC)
}
func zapFullStacktrace(frames []runtime.Frame) zap.Field {
	return zap.Reflect(logFullStacktraceField, frames)
}
func zapCaller(frames []runtime.Frame) zap.Field {
	f := frames[0]
	return zap.Reflect(logCallerField, zapcore.EntryCaller{
		Defined:  true,
		PC:       f.PC,
		File:     f.File,
		Line:     f.Line,
		Function: f.Function,
	})
}

func formError(err error, reason string) (string, zap.Field, zap.Field, zap.Field, zap.Field, zap.Field) {
	frames := getFunctionsCallTrace(1)

	return err.Error(), setZapMethod(frames), zapReason(reason), zapStacktrace(frames), zapCaller(frames), zapFullStacktrace(frames)
}

// func getSentryCallerFrames(pc uintptr) []sentry.Frame {
// 	var frames = []sentry.Frame{}
// 	callersFrames := runtime.CallersFrames([]uintptr{pc})

// 	for {
// 		callerFrame, more := callersFrames.Next()
// 		frames = append(frames, sentry.NewFrame(callerFrame))
// 		if !more {
// 			break
// 		}
// 	}

// 	for i, j := 0, len(frames)-1; i < j; i, j = i+1, j-1 {
// 		frames[i], frames[j] = frames[j], frames[i]
// 	}

// 	return frames
// }
