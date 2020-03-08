FROM golang:1.13 AS builder
WORKDIR /go/src/github.com/yfuruyama/crzerolog
ADD . .
# build server
RUN cd example/grpc/server && CGO_ENABLED=0 go build -o server && cd ../../
# build client
RUN cd example/grpc/client && CGO_ENABLED=0 go build -o client && cd ../../

FROM alpine:latest
RUN apk update && apk add ca-certificates && update-ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder /go/src/github.com/yfuruyama/crzerolog/example/grpc/server/server /server
COPY --from=builder /go/src/github.com/yfuruyama/crzerolog/example/grpc/client/client /client
ENV PORT 8080
CMD ["/server"]
