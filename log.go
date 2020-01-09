package cloudrunlog

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var (
	// CallerSkipFrameCount is the number of stack frames to skip to find the caller.
	CallerSkipFrameCount = 3

	sourceLocationHook = &callerHook{}
)

func init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.LevelFieldName = "severity"
	zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
		// mapping to Stackdriver LogSeverity
		// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogSeverity
		switch l {
		case zerolog.TraceLevel:
			return "DEFAULT"
		case zerolog.DebugLevel:
			return "DEBUG"
		case zerolog.InfoLevel:
			return "INFO"
		case zerolog.WarnLevel:
			return "WARNING"
		case zerolog.ErrorLevel:
			return "ERROR"
		case zerolog.FatalLevel:
			return "CRITICAL"
		case zerolog.PanicLevel:
			return "ALERT"
		case zerolog.NoLevel:
			return "DEFAULT"
		default:
			return "DEFAULT"
		}
	}
}

type callerHook struct{}

func (h *callerHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	var file, line, function string
	if pc, filePath, lineNum, ok := runtime.Caller(CallerSkipFrameCount); ok {
		if f := runtime.FuncForPC(pc); f != nil {
			function = f.Name()
		}
		line = fmt.Sprintf("%d", lineNum)
		parts := strings.Split(filePath, "/")
		file = parts[len(parts)-1] // use short file name
	}
	e.Dict("logging.googleapis.com/sourceLocation", zerolog.Dict().Str("file", file).Str("line", line).Str("function", function))
}

// TODO: Support handleFunc
func HandleWithLogger(rootLogger *zerolog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: platform-dependent
		projectId := os.Getenv("PROJECT_ID")
		traceContext := r.Header.Get("X-Cloud-Trace-Context")
		traceID := strings.Split(traceContext, "/")[0]
		trace := fmt.Sprintf("projects/%s/traces/%s", projectId, traceID)

		l := rootLogger.With().Timestamp().Logger().Hook(sourceLocationHook)
		l.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("logging.googleapis.com/trace", trace)
		})
		r = r.WithContext(l.WithContext(r.Context()))
		next.ServeHTTP(w, r)
	})
}
