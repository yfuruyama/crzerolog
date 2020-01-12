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
	"github.com/rs/zerolog/hlog"
)

type stackdriverLog struct {
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

func TestServeHTTP(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	projectID = "myproject"
	req.Header.Add("X-Cloud-Trace-Context", "0123456789abcdef0123456789abcdef/123;o=1")

	buf := &bytes.Buffer{}
	rootLogger := zerolog.New(buf)
	middleware := InjectLogger(&rootLogger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := hlog.FromRequest(r)
		logger.Info().Msg("hello")
	}))
	resprec := httptest.NewRecorder()

	middleware.ServeHTTP(resprec, req)

	var got stackdriverLog
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	want := stackdriverLog{
		Time:     "ignore",
		Severity: "INFO",
		SourceLocation: sourceLocation{
			File:     "log_test.go",
			Line:     "ignore",
			Function: "ignore",
		},
		Trace:   "projects/myproject/traces/0123456789abcdef0123456789abcdef",
		SpanID:  "123",
		Message: "hello",
	}

	opt := cmpopts.IgnoreFields(stackdriverLog{}, "Time", "SourceLocation.Line", "SourceLocation.Function")
	if diff := cmp.Diff(got, want, opt); diff != "" {
		t.Errorf("Log output diff: %s", diff)
	}
}

func TestTraceContextFromHeader(t *testing.T) {
	for _, tt := range []struct {
		header      string
		wantTraceID string
		wantSpanID  string
	}{
		{"0123456789abcdef0123456789abcdef/123;o=1", "0123456789abcdef0123456789abcdef", "123"},
		{"0123456789abcdef0123456789abcdef/123;o=0", "0123456789abcdef0123456789abcdef", "123"},
		{"0123456789abcdef0123456789abcdef/123", "0123456789abcdef0123456789abcdef", "123"},
		{"0123456789abcdef0123456789abcdef", "0123456789abcdef0123456789abcdef", ""},
		{"0123456789abcdef0123456789abcdef/invalid", "", ""},
		{"invalid", "", ""},
		{"", "", ""},
	} {
		traceID, spanID := traceContextFromHeader(tt.header)
		if traceID != tt.wantTraceID || spanID != tt.wantSpanID {
			t.Errorf("traceContextFromHeader(%q) = (%q, %q), want = (%q, %q)", tt.header, traceID, spanID, tt.wantTraceID, tt.wantSpanID)
		}
	}
}
