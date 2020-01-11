FROM golang:1.13 AS builder
WORKDIR /go/src/app
ADD . .
RUN CGO_ENABLED=0 go build -o main

FROM alpine:latest
COPY --from=builder /go/src/app/main /main
ENV PORT 8080
CMD ["/main"]
