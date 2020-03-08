module github.com/yfuruyama/crzerolog/example/grpc/client

go 1.13

replace github.com/yfuruyama/crzerolog/example/grpc/proto => ../proto

require (
	github.com/yfuruyama/crzerolog/example/grpc/proto v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.27.1
)
