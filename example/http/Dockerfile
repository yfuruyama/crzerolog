FROM golang:1.13 AS builder
WORKDIR /go/src/github.com/yfuruyama/crzerolog
ADD . .
RUN cd example/http && CGO_ENABLED=0 go build -o http && cd ../../

FROM alpine:latest
COPY --from=builder /go/src/github.com/yfuruyama/crzerolog/example/http/http /http
ENV PORT 8080
CMD ["/http"]
