module github.com/yfuruyama/crzerolog/example/grpc/server

go 1.13

replace github.com/yfuruyama/crzerolog/example/grpc/proto => ../proto

replace github.com/yfuruyama/crzerolog => ../../../

require (
	github.com/rs/zerolog v1.18.0
	github.com/yfuruyama/crzerolog v0.1.1
	github.com/yfuruyama/crzerolog/example/grpc/proto v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.27.1
)
