package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/yfuruyama/crzerolog/example/grpc/proto"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:8080", "The server address in the format of host:port")
)

func main() {
	flag.Parse()

	creds, err := credentials.NewClientTLSFromFile("/etc/ssl/certs/ca-certificates.crt", "")
	if err != nil {
		log.Fatalf("Failed to load credentials: %v", err)
	}

	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	ctx := context.Background()
	client := pb.NewHelloClient(conn)
	rep, err := client.Echo(ctx, &pb.EchoRequest{Msg: "Hello Cloud Run"})
	if err != nil {
		log.Fatalf("Failed to request: %v", err)
	}
	fmt.Printf("Echo reply: %q\n", rep.GetMsg())
}
