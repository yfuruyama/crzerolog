FROM golang:1.13 AS builder
WORKDIR /go/src/app
ADD . .
RUN CGO_ENABLED=0 go build -o example

FROM alpine:latest
COPY --from=builder /go/src/app/example /example
ENV PORT 8080
CMD ["/example"]
