FROM golang:1.17.11-alpine3.16 AS builder
WORKDIR /go/src/github.com/kffl/speedbump/
COPY ./ ./
RUN go get ./
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -v -o speedbump .

FROM alpine:3.16

WORKDIR /root/
COPY --from=builder /go/src/github.com/kffl/speedbump/speedbump .
ENTRYPOINT ["/root/speedbump"]
CMD ["--help"]