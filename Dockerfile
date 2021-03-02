FROM golang:1.15.3 as builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go clean --modcache
RUN GOOS=linux CGO_ENABLED=0 go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /etc/go-sro-agent-server
COPY --from=builder /app/main .
CMD ["etc/go-sro-agent-server/main"]