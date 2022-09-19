package crzerolog

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// InjectLoggerInterceptor returns a gRPC unary interceptor for injecting zerolog.Logger to the RPC invocation context.
func InjectLoggerInterceptor(rootLogger *zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx = rootLogger.With().Timestamp().Logger().Hook(sourceLocationHook).WithContext(ctx)

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}
		values := md.Get("x-cloud-trace-context")
		if len(values) != 1 {
			return handler(ctx, req)
		}

		traceID, _ := traceContextFromHeader(values[0])
		if traceID == "" {
			return handler(ctx, req)
		}
		trace := fmt.Sprintf("projects/%s/traces/%s", projectID, traceID)

		log.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("logging.googleapis.com/trace", trace)
		})

		return handler(ctx, req)
	}
}
