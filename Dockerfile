FROM golang:1.20-alpine AS builder

WORKDIR /go/src/github.com/at-wat/switchweb

COPY go.mod go.sum /go/src/github.com/at-wat/switchweb/
RUN go mod download

COPY . /go/src/github.com/at-wat/switchweb
ENV CGO_ENABLED=0
RUN go build -ldflags "-s -w" .

RUN apk add --no-cache ca-certificates && update-ca-certificates

FROM scratch

COPY --from=builder /go/src/github.com/at-wat/switchweb/switchweb /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV SWITCHBOT_TOKEN= \
  SWITCHBOT_CLIENT_SECRET= \
  ADDR=

EXPOSE 8080

ENTRYPOINT ["/switchweb"]
