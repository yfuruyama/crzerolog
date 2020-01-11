// Package crzerolog provides zerolog-based logging for Cloud Run
package crzerolog

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var (
	// CallerSkipFrameCount is the number of stack frames to skip to find the caller.
	CallerSkipFrameCount = 3

	projectID          string
	sourceLocationHook = &callerHook{}
	// TODO: test
	traceHeaderRegExp = regexp.MustCompile(`([0-9a-fA-F]+)(?:/(\d+))?(?:;o=[01])?`)
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

	if isCloudRun() || isAppEngineSecond() {
		// For performance, fetching Project ID here only once,
		// rather than fetching it in every request.
		id, err := fetchProjectIDFromMetadata()
		if err != nil {
			log.Fatalf("Failed to fetch mandatory project ID: %v", err)
		}
		projectID = id
	} else {
		projectID = fetchProjectIDFromEnv()
	}
}

// callerHook implements zerolog.Hook interface.
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

func HandleWithLogger(rootLogger *zerolog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := rootLogger.With().Timestamp().Logger().Hook(sourceLocationHook)

		traceContext := traceHeaderRegExp.FindStringSubmatch(r.Header.Get("X-Cloud-Trace-Context"))
		if len(traceContext) < 3 {
			r = r.WithContext(l.WithContext(r.Context()))
			next.ServeHTTP(w, r)
			return
		}

		traceID := traceContext[1]
		spanID := traceContext[2]
		trace := fmt.Sprintf("projects/%s/traces/%s", projectID, traceID)

		l.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("logging.googleapis.com/trace", trace).Str("logging.googleapis.com/spanId", spanID)
		})

		r = r.WithContext(l.WithContext(r.Context()))
		next.ServeHTTP(w, r)
	})
}

func fetchProjectIDFromMetadata() (string, error) {
	req, err := http.NewRequest("GET",
		"http://metadata.google.internal/computeMetadata/v1/project/project-id", nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Metadata-Flavor", "Google")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func fetchProjectIDFromEnv() string {
	return os.Getenv("GOOGLE_CLOUD_PROJECT")
}
