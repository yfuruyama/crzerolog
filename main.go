package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

func main() {
	port := "8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	rootLogger := zerolog.New(os.Stdout)

	server := http.NewServeMux()
	server.Handle("/", HandleWithLogger(&rootLogger, http.HandlerFunc(hello)))
	log.Printf("Server listening on port %s", port)
	err := http.ListenAndServe(":"+port, server)
	log.Fatal(err)
}

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

// TODO: Support handleFunc
func HandleWithLogger(rootLogger *zerolog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projectId := os.Getenv("PROJECT_ID")
		traceContext := r.Header.Get("X-Cloud-Trace-Context")
		traceID := strings.Split(traceContext, "/")[0]
		trace := fmt.Sprintf("projects/%s/traces/%s", projectId, traceID)

		l := rootLogger.With().Timestamp().Logger()
		l.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("logging.googleapis.com/trace", trace)
		})
		r = r.WithContext(l.WithContext(r.Context()))
		next.ServeHTTP(w, r)
	})
}

func hello(w http.ResponseWriter, r *http.Request) {
	logger := hlog.FromRequest(r)
	logger.Info().Str("additionalField", "test").Msg("Hello structured log 1")
	logger.Warn().Str("additionalField", "test").Msg("Hello structured log 2")

	status := r.URL.Query().Get("status")
	statusCode, err := strconv.Atoi(status)
	if status == "" || err != nil {
		statusCode = 200
	}

	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "Hello! %s", http.StatusText(statusCode))
}
