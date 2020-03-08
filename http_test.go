package crzerolog

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type logEntry struct {
	Time           string         `json:"time"`
	Severity       string         `json:"severity"`
	SourceLocation sourceLocation `json:"logging.googleapis.com/sourceLocation"`
	Trace          string         `json:"logging.googleapis.com/trace"`
	SpanID         string         `json:"logging.googleapis.com/spanId"`
	Message        string         `json:"message"`
}

type sourceLocation struct {
	File     string `json:"file"`
	Line     string `json:"line"`
	Function string `json:"function"`
}

func TestInjectLogger(t *testing.T) {
	tests := []struct {
		desc        string
		requestFunc func() *http.Request
		handler     http.Handler
		want        logEntry
	}{
		{
			desc: "With X-Cloud-Trace-Context",
			requestFunc: func() *http.Request {
				req, err := http.NewRequest("GET", "/", nil)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				req.Header.Add("X-Cloud-Trace-Context", "0123456789abcdef0123456789abcdef/123;o=1")
				return req
			},
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				logger := log.Ctx(r.Context())
				logger.Info().Msg("hello")
			}),
			want: logEntry{
				Time:     "ignore",
				Severity: "INFO",
				SourceLocation: sourceLocation{
					File:     "http_test.go",
					Line:     "ignore",
					Function: "ignore",
				},
				Trace:   "projects/myproject/traces/0123456789abcdef0123456789abcdef",
				SpanID:  "123",
				Message: "hello",
			},
		},
		{
			desc: "Without X-Cloud-Trace-Context",
			requestFunc: func() *http.Request {
				req, err := http.NewRequest("GET", "/", nil)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				return req
			},
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				logger := log.Ctx(r.Context())
				logger.Debug().Msg("hi") // Debug log is ignored
				logger.Info().Msg("hello")
			}),
			want: logEntry{
				Time:     "ignore",
				Severity: "INFO",
				SourceLocation: sourceLocation{
					File:     "http_test.go",
					Line:     "ignore",
					Function: "ignore",
				},
				Trace:   "",
				SpanID:  "",
				Message: "hello",
			},
		},
	}

	for _, tt := range tests {
		projectID = "myproject"
		buf := &bytes.Buffer{}
		rootLogger := zerolog.New(buf)
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		resprec := httptest.NewRecorder()

		InjectLogger(&rootLogger)(tt.handler).ServeHTTP(resprec, tt.requestFunc())

		var got logEntry
		if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		opt := cmpopts.IgnoreFields(logEntry{}, "Time", "SourceLocation.Line", "SourceLocation.Function")
		if diff := cmp.Diff(tt.want, got, opt); diff != "" {
			t.Errorf("%s: Log output diff: %s", tt.desc, diff)
		}
	}
}
