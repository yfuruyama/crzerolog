// Package crzerolog provides a zerolog-based logger for Cloud Run.
package crzerolog

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var (
	// CallerSkipFrameCount is the number of stack frames to skip to find the caller.
	CallerSkipFrameCount = 3

	projectID          string
	sourceLocationHook = &callerHook{}
	// For trace header, see https://cloud.google.com/trace/docs/troubleshooting#force-trace
	traceHeaderRegExp = regexp.MustCompile(`^\s*([0-9a-fA-F]+)(?:/(\d+))?(?:;o=[01])?\s*$`)
)

func init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.LevelFieldName = "severity"
	zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
		// mapping to Cloud Logging LogSeverity
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

// Run adds sourceLocation for the log to zerolog.Event.
func (h *callerHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	var file, line, function string
	if pc, filePath, lineNum, ok := runtime.Caller(CallerSkipFrameCount); ok {
		if f := runtime.FuncForPC(pc); f != nil {
			function = f.Name()
		}
		line = fmt.Sprintf("%d", lineNum)
		parts := strings.Split(filePath, "/")
		file = parts[len(parts)-1]
	}
	e.Dict("logging.googleapis.com/sourceLocation",
		zerolog.Dict().Str("file", file).Str("line", line).Str("function", function))
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

func traceContextFromHeader(header string) (string, string) {
	matched := traceHeaderRegExp.FindStringSubmatch(header)
	if len(matched) < 3 {
		return "", ""
	}

	traceID, spanID := matched[1], matched[2]
	if spanID == "" {
		return traceID, ""
	}
	spanIDInt, err := strconv.ParseUint(spanID, 10, 64)
	if err != nil {
		// invalid
		return "", ""
	}
	// spanId for cloud logging must be 16-character hexadecimal number.
	// See: https://cloud.google.com/trace/docs/trace-log-integration#associating
	spanIDHex := fmt.Sprintf("%016x", spanIDInt)
	return traceID, spanIDHex
}
