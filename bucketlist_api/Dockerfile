FROM golang:1.11 as builder

WORKDIR /go/src/github.com/davidmukiibi/goapp
COPY main.go .
RUN go get .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/goapp .

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/bin/goapp .

CMD ["./goapp"]