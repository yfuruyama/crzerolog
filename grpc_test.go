package crzerolog

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestInjectLoggerInterceptor(t *testing.T) {
	tests := []struct {
		desc    string
		md metadata.MD
		handler func(context.Context, interface{}) (interface{}, error)
		want logEntry
	}{
		{
			desc: "With x-cloud-trace-context",
			md: metadata.Pairs("x-cloud-trace-context", "0123456789abcdef0123456789abcdef/123;o=1"),
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				logger := log.Ctx(ctx)
				logger.Debug().Msg("hi") // Debug log is ignored
				logger.Info().Msg("hello")
				return nil, nil
			},
			want: logEntry{
				Time:     "ignore",
				Severity: "INFO",
				SourceLocation: sourceLocation{
					File:     "grpc_test.go",
					Line:     "ignore",
					Function: "ignore",
				},
				Trace:   "projects/myproject/traces/0123456789abcdef0123456789abcdef",
				SpanID:  "123",
				Message: "hello",
			},
		},
		{
			desc: "Without x-cloud-trace-context",
			md: metadata.New(nil),
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				logger := log.Ctx(ctx)
				logger.Debug().Msg("hi") // Debug log is ignored
				logger.Info().Msg("hello")
				return nil, nil
			},
			want: logEntry{
				Time:     "ignore",
				Severity: "INFO",
				SourceLocation: sourceLocation{
					File:     "grpc_test.go",
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

		unaryInfo := &grpc.UnaryServerInfo{
			FullMethod: "TestService.TestMethod",
		}

		ctx := context.Background()
		ctx = metadata.NewIncomingContext(ctx, tt.md)
		interceptor := InjectLoggerInterceptor(&rootLogger)
		interceptor(ctx, nil, unaryInfo, tt.handler)

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
