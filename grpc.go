package crzerolog

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// TODO: LoggerInterceptor is ...
func LoggerInterceptor(rootLogger *zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		l := rootLogger.With().Timestamp().Logger().Hook(sourceLocationHook)
		ctx = l.WithContext(ctx)

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}
		values := md.Get("x-cloud-trace-context")
		if len(values) != 1 {
			return handler(ctx, req)
		}

		traceID, spanID := traceContextFromHeader(values[0])
		if traceID == "" {
			return handler(ctx, req)
		}
		trace := fmt.Sprintf("projects/%s/traces/%s", projectID, traceID)

		l.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("logging.googleapis.com/trace", trace).Str("logging.googleapis.com/spanId", spanID)
		})

		return handler(ctx, req)
	}
}
