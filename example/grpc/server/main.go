package main

import (
	"context"
	"github.com/yfuruyama/crzerolog"
	"net"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	pb "github.com/yfuruyama/crzerolog/example/grpc/proto"
)

type server struct{}

func (s *server) Echo(ctx context.Context, r *pb.EchoRequest) (*pb.EchoReply, error) {
	logger := log.Ctx(ctx)

	logger.Info().Msg("Hi")
	logger.Warn().Str("foo", "bar").Msg("This is")
	logger.Error().Int("num", 123).Msg("Structured Log")

	return &pb.EchoReply{Msg: r.GetMsg() + "!"}, nil
}

func main() {
	port := "8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal().Msgf("Failed to listen: %v", err)
	}

	rootLogger := zerolog.New(os.Stdout)
	s := grpc.NewServer(
		grpc.UnaryInterceptor(crzerolog.InjectLoggerInterceptor(&rootLogger)),
	)
	pb.RegisterHelloServer(s, &server{})
	if err := s.Serve(l); err != nil {
		log.Fatal().Msgf("Failed to serve: %v", err)
	}
}
